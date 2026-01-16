---
name: charmbracelet
version: "1.0"
description: This skill should be used when building terminal UIs with the Charmbracelet ecosystem. Covers bubbletea (Elm architecture TUI framework), lipgloss (styling), bubbles (components like list, viewport, spinner, progress, table, textarea, help), huh (forms), glamour (markdown), harmonica (animations), log (structured logging), and x/ experimental packages (teatest, ansi, term). Use when the user asks about "bubble tea model", "tea.Cmd", "tea.Msg", "lipgloss styling", "terminal colors", "TUI components", "viewport scrolling", "text input component", "spinner animation", "list.Model", "form validation with huh", "render markdown in terminal", "spring animation", "testing TUI apps with teatest", "REDUCE_MOTION", "accessibility in TUI", or "focus indicators".
---

# Charmbracelet Ecosystem Skill

Complete guidance for building elegant terminal applications with the Charmbracelet ecosystem.

## Quick Reference

| Package | Purpose | When to Use |
|---------|---------|-------------|
| **bubbletea** | TUI framework | Every TUI app - Model/Update/View pattern |
| **lipgloss** | Terminal styling | Colors, borders, layout, spacing |
| **bubbles** | Pre-built components | Lists, viewports, inputs, tables, spinners |
| **huh** | Forms & prompts | User input forms, confirmations |
| **glamour** | Markdown rendering | Display formatted text, documentation |
| **harmonica** | Spring animations | Smooth transitions, scroll, slide effects |
| **log** | Structured logging | File-based logging for TUI apps |
| **x/teatest** | TUI testing | Unit tests for Bubble Tea models |

## Core Concepts

### Elm Architecture (Bubble Tea)

Every Bubble Tea app follows this pattern:

```go
type Model struct {
    // All application state
}

func (m Model) Init() tea.Cmd {
    // Return initial commands (async loading, etc.)
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages, return new state and commands
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m Model) View() string {
    // Render UI as string
    return "Hello, World!"
}
```

### Commands and Messages

Commands perform side effects, messages carry results:

```go
// Command function
func loadData() tea.Msg {
    data, err := fetch()
    if err != nil {
        return errMsg{err}
    }
    return dataMsg{data}
}

// Typed message
type dataMsg struct{ data []Item }

// Return command from Update
return m, loadData

// Batch multiple commands
return m, tea.Batch(cmd1, cmd2, cmd3)
```

### Styling with Lip Gloss

```go
style := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(1, 2).
    Border(lipgloss.RoundedBorder())

output := style.Render("Styled text")

// Layout
row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
col := lipgloss.JoinVertical(lipgloss.Left, top, bottom)
```

## Reference Files

For detailed API documentation, read the appropriate reference file:

| Need | Read |
|------|------|
| Bubble Tea patterns, commands, messages | [references/bubbletea.md](references/bubbletea.md) |
| Styling, colors, borders, layout | [references/lipgloss.md](references/lipgloss.md) |
| List, viewport, input, spinner, table | [references/bubbles.md](references/bubbles.md) |
| Forms, select, confirm, validation | [references/huh.md](references/huh.md) |
| Markdown rendering, styles | [references/glamour.md](references/glamour.md) |
| Spring animations, smooth motion | [references/harmonica.md](references/harmonica.md) |
| Structured logging, formatters | [references/log.md](references/log.md) |
| TUI testing with teatest | [references/teatest.md](references/teatest.md) |
| ansi, term, other x/ utilities | [references/experimental.md](references/experimental.md) |

## Installation

```bash
# Core packages
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/huh
go get github.com/charmbracelet/glamour
go get github.com/charmbracelet/harmonica
go get github.com/charmbracelet/log

# Experimental
go get github.com/charmbracelet/x/exp/teatest
```

## Best Practices

1. **Always handle WindowSizeMsg** - Store dimensions for responsive layouts
2. **Never block in Init/Update** - Use commands for async operations
3. **Define styles as package variables** - Avoid recreating in View()
4. **Use typed messages** - Avoid interface{} for clarity
5. **Forward messages to sub-components** - They need tick/key messages
6. **Log to file for TUIs** - Avoid stdout interference
7. **Test with teatest** - Unit test your models

## Accessibility

Build accessible TUI applications:

### REDUCE_MOTION Support

Check for reduced motion preference and disable animations:

```go
import "os"

func shouldReduceMotion() bool {
    return os.Getenv("REDUCE_MOTION") == "1"
}

// In animation code
if shouldReduceMotion() {
    // Skip animation, apply final state immediately
    m.position = targetPosition
    return m, nil
}
// Normal animation with harmonica
```

See [references/harmonica.md](references/harmonica.md) for spring animation accessibility patterns.

### Focus Indicators

Always show visible focus indicators on interactive elements:

```go
func (m Model) View() string {
    style := m.styles.Item

    if m.focused {
        style = style.
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("#7D56F4"))
    }

    return style.Render(content)
}
```

### Color Independence

Never rely on color alone - use symbols alongside colors:

```go
// Good: Symbol + color
func statusView(isError bool) string {
    if isError {
        return lipgloss.NewStyle().
            Foreground(lipgloss.Color("9")).
            Render("[!] Error")
    }
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("10")).
        Render("[+] Success")
}
```

### High Contrast Support

Support high contrast mode via environment variable:

```go
func loadTheme() Theme {
    if os.Getenv("HIGH_CONTRAST") == "1" {
        return HighContrastTheme
    }
    return DefaultTheme
}
```

## Common Patterns

### Sub-Component Pattern
```go
type Model struct {
    list   list.Model
    input  textinput.Model
    active string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch m.active {
    case "list":
        m.list, cmd = m.list.Update(msg)
    case "input":
        m.input, cmd = m.input.Update(msg)
    }
    return m, cmd
}
```

### Responsive Layout Pattern
```go
func (m Model) View() string {
    style := lipgloss.NewStyle().Width(m.width - 4)
    return style.Render(content)
}
```

### Async Loading Pattern
```go
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.spinner.Tick,
        loadData,
    )
}
```
