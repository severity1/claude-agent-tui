// Package styles provides shared styling utilities for TUI components.
package styles

import "github.com/charmbracelet/lipgloss"

// BorderConfig holds common border configuration used by multiple components.
// This reduces duplication of border field definitions and styling logic.
type BorderConfig struct {
	Style         lipgloss.Border        // border style (Thick, Rounded, Normal, etc.)
	Color         lipgloss.TerminalColor // border color when not focused
	FocusColor    lipgloss.TerminalColor // border color when focused (optional)
	Left          bool                   // show left border
	Right         bool                   // show right border
	PaddingTop    int                    // top padding inside border
	PaddingBottom int                    // bottom padding inside border
	PaddingLeft   int                    // left padding inside border
	PaddingRight  int                    // right padding inside border
}

// DefaultBorderConfig returns sensible defaults for border configuration.
// Default: ThickBorder with right border visible, left padding of 1.
func DefaultBorderConfig() BorderConfig {
	return BorderConfig{
		Style:       lipgloss.ThickBorder(),
		Left:        false,
		Right:       true,
		PaddingLeft: 1,
	}
}

// ModifyBorderVisibility returns a copy of the border with invisible characters
// for hidden sides. This allows space reservation while hiding the border visually.
func ModifyBorderVisibility(border lipgloss.Border, showLeft, showRight bool) lipgloss.Border {
	result := border
	if !showLeft {
		result.Left = " "
		result.MiddleLeft = " "
	}
	if !showRight {
		result.Right = " "
		result.MiddleRight = " "
	}
	return result
}

// ApplyToStyle applies border config to a lipgloss.Style.
// When focused is true and FocusColor is set, FocusColor is used for the border.
func (c BorderConfig) ApplyToStyle(s lipgloss.Style, focused bool) lipgloss.Style {
	color := c.Color
	if focused && c.FocusColor != nil {
		color = c.FocusColor
	}

	borderStyle := ModifyBorderVisibility(c.Style, c.Left, c.Right)

	return s.
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(borderStyle).
		BorderForeground(color).
		PaddingTop(c.PaddingTop).
		PaddingBottom(c.PaddingBottom).
		PaddingLeft(c.PaddingLeft).
		PaddingRight(c.PaddingRight)
}

// SetPaddingIfValid sets padding values only if they are non-negative.
// Returns a new BorderConfig with validated padding values applied.
func (c BorderConfig) SetPaddingIfValid(top, bottom, left, right int) BorderConfig {
	if top >= 0 {
		c.PaddingTop = top
	}
	if bottom >= 0 {
		c.PaddingBottom = bottom
	}
	if left >= 0 {
		c.PaddingLeft = left
	}
	if right >= 0 {
		c.PaddingRight = right
	}
	return c
}

// ApplyToStyleWithColor applies border config with a specific color override.
// Use for dynamic color selection (e.g., flash animation in chatinput).
// The color parameter overrides c.Color and c.FocusColor.
func (c BorderConfig) ApplyToStyleWithColor(s lipgloss.Style, color lipgloss.TerminalColor) lipgloss.Style {
	borderStyle := ModifyBorderVisibility(c.Style, c.Left, c.Right)

	return s.
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(borderStyle).
		BorderForeground(color).
		PaddingTop(c.PaddingTop).
		PaddingBottom(c.PaddingBottom).
		PaddingLeft(c.PaddingLeft).
		PaddingRight(c.PaddingRight)
}

// ApplyToStyleCentered applies config for centered/modal windows (all 4 borders).
// Use for keyhints center mode styling. Unlike ApplyToStyle, this enables all four
// borders and ignores the Left/Right visibility settings.
func (c BorderConfig) ApplyToStyleCentered(s lipgloss.Style, color lipgloss.TerminalColor) lipgloss.Style {
	return s.
		Border(c.Style, true). // all 4 borders
		BorderForeground(color).
		PaddingTop(c.PaddingTop).
		PaddingBottom(c.PaddingBottom).
		PaddingLeft(c.PaddingLeft).
		PaddingRight(c.PaddingRight)
}

// CenterBorderConfig returns a BorderConfig suitable for centered/modal windows.
// Default: NormalBorder, all sides visible (not used by ApplyToStyleCentered),
// uniform padding of 1.
func CenterBorderConfig() BorderConfig {
	return BorderConfig{
		Style:         lipgloss.NormalBorder(),
		Left:          true, // not used by ApplyToStyleCentered, but semantically correct
		Right:         true,
		PaddingTop:    1,
		PaddingBottom: 1,
		PaddingLeft:   1,
		PaddingRight:  1,
	}
}

// BorderConfigOption is a function that modifies a BorderConfig.
type BorderConfigOption func(*BorderConfig)

// WithBorderStyle sets the border style.
func WithBorderStyle(style lipgloss.Border) BorderConfigOption {
	return func(c *BorderConfig) { c.Style = style }
}

// WithBorderLeft enables or disables left border visibility.
func WithBorderLeft(enabled bool) BorderConfigOption {
	return func(c *BorderConfig) { c.Left = enabled }
}

// WithBorderRight enables or disables right border visibility.
func WithBorderRight(enabled bool) BorderConfigOption {
	return func(c *BorderConfig) { c.Right = enabled }
}

// WithBorderColor sets the border color.
func WithBorderColor(color lipgloss.TerminalColor) BorderConfigOption {
	return func(c *BorderConfig) { c.Color = color }
}

// WithBorderFocusColor sets the border color when focused.
func WithBorderFocusColor(color lipgloss.TerminalColor) BorderConfigOption {
	return func(c *BorderConfig) { c.FocusColor = color }
}

// WithBorderPadding sets padding values (negative values ignored).
func WithBorderPadding(top, bottom, left, right int) BorderConfigOption {
	return func(c *BorderConfig) {
		*c = c.SetPaddingIfValid(top, bottom, left, right)
	}
}

// NewBorderConfig creates a BorderConfig with the given options applied to defaults.
func NewBorderConfig(opts ...BorderConfigOption) BorderConfig {
	c := DefaultBorderConfig()
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
