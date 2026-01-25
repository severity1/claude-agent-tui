// Package main demonstrates the ChatInput component with theme switching.
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
//   - 1-7: Switch theme (dynamically discovered from theme.List())
//   - Ctrl+C: Quit
//
// Features:
//   - Theme switching at runtime with 7 built-in themes
//   - Subtle background styling with focus highlight
//   - Responsive width that adapts to terminal size
//   - Character/line counter
//   - History navigation with draft preservation
//   - Keyhints component with collapsible help window
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/component/input/chatinput"
	"github.com/severity1/claude-agent-tui/component/shared/keyhints"
	"github.com/severity1/claude-agent-tui/theme"
	_ "github.com/severity1/claude-agent-tui/theme/themes" // Register built-in themes
)

// model holds the application state.
type model struct {
	input      chatinput.Model
	hints      keyhints.Model
	messages   []string
	width      int          // terminal width
	height     int          // terminal height
	initCmd    tea.Cmd      // command to run on Init (cursor blink)
	themeName  string       // current theme name
	themeList  []string     // dynamically discovered theme names
	themeInfos []theme.Info // theme metadata for descriptions
	flashMsg   string       // temporary message shown after theme switch
}

// newModel creates a new model with default state.
func newModel() model {
	// Use DefaultName() instead of hardcoded string
	themeName := theme.DefaultName()
	p := theme.Get(themeName)

	// Dynamically discover available themes
	themeList := theme.List()
	themeInfos := theme.ListInfo()

	// ChatInput with theme-based styling
	input := chatinput.New(
		chatinput.WithPalette(p),
		chatinput.WithPlaceholder("Type a message..."),
		chatinput.WithMinHeight(1),
		chatinput.WithMaxHeight(5),
		chatinput.WithShowCounter(true),
		chatinput.WithHistorySize(100),
		chatinput.WithPrompt("> "),
		chatinput.WithMode("Build"),
		chatinput.WithInfoBar("Claude Sonnet 4.5"),
		chatinput.WithFlashDuration(150*time.Millisecond),
		chatinput.WithBorderPadding(1, 1),
	)
	focusCmd := input.Focus()

	// Define keybindings for help display
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Alt+Enter", Desc: "Newline"},
		{Key: "Up/Down", Desc: "History"},
		{Key: "Ctrl+Y", Desc: "Copy"},
		{Key: "Esc", Desc: "Clear/Close"},
		{Key: "Tab", Desc: "Focus hints"},
		{Key: "1-7", Desc: "Theme"},
		{Key: "H", Desc: "Cycle help mode"},
		{Key: "?", Desc: "Help"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}

	// KeyHints with theme-based styling
	hints := keyhints.New(bindings,
		keyhints.WithPalette(p),
		keyhints.WithDefaultVisible(5),
		keyhints.WithCollapsed(),
		keyhints.WithSeparator(" | "),
	)

	return model{
		input:      input,
		hints:      hints,
		messages:   []string{},
		initCmd:    focusCmd,
		themeName:  themeName,
		themeList:  themeList,
		themeInfos: themeInfos,
	}
}

// getThemeDescription returns the description for a theme name.
func (m model) getThemeDescription(name string) string {
	for _, info := range m.themeInfos {
		if info.Name == name {
			return info.Description
		}
	}
	return ""
}

