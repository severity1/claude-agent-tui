// Package keyhints provides a reusable component for displaying keybinding hints.
package keyhints

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/theme"
)

// Binding represents a single keybinding with its description.
type Binding struct {
	Key  string // The key or key combination (e.g., "Enter", "Ctrl+C")
	Desc string // Description of what the key does (e.g., "Submit", "Quit")
}

// DefaultVisibleCount is the default number of hints shown when collapsed.
const DefaultVisibleCount = 3

// DefaultPadding is the default left padding for inline hints view.
// Empty by default since the left border provides visual separation.
const DefaultPadding = ""

// borderWidth accounts for the left border character and internal padding
// used in View() width calculations.
const borderWidth = 2

// HelpMode defines how the help window expands when toggled.
type HelpMode int

const (
	HelpModeDown   HelpMode = iota // Expands downward (default)
	HelpModeUp                     // Expands upward
	HelpModeCenter                 // Floating window at center
	HelpModeTop                    // Expands from top
)

// Model represents the keyhints component state.
type Model struct {
	bindings       []Binding
	separator      string
	keyStyle       lipgloss.Style
	descStyle      lipgloss.Style
	sepStyle       lipgloss.Style
	focusStyle     lipgloss.Style // style applied when focused
	helpStyle      lipgloss.Style // style for the help window border
	padding        string         // left padding for inline views
	width          int
	expanded       bool // whether all bindings are visible in inline view
	defaultVisible int  // number of bindings to show when collapsed
	focused        bool // whether the component has focus
	showHelp       bool // whether to show the help window

	// Colors for dynamic styling
	borderColor  lipgloss.TerminalColor // border color when not focused
	primaryColor lipgloss.TerminalColor // accent color for focus indication
	contentBg    lipgloss.TerminalColor // background color for help content

	// Help window configuration
	helpMode HelpMode // how the help window expands
}

// Option configures the Model.
type Option func(*Model)

// New creates a new keyhints model with the given bindings and options.
func New(bindings []Binding, opts ...Option) Model {
	// Use default styles
	defaults := DefaultStyles()

	m := Model{
		bindings:       bindings,
		separator:      " | ",
		keyStyle:       defaults.KeyStyle,
		descStyle:      defaults.DescStyle,
		sepStyle:       defaults.SepStyle,
		focusStyle:     defaults.FocusStyle,
		helpStyle:      defaults.HelpStyle,
		padding:        DefaultPadding,
		expanded:       true, // show all by default for backwards compatibility
		defaultVisible: DefaultVisibleCount,
		borderColor:    defaults.BorderColor,
		primaryColor:   defaults.PrimaryColor,
		contentBg:      defaults.ContentBg,
		helpMode:       HelpModeDown, // default: expand downward
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// WithSeparator sets the separator between keybinding hints.
func WithSeparator(sep string) Option {
	return func(m *Model) {
		m.separator = sep
	}
}

// WithKeyStyle sets the style for key names.
func WithKeyStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.keyStyle = style
	}
}

// WithDescStyle sets the style for descriptions.
func WithDescStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.descStyle = style
	}
}

// WithSepStyle sets the style for the separator.
func WithSepStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.sepStyle = style
	}
}

// WithWidth sets the maximum width for the hints display.
// When set, View() will truncate hints that exceed this width, showing only
// hints that fit with a "+N more" indicator for hidden ones.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// WithDefaultVisible sets the number of bindings shown when collapsed.
// Values <= 0 are ignored and DefaultVisibleCount (3) is used instead.
// This is intentional: zero would show nothing, and negative values have no meaning.
func WithDefaultVisible(n int) Option {
	return func(m *Model) {
		if n > 0 {
			m.defaultVisible = n
		}
		// Note: Invalid values (n <= 0) fall through to default.
		// This is intentional per the documented behavior.
	}
}

// WithCollapsed starts the keyhints in collapsed state.
func WithCollapsed() Option {
	return func(m *Model) {
		m.expanded = false
	}
}

// WithFocusStyle sets the style used when the component is focused.
func WithFocusStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.focusStyle = style
	}
}

// WithHelpStyle sets the style for the help window border.
func WithHelpStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.helpStyle = style
	}
}

// WithHelpMode sets how the help window expands when toggled.
func WithHelpMode(mode HelpMode) Option {
	return func(m *Model) {
		m.helpMode = mode
	}
}

// WithContentBackground sets the background color for help content.
func WithContentBackground(color lipgloss.TerminalColor) Option {
	return func(m *Model) {
		m.contentBg = color
	}
}

