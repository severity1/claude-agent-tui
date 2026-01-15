# Architecture

## Project Structure

```
claude-agent-tui/
├── component/
│   ├── output/                 # Display components
│   │   ├── streamtext/         # Live streaming text with cursor
│   │   │   ├── streamtext.go   # Component implementation
│   │   │   └── styles.go       # Lip Gloss style definitions
│   │   ├── message/            # Message bubble (user/assistant)
│   │   ├── thinking/           # Collapsible thinking block
│   │   ├── tooluse/            # Tool invocation card
│   │   ├── codeblock/          # Syntax highlighted code
│   │   └── status/             # Token counter, model info
│   │
│   └── input/                  # Input components
│       ├── chat/               # Main chat input + mode toggles
│       ├── permission/         # Tool permission prompt
│       ├── question/           # AskUserQuestion multi-choice
│       ├── planenter/          # EnterPlanMode confirmation
│       ├── planexit/           # ExitPlanMode plan review
│       └── interrupt/          # Stop confirmation modal
│
├── layout/                     # High-level composite screens
│   ├── chat/                   # Full chat interface
│   │   ├── chat.go             # Layout composition
│   │   ├── state.go            # Input state machine
│   │   └── history.go          # Message history management
│   └── agent/                  # Agent dashboard with tools view
│
├── adapter/                    # SDK integration layer
│   ├── stream.go               # StreamEvent -> tea.Msg
│   ├── control.go              # Hook events -> input prompts
│   ├── message.go              # Message type adapters
│   └── client.go               # Client lifecycle management
│
├── style/                      # Theming
│   └── style.go                # Default Lip Gloss styles
│
└── tui.go                      # Public API exports
```

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        Bubble Tea Program                        │
│                                                                  │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│  │  Model   │───►│  Update  │───►│  View    │───►│ Terminal │  │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘  │
│       ▲               │                                         │
│       │               │                                         │
│       │          tea.Msg                                        │
│       │               │                                         │
│  ┌────┴───────────────┴────────────────────────────────────┐   │
│  │                    Adapter Layer                         │   │
│  │                                                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │ StreamCmd() │  │ ControlCmd()│  │ ClientCmd() │      │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘      │   │
│  └─────────┼────────────────┼────────────────┼─────────────┘   │
│            │                │                │                  │
└────────────┼────────────────┼────────────────┼──────────────────┘
             │                │                │
             ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    claude-agent-sdk-go                           │
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ StreamEvent │  │ Hook Input  │  │   Client    │              │
│  │   Channel   │  │  Callbacks  │  │  Lifecycle  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

## Adapter Pattern

The adapter layer bridges the SDK's channel-based streaming to Bubble Tea's message-passing architecture.

### Stream Adapter

Converts SDK `StreamEvent` messages to Bubble Tea messages:

```go
// StreamCmd returns a tea.Cmd that reads from SDK channel
func StreamCmd(ctx context.Context, ch <-chan claudecode.Message) tea.Cmd {
    return func() tea.Msg {
        select {
        case msg, ok := <-ch:
            if !ok {
                return StreamDoneMsg{}
            }
            return adaptMessage(msg)
        case <-ctx.Done():
            return StreamErrorMsg{Err: ctx.Err()}
        }
    }
}

// Message types emitted by adapter
type StreamDeltaMsg struct {
    Text string
}

type StreamDoneMsg struct{}

type StreamErrorMsg struct {
    Err error
}
```

### Control Adapter

Converts SDK hook inputs to input prompt messages:

```go
// Control event messages
type ShowPermissionPromptMsg struct {
    ToolName  string
    ToolInput map[string]any
    Respond   func(PermissionResult)
}

type ShowQuestionPromptMsg struct {
    Question    string
    Options     []QuestionOption
    MultiSelect bool
    Respond     func([]string)
}

type ShowPlanEnterPromptMsg struct {
    Message string
    Respond func(bool)
}

type ShowPlanExitPromptMsg struct {
    Plan        string
    Permissions []RequestedPermission
    Respond     func(PlanExitResult)
}
```

## Input State Machine

The chat layout manages which input component has focus:

```
                         ┌─────────────────┐
                         │   ChatInput     │◄── Default state
                         └────────┬────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│PermissionPrompt │    │ QuestionPrompt  │    │ PlanEnterConfirm│
└────────┬────────┘    └────────┬────────┘    └────────┬────────┘
         │                      │                      │
         │                      │                      │ (confirmed)
         │                      │                      ▼
         │                      │             ┌─────────────────┐
         │                      │             │  Planning...    │
         │                      │             │  (streaming)    │
         │                      │             └────────┬────────┘
         │                      │                      │
         │                      │                      ▼
         │                      │             ┌─────────────────┐
         │                      │             │ PlanExitPrompt  │
         │                      │             └────────┬────────┘
         │                      │                      │
         └──────────────────────┴──────────────────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │   ChatInput     │
                         └─────────────────┘
```

### State Enum

```go
type InputState int

const (
    InputStateChat InputState = iota
    InputStatePermission
    InputStateQuestion
    InputStatePlanEnter
    InputStatePlanning
    InputStatePlanExit
    InputStateInterrupt
)
```

## Component Interface

All components follow the Bubble Tea Model interface:

```go
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

### With Focus Management

Input components include focus tracking:

```go
type FocusableComponent interface {
    Component
    Focus() tea.Cmd
    Blur()
    Focused() bool
}
```

## Styling System

Components accept Lip Gloss styles via functional options:

```go
// Each component defines its own Styles struct
type StreamTextStyles struct {
    Text   lipgloss.Style
    Cursor lipgloss.Style
}

// Functional option pattern
func WithStyles(s StreamTextStyles) Option {
    return func(m *Model) {
        m.styles = s
    }
}

// Usage
component := streamtext.New(
    streamtext.WithStyles(StreamTextStyles{
        Text:   lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
        Cursor: lipgloss.NewStyle().Background(lipgloss.Color("205")),
    }),
)
```

## SDK Integration Points

### Message Types Handled

| SDK Type | TUI Response |
|----------|--------------|
| `StreamEvent` (content_block_delta) | Update `StreamText` component |
| `AssistantMessage` | Create/update `MessageBubble` |
| `ResultMessage` | Update `StatusBar`, end stream |
| `PreToolUseHookInput` | Show `PermissionPrompt` |
| `PostToolUseHookInput` | Update `ToolUseCard` with result |

### Tool-Based Inputs

Some SDK tool calls trigger input prompts:

| Tool | Input Component |
|------|-----------------|
| `AskUserQuestion` | `QuestionPrompt` |
| `EnterPlanMode` | `PlanEnterConfirm` |
| `ExitPlanMode` | `PlanExitPrompt` |

### Permission Modes

`ChatInput` mode toggles map to SDK permission modes:

| Toggle | SDK Mode |
|--------|----------|
| Auto-accept ON | `PermissionModeAcceptEdits` |
| Plan mode ON | `PermissionModePlan` |
| Both OFF | `PermissionModeDefault` |
