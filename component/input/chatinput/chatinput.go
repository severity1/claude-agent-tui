package chatinput

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/style"
)

// reduceMotion checks if the user prefers reduced motion.
// Respects the REDUCE_MOTION environment variable for accessibility.
func reduceMotion() bool {
	return os.Getenv("REDUCE_MOTION") != ""
}

// DefaultHistorySize is the default maximum number of history entries.
const DefaultHistorySize = 100

// SubmitMsg is emitted when the user submits the input (presses Enter).
type SubmitMsg struct {
	Text string
}

// CopiedMsg is emitted when the input content is copied to clipboard.
// Uses OSC 52 escape sequences which work cross-platform in modern terminals
// including over SSH, tmux, and Wayland sessions.
//
// Note: Err captures write failures to stderr, but a nil error does NOT guarantee
// the clipboard was updated - OSC 52 is "fire and forget" and terminal support varies.
type CopiedMsg struct {
	Text string
	Err  error
}

// flashDoneMsg signals the end of the copy flash animation.
type flashDoneMsg struct{}

// DefaultPrompt is the default prompt prefix shown before the input.
// Empty by default; use WithPrompt() to add a prompt prefix.
const DefaultPrompt = ""

// Model represents the chat input component state.
type Model struct {
	textarea    textarea.Model
	focused     bool
	history     []string
	historyIdx  int // -1 means "new input", 0+ means position in history (0 = most recent)
	historySize int
	draft       string // saved draft when navigating history
	showCounter bool   // whether to show character/line count
	width       int    // stored width for counter display
	prompt      string // prompt prefix shown before textarea
	promptStyle lipgloss.Style
	minHeight   int            // minimum visible lines (default: 1)
	maxHeight   int            // maximum lines before scroll (default: 10)
	mode        string         // mode indicator (e.g., "Build", "Shell") - left of infoBar
	modeStyle   lipgloss.Style // style for mode indicator
	infoBar     string         // bottom info text (e.g., "Claude Sonnet 4.5")
	infoStyle   lipgloss.Style // style for info bar
	flashing    bool           // true during copy flash animation
	placeholder string         // placeholder text (rendered by us, not bubbles textarea)

	// Visual styling options
	contentBg           lipgloss.Color  // background color for all content inside borders
	borderStyle         lipgloss.Border // border style (Thick, Rounded, Normal, etc.)
	borderColor         lipgloss.Color  // border color when not flashing
	flashBorderColor    lipgloss.Color  // border color during copy flash
	flashDuration       time.Duration   // duration of copy flash animation
	borderPaddingTop    int             // top padding inside border
	borderPaddingBottom int             // bottom padding inside border
}

// Option configures the Model.
type Option func(*Model)

// configureTextareaStyle sets consistent styling for a textarea style struct.
func configureTextareaStyle(s *textarea.Style, fg, muted, bg lipgloss.TerminalColor) {
	s.Base = lipgloss.NewStyle().Foreground(fg).Background(bg)
	s.CursorLine = lipgloss.NewStyle().Background(bg)
	s.EndOfBuffer = lipgloss.NewStyle().Background(bg)
	s.Placeholder = lipgloss.NewStyle().Foreground(muted).Background(bg)
	s.Prompt = lipgloss.NewStyle().Background(bg)
	s.Text = lipgloss.NewStyle().Foreground(fg).Background(bg)
}

// New creates a new ChatInput model with the given options.
func New(opts ...Option) Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Prompt = "" // Hide built-in prompt, we render our own externally

	// Style textarea with Surface background to match Place() container.
	// All inner content uses consistent Surface background for seamless appearance.
	textColor := style.Text
	mutedColor := style.TextMuted
	bgColor := style.Surface

	configureTextareaStyle(&ta.BlurredStyle, textColor, mutedColor, bgColor)
	configureTextareaStyle(&ta.FocusedStyle, textColor, mutedColor, bgColor)

	m := Model{
		textarea:    ta,
		history:     []string{},
		historyIdx:  -1,
		historySize: DefaultHistorySize,
		prompt:      DefaultPrompt,
		promptStyle: style.InputPrompt,
		minHeight:   1,  // single line by default
		maxHeight:   10, // reasonable scroll point
		modeStyle:   lipgloss.NewStyle().Foreground(style.Primary).Bold(true).Background(style.Surface),
		infoStyle:   lipgloss.NewStyle().Foreground(style.TextMuted).Background(style.Surface),
		// Visual styling defaults
		contentBg:           style.Surface,
		borderStyle:         lipgloss.ThickBorder(),
		borderColor:         style.Border,
		flashBorderColor:    style.Primary,
		flashDuration:       150 * time.Millisecond,
		borderPaddingTop:    1,
		borderPaddingBottom: 1,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// WithPlaceholder sets the placeholder text.
// Note: We render the placeholder ourselves (not using bubbles textarea's placeholder)
// to work around a bug where the first placeholder line doesn't fill the background.
func WithPlaceholder(text string) Option {
	return func(m *Model) {
		m.placeholder = text
		// Don't set textarea.Placeholder - we render it ourselves
	}
}

// WithWidth sets the input width.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.textarea.SetWidth(w)
		m.width = w
	}
}

