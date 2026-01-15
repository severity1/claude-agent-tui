# Adapter Layer

The adapter layer bridges `claude-agent-sdk-go` to Bubble Tea's message-passing architecture.

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Bubble Tea Program                        │
│                                                              │
│    Model.Update(msg tea.Msg) (tea.Model, tea.Cmd)           │
│         ▲                                                    │
│         │ tea.Msg                                            │
│         │                                                    │
├─────────┴────────────────────────────────────────────────────┤
│                     Adapter Layer                            │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │  stream.go   │  │  control.go  │  │  client.go   │       │
│  │              │  │              │  │              │       │
│  │ StreamCmd()  │  │ HookAdapter  │  │ ClientCmd()  │       │
│  │ adaptMsg()   │  │ ToolAdapter  │  │ Lifecycle    │       │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘       │
│         │                 │                 │                │
└─────────┼─────────────────┼─────────────────┼────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────┐
│                   claude-agent-sdk-go                        │
│                                                              │
│  ReceiveMessages()    Hooks/Callbacks      Client API       │
│  StreamEvent          PreToolUse           Connect()        │
│  AssistantMessage     PostToolUse          Query()          │
│  ResultMessage        etc.                 Disconnect()     │
└─────────────────────────────────────────────────────────────┘
```

## stream.go - Message Stream Adapter

Converts SDK streaming messages to Bubble Tea messages.

### Types

```go
package adapter

import (
    tea "github.com/charmbracelet/bubbletea"
    claudecode "github.com/severity1/claude-agent-sdk-go"
)

// Output messages (emitted to Bubble Tea)

// StreamDeltaMsg contains incremental text from streaming
type StreamDeltaMsg struct {
    Text string
}

// StreamBlockStartMsg signals a new content block
type StreamBlockStartMsg struct {
    BlockType string
    Index     int
}

// StreamBlockStopMsg signals content block completion
type StreamBlockStopMsg struct {
    Index int
}

// AssistantMsg contains a complete assistant message
type AssistantMsg struct {
    Message *claudecode.AssistantMessage
}

// ResultMsg contains the final result
type ResultMsg struct {
    Message *claudecode.ResultMessage
}

// StreamDoneMsg signals stream completion
type StreamDoneMsg struct{}

// StreamErrorMsg signals a stream error
type StreamErrorMsg struct {
    Err error
}
```

### Functions

```go
// StreamCmd returns a tea.Cmd that reads from SDK message channel.
// Call this in a loop to continuously receive messages.
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

// adaptMessage converts SDK message to appropriate tea.Msg
func adaptMessage(msg claudecode.Message) tea.Msg {
    switch m := msg.(type) {
    case *claudecode.StreamEvent:
        return adaptStreamEvent(m)
    case *claudecode.AssistantMessage:
        return AssistantMsg{Message: m}
    case *claudecode.ResultMessage:
        return ResultMsg{Message: m}
    default:
        return nil
    }
}

// adaptStreamEvent extracts delta text from StreamEvent
func adaptStreamEvent(event *claudecode.StreamEvent) tea.Msg {
    eventType, _ := event.Event["type"].(string)

    switch eventType {
    case claudecode.StreamEventTypeContentBlockStart:
        blockType := ""
        if cb, ok := event.Event["content_block"].(map[string]any); ok {
            blockType, _ = cb["type"].(string)
        }
        index, _ := event.Event["index"].(float64)
        return StreamBlockStartMsg{
            BlockType: blockType,
            Index:     int(index),
        }

    case claudecode.StreamEventTypeContentBlockDelta:
        if delta, ok := event.Event["delta"].(map[string]any); ok {
            if text, ok := delta["text"].(string); ok {
                return StreamDeltaMsg{Text: text}
            }
        }
        return nil

    case claudecode.StreamEventTypeContentBlockStop:
        index, _ := event.Event["index"].(float64)
        return StreamBlockStopMsg{Index: int(index)}

    case claudecode.StreamEventTypeMessageStop:
        return StreamDoneMsg{}

    default:
        return nil
    }
}
```

### Usage Example

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.StreamDeltaMsg:
        // Append to streaming text component
        m.streamText, cmd = m.streamText.Update(streamtext.DeltaMsg{Text: msg.Text})
        // Continue reading stream
        return m, tea.Batch(cmd, adapter.StreamCmd(m.ctx, m.msgChan))

    case adapter.AssistantMsg:
        // Add complete message to history
        m.messages = append(m.messages, message.FromAssistantMessage(msg.Message))
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.ResultMsg:
        // Update status bar, streaming complete
        m.status = status.FromResultMessage(msg.Message)
        m.streaming = false
        return m, nil

    case adapter.StreamDoneMsg:
        m.streaming = false
        return m, nil

    case adapter.StreamErrorMsg:
        m.err = msg.Err
        return m, nil
    }
    return m, nil
}
```

