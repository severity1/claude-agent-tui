# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

<!-- AUTO-MANAGED: project-description -->
## Overview

**Claude Agent TUI** - Reusable terminal UI component library for building Claude Code interfaces using the Charmbracelet ecosystem. Provides Bubble Tea components that integrate with `claude-agent-sdk-go`.

- **Module**: `github.com/severity1/claude-agent-tui`
- **Go Version**: 1.24+
- **Core Dependencies**: Bubble Tea, Lip Gloss, Bubbles, claude-agent-sdk-go

<!-- END AUTO-MANAGED -->

<!-- AUTO-MANAGED: build-commands -->
## Build & Development Commands

```bash
# Build and test
go build ./...                    # Build all packages
go test ./...                     # Run all tests
go test -race ./...               # Race condition detection

# Code quality
go fmt ./...                      # Format code
go vet ./...                      # Static analysis

# Run examples
go run example/chat/main.go      # Run chat example with prompt
go run example/chat/main.go --validate  # Validate all adapter message types
```

<!-- END AUTO-MANAGED -->

<!-- AUTO-MANAGED: architecture -->
## Architecture

```
claude-agent-tui/
├── adapter/                      # SDK integration layer
│   ├── stream.go                 # StreamEvent -> tea.Msg conversion
│   ├── stream_test.go            # Comprehensive adapter tests
│   └── doc.go                    # Package documentation
│
├── component/                    # TUI components (planned)
│   ├── input/                    # Input components
│   └── output/                   # Display components
│
├── layout/                       # Composite screens
│   └── chat/                     # Full chat interface
│       └── doc.go                # Package documentation
│
├── example/                      # Example applications
│   └── chat/                     # Stream adapter validation example
│       └── main.go               # Demonstrates all message types
│
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md           # System design and data flow
│   ├── SDK-MAPPING.md            # SDK events to component mapping
│   ├── COMPONENTS.md             # Component specifications
│   └── ...                       # Additional docs
│
└── tui.go                        # Public API exports
```

### Data Flow

1. SDK `ReceiveMessages()` -> `adapter.StreamCmd()` -> `tea.Msg`
2. Component receives `tea.Msg` -> `Update()` -> state change
3. `canUseTool` callback -> control adapter -> prompt `tea.Msg`
4. `View()` renders with Lip Gloss styles

<!-- END AUTO-MANAGED -->

<!-- AUTO-MANAGED: conventions -->
## Code Conventions

### Bubble Tea Patterns
- **Model Interface**: All components implement `Init()`, `Update()`, `View()`
- **Message Types**: Each package defines its own `XxxMsg` types (e.g., `StreamDeltaMsg`, `AssistantMsg`)
- **Commands**: Use `tea.Cmd` for async operations; return from `Update()` or `Init()`
- **Functional Options**: Use `WithXxx()` pattern for configuration (e.g., `WithStyles()`, `WithTheme()`)

### Adapter Pattern
- Bridge SDK channels to Bubble Tea messages via `StreamCmd()`
- Use `AdaptMessage()` for synchronous SDK message conversion
- Handle context cancellation with `StreamErrorMsg`
- Return `StreamDoneMsg` on channel close
- **Error handling**: Return `UnknownMessageMsg` instead of nil for unrecognized types (fail informatively, not silently)

### Testing
- **External test packages**: Use `package adapter_test` (not `package adapter`)
- **Table-driven tests**: Group related test cases in `[]struct{...}` slices
- **Test naming**: `TestAdaptMessage_StreamEvent_ContentBlockDelta_Text`
- **Coverage**: Test nil inputs, missing fields, edge cases explicitly

### Naming
- **Constants**: Group related constants with const blocks
- **Block types**: `BlockTypeText`, `BlockTypeThinking`, `BlockTypeToolUse`
- **Delta types**: `DeltaTypeText`, `DeltaTypeThinking`, `DeltaTypeInputJSON`

### Imports
- Group: stdlib, third-party, local
- Alias bubbletea as `tea`: `tea "github.com/charmbracelet/bubbletea"`
- Alias SDK as `claudecode`: `claudecode "github.com/severity1/claude-agent-sdk-go"`

<!-- END AUTO-MANAGED -->

<!-- AUTO-MANAGED: patterns -->
## Detected Patterns

### Message Type Adaptation
The adapter converts SDK messages to TUI-specific message types:
- `*claudecode.StreamEvent` -> `StreamDeltaMsg`, `StreamBlockStartMsg`, `StreamBlockStopMsg`, `ThinkingDeltaMsg`, `ToolUseDeltaMsg`, `MessageStartMsg`, `MessageDeltaMsg`
- `*claudecode.AssistantMessage` -> `AssistantMsg`
- `*claudecode.ResultMessage` -> `ResultMsg` (with extracted usage stats)
- `*claudecode.UserMessage` -> `UserMsg`
- `*claudecode.SystemMessage` -> `SystemInitMsg`, `SystemHookResponseMsg`
- `*claudecode.RawControlMessage` -> `ControlRequestMsg`, `ControlResponseMsg`
- Unknown types -> `UnknownMessageMsg` (for debugging)

### Error Handling Pattern
- Never return nil for unrecognized message types
- Return `UnknownMessageMsg` with descriptive `TypeName` field for debugging
- TypeName format examples: `"StreamEvent/nil"`, `"StreamEvent/unknown_event"`, `"SystemMessage/unknown_subtype"`
- Pattern: Fail informatively with context, not silently

### Stream Event Types
Content block events: `content_block_start`, `content_block_delta`, `content_block_stop`
Message events: `message_start`, `message_delta`, `message_stop`

### State Tracking Pattern
- `StreamBlockStartMsg` provides `ToolName` and `ToolID` for tool_use blocks
- `ToolUseDeltaMsg` contains only `PartialJSON` and `Index`
- TUI layer tracks state to associate tool info with deltas by index
- This separates concerns: adapter emits granular events, TUI manages state

### Type Extraction Helpers
- `toInt(v any)`: Handles int, float64, int64 (JSON decodes numbers as float64)
- `toString(v any)`: Safe string extraction
- `toStringSlice(v any)`: Convert `[]any` to `[]string`

<!-- END AUTO-MANAGED -->

<!-- AUTO-MANAGED: dependencies -->
## Key Dependencies

| Package | Import Alias | Usage |
|---------|--------------|-------|
| bubbletea | `tea` | Framework, Model interface, Cmd, Msg |
| lipgloss | `lipgloss` | Styling, AdaptiveColor |
| bubbles | various | viewport, textarea, list, spinner |
| claude-agent-sdk-go | `claudecode` | SDK integration, Message types |

<!-- END AUTO-MANAGED -->

<!-- MANUAL -->
## Custom Notes

Add project-specific notes here. This section is never auto-modified.

<!-- END MANUAL -->
