# Bubble Tea Reference

**Package**: `github.com/charmbracelet/bubbletea`
**Documentation**: https://github.com/charmbracelet/bubbletea
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/bubbletea

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Core Types](#core-types) - Model, Cmd, Msg interfaces
- [The Elm Architecture](#the-elm-architecture) - Init, Update, View pattern
- [Commands](#commands) - Async operations and side effects
- [Built-in Messages](#built-in-messages) - WindowSizeMsg, KeyMsg, etc.
- [Key Message Handling](#key-message-handling) - Keyboard input processing
- [Program Options](#program-options) - WithAltScreen, WithMouseCellMotion, etc.
- [Sending Messages from Outside](#sending-messages-from-outside) - Program.Send
- [Common Patterns](#common-patterns) - Sub-models, batching, ticks

## Overview

Bubble Tea is a TUI framework based on the Elm architecture. Every application follows the Model-Update-View pattern with unidirectional data flow.

## Core Types

```go
// Model interface - implement this for your application state
type Model interface {
    Init() Cmd           // Initial command to run
    Update(Msg) (Model, Cmd)  // Handle messages, return new state
    View() string        // Render UI as string
}

// Msg is any message type (use custom structs)
type Msg interface{}

// Cmd is a function that returns a message
type Cmd func() Msg
```

## The Elm Architecture

```go
type Model struct {
    items    []string
    cursor   int
    width    int
    height   int
    ready    bool
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        loadData,            // Start async loading
        tea.EnterAltScreen,  // Full-screen mode
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.ready = true
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "j", "down":
            m.cursor++
        case "k", "up":
            m.cursor--
        }
    case dataLoadedMsg:
        m.items = msg.items
    }
    return m, nil
}

func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    // Render your UI
    return renderList(m.items, m.cursor)
}
```

## Commands

Commands perform side effects and return messages:

```go
// Simple command
func loadData() tea.Msg {
    data, err := fetchFromAPI()
    if err != nil {
        return errMsg{err}
    }
    return dataLoadedMsg{data}
}

// Command with arguments (use closure)
func loadItem(id string) tea.Cmd {
    return func() tea.Msg {
        item, err := fetchItem(id)
        if err != nil {
            return errMsg{err}
        }
        return itemLoadedMsg{item}
    }
}

// Batch multiple commands (run in parallel)
return m, tea.Batch(cmd1, cmd2, cmd3)

// Sequence commands (run in order)
return m, tea.Sequence(first, second, third)
```

## Built-in Messages

```go
tea.KeyMsg{Type: tea.KeyEnter}      // Keyboard input
tea.MouseMsg{X: 10, Y: 5}           // Mouse events
tea.WindowSizeMsg{Width, Height}   // Terminal resize
tea.FocusMsg / tea.BlurMsg         // Window focus
tea.QuitMsg{}                       // Quit signal
```

## Key Message Handling

```go
case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyEnter:
        // Enter pressed
    case tea.KeyEsc:
        // Escape pressed
    case tea.KeyCtrlC:
        return m, tea.Quit
    case tea.KeyRunes:
        // Character input: msg.Runes
        char := string(msg.Runes)
    }

    // Or use String() for simple matching
    switch msg.String() {
    case "q":
        return m, tea.Quit
    case "j", "down":
        m.cursor++
    case "ctrl+s":
        return m, save
    }
```

## Program Options

```go
p := tea.NewProgram(
    initialModel,
    tea.WithAltScreen(),           // Full-screen mode
    tea.WithMouseCellMotion(),     // Mouse click/motion
    tea.WithMouseAllMotion(),      // All mouse (for drag)
    tea.WithOutput(os.Stderr),     // Custom output
    tea.WithInput(os.Stdin),       // Custom input
    tea.WithoutCatchPanics(),      // Debug panics
    tea.WithFilter(filterFunc),    // Filter messages
)

// Run blocking
if _, err := p.Run(); err != nil {
    log.Fatal(err)
}
```

## Sending Messages from Outside

```go
// Get program reference
p := tea.NewProgram(model)

// Send message from another goroutine
go func() {
    time.Sleep(5 * time.Second)
    p.Send(timeoutMsg{})
}()

p.Run()
```

## Common Patterns

### Async Loading Pattern
```go
func (m Model) Init() tea.Cmd {
    return m.loadData
}

func (m Model) loadData() tea.Msg {
    // This runs in a goroutine
    data, err := slowOperation()
    if err != nil {
        return errMsg{err}
    }
    return dataMsg{data}
}
```

### Tick/Timer Pattern
```go
type tickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        m.elapsed++
        return m, tick() // Schedule next tick
    }
    return m, nil
}
```

### Sub-Model Pattern
```go
type Model struct {
    list   list.Model  // Sub-component
    input  textinput.Model
    active string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Forward to active component
    switch m.active {
    case "list":
        newList, cmd := m.list.Update(msg)
        m.list = newList
        cmds = append(cmds, cmd)
    case "input":
        newInput, cmd := m.input.Update(msg)
        m.input = newInput
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}
```

## Best Practices

1. **Always handle WindowSizeMsg** - Store dimensions for responsive layouts
2. **Never block in Init/Update** - Use commands for async work
3. **Type your messages** - Avoid `interface{}` for clarity
4. **Batch independent commands** - More efficient than sequential
5. **Handle context cancellation** - Respect `ctx.Done()` in commands

## Common Pitfalls

- Forgetting to return `tea.Quit` on exit
- Blocking operations in Update (use Cmd instead)
- Not handling WindowSizeMsg (breaks on resize)
- Mutating model without returning it
- Infinite loops in tick commands
