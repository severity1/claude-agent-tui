# Query API Reference

The Query API provides stateless, one-shot queries to Claude. Each call creates a new session.

## Function Signature

```go
func Query(ctx context.Context, prompt string, opts ...Option) (MessageIterator, error)
```

## MessageIterator

```go
type MessageIterator interface {
    Next() bool           // Advances to next message
    Message() Message     // Returns current message
    Err() error          // Returns any error
    Close() error        // Cleans up resources
}
```

## Basic Usage

```go
iter, err := claudecode.Query(ctx, "Explain Go interfaces")
if err != nil {
    return err
}
defer iter.Close()

for iter.Next() {
    msg := iter.Message()
    switch m := msg.(type) {
    case *claudecode.AssistantMessage:
        fmt.Println(m.Content())
    case *claudecode.ResultMessage:
        fmt.Printf("Cost: %s\n", m.CostUSD)
    }
}
if err := iter.Err(); err != nil {
    return err
}
```

## With Options

```go
iter, err := claudecode.Query(ctx, prompt,
    claudecode.WithModel("sonnet"),
    claudecode.WithMaxTurns(5),
    claudecode.WithSystemPrompt("You are a Go expert"),
    claudecode.WithAllowedTools([]string{"Read", "Write"}),
)
```

## TUI Integration Pattern

```go
// Convert Query to tea.Cmd for Bubble Tea
func QueryCmd(ctx context.Context, prompt string) tea.Cmd {
    return func() tea.Msg {
        iter, err := claudecode.Query(ctx, prompt)
        if err != nil {
            return QueryErrorMsg{Err: err}
        }

        var messages []Message
        for iter.Next() {
            messages = append(messages, iter.Message())
        }
        iter.Close()

        if err := iter.Err(); err != nil {
            return QueryErrorMsg{Err: err}
        }
        return QueryCompleteMsg{Messages: messages}
    }
}
```

## Use Cases

- Automation scripts
- CI/CD pipelines
- One-off code generation
- Simple question answering

For multi-turn conversations with context preservation, use the Client API instead.

## Edge Cases

### Empty Prompt

```go
// Empty prompt returns error
iter, err := claudecode.Query(ctx, "")
if err != nil {
    // err: prompt cannot be empty
}
```

### Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

iter, err := claudecode.Query(ctx, prompt)
if err != nil {
    return err
}
defer iter.Close()

for iter.Next() {
    // Process messages
}

// Check if canceled
if err := iter.Err(); err != nil {
    if errors.Is(err, context.Canceled) {
        // User canceled
    }
    if errors.Is(err, context.DeadlineExceeded) {
        // Timeout
    }
}
```

### Large Responses

```go
var totalContent strings.Builder

for iter.Next() {
    msg := iter.Message()
    if m, ok := msg.(*claudecode.AssistantMessage); ok {
        // Stream content incrementally
        totalContent.WriteString(m.Content())

        // Check if too large
        if totalContent.Len() > maxSize {
            cancel()  // Cancel context
            break
        }
    }
}
```

### Error in Response

```go
for iter.Next() {
    msg := iter.Message()
    if m, ok := msg.(*claudecode.AssistantMessage); ok {
        if m.HasError() {
            // Handle error from Claude
            fmt.Printf("Claude error: %s\n", m.Error.Message)
            continue
        }
    }
}
```

## Streaming with TUI

For live updates in TUI, process messages incrementally:

```go
func StreamQueryCmd(ctx context.Context, prompt string) tea.Cmd {
    return func() tea.Msg {
        iter, err := claudecode.Query(ctx, prompt)
        if err != nil {
            return QueryErrorMsg{Err: err}
        }

        // Return first message, schedule next read
        if iter.Next() {
            return QueryMessageMsg{
                Message:  iter.Message(),
                Iterator: iter,  // Pass iterator for next read
                HasMore:  true,
            }
        }

        iter.Close()
        if err := iter.Err(); err != nil {
            return QueryErrorMsg{Err: err}
        }
        return QueryDoneMsg{}
    }
}

// Continue reading in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case QueryMessageMsg:
        m.processMessage(msg.Message)

        if msg.HasMore {
            // Schedule next read
            return m, continueQuery(msg.Iterator)
        }
        msg.Iterator.Close()
        return m, nil
    }
    return m, nil
}
```

## Resource Cleanup

Always close the iterator to prevent resource leaks:

```go
iter, err := claudecode.Query(ctx, prompt)
if err != nil {
    return err
}

// Use defer for guaranteed cleanup
defer func() {
    if closeErr := iter.Close(); closeErr != nil {
        log.Printf("failed to close iterator: %v", closeErr)
    }
}()

for iter.Next() {
    // Process
}
```

## Comparison: Query vs Client

| Aspect | Query API | Client API |
|--------|-----------|------------|
| State | Stateless | Stateful |
| Session | New each call | Persistent |
| Context | None | Full history |
| Use case | One-off queries | Conversations |
| Overhead | Lower | Higher |
| Streaming | Via iterator | Via channel |

## Live Documentation

For latest API details:
- SDK source: github.com/severity1/claude-agent-sdk-go
- Query godoc: pkg.go.dev/github.com/severity1/claude-agent-sdk-go#Query