// WithHeight sets the input height.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.textarea.SetHeight(h)
	}
}

// WithHistorySize sets the maximum number of history entries.
// Values <= 0 are ignored and DefaultHistorySize (100) is used instead.
// This is intentional: negative values have no semantic meaning for history size,
// and zero would disable history entirely which should be explicit via a separate option.
func WithHistorySize(n int) Option {
	return func(m *Model) {
		if n > 0 {
			m.historySize = n
		}
		// Note: Invalid values (n <= 0) fall through to default.
		// This is intentional per the documented behavior.
	}
}

// WithShowCounter enables the character/line counter display.
func WithShowCounter(show bool) Option {
	return func(m *Model) {
		m.showCounter = show
	}
}

// WithLineNumbers enables or disables line numbers in the textarea.
func WithLineNumbers(show bool) Option {
	return func(m *Model) {
		m.textarea.ShowLineNumbers = show
	}
}

// WithPrompt sets the prompt prefix shown before the textarea.
// Default is "> ".
func WithPrompt(prompt string) Option {
	return func(m *Model) {
		m.prompt = prompt
	}
}

// WithPromptStyle sets the style for the prompt prefix.
func WithPromptStyle(s lipgloss.Style) Option {
	return func(m *Model) {
		m.promptStyle = s
	}
}

// WithMinHeight sets the minimum visible lines (default 1).
// The input will never shrink below this height.
func WithMinHeight(h int) Option {
	return func(m *Model) {
		if h > 0 {
			m.minHeight = h
		}
	}
}

// WithMaxHeight sets the maximum lines before scrolling (default 10).
// Set to 0 for unlimited height (not recommended).
func WithMaxHeight(h int) Option {
	return func(m *Model) {
		if h >= 0 {
			m.maxHeight = h
		}
	}
}

// WithMode sets the mode indicator (e.g., "Build", "Shell").
// Displayed on the left side of the info bar, before the model name.
func WithMode(mode string) Option {
	return func(m *Model) {
		m.mode = mode
	}
}

// WithModeStyle sets the style for the mode indicator.
func WithModeStyle(s lipgloss.Style) Option {
	return func(m *Model) {
		m.modeStyle = s
	}
}

// WithInfoBar sets the bottom info text (e.g., "Claude Sonnet 4.5").
func WithInfoBar(text string) Option {
	return func(m *Model) {
		m.infoBar = text
	}
}

// WithInfoStyle sets the style for the info bar.
func WithInfoStyle(s lipgloss.Style) Option {
	return func(m *Model) {
		m.infoStyle = s
	}
}

// WithContentBackground sets the background color for all content inside the borders.
// Default is style.Surface.
func WithContentBackground(color lipgloss.Color) Option {
	return func(m *Model) {
		m.contentBg = color
	}
}

// WithBorderStyle sets the border style (e.g., ThickBorder, RoundedBorder, NormalBorder).
// Default is lipgloss.ThickBorder().
func WithBorderStyle(border lipgloss.Border) Option {
	return func(m *Model) {
		m.borderStyle = border
	}
}

// WithBorderColor sets the border color when not flashing.
// Default is style.Border.
func WithBorderColor(color lipgloss.Color) Option {
	return func(m *Model) {
		m.borderColor = color
	}
}

// WithFlashBorderColor sets the border color during the copy flash animation.
// Default is style.Primary.
func WithFlashBorderColor(color lipgloss.Color) Option {
	return func(m *Model) {
		m.flashBorderColor = color
	}
}

// WithFlashDuration sets the duration of the copy flash animation.
// Default is 150ms.
func WithFlashDuration(d time.Duration) Option {
	return func(m *Model) {
		if d > 0 {
			m.flashDuration = d
		}
	}
}

