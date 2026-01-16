# Message Types Reference

The SDK provides several message types for different stages of the conversation.

## Message Interface

```go
type Message interface {
    // Common interface for all messages
}
```

## UserMessage

Represents user input.

```go
type UserMessage struct {
    ID      string
    Content string
}
```

## AssistantMessage

Claude's response, which may contain multiple content blocks.

```go
type AssistantMessage struct {
    ID            string
    StopReason    string
    ContentBlocks []ContentBlock
    Error         *ErrorInfo  // Non-nil if response contains error
}

// Methods
func (m *AssistantMessage) Content() string     // Combined text content
func (m *AssistantMessage) HasError() bool      // Check for errors
```

## ResultMessage

Final message with metadata about the conversation.

```go
type ResultMessage struct {
    SessionID      string
    CostUSD        string
    DurationMS     int64
    InputTokens    int64
    OutputTokens   int64
    TotalTokens    int64
    IsError        bool
    ErrorMessage   string
}
```

## StreamMessage

Raw streaming protocol message (usually processed internally).

```go
type StreamMessage struct {
    Type    string
    Payload json.RawMessage
}
```

## SystemMessage

System-level information.

```go
type SystemMessage struct {
    Content string
}
```

## Content Blocks

AssistantMessage contains content blocks:

```go
type ContentBlock interface {
    BlockType() string
}

type TextBlock struct {
    Text string
}

type ThinkingBlock struct {
    Thinking string
}

type ToolUseBlock struct {
    ID    string
    Name  string
    Input map[string]interface{}
}

type ToolResultBlock struct {
    ToolUseID string
    Content   string
    IsError   bool
}
```

## TUI Message Mapping

Map SDK messages to TUI messages:

```go
func convertToTeaMsg(msg claudecode.Message) tea.Msg {
    switch m := msg.(type) {
    case *claudecode.AssistantMessage:
        if m.HasError() {
            return ErrorMsg{Err: m.Error}
        }
        return AssistantMsg{
            ID:      m.ID,
            Content: m.Content(),
            Blocks:  convertBlocks(m.ContentBlocks),
        }

    case *claudecode.ResultMessage:
        return ResultMsg{
            SessionID:    m.SessionID,
            Cost:         m.CostUSD,
            Duration:     time.Duration(m.DurationMS) * time.Millisecond,
            InputTokens:  m.InputTokens,
            OutputTokens: m.OutputTokens,
        }

    case *claudecode.UserMessage:
        return UserMsg{
            ID:      m.ID,
            Content: m.Content,
        }

    default:
        return UnknownMsg{Raw: msg}
    }
}
```

## Extracting Content Blocks

```go
func extractBlocks(msg *claudecode.AssistantMessage) []Block {
    var blocks []Block
    for _, cb := range msg.ContentBlocks {
        switch b := cb.(type) {
        case *claudecode.TextBlock:
            blocks = append(blocks, TextUIBlock{Text: b.Text})
        case *claudecode.ThinkingBlock:
            blocks = append(blocks, ThinkingUIBlock{Content: b.Thinking})
        case *claudecode.ToolUseBlock:
            blocks = append(blocks, ToolUseUIBlock{
                Name:  b.Name,
                Input: b.Input,
            })
        }
    }
    return blocks
}
```

## Error Handling

```go
// Check error types
if claudecode.IsConnectionError(err) {
    // Handle connection failure
}
if claudecode.IsCLINotFoundError(err) {
    // Claude CLI not installed
}
if claudecode.IsProcessError(err) {
    // CLI process error
}
if claudecode.IsJSONDecodeError(err) {
    // Parse error
}

// Extract error details
var connErr *claudecode.ConnectionError
if errors.As(err, &connErr) {
    fmt.Println(connErr.Message)
}
```