// WithPadding sets the left padding for inline views (View and ViewCompact).
// Use empty string "" to disable padding. Default is DefaultPadding (empty string).
func WithPadding(padding string) Option {
	return func(m *Model) {
		m.padding = padding
	}
}

// WithPalette applies a theme palette to the component.
// This sets all styling from the palette, which can be overridden by subsequent options.
func WithPalette(p *theme.Palette) Option {
	return func(m *Model) {
		if p == nil {
			return
		}
		s := StylesFromPalette(p)
		m.keyStyle = s.KeyStyle
		m.descStyle = s.DescStyle
		m.sepStyle = s.SepStyle
		m.focusStyle = s.FocusStyle
		m.helpStyle = s.HelpStyle
		m.borderColor = s.BorderColor
		m.primaryColor = s.PrimaryColor
		m.contentBg = s.ContentBg
	}
}

// WithStyles applies a custom Styles struct to the component.
// Use StylesFromPalette() to create a Styles from a theme, then customize individual fields.
func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.keyStyle = s.KeyStyle
		m.descStyle = s.DescStyle
		m.sepStyle = s.SepStyle
		m.focusStyle = s.FocusStyle
		m.helpStyle = s.HelpStyle
		m.borderColor = s.BorderColor
		m.primaryColor = s.PrimaryColor
		m.contentBg = s.ContentBg
	}
}

// SetBindings replaces the current bindings.
// The input slice is cloned to prevent external mutation.
func (m *Model) SetBindings(bindings []Binding) {
	m.bindings = slices.Clone(bindings)
}

// AddBinding adds a new binding to the list.
func (m *Model) AddBinding(b Binding) {
	m.bindings = append(m.bindings, b)
}

// Bindings returns a copy of the current list of bindings.
// The returned slice is a defensive copy to prevent external mutation.
func (m Model) Bindings() []Binding {
	return slices.Clone(m.bindings)
}

// ToggleExpanded switches between expanded and collapsed states.
// Named explicitly to distinguish from ToggleHelp().
func (m *Model) ToggleExpanded() {
	m.expanded = !m.expanded
}

// Expanded returns whether all bindings are visible.
func (m Model) Expanded() bool {
	return m.expanded
}

// SetExpanded sets the expanded state directly.
func (m *Model) SetExpanded(expanded bool) {
	m.expanded = expanded
}

// SetWidth sets the width for the component.
// This affects both the inline View() truncation and ViewHelp() width in attached modes.
func (m *Model) SetWidth(width int) {
	m.width = width
}

// Focus sets the focused state and returns nil (no cursor blink needed for this component).
// The return type matches Bubble Tea conventions for composability with other components.
func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// Blur removes focus from the component.
func (m *Model) Blur() {
	m.focused = false
}

// Focused returns whether the component has focus.
func (m Model) Focused() bool {
	return m.focused
}

// ShowHelp opens the help window.
func (m *Model) ShowHelp() {
	m.showHelp = true
}

// HideHelp closes the help window.
func (m *Model) HideHelp() {
	m.showHelp = false
}

// ToggleHelp toggles the help window visibility.
func (m *Model) ToggleHelp() {
	m.showHelp = !m.showHelp
}

// ShowingHelp returns whether the help window is visible.
func (m Model) ShowingHelp() bool {
	return m.showHelp
}

// HelpMode returns how the help window expands.
func (m Model) HelpMode() HelpMode {
	return m.helpMode
}

// HelpHeight returns the line count of the help window for parent layout calculations.
// Returns 0 if there are no bindings.
func (m Model) HelpHeight() int {
	if len(m.bindings) == 0 {
		return 1 // "No keybindings defined" message
	}
	// Title + blank line + bindings + blank line + close hint
	return 1 + 1 + len(m.bindings) + 1 + 1
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
// When focused, handles "?" key to toggle help window, Esc to close it.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "?":
			m.showHelp = !m.showHelp
		case "esc":
			if m.showHelp {
				m.showHelp = false
			}
		}
	}
	return m, nil
}

