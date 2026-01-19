# claude-agent-tui

Reusable TUI components for building Claude Code interfaces with the Charmbracelet ecosystem.

## Overview

`claude-agent-tui` provides a comprehensive library of terminal UI components built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) that integrate with [claude-agent-sdk-go](https://github.com/severity1/claude-agent-sdk-go). It enables Go developers to build custom interactive CLI applications that leverage Claude Code's agentic capabilities.

## Features

### Implemented

- **Stream Adapter** - Converts SDK messages to Bubble Tea messages with 16+ message types
- **Validation Example** - Demonstrates all adapter message types with live and synthetic modes

### Planned

- **Component Library** - 27+ reusable Bubble Tea components
- **Multiple Variants** - Each component supports multiple display modes (fullscreen, modal, inline, etc.)
- **Theme System** - 12 built-in themes including Catppuccin, Dracula, Nord, High Contrast, and more
- **Animation System** - Spring-based animations with harmonica physics
- **Mouse Support** - Full mouse interaction with click zones, hover states, and scroll
- **Keybinding System** - Customizable keybindings with vim-style defaults
- **Responsive Layout** - Adaptive layouts for different terminal sizes
- **Full SDK Coverage** - Components for all SDK tools including background tasks and sessions
- **Error Recovery** - Automatic reconnection with exponential backoff
- **Accessibility** - Screen reader support, reduced motion, high contrast themes

## Architecture

The library uses a layered architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                      Your Application                        │
├─────────────────────────────────────────────────────────────┤
│  layout/          High-level composite screens     [PLANNED] │
│  ├── chat/        Full chat interface                       │
│  └── agent/       Agent dashboard with tools view           │
├─────────────────────────────────────────────────────────────┤
│  component/       Low-level Bubble Tea components  [PLANNED] │
│  ├── output/      StreamText, Message, Thinking, ToolUse... │
│  ├── input/       ChatInput, Permission, Question, Todo...  │
│  └── shared/      Button, Chip, Card, Overlay               │
├─────────────────────────────────────────────────────────────┤
│  system/          Cross-cutting systems            [PLANNED] │
│  ├── theme/       Theming engine with built-in themes       │
│  ├── animation/   Spring-based animations                   │
│  ├── keymap/      Keybinding management                     │
│  ├── mouse/       Mouse event routing                       │
│  ├── responsive/  Breakpoints and layout constraints        │
│  ├── focus/       Focus ring and traversal                  │
│  ├── recovery/    Error recovery and reconnection           │
│  └── accessibility/ Screen reader and reduced motion        │
├─────────────────────────────────────────────────────────────┤
│  adapter/         SDK integration layer                      │
│  ├── stream.go    StreamEvent -> tea.Msg       [IMPLEMENTED] │
│  ├── control.go   canUseTool -> input prompts      [PLANNED] │
│  └── client.go    Client lifecycle management      [PLANNED] │
├─────────────────────────────────────────────────────────────┤
│  claude-agent-sdk-go                                         │
└─────────────────────────────────────────────────────────────┘
```

## Components (Planned)

### Output Components

| Component | Purpose | Variants |
|-----------|---------|----------|
| `StreamText` | Live streaming text with cursor | View, Inline |
| `Message` | Conversation message bubble | View, Inline, Card |
| `Thinking` | Collapsible extended thinking | View, Collapsible, Modal |
| `ToolUse` | Tool invocation with I/O display | View, Card, Inline, Modal |
| `CodeBlock` | Syntax-highlighted code | View, Inline, Modal, Window |
| `Status` | Token count, cost, model info | Bar, Inline, Compact |
| `Toast` | Notification messages | Toast |
| `DiffView` | Side-by-side or unified diff | View, Modal, Window |
| `GlobResults` | File pattern match results | View, Sidebar, Modal |
| `GrepResults` | Text search match results | View, Sidebar, Modal |
| `TaskStatus` | Background task progress | Bar, Compact |

### Input Components

| Component | SDK Trigger | Purpose | Variants |
|-----------|-------------|---------|----------|
| `ChatInput` | Default | Main input with mode toggles | - |
| `Permission` | canUseTool | Tool approval (Allow/Deny/Always) | Modal, Overlay, Inline |
| `Question` | AskUserQuestion | Multi-choice question response | Modal, Overlay, View |
| `PlanEnter` | EnterPlanMode | Confirm entering plan mode | Modal, Overlay |
| `PlanExit` | ExitPlanMode | Review plan and decide action | Modal, Fullscreen, Window |
| `Todo` | TodoWrite | Task list display | View, Sidebar, Collapsible |
| `Interrupt` | Ctrl+C | Stop confirmation | Modal, Overlay |
| `Confirm` | Generic | Yes/no prompt | Modal, Inline |
| `FilePicker` | File selection | File browser dialog | Modal, Overlay, Window |
| `TaskManager` | TaskOutput/KillShell | Background task management | Sidebar, Overlay, Modal |
| `Session` | Session control | Session picker and management | Modal, Overlay, Sidebar |
| `Search` | Message search | Search message history | Overlay, Modal, Inline |

### Shared Components

| Component | Purpose |
|-----------|---------|
| `Button` | Clickable button with states |
| `Chip` | Toggle pill/chip |
| `Card` | Container with border |
| `Overlay` | Modal backdrop wrapper |

## Themes (Planned)

Built-in themes with light/dark variants:

| Theme | Variants |
|-------|----------|
| Catppuccin | Mocha, Latte, Frappe, Macchiato |
| Dracula | - |
| Nord | - |
| Gruvbox | Dark, Light |
| Tokyo Night | - |
| One Dark | - |
| High Contrast | WCAG AAA compliant |
| System | Adaptive |

```go
import "github.com/severity1/claude-agent-tui/system/theme"

// Use a built-in theme
t := theme.Get("catppuccin-mocha")

// Apply to components
input := chatinput.New(
    chatinput.WithStyles(t.ChatInput()),
)
```

## Quick Start (Planned)

The full chat interface is planned for future implementation:

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/severity1/claude-agent-tui/layout/chat"
    "github.com/severity1/claude-agent-tui/system/theme"
    claudecode "github.com/severity1/claude-agent-sdk-go"
)

func main() {
    // Create chat layout with theme and SDK client
    chatLayout := chat.New(
        chat.WithClient(claudecode.NewClient()),
        chat.WithTheme(theme.Get("dracula")),
    )

    // Run with mouse support
    p := tea.NewProgram(
        chatLayout,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Using the Stream Adapter (Implemented)

The stream adapter converts SDK messages to Bubble Tea messages, enabling reactive TUI updates:

```go
package main

import (
    "context"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/severity1/claude-agent-tui/adapter"
    claudecode "github.com/severity1/claude-agent-sdk-go"
)

type Model struct {
    client   *claudecode.Client
    messages chan any
    ctx      context.Context
    content  string
}

func (m Model) Init() tea.Cmd {
    // Start streaming SDK messages
    return adapter.StreamCmd(m.ctx, m.messages)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.StreamDeltaMsg:
        // Incremental text content
        m.content += msg.Delta
        return m, adapter.StreamCmd(m.ctx, m.messages)

    case adapter.StreamBlockStartMsg:
        // New content block (text, thinking, tool_use)
        if msg.BlockType == adapter.BlockTypeToolUse {
            // Tool use started: msg.ToolName, msg.ToolID
        }
        return m, adapter.StreamCmd(m.ctx, m.messages)

    case adapter.ThinkingDeltaMsg:
        // Extended thinking content
        return m, adapter.StreamCmd(m.ctx, m.messages)

    case adapter.ToolUseDeltaMsg:
        // Incremental tool input JSON
        return m, adapter.StreamCmd(m.ctx, m.messages)

    case adapter.ResultMsg:
        // Final result with usage stats
        // msg.InputTokens, msg.OutputTokens, msg.CacheReadTokens
        return m, nil

    case adapter.StreamDoneMsg:
        // Stream completed
        return m, nil

    case adapter.StreamErrorMsg:
        // Handle error: msg.Err
        return m, nil
    }
    return m, nil
}
```

### Adapter Message Types

| Message Type | Purpose |
|--------------|---------|
| `StreamDeltaMsg` | Incremental text content |
| `StreamBlockStartMsg` | New content block (text, thinking, tool_use) |
| `StreamBlockStopMsg` | Content block completion |
| `ThinkingDeltaMsg` | Incremental extended thinking |
| `ToolUseDeltaMsg` | Incremental tool input JSON |
| `MessageStartMsg` | Initial message metadata with model info |
| `MessageDeltaMsg` | Final usage stats and stop reason |
| `AssistantMsg` | Complete assistant message |
| `ResultMsg` | Final result with usage statistics |
| `UserMsg` | User message (typically tool results) |
| `SystemInitMsg` | Session initialization data |
| `SystemHookResponseMsg` | Hook execution results |
| `ControlRequestMsg` | Control protocol request |
| `ControlResponseMsg` | Control protocol response |
| `StreamDoneMsg` | Stream completion signal |
| `StreamErrorMsg` | Stream error with context |
| `UnknownMessageMsg` | Unrecognized message (for debugging) |

### Running the Validation Example

```bash
# Run with a prompt (live mode - requires Claude CLI)
go run example/chat/main.go "Hello, Claude!"

