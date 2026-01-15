# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Overview

**Claude Agent TUI** - Terminal UI components for Claude Agent SDK using Charmbracelet ecosystem. Provides reusable Bubble Tea components that integrate with `claude-agent-sdk-go`.

- **Module**: `github.com/severity1/claude-agent-tui`
- **Go Version**: 1.21+
- **Dependencies**: Bubble Tea, Lip Gloss, Bubbles, claude-agent-sdk-go

## Build & Development Commands

```bash
# Build and test
go build ./...                    # Build all packages
go test ./...                     # Run all tests
go test -race ./...               # Race condition detection

# Code quality
go fmt ./...                      # Format code
go vet ./...                      # Static analysis
golangci-lint run                 # Comprehensive linting

# Run examples
go run examples/chat/main.go     # Run chat example
```

## Architecture

```
claude-agent-tui/
├── component/
│   ├── output/           # Display components
│   │   ├── streamtext/   # Live streaming text
│   │   ├── message/      # Message bubble
│   │   ├── thinking/     # Collapsible thinking
│   │   ├── tooluse/      # Tool invocation card
│   │   ├── codeblock/    # Syntax highlighted code
│   │   └── status/       # Token/cost/model bar
│   │
│   └── input/            # Input components
│       ├── chat/         # Main input + mode toggles
│       ├── permission/   # Tool permission prompt
│       ├── question/     # AskUserQuestion prompt
│       ├── planenter/    # EnterPlanMode confirm
│       ├── planexit/     # ExitPlanMode review
│       └── interrupt/    # Stop confirmation
│
├── layout/               # Composite screens
│   ├── chat/             # Full chat interface
│   └── agent/            # Agent dashboard
│
├── adapter/              # SDK integration
│   ├── stream.go         # StreamEvent -> tea.Msg
│   ├── control.go        # Hooks -> prompts
│   └── client.go         # Client lifecycle
│
└── style/                # Default styles
    └── style.go          # Lip Gloss defaults
```

**Data Flow**:
1. SDK `ReceiveMessages()` -> `adapter.StreamCmd()` -> `tea.Msg`
2. Component receives `tea.Msg` -> `Update()` -> state change
3. `View()` renders with Lip Gloss styles

## Code Conventions

- **Bubble Tea Model**: All components implement `Init()`, `Update()`, `View()`
- **Functional options**: `WithStyles()`, `WithXxx()` pattern for configuration
- **Messages**: Each component defines its own `XxxMsg` types
- **Styles**: All components accept `Styles` struct with Lip Gloss styles
- **Focus management**: Input components implement `Focus()`, `Blur()`, `Focused()`

## Patterns

- **Adapter pattern**: Bridge SDK channels to Bubble Tea messages
- **State machine**: Layout manages input focus transitions
- **Response channels**: Hook adapters use channels for sync response

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - Pre-built components
- `github.com/severity1/claude-agent-sdk-go` - Claude SDK

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design
- [Components](docs/COMPONENTS.md) - Component specifications
- [Adapters](docs/ADAPTERS.md) - SDK integration layer
- [Roadmap](docs/ROADMAP.md) - Implementation phases
