---
name: claude-agent-sdk-go
description: Claude Agent SDK for Go - API patterns and integration guidance for building TUI components that interface with the Claude Code CLI
---

# Claude Agent SDK for Go

This skill provides reference documentation for the claude-agent-sdk-go library, which enables Go applications to interact with the Claude Code CLI.

## Overview

The SDK provides two main APIs:
- **Query API**: Stateless, one-shot queries for automation
- **Client API**: Stateful, multi-turn conversations with session management

## Key Concepts

### Message Flow
```
SDK ReceiveMessages() channel -> Adapter converts to tea.Msg -> TUI component Update()
```

### Integration Points for TUI
1. **Stream Events**: Convert to `StreamDeltaMsg` for live text updates
2. **Assistant Messages**: Convert to `AssistantMsg` for message bubbles
3. **Result Messages**: Extract metadata for status bar
4. **Hook Callbacks**: Trigger permission/question prompts

## Reference Files

- `references/query-api.md` - Query() function and MessageIterator
- `references/client-api.md` - Client interface and WithClient() pattern
- `references/options.md` - Functional options (60+ With* functions)
- `references/messages.md` - Message types and content blocks

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

### Error Handling
```go
if claudecode.IsConnectionError(err) {
    return ConnectionErrorMsg{Err: err}
}
if claudecode.IsCLINotFoundError(err) {
    return CLINotFoundMsg{}
}
```
