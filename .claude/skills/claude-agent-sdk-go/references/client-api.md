# Client API Reference

The Client API provides stateful, multi-turn conversations with session management.

## Client Interface

```go
type Client interface {
    Connect(ctx context.Context) error
    Query(ctx context.Context, prompt string) error
    QueryWithSession(ctx context.Context, prompt string, sessionID string) error
    ReceiveMessages(ctx context.Context) (<-chan Message, <-chan error)
    Interrupt(ctx context.Context) error
    SetModel(ctx context.Context, model *string) error
    SetPermissionMode(ctx context.Context, mode string) error
    RewindFiles(ctx context.Context, userMessageID string) error
    Disconnect() error
}
```

## WithClient Pattern (Recommended)

Go-idiomatic resource management with automatic cleanup:

```go
err := claudecode.WithClient(ctx, func(client claudecode.Client) error {
    // Client is connected and ready

    if err := client.Query(ctx, "Hello"); err != nil {
        return err
    }

    msgCh, errCh := client.ReceiveMessages(ctx)
    for {
        select {
        case msg, ok := <-msgCh:
            if !ok {
                return nil // Stream complete
            }
            handleMessage(msg)
        case err := <-errCh:
            return err
        }
    }
}, claudecode.WithModel("sonnet"))
// Client automatically disconnected here
```

## Manual Client Management

```go
client, err := claudecode.NewClient(claudecode.WithModel("sonnet"))
if err != nil {
    return err
}
defer client.Disconnect()

if err := client.Connect(ctx); err != nil {
    return err
}

// Use client...
```

## Receiving Messages

```go
msgCh, errCh := client.ReceiveMessages(ctx)

for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case msg, ok := <-msgCh:
        if !ok {
            return nil // Stream ended
        }
        // Process message
    case err := <-errCh:
        return err
    }
}
```

## TUI Integration Pattern

```go
type ClientAdapter struct {
    client claudecode.Client
    ctx    context.Context
    cancel context.CancelFunc
}

func (a *ClientAdapter) Connect() tea.Cmd {
    return func() tea.Msg {
        if err := a.client.Connect(a.ctx); err != nil {
            return ConnectErrorMsg{Err: err}
        }
        return ConnectedMsg{}
    }
}

func (a *ClientAdapter) StreamMessages() tea.Cmd {
    return func() tea.Msg {
        msgCh, errCh := a.client.ReceiveMessages(a.ctx)
        select {
        case msg, ok := <-msgCh:
            if !ok {
                return StreamEndMsg{}
            }
            return convertMessage(msg)
        case err := <-errCh:
            return StreamErrorMsg{Err: err}
        }
    }
}
```

## Session Management

```go
// Resume existing session
err := client.QueryWithSession(ctx, "Continue our discussion", sessionID)

// Rewind file changes to a specific point
err := client.RewindFiles(ctx, userMessageID)
```

## Use Cases

- Interactive chat interfaces
- Multi-turn conversations
- Session persistence
- File change tracking