// switchTheme applies a new theme to all components.
func (m model) switchTheme(name string) model {
	p := theme.Get(name)
	if p == nil {
		return m
	}

	m.themeName = name
	m.flashMsg = fmt.Sprintf("Switched to: %s - %s", name, m.getThemeDescription(name))

	// Recreate input with new theme
	m.input = chatinput.New(
		chatinput.WithPalette(p),
		chatinput.WithPlaceholder("Type a message..."),
		chatinput.WithMinHeight(1),
		chatinput.WithMaxHeight(5),
		chatinput.WithShowCounter(true),
		chatinput.WithHistorySize(100),
		chatinput.WithPrompt("> "),
		chatinput.WithMode("Build"),
		chatinput.WithInfoBar("Claude Sonnet 4.5"),
		chatinput.WithWidth(m.width-2),
	)
	m.input.Focus()

	// Recreate keyhints bindings
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Alt+Enter", Desc: "Newline"},
		{Key: "Up/Down", Desc: "History"},
		{Key: "Ctrl+Y", Desc: "Copy"},
		{Key: "Esc", Desc: "Clear/Close"},
		{Key: "Tab", Desc: "Focus hints"},
		{Key: "1-7", Desc: "Theme"},
		{Key: "H", Desc: "Cycle help mode"},
		{Key: "?", Desc: "Help"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	// Calculate input width for consistency
	inputWidth := m.width - 2
	if inputWidth < 20 {
		inputWidth = 20
	}
	m.hints = keyhints.New(bindings,
		keyhints.WithPalette(p),
		keyhints.WithDefaultVisible(5),
		keyhints.WithCollapsed(),
		keyhints.WithSeparator(" | "),
		keyhints.WithWidth(inputWidth), // Match ChatInput width for border alignment
	)

	return m
}

// Init returns the cursor blink command for the already-focused input.
func (m model) Init() tea.Cmd {
	return m.initCmd
}

// Update handles messages and updates the model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputWidth := msg.Width - 2
		if inputWidth < 20 {
			inputWidth = 20
		}
		m.input.SetWidth(inputWidth)
		m.hints.SetWidth(inputWidth) // Match ChatInput width for border alignment
		return m, nil

	case tea.KeyMsg:
		// Clear flash message on any key press
		m.flashMsg = ""

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7":
			idx := int(msg.String()[0] - '1')
			if idx >= 0 && idx < len(m.themeList) {
				return m.switchTheme(m.themeList[idx]), nil
			}
		case "tab":
			if m.hints.ShowingHelp() {
				return m, nil
			}
			if m.input.Focused() {
				m.input.Blur()
				m.hints.Focus()
				return m, nil
			}
			m.hints.Blur()
			return m, m.input.Focus()
		case "esc":
			if m.hints.ShowingHelp() {
				m.hints.HideHelp()
				return m, nil
			}
			if m.hints.Focused() {
				m.hints.Blur()
				m.input.Focus()
				return m, nil
			}
		case "H":
			// Only cycle help modes when keyhints is focused or help is showing
			if m.hints.Focused() || m.hints.ShowingHelp() {
				// Cycle through help modes: Down -> Up -> Top -> Center -> Down
				modes := []keyhints.HelpMode{
					keyhints.HelpModeDown,
					keyhints.HelpModeUp,
					keyhints.HelpModeTop,
					keyhints.HelpModeCenter,
				}
				current := m.hints.HelpMode()
				nextIdx := 0
				for i, mode := range modes {
					if mode == current {
						nextIdx = (i + 1) % len(modes)
						break
					}
				}
				m.hints.SetHelpMode(modes[nextIdx])
				modeNames := map[keyhints.HelpMode]string{
					keyhints.HelpModeDown:   "Down",
					keyhints.HelpModeUp:     "Up",
					keyhints.HelpModeTop:    "Top",
					keyhints.HelpModeCenter: "Center",
				}
				m.flashMsg = fmt.Sprintf("Help mode: %s", modeNames[modes[nextIdx]])
				return m, nil
			}
		}

	case chatinput.SubmitMsg:
		m.messages = append(m.messages, msg.Text)
		return m, nil

	case chatinput.CopiedMsg:
		return m, nil
	}

	// Pass messages to focused component
	if m.hints.Focused() {
		updated, cmd := m.hints.Update(msg)
		if model, ok := updated.(keyhints.Model); ok {
			m.hints = model
		}
		return m, cmd
	}

	var cmd tea.Cmd
	updated, cmd := m.input.Update(msg)
	if model, ok := updated.(chatinput.Model); ok {
		m.input = model
	}
	return m, cmd
}

// View renders the current state.
func (m model) View() string {
	var sb strings.Builder

	// Map help modes to display names
	modeNames := map[keyhints.HelpMode]string{
		keyhints.HelpModeDown:   "Down",
		keyhints.HelpModeUp:     "Up",
		keyhints.HelpModeTop:    "Top",
		keyhints.HelpModeCenter: "Center",
	}

	// Header with theme indicator and description
	titleStyle := lipgloss.NewStyle().Bold(true)
	sb.WriteString(titleStyle.Render("ChatInput Theme Demo"))
	sb.WriteString(" - Theme: ")
	themeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	sb.WriteString(themeStyle.Render(m.themeName))
	sb.WriteString(" - Help Mode: ")
	modeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	sb.WriteString(modeStyle.Render(modeNames[m.hints.HelpMode()]))
	sb.WriteString("\n")

	// Show current theme description
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true)
	sb.WriteString(descStyle.Render(m.getThemeDescription(m.themeName)))
	sb.WriteString("\n")

	// Theme list - dynamically generated from theme.List()
	sb.WriteString("\nThemes: ")
	for i, name := range m.themeList {
		if name == m.themeName {
			fmt.Fprintf(&sb, "[%d:%s] ", i+1, name)
		} else {
			fmt.Fprintf(&sb, "%d:%s ", i+1, name)
		}
	}
	sb.WriteString("\n")

	// Show flash message if present
	if m.flashMsg != "" {
		flashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
		sb.WriteString(flashStyle.Render(m.flashMsg))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Display submitted messages
	if len(m.messages) > 0 {
		sb.WriteString("Submitted:\n")
		for i, msg := range m.messages {
			display := strings.ReplaceAll(msg, "\n", "\\n")
			fmt.Fprintf(&sb, "  %d. %s\n", i+1, display)
		}
		sb.WriteString("\n")
	}

	// Display input and keyhints/help based on mode
	if m.hints.ShowingHelp() {
		switch m.hints.HelpMode() {
		case keyhints.HelpModeUp:
			// Help above input
			sb.WriteString(m.hints.ViewHelp())
			sb.WriteString("\n")
			sb.WriteString(m.input.View())
		case keyhints.HelpModeTop:
			// Help at very top (already at top, same as Up but semantic)
			sb.WriteString(m.hints.ViewHelp())
			sb.WriteString("\n\n")
			sb.WriteString(m.input.View())
		case keyhints.HelpModeCenter:
			// Centered overlay (render after input)
			sb.WriteString(m.input.View())
			sb.WriteString("\n\n")
			sb.WriteString(m.hints.ViewHelp())
		default: // HelpModeDown
			sb.WriteString(m.input.View())
			sb.WriteString("\n\n")
			sb.WriteString(m.hints.ViewHelp())
		}
	} else {
		sb.WriteString(m.input.View())
		sb.WriteString("\n\n")
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
