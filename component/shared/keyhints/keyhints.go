// Package keyhints provides a reusable component for displaying keybinding hints.
package keyhints

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Binding represents a single keybinding with its description.
type Binding struct {
	Key  string // The key or key combination (e.g., "Enter", "Ctrl+C")
	Desc string // Description of what the key does (e.g., "Submit", "Quit")
}

// DefaultVisibleCount is the default number of hints shown when collapsed.
const DefaultVisibleCount = 3

// Model represents the keyhints component state.
type Model struct {
	bindings       []Binding
	separator      string
	keyStyle       lipgloss.Style
	descStyle      lipgloss.Style
	sepStyle       lipgloss.Style
	focusStyle     lipgloss.Style // style applied when focused
	helpStyle      lipgloss.Style // style for the help window border
	width          int
	expanded       bool // whether all bindings are visible in inline view
	defaultVisible int  // number of bindings to show when collapsed
	focused        bool // whether the component has focus
	showHelp       bool // whether to show the help window
}

// Option configures the Model.
type Option func(*Model)

// New creates a new keyhints model with the given bindings and options.
func New(bindings []Binding, opts ...Option) Model {
	m := Model{
		bindings:       bindings,
		separator:      " | ",
		keyStyle:       lipgloss.NewStyle().Bold(true),
		descStyle:      lipgloss.NewStyle(),
		sepStyle:       lipgloss.NewStyle().Faint(true),
		focusStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("212")), // pink highlight
		helpStyle:      lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2),
		expanded:       true, // show all by default for backwards compatibility
		defaultVisible: DefaultVisibleCount,
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
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// WithDefaultVisible sets the number of bindings shown when collapsed.
// If n <= 0, defaults to DefaultVisibleCount.
func WithDefaultVisible(n int) Option {
	return func(m *Model) {
		if n > 0 {
			m.defaultVisible = n
		}
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

// SetBindings replaces the current bindings.
func (m *Model) SetBindings(bindings []Binding) {
	m.bindings = bindings
}

// AddBinding adds a new binding to the list.
func (m *Model) AddBinding(b Binding) {
	m.bindings = append(m.bindings, b)
}

// Bindings returns the current list of bindings.
func (m Model) Bindings() []Binding {
	return m.bindings
}

// Toggle switches between expanded and collapsed states.
func (m *Model) Toggle() {
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

// Focus sets the focused state.
func (m *Model) Focus() {
	m.focused = true
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

	var parts []string
	for _, b := range bindings {
		key := m.keyStyle.Render(b.Key)
		desc := m.descStyle.Render(b.Desc)
		parts = append(parts, key+": "+desc)
	}

	// Add "? more" indicator when collapsed with hidden hints
	if hiddenCount > 0 {
		moreText := m.descStyle.Render(fmt.Sprintf("? %d more", hiddenCount))
		parts = append(parts, moreText)
	}

	sep := m.sepStyle.Render(m.separator)
	result := strings.Join(parts, sep)

	// Add focus indicator
	if m.focused {
		prefix := m.focusStyle.Render("> ")
		result = prefix + result
	}

	return result
}

// ViewCompact renders a more compact version with just keys.
func (m Model) ViewCompact() string {
	if len(m.bindings) == 0 {
		return ""
	}

	var parts []string
	for _, b := range m.bindings {
		parts = append(parts, m.keyStyle.Render(b.Key))
	}

	sep := m.sepStyle.Render(m.separator)
	return strings.Join(parts, sep)
}

// ViewHelp renders a boxed help window showing all keybindings.
func (m Model) ViewHelp() string {
	if len(m.bindings) == 0 {
		return m.helpStyle.Render("No keybindings defined")
	}

	var sb strings.Builder
	sb.WriteString(m.keyStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	// Find max key length for alignment
	maxKeyLen := 0
	for _, b := range m.bindings {
		if len(b.Key) > maxKeyLen {
			maxKeyLen = len(b.Key)
		}
	}

	// Render each binding
	for _, b := range m.bindings {
		key := m.keyStyle.Render(fmt.Sprintf("%-*s", maxKeyLen, b.Key))
		desc := m.descStyle.Render(b.Desc)
		sb.WriteString(fmt.Sprintf("  %s  %s\n", key, desc))
	}

	sb.WriteString("\n")
	sb.WriteString(m.sepStyle.Render("Press ? or Esc to close"))

	return m.helpStyle.Render(sb.String())
}