// WithBorderPadding sets the top and bottom padding inside the border.
// Default is (1, 1).
func WithBorderPadding(top, bottom int) Option {
	return func(m *Model) {
		if top >= 0 {
			m.borderPaddingTop = top
		}
		if bottom >= 0 {
			m.borderPaddingBottom = bottom
		}
	}
}

// Focus sets the focused state and returns a command for cursor blink.
func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return m.textarea.Focus()
}

// Blur removes focus from the input.
func (m *Model) Blur() {
	m.focused = false
	m.textarea.Blur()
}

// Focused returns whether the input is focused.
func (m Model) Focused() bool {
	return m.focused
}

// Value returns the current input text.
func (m Model) Value() string {
	return m.textarea.Value()
}

// SetValue sets the input text.
func (m *Model) SetValue(s string) {
	m.textarea.SetValue(s)
}

// Prompt returns the current prompt prefix.
func (m Model) Prompt() string {
	return m.prompt
}

// SetWidth updates the input width at runtime.
// Useful for responsive layouts that resize components when terminal size changes.
func (m *Model) SetWidth(w int) {
	m.textarea.SetWidth(w)
	m.width = w
}

// Width returns the current configured width.
// Returns 0 if no width was explicitly set.
func (m Model) Width() int {
	return m.width
}

// SetHeight updates the input height at runtime.
// Note: When using auto-expanding behavior (minHeight/maxHeight), prefer
// SetMinHeight and SetMaxHeight instead, as SetHeight sets a static height.
func (m *Model) SetHeight(h int) {
	m.textarea.SetHeight(h)
}

// Height returns the current calculated display height based on content.
// This is clamped between minHeight and maxHeight.
func (m Model) Height() int {
	lineCount := m.LineCount()
	displayHeight := lineCount
	if displayHeight < m.minHeight {
		displayHeight = m.minHeight
	}
	if m.maxHeight > 0 && displayHeight > m.maxHeight {
		displayHeight = m.maxHeight
	}
	return displayHeight
}

// MinHeight returns the configured minimum height.
func (m Model) MinHeight() int {
	return m.minHeight
}

// MaxHeight returns the configured maximum height.
func (m Model) MaxHeight() int {
	return m.maxHeight
}

// SetMinHeight updates the minimum height at runtime.
func (m *Model) SetMinHeight(h int) {
	if h > 0 {
		m.minHeight = h
	}
}

// SetMaxHeight updates the maximum height at runtime.
func (m *Model) SetMaxHeight(h int) {
	if h >= 0 {
		m.maxHeight = h
	}
}

// Mode returns the current mode indicator.
func (m Model) Mode() string {
	return m.mode
}

// SetMode updates the mode indicator at runtime.
func (m *Model) SetMode(mode string) {
	m.mode = mode
}

// InfoBar returns the current info bar text.
func (m Model) InfoBar() string {
	return m.infoBar
}

// SetInfoBar updates the info bar text at runtime.
func (m *Model) SetInfoBar(text string) {
	m.infoBar = text
}

// Flashing returns whether the copy flash animation is currently active.
func (m Model) Flashing() bool {
	return m.flashing
}

// SetSize updates both width and height at runtime.
// Matches ResponsiveComponent interface pattern from RESPONSIVE.md.
func (m *Model) SetSize(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle flash done regardless of focus state
	if _, ok := msg.(flashDoneMsg); ok {
		m.flashing = false
		return m, nil
	}

	if !m.focused {
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if handled, cmd := m.handleKeyMsg(keyMsg); handled {
			return m, cmd
		}
	}

	// Pass other messages to textarea
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// handleKeyMsg processes key messages. Returns (handled, cmd).
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEnter && !msg.Alt:
		return true, m.handleSubmit()

	case msg.Type == tea.KeyEnter && msg.Alt, msg.Type == tea.KeyCtrlJ:
		m.textarea.SetValue(m.textarea.Value() + "\n")
		return true, nil

	case msg.Type == tea.KeyEsc:
		m.textarea.SetValue("")
		m.historyIdx = -1
		return true, nil

	case msg.Type == tea.KeyUp:
		m.historyUp()
		return true, nil

	case msg.Type == tea.KeyDown:
		m.historyDown()
		return true, nil

	case msg.Type == tea.KeyCtrlY, msg.String() == "ctrl+shift+c":
		// Ctrl+Y (reliable, traditional Unix "yank") or Ctrl+Shift+C (if terminal supports it)
		return true, m.CopyCmd()
	}
	return false, nil
}

