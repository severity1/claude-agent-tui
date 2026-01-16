---
name: claude-agent-sdk-go
version: "1.0"
description: This skill should be used when integrating with the Claude Agent SDK for Go (claude-agent-sdk-go). Provides API patterns for Query API (stateless one-shot queries), Client API (stateful multi-turn conversations), functional options (With* pattern), and message types (AssistantMessage, ResultMessage, ContentBlocks). Use when the user asks about "Claude SDK Go", "query claude programmatically", "stream messages from Claude", "TUI Claude integration", "Claude client pattern", "WithModel option", "ReceiveMessages channel", "tool permission callbacks", "canUseTool callback", "claude-agent-sdk-go API", "StreamDeltaMsg", "convert SDK to tea.Msg", "permission adapter pattern", "background tasks", "TaskOutput", or "KillShell".
---

# Claude Agent SDK for Go

Reference documentation for the claude-agent-sdk-go library, enabling Go applications to interact with the Claude Code CLI.

## Quick Reference

| Need | Use | Reference |
|------|-----|-----------|
| One-shot query, no state | `claudecode.Query()` | [query-api.md](references/query-api.md) |
| Multi-turn conversation | `claudecode.NewClient()` | [client-api.md](references/client-api.md) |
| Configure behavior | `claudecode.With*()` options | [options.md](references/options.md) |
| Parse response content | `AssistantMessage`, `ContentBlock` | [messages.md](references/messages.md) |
| Stream to TUI | `ReceiveMessages()` channel | See patterns below |
| Handle tool permissions | `canUseTool` callback | See patterns below |

## When to Use Each API

| Scenario | API | Why |
|----------|-----|-----|
| CLI tool, single prompt | Query | Stateless, simple, returns iterator |
| Chat interface | Client | Maintains session, conversation history |
| Automation script | Query | No session management overhead |
| Interactive TUI | Client | Need streaming, tool callbacks |

## Key Concepts

### Message Flow
```
User Input -> Client.SendMessage() -> SDK processes
                                          |
                                          v
TUI Update() <- tea.Msg <- Adapter <- ReceiveMessages() channel
```

### Core Types

| Type | Purpose |
|------|---------|
| `Message` | Base interface for all messages |
| `AssistantMessage` | Claude's response with content blocks |
| `ResultMessage` | Final result with token usage, cost |
| `ContentBlock` | Text, tool use, or tool result |
| `StreamEvent` | Real-time streaming delta |

### Integration Points for TUI

1. **Stream Events** - Convert to `StreamDeltaMsg` for live text updates
2. **Assistant Messages** - Convert to `AssistantMsg` for message bubbles
3. **Result Messages** - Extract metadata for status bar (tokens, cost)
4. **Tool Callbacks** - `canUseTool` triggers permission prompts
5. **Error Events** - Convert to appropriate error messages

## Common Patterns for TUI Integration

### Converting SDK Channel to tea.Cmd
```go
func StreamCmd(ctx context.Context, msgCh <-chan claudecode.Message) tea.Cmd {
    return func() tea.Msg {
        select {
        case <-ctx.Done():
            return StreamEndMsg{}
        case msg, ok := <-msgCh:
            if !ok {
                return StreamEndMsg{}
            }
            return convertToTeaMsg(msg)
        }
    }
}
```

### Permission Callback Integration
```go
func NewControlAdapter() *ControlAdapter {
    return &ControlAdapter{
        responseCh: make(chan PermissionResponse),
    }
}

func (a *ControlAdapter) CanUseTool(input claudecode.CanUseToolInput) claudecode.PermissionResult {
    // Send prompt to TUI
    a.promptCh <- ShowPermissionPromptMsg{Input: input}

    // Block until TUI responds
    response := <-a.responseCh

    if response.Allow {
        return claudecode.PermissionResultAllow(nil)
    }
    return claudecode.PermissionResultDeny("User denied")
}
```

### Message Type Conversion
```go
func convertToTeaMsg(msg claudecode.Message) tea.Msg {
    switch m := msg.(type) {
    case *claudecode.AssistantMessage:
        return AssistantMsg{Content: m.Content}
    case *claudecode.ResultMessage:
        return ResultMsg{
            Tokens: m.Usage.InputTokens + m.Usage.OutputTokens,
            Cost:   m.Cost,
        }
    case *claudecode.StreamEvent:
        return StreamDeltaMsg{Delta: m.Delta}
    default:
        return nil
    }
}
```

### Error Handling
```go
if claudecode.IsConnectionError(err) {
    return ConnectionErrorMsg{Err: err}
}
if claudecode.IsCLINotFoundError(err) {
    return CLINotFoundMsg{}
}
```

## Best Practices

1. **Always use context** - Pass `context.Context` for cancellation support
2. **Handle channel closure** - Check `ok` when receiving from `ReceiveMessages()`
3. **Don't block callbacks** - Use channels to communicate between callback and TUI
4. **Batch option application** - Pass all `With*` options in single call
5. **Check error types** - Use `IsConnectionError`, `IsCLINotFoundError` for specific handling
6. **Close clients** - Call `client.Close()` when done to clean up resources
7. **Forward all messages** - Don't filter stream events; components may need them

## Common Pitfalls

- Blocking in `canUseTool` callback without channel communication
- Not handling `StreamEndMsg` / channel closure
- Forgetting to check context cancellation in stream loops
- Creating new client per message instead of reusing
- Not passing required options like `WithModel()` for Query API

## Reference Files

| File | Contents |
|------|----------|
| [query-api.md](references/query-api.md) | `Query()` function, `MessageIterator`, one-shot usage |
| [client-api.md](references/client-api.md) | `Client` interface, session management, `SendMessage()` |
| [options.md](references/options.md) | 60+ functional options (`With*` pattern) |
| [messages.md](references/messages.md) | Message types, content blocks, stream events |
| [error-handling.md](references/error-handling.md) | Error types, detection patterns, recovery strategies |
| [background-tasks.md](references/background-tasks.md) | Background tasks, TaskOutput, KillShell integration |
