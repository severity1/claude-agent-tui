# Error Handling Reference

Error types, detection patterns, and recovery strategies for SDK integration.

## SDK Error Types

The SDK provides typed errors for different failure scenarios.

### Error Detection

```go
import "github.com/severity1/claudecode"

// Check error types using provided helpers
if claudecode.IsConnectionError(err) {
    // CLI connection failed
}
if claudecode.IsCLINotFoundError(err) {
    // Claude CLI not installed
}
if claudecode.IsProcessError(err) {
    // CLI process error (crash, exit code)
}
if claudecode.IsJSONDecodeError(err) {
    // Response parsing failed
}
if claudecode.IsTimeoutError(err) {
    // Operation timed out
}
```

### Extract Error Details

```go
// Connection errors
var connErr *claudecode.ConnectionError
if errors.As(err, &connErr) {
    fmt.Printf("Connection failed: %s\n", connErr.Message)
    fmt.Printf("Cause: %v\n", connErr.Cause)
}

// Process errors
var procErr *claudecode.ProcessError
if errors.As(err, &procErr) {
    fmt.Printf("Process exited with code %d\n", procErr.ExitCode)
    fmt.Printf("Stderr: %s\n", procErr.Stderr)
}

// Parse errors
var parseErr *claudecode.JSONDecodeError
if errors.As(err, &parseErr) {
    fmt.Printf("Failed to parse at offset %d\n", parseErr.Offset)
    fmt.Printf("Raw: %s\n", parseErr.Raw)
}
```

## Error Categories for TUI

Map SDK errors to TUI behavior.

### Fatal Errors (Don't Retry)

```go
func isFatalError(err error) bool {
    // CLI not installed - user action required
    if claudecode.IsCLINotFoundError(err) {
        return true
    }

    // Version mismatch - user action required
    if claudecode.IsVersionError(err) {
        return true
    }

    // Authentication issues
    if claudecode.IsAuthError(err) {
        return true
    }

    return false
}
```

### Transient Errors (Safe to Retry)

```go
func isTransientError(err error) bool {
    // Network issues usually recover
    if claudecode.IsConnectionError(err) {
        var connErr *claudecode.ConnectionError
        if errors.As(err, &connErr) {
            // Timeout is transient
            return connErr.IsTimeout
        }
    }

    // Process crash might be transient
    if claudecode.IsProcessError(err) {
        return true
    }

    // Timeout is transient
    if claudecode.IsTimeoutError(err) {
        return true
    }

    return false
}
```

## TUI Message Conversion

Convert SDK errors to tea.Msg types.

```go
// Define TUI error messages
type ConnectionErrorMsg struct {
    Err       error
    Retryable bool
    Message   string
}

type CLINotFoundMsg struct{}

type ParseErrorMsg struct {
    Err    error
    Offset int
}

// Convert SDK errors
func convertError(err error) tea.Msg {
    if claudecode.IsCLINotFoundError(err) {
        return CLINotFoundMsg{}
    }

    if claudecode.IsJSONDecodeError(err) {
        var parseErr *claudecode.JSONDecodeError
        errors.As(err, &parseErr)
        return ParseErrorMsg{
            Err:    parseErr,
            Offset: parseErr.Offset,
        }
    }

    return ConnectionErrorMsg{
        Err:       err,
        Retryable: isTransientError(err),
        Message:   userFriendlyMessage(err),
    }
}
```

## Stream Error Handling

Handle errors during streaming.

```go
func StreamCmd(ctx context.Context, ch <-chan claudecode.Message) tea.Cmd {
    return func() tea.Msg {
        select {
        case <-ctx.Done():
            err := ctx.Err()
            if errors.Is(err, context.Canceled) {
                // User canceled - not an error
                return StreamCanceledMsg{}
            }
            // Timeout
            return StreamErrorMsg{Err: err, Recoverable: true}

        case msg, ok := <-ch:
            if !ok {
                return StreamDoneMsg{}
            }

            // Check for error in message
            if errMsg, isErr := msg.(*claudecode.ErrorMessage); isErr {
                return StreamErrorMsg{
                    Err:         errors.New(errMsg.Message),
                    Recoverable: errMsg.Recoverable,
                }
            }

            return convertToTeaMsg(msg)
        }
    }
}
```

## Client Recovery Pattern

Reconnect after connection loss.

```go
type ClientAdapter struct {
    client    claudecode.Client
    ctx       context.Context
    cancel    context.CancelFunc
    connected bool
    lastErr   error
}

func (a *ClientAdapter) ReconnectCmd() tea.Cmd {
    return func() tea.Msg {
        // Cancel existing context
        if a.cancel != nil {
            a.cancel()
        }

        // Create new context
        a.ctx, a.cancel = context.WithCancel(context.Background())

        // Attempt to create new client
        client, err := claudecode.NewClient(a.ctx, a.opts...)
        if err != nil {
            a.lastErr = err
            return ReconnectFailedMsg{Err: err}
        }

        a.client = client
        a.connected = true
        a.lastErr = nil

        return ReconnectedMsg{}
    }
}
```

## Permission Callback Error Handling

Handle errors in canUseTool callback.

```go
func (a *ControlAdapter) CanUseTool(input claudecode.CanUseToolInput) claudecode.PermissionResult {
    // Send prompt to TUI
    a.promptCh <- ShowPermissionMsg{Input: input}

    // Wait for response with timeout
    select {
    case response := <-a.responseCh:
        if response.Allow {
            return claudecode.PermissionResultAllow(nil)
        }
        return claudecode.PermissionResultDeny(response.Reason)

    case <-time.After(5 * time.Minute):
        // Timeout - deny by default
        return claudecode.PermissionResultDeny("Permission prompt timed out")

    case <-a.ctx.Done():
        // Context canceled
        return claudecode.PermissionResultDeny("Operation canceled")
    }
}
```

## Best Practices

1. **Always check error type** - Use Is/As helpers for proper handling
2. **Distinguish fatal vs transient** - Only retry transient errors
3. **Preserve context** - Include original error in wrapped errors
4. **Timeout all operations** - Prevent indefinite hangs
5. **Log errors appropriately** - Debug level for transient, error level for fatal
6. **User-friendly messages** - Convert technical errors to clear messages

## User-Friendly Messages

```go
func userFriendlyMessage(err error) string {
    if claudecode.IsCLINotFoundError(err) {
        return "Claude CLI is not installed. Run 'claude install' to set it up."
    }
    if claudecode.IsVersionError(err) {
        return "Claude CLI needs to be updated. Run 'claude update' to update."
    }
    if claudecode.IsAuthError(err) {
        return "Authentication failed. Run 'claude login' to authenticate."
    }
    if claudecode.IsConnectionError(err) {
        return "Connection lost. Attempting to reconnect..."
    }
    if claudecode.IsTimeoutError(err) {
        return "Request timed out. Please try again."
    }
    return "An unexpected error occurred. Please try again."
}
```

## Related Documentation

- [messages.md](messages.md) - Message types and conversion
- [client-api.md](client-api.md) - Client lifecycle management
- docs/ERROR-HANDLING.md - TUI-side error handling patterns