// CopyCmd returns a command that copies the current value to the system clipboard
// using OSC 52 escape sequences. This works cross-platform in modern terminals
// including over SSH, tmux, and Wayland sessions.
//
// Can be triggered programmatically or via Ctrl+Y (recommended) or Ctrl+Shift+C.
// Note: Ctrl+Shift+C may not work in all terminals as some intercept it for their
// own copy function. Ctrl+Y is the traditional Unix "yank" key and works reliably.
//
// This also triggers a brief flash animation on the border to provide visual feedback.
// The flash animation respects REDUCE_MOTION environment variable for accessibility.
func (m *Model) CopyCmd() tea.Cmd {
	text := m.textarea.Value()
	m.flashing = true

	copyCmd := func() tea.Msg {
		// OSC 52 writes to terminal's clipboard via escape sequence.
		// Note: Even with a nil error, clipboard success cannot be verified
		// as OSC 52 is "fire and forget" - the terminal handles the actual copy.
		_, err := osc52.New(text).WriteTo(os.Stderr)
		return CopiedMsg{Text: text, Err: err}
	}

	// Respect REDUCE_MOTION preference for accessibility
	if reduceMotion() {
		return copyCmd
	}

	return tea.Batch(
		copyCmd,
		tea.Tick(m.flashDuration, func(time.Time) tea.Msg {
			return flashDoneMsg{}
		}),
	)
}

// handleSubmit processes Enter key for submission.
func (m *Model) handleSubmit() tea.Cmd {
	text := m.textarea.Value()
	if text == "" {
		return nil
	}
	m.addToHistory(text)
	m.historyIdx = -1
	m.textarea.SetValue("")
	return func() tea.Msg {
		return SubmitMsg{Text: text}
	}
}

// historyUp navigates to an older history entry.
func (m *Model) historyUp() {
	if len(m.history) == 0 {
		return
	}
	// Save current input as draft when first entering history
	if m.historyIdx == -1 {
		m.draft = m.textarea.Value()
	}
	if m.historyIdx < len(m.history)-1 {
		m.historyIdx++
		m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
	}
}

// historyDown navigates to a newer history entry.
func (m *Model) historyDown() {
	if m.historyIdx > 0 {
		m.historyIdx--
		m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
	} else if m.historyIdx == 0 {
		// Restore draft when exiting history
		m.historyIdx = -1
		m.textarea.SetValue(m.draft)
		m.draft = ""
	}
}

// addToHistory adds text to the history buffer, respecting max size.
func (m *Model) addToHistory(text string) {
	m.history = append(m.history, text)
	// Trim if over max size
	if len(m.history) > m.historySize {
		m.history = m.history[len(m.history)-m.historySize:]
	}
}

