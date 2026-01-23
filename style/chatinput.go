package style

import "github.com/charmbracelet/lipgloss"

// ChatInput component styles.
var (
	// InputPrompt is the style for the ">" prompt prefix.
	// Uses Surface background to match the textarea and container.
	// The prompt uses bold text with the primary accent color.
	InputPrompt = lipgloss.NewStyle().
		Padding(0, 0, 0, 1).
		Bold(true).
		Foreground(Primary).
		Background(Surface)
)