---

## control.go - Control Event Adapter

Converts SDK hooks and tool calls to input prompt messages.

### Types

```go
package adapter

// ShowPermissionPromptMsg triggers permission input component
type ShowPermissionPromptMsg struct {
    ToolName  string
    ToolInput map[string]any
    ToolUseID string
    Respond   func(decision string, alwaysAllow bool)
}

// ShowQuestionPromptMsg triggers question input component
type ShowQuestionPromptMsg struct {
    Question    string
    Header      string
    Options     []QuestionOption
    MultiSelect bool
    Respond     func(selections []string, customText string)
}

type QuestionOption struct {
    Label       string
    Description string
}

// ShowPlanEnterPromptMsg triggers plan entry confirmation
type ShowPlanEnterPromptMsg struct {
    Message string
    Respond func(confirmed bool)
}

// ShowPlanExitPromptMsg triggers plan review component
type ShowPlanExitPromptMsg struct {
    PlanContent string
    Permissions []RequestedPermission
    Respond     func(action string, feedback string, approvedPerms []string)
}

type RequestedPermission struct {
    Tool   string
    Prompt string
}
```

### Hook Adapter

```go
// HookAdapter wraps SDK hooks to emit TUI messages
type HookAdapter struct {
    program *tea.Program
}

func NewHookAdapter(p *tea.Program) *HookAdapter {
    return &HookAdapter{program: p}
}

// PreToolUseHook returns an SDK hook that emits ShowPermissionPromptMsg
func (h *HookAdapter) PreToolUseHook() claudecode.Option {
    return claudecode.WithPreToolUseHook("", func(
        ctx context.Context,
        input any,
        toolUseID *string,
        hookCtx claudecode.HookContext,
    ) (claudecode.HookJSONOutput, error) {
        preInput, ok := input.(*claudecode.PreToolUseHookInput)
        if !ok {
            return claudecode.HookJSONOutput{}, nil
        }

        // Create response channel
        resultCh := make(chan struct {
            decision    string
            alwaysAllow bool
        }, 1)

        // Send message to TUI
        h.program.Send(ShowPermissionPromptMsg{
            ToolName:  preInput.ToolName,
            ToolInput: preInput.ToolInput,
            ToolUseID: *toolUseID,
            Respond: func(decision string, alwaysAllow bool) {
                resultCh <- struct {
                    decision    string
                    alwaysAllow bool
                }{decision, alwaysAllow}
            },
        })

        // Wait for user response
        select {
        case result := <-resultCh:
            if result.decision == "deny" {
                decision := "block"
                reason := "User denied tool execution"
                return claudecode.HookJSONOutput{
                    Decision: &decision,
                    Reason:   &reason,
                }, nil
            }
            return claudecode.HookJSONOutput{}, nil
        case <-ctx.Done():
            return claudecode.HookJSONOutput{}, ctx.Err()
        }
    })
}
```

### Tool Call Adapter

Detects special tool calls (AskUserQuestion, EnterPlanMode, ExitPlanMode) and converts them to TUI messages.