// View implements tea.Model.
// Note: This method intentionally mutates the receiver's textarea height
// to implement auto-expansion. This is safe because Bubble Tea treats
// models as values and the mutation enables dynamic height adjustment
// based on content without requiring explicit SetHeight calls.
func (m Model) View() string {
	displayHeight := m.Height()
	m.textarea.SetHeight(displayHeight) // Intentional mutation for auto-expansion

	prompt := m.promptStyle.Render(m.prompt)
	promptWidth := lipgloss.Width(prompt)

	// Calculate content dimensions (inside borders)
	contentWidth := max(20, m.width-2)

	// Set textarea width to fit within content area (minus prompt)
	textareaWidth := max(10, contentWidth-promptWidth)
	m.textarea.SetWidth(textareaWidth)

	// Style for content background fill
	bgStyle := lipgloss.NewStyle().Background(m.contentBg)

	// Render textarea or custom placeholder
	var textareaView string
	if m.textarea.Value() == "" && m.placeholder != "" {
		// Custom placeholder rendering (workaround for bubbles textarea bug)
		// The bug: first placeholder line doesn't pad to full width
		placeholderStyle := lipgloss.NewStyle().
			Foreground(style.TextMuted).
			Background(m.contentBg).
			Width(textareaWidth)
		textareaView = placeholderStyle.Render(m.placeholder)
	} else {
		textareaView = m.textarea.View()
	}

	// Use Place() to ensure content background fills entire area
	textareaView = lipgloss.Place(
		textareaWidth,
		displayHeight,
		lipgloss.Left,
		lipgloss.Top,
		textareaView,
		lipgloss.WithWhitespaceBackground(m.contentBg),
	)

	// Full-width style for lines that need to fill contentWidth
	fullWidthStyle := bgStyle.Width(contentWidth)

	// Build inner content - prompt + textarea (already has bg via Place)
	inputRow := lipgloss.JoinHorizontal(lipgloss.Top, prompt, textareaView)
	inputRow = lipgloss.Place(
		contentWidth,
		displayHeight,
		lipgloss.Left,
		lipgloss.Top,
		inputRow,
		lipgloss.WithWhitespaceBackground(m.contentBg),
	)

	var innerContent string
	if m.infoBar != "" || m.showCounter {
		// Spacing line: full-width spaces with Surface background
		spacingLine := fullWidthStyle.Render("")
		// Info bar with Surface background and explicit width
		infoLine := fullWidthStyle.Render(m.buildInfoBar(prompt, contentWidth))
		innerContent = lipgloss.JoinVertical(lipgloss.Left, inputRow, spacingLine, infoLine)
	} else {
		innerContent = inputRow
	}

	// Calculate total height for Place:
	// - displayHeight lines for textarea
	// - 1 spacing line (if info bar shown)
	// - 1 info bar line (if shown)
	totalHeight := displayHeight
	if m.infoBar != "" || m.showCounter {
		totalHeight += 2 // spacing + info bar
	}

	// Use Place to fill background with content background color
	innerContent = lipgloss.Place(
		contentWidth,
		totalHeight,
		lipgloss.Left,
		lipgloss.Top,
		innerContent,
		lipgloss.WithWhitespaceBackground(m.contentBg),
	)

	currentBorderColor := m.borderColor
	if m.flashing {
		currentBorderColor = m.flashBorderColor
	}

	// Border with padding and content background for padding areas
	// Width ensures content fits exactly, making right border visible
	borderStyleRendered := lipgloss.NewStyle().
		Width(contentWidth).
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(m.borderStyle).
		BorderForeground(currentBorderColor).
		PaddingTop(m.borderPaddingTop).
		PaddingBottom(m.borderPaddingBottom).
		Background(m.contentBg)

	return borderStyleRendered.Render(innerContent)
}

// buildInfoBar constructs the info bar line with mode and info on left, counter on right.
// Layout: [indent][Mode] [InfoBar/Model]  ...spacing...  [Counter]
// The contentWidth is the width inside the borders (total width minus border chars).
// When a prompt is configured, the info bar aligns with the textarea (indented by prompt width).
// When no prompt is configured, the info bar starts at the left edge.
// All content uses the configured content background for consistent appearance.
func (m Model) buildInfoBar(prompt string, contentWidth int) string {
	promptWidth := lipgloss.Width(prompt)

	// Style for background-colored spaces
	bgStyle := lipgloss.NewStyle().Background(m.contentBg)

	// Left side: mode (if set) followed by info bar text
	left := ""
	if m.mode != "" {
		left = m.modeStyle.Render(m.mode)
		if m.infoBar != "" {
			left += bgStyle.Render(" ") + m.infoStyle.Render(m.infoBar)
		}
	} else if m.infoBar != "" {
		left = m.infoStyle.Render(m.infoBar)
	}

	// Right side: counter (if enabled)
	right := ""
	if m.showCounter {
		right = m.infoStyle.Render(fmt.Sprintf("%d chars, %d lines", m.CharCount(), m.LineCount()))
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	// Only indent if there's a prompt to align with
	indent := ""
	if promptWidth > 0 {
		indent = bgStyle.Render(strings.Repeat(" ", promptWidth))
	}

	// Right padding matches left indent pattern for visual symmetry
	rightPad := ""
	rightPadWidth := 0
	if promptWidth > 0 {
		rightPad = bgStyle.Render(strings.Repeat(" ", promptWidth))
		rightPadWidth = promptWidth
	}

	spacing := max(1, contentWidth-promptWidth-leftWidth-rightWidth-rightPadWidth)

	return indent + left + bgStyle.Render(strings.Repeat(" ", spacing)) + right + rightPad
}

// CharCount returns the current character count (Unicode runes, not bytes).
func (m Model) CharCount() int {
	return utf8.RuneCountInString(m.textarea.Value())
}

// LineCount returns the current line count.
func (m Model) LineCount() int {
	text := m.textarea.Value()
	if text == "" {
		return 1
	}
	return strings.Count(text, "\n") + 1
}
