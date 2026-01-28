// Package keyhints provides a reusable component for displaying keybinding hints.
package keyhints

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/internal/accessibility"
	"github.com/severity1/claude-agent-tui/internal/focus"
	"github.com/severity1/claude-agent-tui/internal/styles"
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

// HelpMode defines how the help window expands when toggled.
type HelpMode int

const (
	HelpModeDown   HelpMode = iota // Expands downward (default)
	HelpModeUp                     // Expands upward
	HelpModeCenter                 // Floating window at center
	HelpModeTop                    // Expands from top
)

// HeightTickMsg triggers animation frame updates for help window expansion.
type HeightTickMsg struct{}

// heightTick returns a command that sends HeightTickMsg at 60 FPS.
func heightTick() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return HeightTickMsg{}
	})
}

// Model represents the keyhints component state.
type Model struct {
	focus.State      // embedded focus tracking (provides Focus(), Blur(), Focused())
	bindings         []Binding
	reservedBindings []Binding // bindings always shown at end of bar (e.g., "?", "Esc")
	separator        string
	keyStyle         lipgloss.Style
	descStyle        lipgloss.Style
	sepStyle         lipgloss.Style
	focusStyle       lipgloss.Style // style applied when focused
	helpStyle        lipgloss.Style // style for the help window border
	padding          string         // left padding for inline views
	width            int
	expanded         bool // whether all bindings are visible in inline view
	defaultVisible   int  // number of bindings to show when collapsed
	showHelp         bool // whether to show the help window

	// Background color for help content
	contentBg lipgloss.TerminalColor

	// Help window configuration
	helpMode HelpMode // how the help window expands

	// Animation state fields
	animationEnabled bool             // toggleable by users
	animatedHeight   float64          // current animated height (lines)
	heightVelocity   float64          // spring velocity
	targetHeight     float64          // target height (0 or HelpHeight())
	spring           harmonica.Spring // spring instance
	animating        bool             // true while animation in progress
	overlay          bool             // when true, help floats above content instead of pushing

	// Border styling (uses shared border package)
	border styles.BorderConfig // consolidated border configuration for all modes
}

// Option configures the Model.
type Option func(*Model)

