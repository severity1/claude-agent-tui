package styles_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/internal/styles"
	"github.com/severity1/claude-agent-tui/theme"
)

func TestDefaultColors_ReturnsExpectedValues(t *testing.T) {
	colors := styles.DefaultColors()

	// Verify Border is set
	if colors.Border == nil {
		t.Error("DefaultColors().Border should not be nil")
	}

	// Verify Primary is set
	if colors.Primary == nil {
		t.Error("DefaultColors().Primary should not be nil")
	}

	// Verify ContentBg is nil (transparent by default)
	if colors.ContentBg != nil {
		t.Error("DefaultColors().ContentBg should be nil by default")
	}

	// Verify Text is set
	if colors.Text == nil {
		t.Error("DefaultColors().Text should not be nil")
	}

	// Verify TextMuted is set
	if colors.TextMuted == nil {
		t.Error("DefaultColors().TextMuted should not be nil")
	}
}

func TestColorsFromPalette_NilReturnsDefaults(t *testing.T) {
	colors := styles.ColorsFromPalette(nil)
	defaults := styles.DefaultColors()

	// Should return the same as DefaultColors when nil is passed
	if colors.Border == nil {
		t.Error("ColorsFromPalette(nil).Border should not be nil")
	}

	// Check that default border color is the same
	defaultBorder := defaults.Border.(lipgloss.AdaptiveColor)
	resultBorder := colors.Border.(lipgloss.AdaptiveColor)
	if defaultBorder.Light != resultBorder.Light || defaultBorder.Dark != resultBorder.Dark {
		t.Error("ColorsFromPalette(nil) should return DefaultColors() values")
	}
}

func TestColorsFromPalette_ExtractsCorrectFields(t *testing.T) {
	palette := &theme.Palette{
		Border:          lipgloss.AdaptiveColor{Light: "#111111", Dark: "#222222"},
		Primary:         lipgloss.AdaptiveColor{Light: "#333333", Dark: "#444444"},
		BackgroundPanel: lipgloss.AdaptiveColor{Light: "#555555", Dark: "#666666"},
		Text:            lipgloss.AdaptiveColor{Light: "#777777", Dark: "#888888"},
		TextMuted:       lipgloss.AdaptiveColor{Light: "#999999", Dark: "#aaaaaa"},
	}

	colors := styles.ColorsFromPalette(palette)

	// Verify Border is extracted correctly
	border := colors.Border.(lipgloss.AdaptiveColor)
	if border.Light != "#111111" || border.Dark != "#222222" {
		t.Errorf("ColorsFromPalette Border = %v, want {Light:#111111, Dark:#222222}", colors.Border)
	}

	// Verify Primary is extracted correctly
	primary := colors.Primary.(lipgloss.AdaptiveColor)
	if primary.Light != "#333333" || primary.Dark != "#444444" {
		t.Errorf("ColorsFromPalette Primary = %v, want {Light:#333333, Dark:#444444}", colors.Primary)
	}

	// Verify ContentBg is extracted correctly
	contentBg := colors.ContentBg.(lipgloss.AdaptiveColor)
	if contentBg.Light != "#555555" || contentBg.Dark != "#666666" {
		t.Errorf("ColorsFromPalette ContentBg = %v, want {Light:#555555, Dark:#666666}", colors.ContentBg)
	}

	// Verify Text is extracted correctly
	text := colors.Text.(lipgloss.AdaptiveColor)
	if text.Light != "#777777" || text.Dark != "#888888" {
		t.Errorf("ColorsFromPalette Text = %v, want {Light:#777777, Dark:#888888}", colors.Text)
	}

	// Verify TextMuted is extracted correctly
	textMuted := colors.TextMuted.(lipgloss.AdaptiveColor)
	if textMuted.Light != "#999999" || textMuted.Dark != "#aaaaaa" {
		t.Errorf("ColorsFromPalette TextMuted = %v, want {Light:#999999, Dark:#aaaaaa}", colors.TextMuted)
	}
}