// View implements tea.Model - renders keybindings as inline hints.
// When collapsed, shows only defaultVisible bindings with a "? more" indicator.
// When width is set, truncates to fit within width with a "+N more" indicator.
// When focused, shows a visual indicator.
func (m Model) View() string {
	if len(m.bindings) == 0 {
		return ""
	}

	bindings := m.bindings
	hiddenCount := 0

	// When collapsed, limit visible bindings
	if !m.expanded && len(m.bindings) > m.defaultVisible {
		hiddenCount = len(m.bindings) - m.defaultVisible
		bindings = m.bindings[:m.defaultVisible]
	}

	sep := m.sepStyle.Render(m.separator)
	sepWidth := lipgloss.Width(sep)

	// Build parts, respecting width constraint if set
	var parts []string
	currentWidth := 0

	// Account for border width (1 character + border padding)
	currentWidth = borderWidth

	for i, b := range bindings {
		key := m.keyStyle.Render(b.Key)
		desc := m.descStyle.Render(b.Desc)
		part := key + ": " + desc
		partWidth := lipgloss.Width(part)

		// Check if adding this part would exceed width (if width is set)
		if m.width > 0 && len(parts) > 0 {
			// Reserve space for potential "+N more" indicator
			moreIndicator := fmt.Sprintf("+%d more", len(bindings)-i+hiddenCount)
			moreWidth := lipgloss.Width(m.descStyle.Render(moreIndicator))
			neededWidth := currentWidth + sepWidth + partWidth

			// If this part won't fit, stop and show remaining count
			if neededWidth+sepWidth+moreWidth > m.width {
				hiddenCount += len(bindings) - i
				break
			}
		}

		if len(parts) > 0 {
			currentWidth += sepWidth
		}
		parts = append(parts, part)
		currentWidth += partWidth
	}

	// Add "? more" or "+N more" indicator when there are hidden hints
	if hiddenCount > 0 {
		var moreText string
		if m.width > 0 {
			// Width-based truncation uses "+N more" format
			moreText = m.descStyle.Render(fmt.Sprintf("+%d more", hiddenCount))
		} else {
			// Collapsed state uses "? N more" format
			moreText = m.descStyle.Render(fmt.Sprintf("? %d more", hiddenCount))
		}
		parts = append(parts, moreText)
	}

	result := strings.Join(parts, sep)

	// Apply left padding
	result = m.padding + result

	// Wrap with left border - use focusStyle color when focused, borderColor when unfocused
	currentBorderColor := m.borderColor
	if m.focused {
		currentBorderColor = m.focusStyle.GetForeground()
	}
	borderStyle := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(currentBorderColor).
		PaddingLeft(1)

	// Apply width for consistent full-width rendering (like ChatInput)
	if m.width > 0 {
		borderStyle = borderStyle.Width(m.width - borderWidth)
	}

	return borderStyle.Render(result)
}

// ViewCompact renders a more compact version with just keys.
// Left padding is applied consistently with View().
func (m Model) ViewCompact() string {
	if len(m.bindings) == 0 {
		return ""
	}

	var parts []string
	for _, b := range m.bindings {
		parts = append(parts, m.keyStyle.Render(b.Key))
	}

	sep := m.sepStyle.Render(m.separator)
	return m.padding + strings.Join(parts, sep)
}

// ViewHelp renders the help window showing all keybindings.
// Uses left border style to match View() and applies background color if set.
// In Down/Up/Top modes, width matches the keyhints bar; Center uses content-width.
func (m Model) ViewHelp() string {
	if len(m.bindings) == 0 {
		return m.renderHelpContent("No keybindings defined", m.width)
	}

	var sb strings.Builder
	sb.WriteString(m.keyStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	// Find max key display width for alignment (handles Unicode)
	maxKeyWidth := 0
	for _, b := range m.bindings {
		if w := lipgloss.Width(b.Key); w > maxKeyWidth {
			maxKeyWidth = w
		}
	}

	// Render each binding with proper Unicode-aware padding
	for _, b := range m.bindings {
		// Use lipgloss Width for proper display width padding (handles CJK, etc.)
		keyWidth := lipgloss.Width(b.Key)
		padding := ""
		if keyWidth < maxKeyWidth {
			padding = strings.Repeat(" ", maxKeyWidth-keyWidth)
		}
		key := m.keyStyle.Render(b.Key + padding)
		desc := m.descStyle.Render(b.Desc)
		fmt.Fprintf(&sb, "%s  %s\n", key, desc)
	}

	sb.WriteString("\n")
	sb.WriteString(m.sepStyle.Render("Press ? or Esc to close"))

	return m.renderHelpContent(sb.String(), m.width)
}

// renderHelpContent wraps content with left border styling for the help window.
// Width is applied for non-Center modes to match the keyhints bar width.
func (m Model) renderHelpContent(content string, width int) string {
	// Use focusStyle color for border (help window is always "active" when shown)
	borderColor := m.focusStyle.GetForeground()

	style := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(borderColor).
		PaddingLeft(1)

	// Apply background color if set
	if m.contentBg != nil {
		style = style.Background(m.contentBg)
	}

	// Apply width for non-Center modes to match keyhints bar width
	if m.helpMode != HelpModeCenter && width > 0 {
		style = style.Width(width - borderWidth)
	}

	return style.Render(content)
}