// New creates a new keyhints model with the given options.
// Use WithBindings() to set initial bindings.
func New(opts ...Option) Model {
	// Use default styles
	defaults := DefaultStyles()

	m := Model{
		separator:      " | ",
		keyStyle:       defaults.KeyStyle,
		descStyle:      defaults.DescStyle,
		sepStyle:       defaults.SepStyle,
		focusStyle:     defaults.FocusStyle,
		helpStyle:      defaults.HelpStyle,
		padding:        DefaultPadding,
		expanded:       true, // show all by default for backwards compatibility
		defaultVisible: DefaultVisibleCount,
		contentBg:      defaults.ContentBg,
		helpMode:       HelpModeDown, // default: expand downward

		// Border defaults - use NewBorderConfig for consistent construction
		border: styles.NewBorderConfig(
			styles.WithBorderLeft(true),
			styles.WithBorderRight(false),
			styles.WithBorderColor(defaults.BorderColor),
			styles.WithBorderFocusColor(defaults.PrimaryColor),
		),
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

// WithBindings sets the initial keybindings.
func WithBindings(bindings ...Binding) Option {
	return func(m *Model) {
		m.bindings = bindings
	}
}

// WithReservedBindings sets bindings that are always shown at the end of the bar.
// These bindings are displayed after the regular bindings and before the "+N more"
// indicator, ensuring they are always visible regardless of width constraints.
// Common use: reserved spots for "?" (help) and "Esc" (close/cancel).
func WithReservedBindings(bindings ...Binding) Option {
	return func(m *Model) {
		m.reservedBindings = bindings
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

// WithBorderStyle sets the border style for the inline bar and attached help modes.
// Default is lipgloss.ThickBorder().
func WithBorderStyle(style lipgloss.Border) Option {
	return func(m *Model) {
		m.border.Style = style
	}
}

// WithBorderLeft enables or disables the left border visibility.
// When false, the border space is still reserved but the border character is invisible (space).
// Default is true (visible).
func WithBorderLeft(enabled bool) Option {
	return func(m *Model) {
		m.border.Left = enabled
	}
}

// WithBorderRight enables or disables the right border visibility.
// When false, the border space is still reserved but the border character is invisible (space).
// Default is false (invisible but space reserved).
func WithBorderRight(enabled bool) Option {
	return func(m *Model) {
		m.border.Right = enabled
	}
}

// WithBorderPadding sets the padding inside the border for inline and attached modes.
// Default is (0, 0, 1, 0) for top, bottom, left, right.
// Negative values are ignored for each parameter independently.
func WithBorderPadding(top, bottom, left, right int) Option {
	return func(m *Model) {
		m.border = m.border.SetPaddingIfValid(top, bottom, left, right)
	}
}

// WithBorderColor sets the border color when not focused.
// When focused, the FocusColor from BorderConfig takes precedence.
func WithBorderColor(color lipgloss.TerminalColor) Option {
	return func(m *Model) {
		m.border.Color = color
	}
}

// WithBorderConfig sets the inline border configuration directly.
// Use for bulk configuration when you need to set many border properties at once.
func WithBorderConfig(cfg styles.BorderConfig) Option {
	return func(m *Model) {
		m.border = cfg
	}
}

// applyStylesToModel applies all styles from a Styles struct to a Model.
// This is the single point of style application used by WithPalette and WithStyles.
func applyStylesToModel(m *Model, s Styles) {
	m.keyStyle = s.KeyStyle
	m.descStyle = s.DescStyle
	m.sepStyle = s.SepStyle
	m.focusStyle = s.FocusStyle
	m.helpStyle = s.HelpStyle
	m.border.Color = s.BorderColor
	m.border.FocusColor = s.PrimaryColor
	m.contentBg = s.ContentBg
}

// WithPalette applies a theme palette to the component.
// This sets all styling from the palette, which can be overridden by subsequent options.
// If palette is nil, the component retains its current styles (no-op).
// Note: A nil palette typically indicates theme.Get() failed for an unknown theme.
// Consider using theme.GetWithDefault() for guaranteed non-nil palette.
func WithPalette(p *theme.Palette) Option {
	return func(m *Model) {
		if p == nil {
			return
		}
		applyStylesToModel(m, StylesFromPalette(p))
	}
}

// WithStyles applies a custom Styles struct to the component.
// Use StylesFromPalette() to create a Styles from a theme, then customize individual fields.
func WithStyles(s Styles) Option {
	return func(m *Model) {
		applyStylesToModel(m, s)
	}
}

// WithAnimation enables/disables physics animation for help window expansion.
// Default is false (no animation). When enabled, respects REDUCE_MOTION env var.
// Animation is skipped for HelpModeCenter (instant show/hide for floating modal).
func WithAnimation(enabled bool) Option {
	return func(m *Model) {
		m.animationEnabled = enabled
		if enabled {
			m.spring = harmonica.NewSpring(harmonica.FPS(60), 12.0, 1.0)
		}
	}
}

// WithOverlay enables overlay mode where help floats above content.
// When false (default), Down/Up/Top modes push content.
// When true, they float like Center mode but at their respective positions.
// Overlay mode uses a snappier animation for quicker show/hide.
func WithOverlay(overlay bool) Option {
	return func(m *Model) {
		m.overlay = overlay
	}
}

// WithAnimationSpring sets custom spring parameters (frequency, damping).
// Only effective when animation is enabled.
func WithAnimationSpring(frequency, damping float64) Option {
	return func(m *Model) {
		m.spring = harmonica.NewSpring(harmonica.FPS(60), frequency, damping)
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

// Width returns the current configured width.
// Returns 0 if no width was explicitly set.
func (m Model) Width() int {
	return m.width
}

// Focus, Blur, and Focused methods are inherited from embedded focus.State

// ShowHelp opens the help window.
// When animation is enabled and REDUCE_MOTION is not set, starts animation.
// Animation is skipped for HelpModeCenter (instant show for floating modal).
// Overlay mode uses a snappier spring (frequency=18.0) for quicker animation.
func (m *Model) ShowHelp() {
	m.showHelp = true
	if m.animationEnabled && !accessibility.ReduceMotion() && m.helpMode != HelpModeCenter {
		// Use snappier spring for overlay mode
		if m.overlay {
			m.spring = harmonica.NewSpring(harmonica.FPS(60), 18.0, 1.0)
		}
		m.animatedHeight = 0
		m.heightVelocity = 0
		m.targetHeight = float64(m.HelpHeight())
		m.animating = true
	} else {
		m.animatedHeight = float64(m.HelpHeight())
	}
}

// HideHelp closes the help window.
// When animation is enabled and REDUCE_MOTION is not set, starts animation.
// Animation is skipped for HelpModeCenter (instant hide for floating modal).
// Overlay mode uses a snappier spring (frequency=18.0) for quicker animation.
func (m *Model) HideHelp() {
	if m.animationEnabled && !accessibility.ReduceMotion() && m.helpMode != HelpModeCenter {
		// Use snappier spring for overlay mode
		if m.overlay {
			m.spring = harmonica.NewSpring(harmonica.FPS(60), 18.0, 1.0)
		}
		m.targetHeight = 0
		m.animating = true
		// Don't set showHelp = false yet; wait until animation completes
	} else {
		m.showHelp = false
		m.animatedHeight = 0
	}
}

// ToggleHelp toggles the help window visibility.
// When animation is enabled and REDUCE_MOTION is not set, starts animation.
func (m *Model) ToggleHelp() {
	if m.showHelp {
		m.HideHelp()
	} else {
		m.ShowHelp()
	}
}

// ShowingHelp returns whether the help window is visible.
func (m Model) ShowingHelp() bool {
	return m.showHelp
}

// HelpMode returns how the help window expands.
func (m Model) HelpMode() HelpMode {
	return m.helpMode
}

// SetHelpMode sets how the help window expands for runtime mode changes.
func (m *Model) SetHelpMode(mode HelpMode) {
	m.helpMode = mode
}

// Overlay returns whether overlay mode is enabled.
func (m Model) Overlay() bool {
	return m.overlay
}

// SetOverlay enables or disables overlay mode.
func (m *Model) SetOverlay(overlay bool) {
	m.overlay = overlay
}

// ToggleOverlay toggles overlay mode on/off.
func (m *Model) ToggleOverlay() {
	m.overlay = !m.overlay
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
// Handles HeightTickMsg for animation, and when focused, handles "?" key to toggle help window, Esc to close it.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case HeightTickMsg:
		if !m.animating {
			return m, nil
		}

		// Update spring
		m.animatedHeight, m.heightVelocity = m.spring.Update(
			m.animatedHeight, m.heightVelocity, m.targetHeight,
		)

		// Check if settled (within threshold)
		if math.Abs(m.animatedHeight-m.targetHeight) < 0.01 &&
			math.Abs(m.heightVelocity) < 0.01 {
			m.animatedHeight = m.targetHeight
			m.animating = false
			if m.targetHeight == 0 {
				m.showHelp = false
			}
			return m, nil
		}

		return m, heightTick()

	case tea.KeyMsg:
		if !m.Focused() {
			return m, nil
		}
		switch msg.String() {
		case "?":
			m.ToggleHelp()
			// Start animation tick if animating
			if m.animating {
				return m, heightTick()
			}
		case "esc":
			if m.showHelp || m.animating {
				m.HideHelp()
				if m.animating {
					return m, heightTick()
				}
			}
		}
	}
	return m, nil
}

// renderBinding formats a single binding as "Key: Desc".
func (m Model) renderBinding(b Binding) string {
	key := m.keyStyle.Render(b.Key)
	desc := m.descStyle.Render(b.Desc)
	return key + ": " + desc
}

// buildReservedParts calculates reserved bindings width and parts.
func (m Model) buildReservedParts(sepWidth int) ([]string, int) {
	var parts []string
	totalWidth := 0
	for _, b := range m.reservedBindings {
		part := m.renderBinding(b)
		parts = append(parts, part)
		if totalWidth > 0 {
			totalWidth += sepWidth
		}
		totalWidth += lipgloss.Width(part)
	}
	return parts, totalWidth
}

// View implements tea.Model - renders keybindings as inline hints.
// When collapsed, shows only defaultVisible bindings with a "? more" indicator.
// When width is set, truncates to fit within width with a "+N more" indicator.
// Reserved bindings (set via WithReservedBindings) are always shown at the end.
// When focused, shows a visual indicator.
func (m Model) View() string {
	if len(m.bindings) == 0 && len(m.reservedBindings) == 0 {
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

	// Always reserve space for both borders (space characters when invisible)
	bw := 2
	currentWidth := bw + m.border.PaddingLeft + m.border.PaddingRight

	// Calculate reserved bindings (always shown at end)
	reservedParts, reservedWidth := m.buildReservedParts(sepWidth)

	for i, b := range bindings {
		part := m.renderBinding(b)
		partWidth := lipgloss.Width(part)

		// Check if adding this part would exceed width (if width is set)
		if m.width > 0 && len(parts) > 0 {
			// Reserve space for potential "+N more" indicator and reserved bindings
			moreIndicator := fmt.Sprintf("+%d more", len(bindings)-i+hiddenCount)
			moreWidth := lipgloss.Width(m.descStyle.Render(moreIndicator))
			neededWidth := currentWidth + sepWidth + partWidth
			totalReserved := reservedWidth
			if totalReserved > 0 {
				totalReserved += sepWidth
			}

			// If this part won't fit, stop and show remaining count
			if neededWidth+sepWidth+moreWidth+totalReserved > m.width {
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
		format := "? %d more"
		if m.width > 0 {
			format = "+%d more"
		}
		parts = append(parts, m.descStyle.Render(fmt.Sprintf(format, hiddenCount)))
	}

	// Add reserved bindings at the end (always visible)
	parts = append(parts, reservedParts...)

	result := strings.Join(parts, sep)

	// Apply left padding
	result = m.padding + result

	// Apply border styling - BorderConfig.Color/FocusColor handles focus state
	style := m.border.ApplyToStyle(lipgloss.NewStyle(), m.Focused())

	// Apply width for consistent full-width rendering (like ChatInput)
	// bw was already calculated at the start of View()
	if m.width > 0 {
		style = style.Width(m.width - bw)
	}

	return style.Render(result)
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

	// Create styles with background for consistent rendering.
	// Without this, styled text would show terminal default background
	// while the container's WithWhitespaceBackground only fills empty space.
	keyStyle := m.keyStyle
	descStyle := m.descStyle
	sepStyle := m.sepStyle
	if m.contentBg != nil {
		keyStyle = keyStyle.Background(m.contentBg)
		descStyle = descStyle.Background(m.contentBg)
		sepStyle = sepStyle.Background(m.contentBg)
	}

	var sb strings.Builder
	sb.WriteString(keyStyle.Render("Keyboard Shortcuts"))
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
		key := keyStyle.Render(b.Key + padding + "  ")
		desc := descStyle.Render(b.Desc)
		fmt.Fprintf(&sb, "%s%s\n", key, desc)
	}

	sb.WriteString("\n")
	sb.WriteString(sepStyle.Render("Press ? or Esc to close"))

	return m.renderHelpContent(sb.String(), m.width)
}

// renderHelpContent wraps content with border styling for the help window.
// Centered mode uses all-sides border styling; attached modes use inline config.
func (m Model) renderHelpContent(content string, width int) string {
	var style lipgloss.Style

	if m.helpMode == HelpModeCenter {
		// Centered: floating modal appearance with all-sides border
		style = m.border.ApplyToStyleCentered(lipgloss.NewStyle(), m.border.FocusColor)
	} else {
		// Attached modes (Down/Up/Top): use inline border config
		style = m.border.ApplyToStyle(lipgloss.NewStyle(), true)

		// Always reserve space for both borders (space characters when invisible)
		bw := 2
		if width > 0 {
			style = style.Width(width - bw)
		}
	}

	// Apply background color to border/padding areas
	if m.contentBg != nil {
		style = style.Background(m.contentBg)
	}

	result := style.Render(content)

	// Use Place() with WhitespaceBackground for consistent background fill
	if m.contentBg != nil {
		result = lipgloss.Place(
			lipgloss.Width(result),
			lipgloss.Height(result),
			lipgloss.Left,
			lipgloss.Top,
			result,
			lipgloss.WithWhitespaceBackground(m.contentBg),
		)
	}

	return result
}

// AnimationCmd returns a tick command if animation is in progress.
// Call this after ShowHelp()/HideHelp() when using animation programmatically.
// To check if animating: cmd := m.AnimationCmd(); isAnimating := cmd != nil
func (m Model) AnimationCmd() tea.Cmd {
	if m.animating {
		return heightTick()
	}
	return nil
}

// AnimatedHelpHeight returns the current animated height for layout calculations.
// Returns full HelpHeight() if animation is disabled or not animating.
// For animation padding (Up/overlay modes): padding = m.HelpHeight() - m.AnimatedHelpHeight()
func (m Model) AnimatedHelpHeight() int {
	if !m.animationEnabled || !m.animating {
		if m.showHelp {
			return m.HelpHeight()
		}
		return 0
	}
	return int(math.Ceil(m.animatedHeight))
}

// ViewHelpPartial renders help window with only visibleLines visible.
// Clips from top for HelpModeUp (bottom-anchored), from bottom for other modes (top-anchored).
// For Up mode, returns only the last N lines - parent handles positioning padding.
// For animated help: m.ViewHelpPartial(m.AnimatedHelpHeight())
func (m Model) ViewHelpPartial(visibleLines int) string {
	if visibleLines <= 0 {
		return ""
	}

	fullContent := m.ViewHelp()
	lines := strings.Split(fullContent, "\n")
	totalLines := len(lines)

	if visibleLines >= totalLines {
		return fullContent
	}

	if m.helpMode == HelpModeUp {
		// Bottom-anchored: show last N lines
		// Parent handles positioning padding for "grow upward" effect
		return strings.Join(lines[totalLines-visibleLines:], "\n")
	}

	// Top-anchored (Down, Top, Center): show first N lines
	return strings.Join(lines[:visibleLines], "\n")
}
