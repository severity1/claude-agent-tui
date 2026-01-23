// Package keyhints provides a reusable component for displaying keybinding hints.
package keyhints

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/style"
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
}

// Option configures the Model.
type Option func(*Model)

// New creates a new keyhints model with the given bindings and options.
func New(bindings []Binding, opts ...Option) Model {
	m := Model{
		bindings:       bindings,
		separator:      " | ",
		keyStyle:       lipgloss.NewStyle().Bold(true).Foreground(style.TextMuted),
		descStyle:      lipgloss.NewStyle().Foreground(style.TextMuted),
		sepStyle:       lipgloss.NewStyle().Faint(true).Foreground(style.TextMuted),
		focusStyle:     lipgloss.NewStyle().Foreground(style.Primary),
		helpStyle:      lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(style.Border).Padding(1, 2),
		padding:        DefaultPadding,
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

// WithPadding sets the left padding for inline views (View and ViewCompact).
// Use empty string "" to disable padding. Default is DefaultPadding (1 space).
func WithPadding(padding string) Option {
	return func(m *Model) {
		m.padding = padding
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
	borderWidth := 2 // left border character + internal padding
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

	// Wrap with left border - use focusStyle color when focused, Border color when unfocused
	var borderColor lipgloss.TerminalColor = style.Border
	if m.focused {
		borderColor = m.focusStyle.GetForeground()
	}
	borderStyle := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(borderColor).
		PaddingLeft(1)

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

// ViewHelp renders a boxed help window showing all keybindings.
func (m Model) ViewHelp() string {
	if len(m.bindings) == 0 {
		return m.helpStyle.Render("No keybindings defined")
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
		fmt.Fprintf(&sb, "  %s  %s\n", key, desc)
	}

	sb.WriteString("\n")
	sb.WriteString(m.sepStyle.Render("Press ? or Esc to close"))

	// Note: ViewHelp does NOT apply padding - the box has its own styling
	// and prepending padding would break the border alignment
	return m.helpStyle.Render(sb.String())
}
