package chatinput

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

// DefaultHistorySize is the default maximum number of history entries.
const DefaultHistorySize = 100

// SubmitMsg is emitted when the user submits the input (presses Enter).
type SubmitMsg struct {
	Text string
}

// CopiedMsg is emitted when the input content is copied to clipboard.
type CopiedMsg struct {
	Text string
	Err  error
}

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
}

// Option configures the Model.
type Option func(*Model)

// New creates a new ChatInput model with the given options.
func New(opts ...Option) Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	m := Model{
		textarea:    ta,
		history:     []string{},
		historyIdx:  -1,
		historySize: DefaultHistorySize,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// WithPlaceholder sets the placeholder text.
func WithPlaceholder(text string) Option {
	return func(m *Model) {
		m.textarea.Placeholder = text
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
func WithHistorySize(n int) Option {
	return func(m *Model) {
		if n > 0 {
			m.historySize = n
		}
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

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case msg.String() == "ctrl+shift+c":
		return true, m.handleCopy()
	}
	return false, nil
}

// handleCopy copies the current value to the system clipboard.
func (m *Model) handleCopy() tea.Cmd {
	text := m.textarea.Value()
	return func() tea.Msg {
		err := clipboard.Init()
		if err != nil {
			return CopiedMsg{Text: text, Err: err}
		}
		clipboard.Write(clipboard.FmtText, []byte(text))
		return CopiedMsg{Text: text, Err: nil}
	}
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
func (m Model) View() string {
	view := m.textarea.View()
	if m.showCounter {
		view += "\n" + m.counterView()
	}
	return view
}

// counterView returns the character and line count display.
func (m Model) counterView() string {
	text := m.textarea.Value()
	chars := len(text)
	lines := 1
	if text != "" {
		lines = strings.Count(text, "\n") + 1
	}
	counter := fmt.Sprintf("%d chars, %d lines", chars, lines)
	// Right-align if width is set
	if m.width > 0 && len(counter) < m.width {
		padding := m.width - len(counter)
		counter = strings.Repeat(" ", padding) + counter
	}
	return counter
}

// CharCount returns the current character count.
func (m Model) CharCount() int {
	return len(m.textarea.Value())
}

// LineCount returns the current line count.
func (m Model) LineCount() int {
	text := m.textarea.Value()
	if text == "" {
		return 1
	}
	return strings.Count(text, "\n") + 1
}