# Run validation mode (tests all message types synthetically)
go run example/chat/main.go --validate
```

## Keybindings (Planned)

Default keybindings use vim-style navigation:

| Action | Keys |
|--------|------|
| Navigation | `j/k` or `up/down` |
| Submit | `enter` |
| Cancel | `esc` |
| Allow tool | `a` or `y` |
| Deny tool | `d` or `n` |
| Always allow | `A` |
| Toggle auto-accept | `ctrl+a` |
| Toggle plan mode | `ctrl+p` |
| Help | `?` |

Custom keybindings:

```go
import "github.com/severity1/claude-agent-tui/system/keymap"

km := keymap.Default()
km.Set("quit", "ctrl+q")
km.Set("allow", "a", "y", "enter")
```

## Accessibility (Planned)

Built-in accessibility support controlled via environment variables:

| Variable | Effect |
|----------|--------|
| `REDUCE_MOTION=1` | Disable all animations |
| `HIGH_CONTRAST=1` | Enable high contrast theme |
| `SCREEN_READER=1` | Enable screen reader announcements |
| `LARGE_TEXT=1` | Increase text sizes |

All components support:
- Keyboard-only navigation
- Visible focus indicators
- Color-blindness safe design (symbols alongside colors)
- Screen reader announcements for state changes

## Error Recovery (Planned)

Automatic error recovery with exponential backoff:

```go
import "github.com/severity1/claude-agent-tui/system/recovery"

