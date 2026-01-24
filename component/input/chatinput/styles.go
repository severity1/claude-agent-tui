package chatinput

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

// Styles contains all visual styling for the ChatInput component.
// Use DefaultStyles() for sensible defaults or StylesFromPalette() to derive
// styles from a theme palette.
type Styles struct {
	// PromptStyle is the style for the prompt prefix (e.g., "> ").
	PromptStyle lipgloss.Style
	// ModeStyle is the style for the mode indicator (e.g., "Build").
	ModeStyle lipgloss.Style
	// InfoStyle is the style for the info bar text.
	InfoStyle lipgloss.Style
	// PlaceholderStyle is the style for placeholder text.
	PlaceholderStyle lipgloss.Style

	// Border is the border style (Thick, Rounded, Normal, etc.).
	Border lipgloss.Border
	// BorderColor is the border color when not flashing.
	BorderColor lipgloss.TerminalColor
	// FlashBorderColor is the border color during copy flash animation.
	FlashBorderColor lipgloss.TerminalColor
	// FlashDuration is the duration of the copy flash animation.
	FlashDuration time.Duration

	// ContentBg is the background color for all content inside borders.
	ContentBg lipgloss.TerminalColor
	// TextColor is the primary text color.
	TextColor lipgloss.TerminalColor
	// TextMutedColor is the secondary/muted text color.
	TextMutedColor lipgloss.TerminalColor
	// PrimaryColor is the accent color for focused elements.
	PrimaryColor lipgloss.TerminalColor

	// BorderPaddingTop is the top padding inside the border.
	BorderPaddingTop int
	// BorderPaddingBottom is the bottom padding inside the border.
	BorderPaddingBottom int
}

// DefaultStyles returns the default styles for ChatInput.
// These use neutral colors that work in most terminals.
func DefaultStyles() Styles {
	primary := lipgloss.AdaptiveColor{Light: "#c96a2a", Dark: "#fab283"}
	text := lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#e0e0e0"}
	textMuted := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#6a6a6a"}
	surface := lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#212121"}
	border := lipgloss.AdaptiveColor{Light: "#c0c0c0", Dark: "#4b4c5c"}

	return Styles{
		PromptStyle: lipgloss.NewStyle().
			Padding(0, 0, 0, 1).
			Bold(true).
			Foreground(primary).
			Background(surface),
		ModeStyle: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Background(surface),
		InfoStyle: lipgloss.NewStyle().
			Foreground(textMuted).
			Background(surface),
		PlaceholderStyle: lipgloss.NewStyle().
			Foreground(textMuted).
			Background(surface),

		Border:           lipgloss.ThickBorder(),
		BorderColor:      border,
		FlashBorderColor: primary,
		FlashDuration:    150 * time.Millisecond,

		ContentBg:      surface,
		TextColor:      text,
		TextMutedColor: textMuted,
		PrimaryColor:   primary,

		BorderPaddingTop:    1,
		BorderPaddingBottom: 1,
	}
}

// StylesFromPalette creates Styles from a theme Palette.
// This allows consistent theming across all components.
func StylesFromPalette(p *theme.Palette) Styles {
	if p == nil {
		return DefaultStyles()
	}

	return Styles{
		PromptStyle: lipgloss.NewStyle().
			Padding(0, 0, 0, 1).
			Bold(true).
			Foreground(p.Primary).
			Background(p.BackgroundPanel),
		ModeStyle: lipgloss.NewStyle().
			Foreground(p.Primary).
			Bold(true).
			Background(p.BackgroundPanel),
		InfoStyle: lipgloss.NewStyle().
			Foreground(p.TextMuted).
			Background(p.BackgroundPanel),
		PlaceholderStyle: lipgloss.NewStyle().
			Foreground(p.TextMuted).
			Background(p.BackgroundPanel),

		Border:           lipgloss.ThickBorder(),
		BorderColor:      p.Border,
		FlashBorderColor: p.Primary,
		FlashDuration:    150 * time.Millisecond,

		ContentBg:      p.BackgroundPanel,
		TextColor:      p.Text,
		TextMutedColor: p.TextMuted,
		PrimaryColor:   p.Primary,

		BorderPaddingTop:    1,
		BorderPaddingBottom: 1,
	}
}
