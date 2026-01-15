# claude-agent-tui

TUI components for Claude Agent SDK using the Charmbracelet ecosystem.

## Overview

`claude-agent-tui` provides reusable terminal UI components built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) that integrate with [claude-agent-sdk-go](https://github.com/severity1/claude-agent-sdk-go). It enables developers to build interactive CLI applications that expose Claude Code capabilities.

## Architecture

The library uses a layered architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                      Your Application                        │
├─────────────────────────────────────────────────────────────┤
│  layout/        High-level composite screens                 │
│  ├── chat/      Full chat interface                         │
│  └── agent/     Agent dashboard with tools view             │
├─────────────────────────────────────────────────────────────┤
│  component/     Low-level Bubbles-style components          │
│  ├── output/    StreamText, MessageBubble, ToolUseCard...   │
│  └── input/     ChatInput, PermissionPrompt, QuestionPrompt │
├─────────────────────────────────────────────────────────────┤
│  adapter/       SDK integration layer                        │
│  └── stream.go  Converts SDK messages to Bubble Tea Msgs    │
├─────────────────────────────────────────────────────────────┤
│  claude-agent-sdk-go                                         │
└─────────────────────────────────────────────────────────────┘
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built components
- [claude-agent-sdk-go](https://github.com/severity1/claude-agent-sdk-go) - Claude Code SDK

## Components

### Output Components

| Component | Purpose |
|-----------|---------|
| `StreamText` | Live-updating streaming text with typing cursor |
| `MessageBubble` | Styled conversation message (user/assistant) |
| `ThinkingBlock` | Collapsible extended thinking display |
| `ToolUseCard` | Tool invocation with input/output display |
| `CodeBlock` | Syntax-highlighted code rendering |
| `StatusBar` | Token count, cost, model info |

### Input Components

| Component | Trigger | Purpose |
|-----------|---------|---------|
| `ChatInput` | Default | Main input with mode toggles (auto-accept, plan mode) |
| `PermissionPrompt` | PreToolUse hook | Tool approval (Allow/Deny/Always) |
| `QuestionPrompt` | AskUserQuestion | Multi-choice question response |
| `PlanEnterConfirm` | EnterPlanMode | Confirm entering plan mode |
| `PlanExitPrompt` | ExitPlanMode | Review plan and decide next action |

## Styling

All components accept Lip Gloss styles for full customization:

```go
styles := chatinput.Styles{
    Container:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
    TextArea:   lipgloss.NewStyle().Padding(1),
    ToggleOn:   lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
    ToggleOff:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
}

input := chatinput.New(chatinput.WithStyles(styles))
```

## Quick Start

```go
package main

import (
    "context"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/severity1/claude-agent-tui/layout/chat"
    claudecode "github.com/severity1/claude-agent-sdk-go"
)

func main() {
    // Create chat layout with SDK client
    chatLayout := chat.New(
        chat.WithClient(claudecode.NewClient()),
    )

    // Run Bubble Tea program
    p := tea.NewProgram(chatLayout, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design and data flow
- [Components](docs/COMPONENTS.md) - Component specifications and interfaces
- [Adapters](docs/ADAPTERS.md) - SDK-to-TUI message conversion
- [Roadmap](docs/ROADMAP.md) - Implementation phases

## Status

This project is in the planning phase. See [ROADMAP.md](docs/ROADMAP.md) for implementation timeline.

## License

MIT License - see LICENSE file for details.
