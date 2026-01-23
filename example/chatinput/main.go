// Package main demonstrates the ChatInput component for text input.
//
// Usage:
//
//	go run example/chatinput/main.go
//
// Controls:
//   - Enter: Submit text
//   - Alt+Enter / Ctrl+J: Insert newline
//   - Up/Down: Navigate history (preserves draft)
//   - Ctrl+Y: Copy to clipboard (recommended)
//   - Ctrl+Shift+C: Copy to clipboard (if terminal supports it)
//   - Esc: Clear input / close help / return to input
//   - Tab: Switch focus between input and keyhints
//   - ?: Open help window (when keyhints is focused)
//   - Ctrl+C: Quit
//
// Features:
//   - Subtle background styling with focus highlight (lighter when focused)
//   - Responsive width that adapts to terminal size
//   - Optional prompt prefix (configured via WithPrompt)
//   - Character/line counter
//   - History navigation with draft preservation
//   - Keyhints component with collapsible help window
//   - Focus switching between input and hints
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/component/input/chatinput"
	"github.com/severity1/claude-agent-tui/component/shared/keyhints"
	"github.com/severity1/claude-agent-tui/style"
)

// model holds the application state.
type model struct {
	input    chatinput.Model
	hints    keyhints.Model
	messages []string
	copied   bool // show copy feedback
	width    int  // terminal width
	height   int  // terminal height
}

// newModel creates a new model with default state.
// This example demonstrates ALL configuration options available for ChatInput and KeyHints.
// Many options use explicit defaults to showcase the API.
func newModel() model {
	// ChatInput with ALL configuration options shown explicitly
	input := chatinput.New(
		// Content behavior options
		chatinput.WithPlaceholder("Type a message..."),
		chatinput.WithMinHeight(1), // Default: 1 - single line minimum
		chatinput.WithMaxHeight(5), // Default: 10 - lines before scrolling
		chatinput.WithShowCounter(true),
		chatinput.WithHistorySize(100), // Default: 100 entries

		// Prompt configuration
		chatinput.WithPrompt("> "),                   // Default: "" (empty) - add prompt prefix
		chatinput.WithPromptStyle(style.InputPrompt), // Default style for prompt

		// Info bar configuration
		chatinput.WithMode("Build"),
		chatinput.WithModeStyle(lipgloss.NewStyle().Foreground(style.Primary).Bold(true).Background(style.Surface)),
		chatinput.WithInfoBar("Claude Sonnet 4.5"),
		chatinput.WithInfoStyle(lipgloss.NewStyle().Foreground(style.TextMuted).Background(style.Surface)),

		// Visual styling options (all showing explicit defaults)
		chatinput.WithContentBackground(style.Surface),    // Background for content inside borders
		chatinput.WithBorderStyle(lipgloss.ThickBorder()), // Border style (Thick, Rounded, Normal)
		chatinput.WithBorderColor(style.Border),           // Normal border color
		chatinput.WithFlashBorderColor(style.Primary),     // Border color during copy flash
		chatinput.WithFlashDuration(150*time.Millisecond), // Copy flash animation duration
		chatinput.WithBorderPadding(1, 1),                 // Internal padding (top, bottom)
	)
	// Focus the input so it accepts keyboard input immediately
	input.Focus()

	// Define keybindings for help display
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Alt+Enter", Desc: "Newline"},
		{Key: "Up/Down", Desc: "History"},
		{Key: "Ctrl+Y", Desc: "Copy"},
		{Key: "Esc", Desc: "Clear/Close"},
		{Key: "Tab", Desc: "Focus hints"},
		{Key: "?", Desc: "Help (focused)"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}

	// KeyHints with ALL configuration options shown explicitly
	hints := keyhints.New(bindings,
		// Display configuration
		keyhints.WithDefaultVisible(4), // Default: 3 - hints shown when collapsed
		keyhints.WithCollapsed(),       // Start collapsed (default: expanded)
		keyhints.WithSeparator(" | "),  // Default separator between hints
		keyhints.WithPadding(""),       // Default: "" - left padding for inline views

		// Style configuration (all showing explicit defaults)
		keyhints.WithKeyStyle(lipgloss.NewStyle().Bold(true).Foreground(style.TextMuted)),
		keyhints.WithDescStyle(lipgloss.NewStyle().Foreground(style.TextMuted)),
		keyhints.WithSepStyle(lipgloss.NewStyle().Faint(true).Foreground(style.TextMuted)),
		keyhints.WithFocusStyle(lipgloss.NewStyle().Foreground(style.Primary)), // Used when component is focused
		keyhints.WithHelpStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(style.Border).Padding(1, 2)),
	)

	return model{
		input:    input,
		hints:    hints,
		messages: []string{},
	}
}

// Init returns the cursor blink command for the already-focused input.
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages and updates the model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Full width minus small padding (no border to account for now)
		inputWidth := msg.Width - 2
		if inputWidth < 20 {
			inputWidth = 20 // minimum width
		}
		m.input.SetWidth(inputWidth)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Don't switch focus if help is showing
			if m.hints.ShowingHelp() {
				return m, nil
			}
			// Toggle focus between input and hints
			if m.input.Focused() {
				m.input.Blur()
				m.hints.Focus()
				return m, nil
			}
			// Return focus to input - return blink command for cursor
			m.hints.Blur()
			return m, m.input.Focus()
		case "esc":
			// If help is showing, close it first
			if m.hints.ShowingHelp() {
				m.hints.HideHelp()
				return m, nil
			}
			// If hints focused, return to input
			if m.hints.Focused() {
				m.hints.Blur()
				m.input.Focus()
				return m, nil
			}
			// Otherwise let chatinput handle it (clear)
		}

	case chatinput.SubmitMsg:
		// Add the submitted text to messages
		m.messages = append(m.messages, msg.Text)
		m.copied = false
		return m, nil

	case chatinput.CopiedMsg:
		// Show copy feedback
		if msg.Err == nil {
			m.copied = true
		}
		return m, nil
	}

	// Pass messages to focused component
	if m.hints.Focused() {
		updated, cmd := m.hints.Update(msg)
		m.hints = updated.(keyhints.Model)
		return m, cmd
	}

	// Pass other messages to chatinput
	var cmd tea.Cmd
	updated, cmd := m.input.Update(msg)
	m.input = updated.(chatinput.Model)
	return m, cmd
}

// View renders the current state.
func (m model) View() string {
	var sb strings.Builder

	sb.WriteString("ChatInput Component Demo\n")
	sb.WriteString("========================\n\n")

	// Display submitted messages
	if len(m.messages) > 0 {
		sb.WriteString("Submitted messages:\n")
		for i, msg := range m.messages {
			// Show message with line breaks escaped for visibility
			display := strings.ReplaceAll(msg, "\n", "\\n")
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, display))
		}
		sb.WriteString("\n")
	}

	// Display the input (now with built-in bordered styling)
	sb.WriteString(m.input.View())
	sb.WriteString("\n")

	// Show copy feedback
	if m.copied {
		sb.WriteString("[Copied to clipboard]\n")
	}
	sb.WriteString("\n")

	// Display keyhints or help window
	// Padding is now built into the keyhints component (default 1 space)
	if m.hints.ShowingHelp() {
		sb.WriteString(m.hints.ViewHelp())
	} else {
		sb.WriteString(m.hints.View())
	}
	sb.WriteString("\n")

	return sb.String()
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
