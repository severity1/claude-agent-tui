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
// Example usage:
//
//	cfg := styles.DefaultBorderConfig()
//	cfg.Left = true
//	cfg.Color = lipgloss.Color("240")
//	style := cfg.ApplyToStyle(lipgloss.NewStyle(), false)
//
// # Colors
//
// Colors provides a consistent base palette for components:
//
//	colors := styles.DefaultColors()
//	// or from a theme palette:
//	colors := styles.ColorsFromPalette(theme.Get("dracula"))
package styles
