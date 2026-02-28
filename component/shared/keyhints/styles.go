package keyhints

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/internal/styles"
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
	// ContentBg is the background color for help content.
	ContentBg lipgloss.TerminalColor
}

// stylesFromColors creates Styles from a Colors struct.
// This is the single source of truth for style construction.
// The useDefaultContentBg parameter controls whether to use nil (transparent)
// or the colors.ContentBg value for the background.
func stylesFromColors(colors styles.Colors, useDefaultContentBg bool) Styles {
	contentBg := colors.ContentBg
	if useDefaultContentBg {
		contentBg = nil // keyhints uses nil by default (transparent)
	}
	return Styles{
		KeyStyle:   lipgloss.NewStyle().Bold(true).Foreground(colors.TextMuted),
		DescStyle:  lipgloss.NewStyle().Foreground(colors.TextMuted),
		SepStyle:   lipgloss.NewStyle().Faint(true).Foreground(colors.TextMuted),
		FocusStyle: lipgloss.NewStyle().Foreground(colors.Primary),
		HelpStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Border).
			Padding(1, 2),

		BorderColor:  colors.Border,
		PrimaryColor: colors.Primary,
		ContentBg:    contentBg,
	}
}

// DefaultStyles returns the default styles for KeyHints.
// These use neutral colors that work in most terminals.
func DefaultStyles() Styles {
	return stylesFromColors(styles.DefaultColors(), true)
}

// StylesFromPalette creates Styles from a theme Palette.
// This allows consistent theming across all components.
func StylesFromPalette(p *theme.Palette) Styles {
	if p == nil {
		return DefaultStyles()
	}
	return stylesFromColors(styles.ColorsFromPalette(p), false)
}