```go
// ToolCallAdapter monitors AssistantMessages for special tools
type ToolCallAdapter struct {
    program *tea.Program
}

func NewToolCallAdapter(p *tea.Program) *ToolCallAdapter {
    return &ToolCallAdapter{program: p}
}

// ProcessMessage checks for special tool calls
func (t *ToolCallAdapter) ProcessMessage(msg *claudecode.AssistantMessage) {
    for _, block := range msg.Content {
        toolUse, ok := block.(*claudecode.ToolUseBlock)
        if !ok {
            continue
        }

        switch toolUse.Name {
        case "AskUserQuestion":
            t.handleAskUserQuestion(toolUse)
        case "EnterPlanMode":
            t.handleEnterPlanMode(toolUse)
        case "ExitPlanMode":
            t.handleExitPlanMode(toolUse)
        }
    }
}

func (t *ToolCallAdapter) handleAskUserQuestion(toolUse *claudecode.ToolUseBlock) {
    // Parse tool input
    questions, _ := toolUse.Input["questions"].([]any)
    if len(questions) == 0 {
        return
    }

    q := questions[0].(map[string]any)
    question, _ := q["question"].(string)
    header, _ := q["header"].(string)
    multiSelect, _ := q["multiSelect"].(bool)

    var options []QuestionOption
    if opts, ok := q["options"].([]any); ok {
        for _, opt := range opts {
            o := opt.(map[string]any)
            options = append(options, QuestionOption{
                Label:       o["label"].(string),
                Description: o["description"].(string),
            })
        }
    }

    t.program.Send(ShowQuestionPromptMsg{
        Question:    question,
        Header:      header,
        Options:     options,
        MultiSelect: multiSelect,
        Respond: func(selections []string, customText string) {
            // Response handling would connect back to SDK
        },
    })
}

func (t *ToolCallAdapter) handleEnterPlanMode(toolUse *claudecode.ToolUseBlock) {
    t.program.Send(ShowPlanEnterPromptMsg{
        Message: "Claude wants to plan the implementation before proceeding.",
        Respond: func(confirmed bool) {
            // Response handling
        },
    })
}

func (t *ToolCallAdapter) handleExitPlanMode(toolUse *claudecode.ToolUseBlock) {
    // Parse plan content and permissions from tool input
    var permissions []RequestedPermission
    if perms, ok := toolUse.Input["allowedPrompts"].([]any); ok {
        for _, p := range perms {
            perm := p.(map[string]any)
            permissions = append(permissions, RequestedPermission{
                Tool:   perm["tool"].(string),
                Prompt: perm["prompt"].(string),
            })
        }
    }

    t.program.Send(ShowPlanExitPromptMsg{
        PlanContent: "Plan content from file...", // Would read from plan file
        Permissions: permissions,
        Respond: func(action string, feedback string, approvedPerms []string) {
            // Response handling
        },
    })
}
```

---

## client.go - Client Lifecycle Adapter

Manages SDK client lifecycle within Bubble Tea.

### Types

```go
package adapter

// ClientState represents the SDK client connection state
type ClientState int

const (
    ClientStateDisconnected ClientState = iota
    ClientStateConnecting
    ClientStateConnected
    ClientStateStreaming
    ClientStateError
)

// ClientStateMsg reports client state changes
type ClientStateMsg struct {
    State ClientState
    Err   error
}

// ClientAdapter manages SDK client lifecycle
type ClientAdapter struct {
    client  claudecode.Client
    ctx     context.Context
    cancel  context.CancelFunc
    state   ClientState
    options []claudecode.Option
}
```

### Functions

