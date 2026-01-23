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
//   - Character/line counter
//   - History navigation with draft preservation
//   - Keyhints component with collapsible help window
//   - Focus switching between input and hints
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/severity1/claude-agent-tui/component/input/chatinput"
	"github.com/severity1/claude-agent-tui/component/shared/keyhints"
)

// model holds the application state.
type model struct {
	input    chatinput.Model
	hints    keyhints.Model
	messages []string
	copied   bool // show copy feedback
}

// newModel creates a new model with default state.
func newModel() model {
	input := chatinput.New(
		chatinput.WithPlaceholder("Type a message..."),
		chatinput.WithWidth(60),
		chatinput.WithHeight(3),
		chatinput.WithShowCounter(true),
	)
	// Focus the input so it accepts keyboard input immediately
	input.Focus()

	hints := keyhints.New([]keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Alt+Enter", Desc: "Newline"},
		{Key: "Up/Down", Desc: "History"},
		{Key: "Ctrl+Y", Desc: "Copy"},
		{Key: "Esc", Desc: "Clear/Close"},
		{Key: "Tab", Desc: "Focus hints"},
		{Key: "?", Desc: "Help (focused)"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}, keyhints.WithCollapsed(), keyhints.WithDefaultVisible(4))

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
			} else {
				m.hints.Blur()
				m.input.Focus()
			}
			return m, nil
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

	// Display the input
	sb.WriteString("Input:\n")
	sb.WriteString(m.input.View())
	sb.WriteString("\n")

	// Show copy feedback
	if m.copied {
		sb.WriteString("[Copied to clipboard]\n")
	}
	sb.WriteString("\n")

	// Display keyhints or help window
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