mgr := recovery.New(
    recovery.WithMaxRetries(5),
    recovery.WithBaseDelay(time.Second),
    recovery.WithOnStateChange(func(state recovery.ConnectionState) {
        // Handle state changes
    }),
)
```

Connection states: `Healthy`, `Unstable`, `Lost`, `Reconnecting`, `Failed`

Features:
- Automatic reconnection with jittered backoff
- Partial message recovery on interruption
- User notification of connection state
- Manual retry option when auto-recovery fails

## Session Management (Planned)

Persist and resume conversations:

```go
import "github.com/severity1/claude-agent-tui/system/session"

store, _ := session.NewFileStore("~/.claude-tui/sessions")
mgr := session.NewManager(store)

// Create new session
sess, _ := mgr.NewSession("claude-opus-4-5-20251101")

// Resume existing session
sess, _ = mgr.Resume(sessionID)

// Export to markdown
data, _ := mgr.Export(sessionID, session.ExportMarkdown)
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built components
- [Harmonica](https://github.com/charmbracelet/harmonica) - Spring animations
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [claude-agent-sdk-go](https://github.com/severity1/claude-agent-sdk-go) - Claude Code SDK

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design and data flow
- [Components](docs/COMPONENTS.md) - Component specifications and interfaces
- [SDK Mapping](docs/SDK-MAPPING.md) - SDK events to component mapping
- [Adapters](docs/ADAPTERS.md) - SDK-to-TUI message conversion
- [Theming](docs/THEMING.md) - Theme system and customization
- [Animation](docs/ANIMATION.md) - Animation system and triggers
- [Keybindings](docs/KEYBINDINGS.md) - Keymap customization
- [Mouse](docs/MOUSE.md) - Mouse support and zones
- [Responsive](docs/RESPONSIVE.md) - Responsive layout system
- [Accessibility](docs/ACCESSIBILITY.md) - Screen reader and reduced motion support
- [Error Handling](docs/ERROR-HANDLING.md) - Recovery and resilience patterns
- [Sessions](docs/SESSIONS.md) - Session persistence and management
- [Roadmap](docs/ROADMAP.md) - Implementation phases

## Status

| Component | Status | Notes |
|-----------|--------|-------|
| **Stream Adapter** | Implemented | Converts SDK messages to Bubble Tea messages |
| **Validation Example** | Implemented | Demonstrates all 16+ message types |
| **Control Adapter** | Planned | Tool permission callbacks |
| **Client Adapter** | Planned | SDK client lifecycle |
| **TUI Components** | Planned | See component tables above |
| **System Packages** | Planned | Theme, animation, keymap, etc. |

See [ROADMAP.md](docs/ROADMAP.md) for implementation phases.

## License

MIT License - see LICENSE file for details.
