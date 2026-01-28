// Package styles provides shared styling utilities for TUI components.
//
// This package consolidates common styling patterns used across components:
//   - BorderConfig: Configures border appearance, visibility, and padding
//   - Colors: Provides common color values shared across components
//
// # Border Configuration
//
// BorderConfig provides a unified way to configure borders with support for:
//   - Border style selection (Thick, Rounded, Normal, etc.)
//   - Color with optional focus color
//   - Independent left/right visibility control
//   - Per-side padding configuration
//
// Direct struct manipulation:
//
//	cfg := styles.DefaultBorderConfig()
//	cfg.Left = true
//	cfg.Color = lipgloss.Color("240")
//	style := cfg.ApplyToStyle(lipgloss.NewStyle(), false)
//
// # Functional Options Constructor
//
// NewBorderConfig creates a BorderConfig with functional options applied to defaults:
//
//	cfg := styles.NewBorderConfig(
//	    styles.WithBorderStyle(lipgloss.RoundedBorder()),
//	    styles.WithBorderLeft(true),
//	    styles.WithBorderColor(lipgloss.Color("240")),
//	    styles.WithBorderFocusColor(lipgloss.Color("39")),
//	    styles.WithBorderPadding(1, 1, 2, 2),
//	)
//
// Available options:
//   - WithBorderStyle(style lipgloss.Border): Set border style
//   - WithBorderLeft(enabled bool): Enable/disable left border visibility
//   - WithBorderRight(enabled bool): Enable/disable right border visibility
//   - WithBorderColor(color lipgloss.TerminalColor): Set border color
//   - WithBorderFocusColor(color lipgloss.TerminalColor): Set focused border color
//   - WithBorderPadding(top, bottom, left, right int): Set padding (negative values ignored)
//
// # Preset Constructors
//
//   - DefaultBorderConfig(): ThickBorder with right border visible, left padding of 1
//   - CenterBorderConfig(): NormalBorder with all sides, uniform padding of 1
//
// # Colors
//
// Colors provides a consistent base palette for components:
//
//	colors := styles.DefaultColors()
//	// or from a theme palette:
//	colors := styles.ColorsFromPalette(theme.Get(theme.Dracula))
package styles
