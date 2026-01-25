package styles

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/theme"
)

// Colors holds common color values shared across components.
type Colors struct {
	Border    lipgloss.TerminalColor
	Primary   lipgloss.TerminalColor
	ContentBg lipgloss.TerminalColor
	Text      lipgloss.TerminalColor
	TextMuted lipgloss.TerminalColor
}

// DefaultColors returns the default color palette.
func DefaultColors() Colors {
	return Colors{
		Border:    lipgloss.AdaptiveColor{Light: "#c0c0c0", Dark: "#4b4c5c"},
		Primary:   lipgloss.AdaptiveColor{Light: "#c96a2a", Dark: "#fab283"},
		ContentBg: nil,
		Text:      lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#e0e0e0"},
		TextMuted: lipgloss.AdaptiveColor{Light: "#666666", Dark: "#6a6a6a"},
	}
}

// ColorsFromPalette creates Colors from a theme palette.
func ColorsFromPalette(p *theme.Palette) Colors {
	if p == nil {
		return DefaultColors()
	}
	return Colors{
		Border:    p.Border,
		Primary:   p.Primary,
		ContentBg: p.BackgroundPanel,
		Text:      p.Text,
		TextMuted: p.TextMuted,
	}
}
