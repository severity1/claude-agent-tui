// Package border provides shared border configuration utilities for TUI components.
package border

import "github.com/charmbracelet/lipgloss"

// Config holds common border configuration used by multiple components.
// This reduces duplication of border field definitions and styling logic.
type Config struct {
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

// DefaultConfig returns sensible defaults for border configuration.
// Default: ThickBorder with right border visible, left padding of 1.
func DefaultConfig() Config {
	return Config{
		Style:       lipgloss.ThickBorder(),
		Left:        false,
		Right:       true,
		PaddingLeft: 1,
	}
}

// ModifyVisibility returns a copy of the border with invisible characters
// for hidden sides. This allows space reservation while hiding the border visually.
func ModifyVisibility(border lipgloss.Border, showLeft, showRight bool) lipgloss.Border {
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
func (c Config) ApplyToStyle(s lipgloss.Style, focused bool) lipgloss.Style {
	color := c.Color
	if focused && c.FocusColor != nil {
		color = c.FocusColor
	}

	borderStyle := ModifyVisibility(c.Style, c.Left, c.Right)

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
// Returns a new Config with validated padding values applied.
func (c Config) SetPaddingIfValid(top, bottom, left, right int) Config {
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
