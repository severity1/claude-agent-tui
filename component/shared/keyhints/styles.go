package keyhints

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

// Styles contains all visual styling for the KeyHints component.
// Use DefaultStyles() for sensible defaults or StylesFromPalette() to derive
// styles from a theme palette.
type Styles struct {
	// KeyStyle is the style for key names (e.g., "Enter", "Ctrl+C").
	KeyStyle lipgloss.Style
	// DescStyle is the style for key descriptions.
	DescStyle lipgloss.Style
	// SepStyle is the style for separators between hints.
	SepStyle lipgloss.Style
	// FocusStyle is the style applied when the component has focus.
	FocusStyle lipgloss.Style
	// HelpStyle is the style for the help window border.
	HelpStyle lipgloss.Style

	// BorderColor is the border color when not focused.
	BorderColor lipgloss.TerminalColor
	// PrimaryColor is the accent color used for focus indication.
	PrimaryColor lipgloss.TerminalColor
}

// DefaultStyles returns the default styles for KeyHints.
// These use neutral colors that work in most terminals.
func DefaultStyles() Styles {
	primary := lipgloss.AdaptiveColor{Light: "#c96a2a", Dark: "#fab283"}
	textMuted := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#6a6a6a"}
	border := lipgloss.AdaptiveColor{Light: "#c0c0c0", Dark: "#4b4c5c"}

	return Styles{
		KeyStyle:   lipgloss.NewStyle().Bold(true).Foreground(textMuted),
		DescStyle:  lipgloss.NewStyle().Foreground(textMuted),
		SepStyle:   lipgloss.NewStyle().Faint(true).Foreground(textMuted),
		FocusStyle: lipgloss.NewStyle().Foreground(primary),
		HelpStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(1, 2),

		BorderColor:  border,
		PrimaryColor: primary,
	}
}

// StylesFromPalette creates Styles from a theme Palette.
// This allows consistent theming across all components.
func StylesFromPalette(p *theme.Palette) Styles {
	if p == nil {
		return DefaultStyles()
	}

	return Styles{
		KeyStyle:   lipgloss.NewStyle().Bold(true).Foreground(p.TextMuted),
		DescStyle:  lipgloss.NewStyle().Foreground(p.TextMuted),
		SepStyle:   lipgloss.NewStyle().Faint(true).Foreground(p.TextMuted),
		FocusStyle: lipgloss.NewStyle().Foreground(p.Primary),
		HelpStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border).
			Padding(1, 2),

		BorderColor:  p.Border,
		PrimaryColor: p.Primary,
	}
}
