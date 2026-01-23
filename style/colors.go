// Package style provides centralized styling for claude-agent-tui components.
package style

import "github.com/charmbracelet/lipgloss"

// Semantic color palette for consistent theming across components.
// Colors match the OpenCode palette for visual consistency.
var (
	// Primary is the accent color used for focused elements and highlights.
	Primary = lipgloss.Color("#fab283") // Orange/gold

	// Secondary is used for subtle accents and interactive elements.
	Secondary = lipgloss.Color("#5c9cf5") // Blue

	// Text is the primary text color.
	Text = lipgloss.Color("#e0e0e0") // Light text

	// TextMuted is used for secondary text, placeholders, and hints.
	TextMuted = lipgloss.Color("#6a6a6a") // Muted text

	// Surface is the background color for containers.
	Surface = lipgloss.Color("#212121") // Background

	// SurfaceHover is the background color for focused/hovered containers.
	SurfaceHover = lipgloss.Color("#252525") // Slightly lighter

	// Border is the color for container borders.
	Border = lipgloss.Color("#4b4c5c") // Border color
)
