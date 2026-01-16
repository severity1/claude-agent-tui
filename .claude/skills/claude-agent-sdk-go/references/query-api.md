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
