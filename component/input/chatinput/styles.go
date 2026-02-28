package chatinput

import (
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/internal/styles"
	"github.com/severity1/claude-agent-tui/theme"
)

// Styles contains all visual styling for the ChatInput component.
// Use DefaultStyles() for sensible defaults or StylesFromPalette() to derive
// styles from a theme palette.
type Styles struct {
	// PromptStyle is the style for the prompt prefix when unfocused (e.g., "> ").
	PromptStyle lipgloss.Style
	// FocusedPromptStyle is the style for the prompt prefix when focused.
	// This allows the prompt to "light up" when the input has focus.
	FocusedPromptStyle lipgloss.Style
	// ModeStyle is the style for the mode indicator (e.g., "Build").
	ModeStyle lipgloss.Style
	// InfoStyle is the style for the info bar text.
	InfoStyle lipgloss.Style
	// PlaceholderStyle is the style for placeholder text.
	// Note: This style's colors (Foreground=TextMutedColor, Background=ContentBg)
	// are applied via configureTextareaStyle() when WithPalette/WithStyles is used.
	// Custom placeholder rendering also uses these same colors directly.
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

// stylesFromColors creates Styles from a Colors struct.
// This is the single source of truth for style construction.
func stylesFromColors(colors styles.Colors) Styles {
	return Styles{
		PromptStyle: lipgloss.NewStyle().
			Padding(0, 0, 0, 1).
			Bold(true).
			Foreground(colors.TextMuted).
			Background(colors.ContentBg),
		FocusedPromptStyle: lipgloss.NewStyle().
			Padding(0, 0, 0, 1).
			Bold(true).
			Foreground(colors.Primary).
			Background(colors.ContentBg),
		ModeStyle: lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true).
			Background(colors.ContentBg),
		InfoStyle: lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Background(colors.ContentBg),
		PlaceholderStyle: lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Background(colors.ContentBg),

		Border:           lipgloss.ThickBorder(),
		BorderColor:      colors.Border,
		FlashBorderColor: colors.Primary,
		FlashDuration:    150 * time.Millisecond,

		ContentBg:      colors.ContentBg,
		TextColor:      colors.Text,
		TextMutedColor: colors.TextMuted,
		PrimaryColor:   colors.Primary,

		BorderPaddingTop:    1,
		BorderPaddingBottom: 1,
	}
}

// DefaultStyles returns the default styles for ChatInput.
// These use neutral colors that work in most terminals.
func DefaultStyles() Styles {
	return stylesFromColors(styles.DefaultColors())
}

// StylesFromPalette creates Styles from a theme Palette.
// This allows consistent theming across all components.
func StylesFromPalette(p *theme.Palette) Styles {
	if p == nil {
		return DefaultStyles()
	}
	return stylesFromColors(styles.ColorsFromPalette(p))
}
