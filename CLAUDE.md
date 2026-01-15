# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Overview

**Claude Agent TUI** - Reusable terminal UI component library for building Claude Code interfaces using the Charmbracelet ecosystem. Provides 20+ Bubble Tea components that integrate with `claude-agent-sdk-go`.

- **Module**: `github.com/severity1/claude-agent-tui`
- **Go Version**: 1.21+
- **Core Dependencies**: Bubble Tea, Lip Gloss, Bubbles, Harmonica, Glamour

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
│   ├── output/                   # Display components
│   │   ├── streamtext/           # Live streaming text
│   │   ├── message/              # Message bubble (user/assistant)
│   │   ├── thinking/             # Collapsible thinking block
│   │   ├── tooluse/              # Tool invocation card
│   │   ├── codeblock/            # Syntax highlighted code
│   │   ├── status/               # Token/cost/model bar
│   │   └── toast/                # Notification toasts
│   │
│   ├── input/                    # Input components
│   │   ├── chat/                 # Main input + mode toggles
│   │   ├── permission/           # Tool permission prompt
│   │   ├── question/             # AskUserQuestion prompt
│   │   ├── planenter/            # EnterPlanMode confirm
│   │   ├── planexit/             # ExitPlanMode review
│   │   ├── interrupt/            # Stop confirmation
│   │   ├── confirm/              # Generic yes/no prompt
│   │   ├── filepicker/           # File selection dialog
│   │   └── todo/                 # TodoWrite task list
│   │
│   └── shared/                   # Shared primitives
│       ├── button/               # Reusable button
│       ├── chip/                 # Toggle chip/pill
│       ├── card/                 # Card container
│       └── overlay/              # Overlay/modal wrapper
│
├── system/                       # Cross-cutting systems
│   ├── theme/                    # Theming engine
│   │   ├── theme.go              # Theme interface + registry
│   │   ├── palette.go            # Color palette definitions
│   │   └── builtin/              # Built-in themes (11 total)
│   │
│   ├── keymap/                   # Keybinding system
│   │   ├── keymap.go             # Keymap registry
│   │   ├── defaults.go           # Default keybindings
│   │   └── help.go               # Help overlay generator
│   │
│   ├── animation/                # Animation engine
│   │   ├── animation.go          # Animation primitives
│   │   ├── spring.go             # Spring-based motion (harmonica)
│   │   └── triggers.go           # Event trigger definitions
│   │
│   ├── responsive/               # Responsive layout
│   │   ├── breakpoint.go         # Breakpoint definitions
│   │   └── constraints.go        # Min/max constraints
│   │
│   ├── mouse/                    # Mouse support
│   │   ├── mouse.go              # Mouse event routing
│   │   └── zones.go              # Click zone registry
│   │
│   └── focus/                    # Focus management
│       └── focus.go              # Focus ring + traversal
│
├── adapter/                      # SDK integration
│   ├── stream.go                 # StreamEvent -> tea.Msg
│   ├── control.go                # canUseTool -> prompts
│   └── client.go                 # Client lifecycle
│
├── layout/                       # Composite screens
│   ├── chat/                     # Full chat interface
│   └── agent/                    # Agent dashboard
│
└── style/                        # Default styles
    └── defaults.go               # Sensible defaults
```

## Data Flow

1. SDK `ReceiveMessages()` -> `adapter.StreamCmd()` -> `tea.Msg`
2. Component receives `tea.Msg` -> `Update()` -> state change
3. `canUseTool` callback -> `adapter.ToolControlAdapter` -> prompt `tea.Msg`
4. `View()` renders with theme-aware Lip Gloss styles

## Code Conventions

- **Bubble Tea Model**: All components implement `Init()`, `Update()`, `View()`
- **Functional options**: `WithStyles()`, `WithTheme()`, `WithKeymap()` pattern
- **Messages**: Each component defines its own `XxxMsg` types
- **Styles**: All components accept `Styles` struct generated from theme
- **Focus management**: Input components implement `Focus()`, `Blur()`, `Focused()`
- **Mouse support**: Components implement `Zones()` for clickable regions
- **Variants**: Components support multiple display modes via `Variant` type
- **Animation**: Components implement `AnimatedComponent` interface for animations

## Patterns

- **Adapter pattern**: Bridge SDK channels to Bubble Tea messages
- **State machine**: Layout manages input focus transitions
- **Response channels**: canUseTool adapter uses channels for sync response
- **Theme injection**: Components receive styles from theme via factories
- **Zone registration**: Mouse zones re-registered on each render frame
- **Spring physics**: Animations use harmonica for natural motion

## Component Variants

```go
type Variant int

const (
    VariantView       Variant = iota  // Embeddable content
    VariantInline                      // Flows within text
    VariantCard                        // Bordered card style
    VariantWindow                      // Bordered, positionable
    VariantModal                       // Centered with backdrop
    VariantOverlay                     // Floats above
    VariantFullscreen                  // Takes entire terminal
    VariantBar                         // Horizontal bar
    VariantCompact                     // Minimal footprint
    VariantCollapsible                 // Expandable/collapsible
    VariantSidebar                     // Vertical side panel
    VariantToast                       // Corner notification
)
```

## Key Dependencies

| Package | Import | Usage |
|---------|--------|-------|
| bubbletea | `tea` | Framework, Model interface |
| lipgloss | `lipgloss` | Styling, AdaptiveColor |
| bubbles | various | viewport, textarea, list, spinner |
| harmonica | `harmonica` | Spring-based animations |
| glamour | `glamour` | Markdown rendering |
| chroma | `chroma` | Syntax highlighting |
| claude-agent-sdk-go | `claudecode` | SDK integration |

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design and data flow
- [Components](docs/COMPONENTS.md) - Component specifications with variants
- [SDK Mapping](docs/SDK-MAPPING.md) - SDK events to component mapping
- [Adapters](docs/ADAPTERS.md) - canUseTool callback integration
- [Theming](docs/THEMING.md) - Theme system with 11 built-in themes
- [Animation](docs/ANIMATION.md) - Spring-based animation system
- [Keybindings](docs/KEYBINDINGS.md) - Customizable vim-style keybindings
- [Mouse](docs/MOUSE.md) - Mouse support with click zones
- [Responsive](docs/RESPONSIVE.md) - Responsive layout breakpoints
- [Roadmap](docs/ROADMAP.md) - 13 milestone implementation plan