```go
// NewClientAdapter creates a new client adapter with options
func NewClientAdapter(opts ...claudecode.Option) *ClientAdapter {
    return &ClientAdapter{
        options: opts,
        state:   ClientStateDisconnected,
    }
}

// ConnectCmd returns a tea.Cmd that establishes SDK connection
func (c *ClientAdapter) ConnectCmd() tea.Cmd {
    return func() tea.Msg {
        c.ctx, c.cancel = context.WithCancel(context.Background())
        c.client = claudecode.NewClient(c.options...)

        if err := c.client.Connect(c.ctx); err != nil {
            c.state = ClientStateError
            return ClientStateMsg{State: ClientStateError, Err: err}
        }

        c.state = ClientStateConnected
        return ClientStateMsg{State: ClientStateConnected}
    }
}

// QueryCmd returns a tea.Cmd that sends a query
func (c *ClientAdapter) QueryCmd(prompt string) tea.Cmd {
    return func() tea.Msg {
        if err := c.client.Query(c.ctx, prompt); err != nil {
            return ClientStateMsg{State: ClientStateError, Err: err}
        }
        c.state = ClientStateStreaming
        return ClientStateMsg{State: ClientStateStreaming}
    }
}

// MessageChannel returns the SDK message channel for streaming
func (c *ClientAdapter) MessageChannel() <-chan claudecode.Message {
    return c.client.ReceiveMessages(c.ctx)
}

// DisconnectCmd returns a tea.Cmd that closes the connection
func (c *ClientAdapter) DisconnectCmd() tea.Cmd {
    return func() tea.Msg {
        if c.cancel != nil {
            c.cancel()
        }
        if c.client != nil {
            c.client.Disconnect(c.ctx)
        }
        c.state = ClientStateDisconnected
        return ClientStateMsg{State: ClientStateDisconnected}
    }
}

// InterruptCmd returns a tea.Cmd that interrupts current operation
func (c *ClientAdapter) InterruptCmd() tea.Cmd {
    return func() tea.Msg {
        if c.client != nil {
            c.client.Interrupt(c.ctx)
        }
        return nil
    }
}

// SetPermissionModeCmd changes permission mode mid-session
func (c *ClientAdapter) SetPermissionModeCmd(mode claudecode.PermissionMode) tea.Cmd {
    return func() tea.Msg {
        if c.client != nil {
            c.client.SetPermissionMode(c.ctx, mode)
        }
        return nil
    }
}
```

### Usage Example

```go
type Model struct {
    client *adapter.ClientAdapter
    // ...
}

func (m Model) Init() tea.Cmd {
    m.client = adapter.NewClientAdapter(
        claudecode.WithPartialStreaming(),
        claudecode.WithMaxTurns(10),
    )
    return m.client.ConnectCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.ClientStateMsg:
        switch msg.State {
        case adapter.ClientStateConnected:
            m.connected = true
        case adapter.ClientStateStreaming:
            // Start reading messages
            return m, adapter.StreamCmd(m.ctx, m.client.MessageChannel())
        case adapter.ClientStateError:
            m.err = msg.Err
        }

    case chatinput.SubmitMsg:
        // Update permission mode based on toggles
        if msg.AutoAccept {
            m.client.SetPermissionModeCmd(claudecode.PermissionModeAcceptEdits)
        }
        return m, m.client.QueryCmd(msg.Text)
    }
    return m, nil
}
```

---

## Integration Summary

| SDK Component | Adapter | TUI Message | Component |
|---------------|---------|-------------|-----------|
| `StreamEvent` (delta) | `stream.go` | `StreamDeltaMsg` | `StreamText` |
| `AssistantMessage` | `stream.go` | `AssistantMsg` | `MessageBubble` |
| `ResultMessage` | `stream.go` | `ResultMsg` | `StatusBar` |
| `PreToolUseHookInput` | `control.go` | `ShowPermissionPromptMsg` | `PermissionPrompt` |
| `AskUserQuestion` tool | `control.go` | `ShowQuestionPromptMsg` | `QuestionPrompt` |
| `EnterPlanMode` tool | `control.go` | `ShowPlanEnterPromptMsg` | `PlanEnterConfirm` |
| `ExitPlanMode` tool | `control.go` | `ShowPlanExitPromptMsg` | `PlanExitPrompt` |
| Client lifecycle | `client.go` | `ClientStateMsg` | Layout state |
