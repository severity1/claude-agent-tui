# Error Handling and Recovery

Resilient error handling and automatic recovery patterns.

## Overview

The error handling system provides:
- Automatic reconnection with exponential backoff
- Partial message recovery
- User notification of connection state
- Graceful degradation
- Error categorization and handling

---

## Error Categories

### Connection Errors

Errors related to SDK client connection.

```go
type ConnectionError struct {
    Type    ConnectionErrorType
    Err     error
    Attempt int
}

type ConnectionErrorType int

const (
    ErrorCLINotFound ConnectionErrorType = iota  // Claude CLI not installed
    ErrorCLIVersion                               // CLI version mismatch
    ErrorProcessStart                             // Failed to start process
    ErrorProcessCrash                             // Process terminated unexpectedly
    ErrorTimeout                                  // Connection timeout
    ErrorNetwork                                  // Network-related error
)
```

### Stream Errors

Errors during message streaming.

```go
type StreamError struct {
    Type     StreamErrorType
    Err      error
    BlockIdx int  // Which content block failed
    Partial  bool // Was partial content received?
}

type StreamErrorType int

const (
    ErrorStreamInterrupted StreamErrorType = iota  // Stream stopped unexpectedly
    ErrorStreamDecode                               // JSON decode error
    ErrorStreamTimeout                              // Read timeout
    ErrorStreamCanceled                             // User canceled
)
```

### Tool Errors

Errors from tool execution.

```go
type ToolError struct {
    ToolName string
    ToolID   string
    Err      error
    Retryable bool
}
```

---

## Recovery Manager

Automatic recovery with exponential backoff.

```go
package recovery

type Manager struct {
    maxRetries    int
    baseDelay     time.Duration
    maxDelay      time.Duration
    backoffFactor float64

    currentState  ConnectionState
    retryCount    int
    lastError     error

    onStateChange func(ConnectionState)
    onRecover     func() tea.Cmd
}

type ConnectionState int

const (
    StateHealthy ConnectionState = iota
    StateUnstable
    StateLost
    StateReconnecting
    StateFailed  // Max retries exceeded
)

func New(opts ...Option) *Manager {
    return &Manager{
        maxRetries:    5,
        baseDelay:     time.Second,
        maxDelay:      30 * time.Second,
        backoffFactor: 2.0,
        currentState:  StateHealthy,
    }
}
```

### Options

```go
func WithMaxRetries(n int) Option {
    return func(m *Manager) { m.maxRetries = n }
}

func WithBaseDelay(d time.Duration) Option {
    return func(m *Manager) { m.baseDelay = d }
}

func WithMaxDelay(d time.Duration) Option {
    return func(m *Manager) { m.maxDelay = d }
}

func WithBackoffFactor(f float64) Option {
    return func(m *Manager) { m.backoffFactor = f }
}

func WithOnStateChange(fn func(ConnectionState)) Option {
    return func(m *Manager) { m.onStateChange = fn }
}
```

### Recovery Flow

```go
// HandleError processes an error and returns recovery command
func (m *Manager) HandleError(err error) tea.Cmd {
    // Categorize error
    if isTransient(err) {
        return m.scheduleRetry()
    }

    if isFatal(err) {
        m.setState(StateFailed)
        return nil
    }

    // Unknown error - attempt recovery
    return m.scheduleRetry()
}

// scheduleRetry calculates delay and returns retry command
func (m *Manager) scheduleRetry() tea.Cmd {
    if m.retryCount >= m.maxRetries {
        m.setState(StateFailed)
        return nil
    }

    // Exponential backoff with jitter
    delay := m.calculateDelay()
    m.retryCount++
    m.setState(StateReconnecting)

    return tea.Tick(delay, func(t time.Time) tea.Msg {
        return RetryMsg{Attempt: m.retryCount}
    })
}

func (m *Manager) calculateDelay() time.Duration {
    delay := float64(m.baseDelay) * math.Pow(m.backoffFactor, float64(m.retryCount))

    // Add jitter (10-20%)
    jitter := delay * (0.1 + rand.Float64()*0.1)
    delay += jitter

    // Cap at max delay
    if delay > float64(m.maxDelay) {
        delay = float64(m.maxDelay)
    }

    return time.Duration(delay)
}

// Reset clears retry state after successful operation
func (m *Manager) Reset() {
    m.retryCount = 0
    m.lastError = nil
    m.setState(StateHealthy)
}
```

---

## Recovery Messages

Messages for TUI integration.

```go
// ConnectionLostMsg signals connection was lost
type ConnectionLostMsg struct {
    Err error
}

// ReconnectingMsg signals retry in progress
type ReconnectingMsg struct {
    Attempt     int
    MaxAttempts int
    NextRetry   time.Duration
}

// ReconnectedMsg signals successful reconnection
type ReconnectedMsg struct{}

// ReconnectFailedMsg signals all retries exhausted
type ReconnectFailedMsg struct {
    Err      error
    Attempts int
}

// RetryMsg triggers a retry attempt
type RetryMsg struct {
    Attempt int
}
```

---

## TUI Integration

### Layout Update Handler

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.StreamErrorMsg:
        return m.handleStreamError(msg.Err)

    case recovery.ConnectionLostMsg:
        m.connectionState = recovery.StateLost
        m.showConnectionLostUI()
        return m, m.recovery.HandleError(msg.Err)

    case recovery.ReconnectingMsg:
        m.connectionState = recovery.StateReconnecting
        m.showReconnectingUI(msg.Attempt, msg.MaxAttempts)
        return m, nil

    case recovery.RetryMsg:
        return m, m.client.ReconnectCmd()

    case recovery.ReconnectedMsg:
        m.recovery.Reset()
        m.connectionState = recovery.StateHealthy
        m.showReconnectedToast()
        return m, nil

    case recovery.ReconnectFailedMsg:
        m.connectionState = recovery.StateFailed
        m.showFailedUI(msg.Err)
        return m, nil
    }
    return m, nil
}
```

### Error Handling Helper

```go
func (m *Model) handleStreamError(err error) (tea.Model, tea.Cmd) {
    // Save partial content if any
    if m.streaming && m.currentContent != "" {
        m.savePartialMessage()
    }

    // Determine error type
    switch {
    case errors.Is(err, context.Canceled):
        // User canceled - not an error
        m.streaming = false
        return m, nil

    case isNetworkError(err):
        // Network issue - attempt recovery
        return m, func() tea.Msg {
            return recovery.ConnectionLostMsg{Err: err}
        }

    case isProcessError(err):
        // CLI process died
        return m, func() tea.Msg {
            return recovery.ConnectionLostMsg{Err: err}
        }

    default:
        // Show error toast, don't attempt recovery
        return m, m.toast.Show(toast.Error, "Error", err.Error())
    }
}
```

---

## Partial Message Recovery

Save and restore partial content during errors.

```go
type PartialMessage struct {
    Content      string
    ContentBlocks []PartialBlock
    Timestamp    time.Time
    LastBlockIdx int
}

type PartialBlock struct {
    Type    string  // "text", "thinking", "tool_use"
    Content string
    Complete bool
}

// Save partial message on error
func (m *Model) savePartialMessage() {
    if m.currentContent == "" {
        return
    }

    m.partialMessage = &PartialMessage{
        Content:      m.currentContent,
        ContentBlocks: m.contentBlocks,
        Timestamp:    time.Now(),
        LastBlockIdx: m.currentBlockIdx,
    }
}

// Restore partial message on reconnect
func (m *Model) restorePartialMessage() {
    if m.partialMessage == nil {
        return
    }

    // Only restore if recent (within 5 minutes)
    if time.Since(m.partialMessage.Timestamp) > 5*time.Minute {
        m.partialMessage = nil
        return
    }

    // Show partial content with indicator
    m.showPartialContent(m.partialMessage)
}
```

---

## UI Components

### Connection Status Indicator

```go
type ConnectionIndicator struct {
    state   recovery.ConnectionState
    attempt int
    maxAttempts int
    styles  ConnectionIndicatorStyles
}

func (c *ConnectionIndicator) View() string {
    switch c.state {
    case recovery.StateHealthy:
        return c.styles.Healthy.Render("[Connected]")

    case recovery.StateUnstable:
        return c.styles.Unstable.Render("[Unstable]")

    case recovery.StateLost:
        return c.styles.Lost.Render("[Disconnected]")

    case recovery.StateReconnecting:
        return c.styles.Reconnecting.Render(
            fmt.Sprintf("[Reconnecting %d/%d...]", c.attempt, c.maxAttempts),
        )

    case recovery.StateFailed:
        return c.styles.Failed.Render("[Connection Failed]")
    }
    return ""
}
```

### Error Toast Messages

```go
var ErrorMessages = map[ConnectionErrorType]string{
    ErrorCLINotFound:  "Claude CLI not found. Please install it first.",
    ErrorCLIVersion:   "Claude CLI version mismatch. Please update.",
    ErrorProcessStart: "Failed to start Claude process.",
    ErrorProcessCrash: "Claude process terminated unexpectedly.",
    ErrorTimeout:      "Connection timed out.",
    ErrorNetwork:      "Network error. Check your connection.",
}

func ShowConnectionError(t *toast.Model, errType ConnectionErrorType) tea.Cmd {
    msg := ErrorMessages[errType]
    return t.Show(toast.Error, "Connection Error", msg)
}
```

### Retry Dialog

```go
type RetryDialog struct {
    error      error
    canRetry   bool
    retrying   bool
    respond    func(retry bool)
    styles     RetryDialogStyles
}

func (d *RetryDialog) View() string {
    title := "Connection Lost"
    message := fmt.Sprintf("Error: %s", d.error)

    var buttons string
    if d.canRetry {
        buttons = lipgloss.JoinHorizontal(
            lipgloss.Center,
            d.styles.Button.Render("Retry"),
            d.styles.Button.Render("Cancel"),
        )
    } else {
        buttons = d.styles.Button.Render("OK")
    }

    return d.styles.Container.Render(
        lipgloss.JoinVertical(
            lipgloss.Center,
            d.styles.Title.Render(title),
            d.styles.Message.Render(message),
            buttons,
        ),
    )
}
```

---

## Error Classification

Determine if errors are retryable.

```go
func isTransient(err error) bool {
    // Network errors are usually transient
    if isNetworkError(err) {
        return true
    }

    // Timeout errors are transient
    if errors.Is(err, context.DeadlineExceeded) {
        return true
    }

    // Process crash might be transient
    var processErr *ProcessError
    if errors.As(err, &processErr) {
        return processErr.ExitCode != 0
    }

    return false
}

func isFatal(err error) bool {
    // CLI not found is fatal
    var notFound *CLINotFoundError
    if errors.As(err, &notFound) {
        return true
    }

    // Version mismatch is fatal
    if strings.Contains(err.Error(), "version") {
        return true
    }

    return false
}

func isNetworkError(err error) bool {
    var netErr net.Error
    return errors.As(err, &netErr)
}

func isProcessError(err error) bool {
    var processErr *ProcessError
    return errors.As(err, &processErr)
}
```

---

## Data Flow Edge Cases

Handling special scenarios.

### Very Long Messages

```go
// Truncate with "Show more" option
func (m *Message) renderWithTruncation(maxLines int) string {
    lines := strings.Split(m.content, "\n")
    if len(lines) <= maxLines {
        return m.content
    }

    truncated := strings.Join(lines[:maxLines], "\n")
    return truncated + "\n" + m.styles.ShowMore.Render("[Show more...]")
}
```

### Horizontal Overflow

```go
// Horizontal scroll or wrap option
type OverflowMode int

const (
    OverflowWrap OverflowMode = iota
    OverflowScroll
    OverflowTruncate
)

func (c *CodeBlock) View() string {
    switch c.overflowMode {
    case OverflowWrap:
        return wordwrap.String(c.code, c.width)
    case OverflowScroll:
        return c.viewport.View()  // Horizontal viewport
    case OverflowTruncate:
        return truncateLines(c.code, c.width)
    }
    return c.code
}
```

### Streaming Interruption

```go
func (m *Model) handleInterrupt() (tea.Model, tea.Cmd) {
    // Save current content
    if m.streaming {
        m.savePartialMessage()
    }

    // Mark message as interrupted
    if m.currentMessage != nil {
        m.currentMessage.Interrupted = true
    }

    // Cancel context
    m.cancel()

    // Show interrupted indicator
    return m, m.toast.Show(toast.Warning, "Interrupted", "Response was interrupted")
}
```

### Concurrent Tool Executions

```go
type ToolQueue struct {
    active    []*ToolExecution
    pending   []*ToolExecution
    maxActive int
}

func (q *ToolQueue) View() string {
    var parts []string

    // Show active tools
    for _, tool := range q.active {
        parts = append(parts, tool.View())
    }

    // Show pending count
    if len(q.pending) > 0 {
        parts = append(parts, fmt.Sprintf("[%d more pending...]", len(q.pending)))
    }

    return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
```

### Large Tool Output

```go
// Paginated or virtualized output
type OutputViewer struct {
    content    string
    viewport   viewport.Model
    pageSize   int
    virtualize bool  // Only render visible portion
}

func (o *OutputViewer) View() string {
    if o.virtualize {
        return o.renderVirtualized()
    }
    return o.viewport.View()
}

func (o *OutputViewer) renderVirtualized() string {
    lines := strings.Split(o.content, "\n")
    start := o.viewport.YOffset
    end := min(start+o.viewport.Height, len(lines))

    visible := lines[start:end]
    return strings.Join(visible, "\n")
}
```

### Binary File Content

```go
func (t *ToolUse) renderFileContent(path string, content []byte) string {
    if isBinary(content) {
        size := len(content)
        return t.styles.Binary.Render(
            fmt.Sprintf("[Binary file: %s, %d bytes]", path, size),
        )
    }
    return string(content)
}

func isBinary(data []byte) bool {
    // Check for null bytes or high proportion of non-printable
    for _, b := range data[:min(512, len(data))] {
        if b == 0 {
            return true
        }
    }
    return false
}
```

### ANSI in Bash Output

```go
type ANSIMode int

const (
    ANSIStrip ANSIMode = iota  // Remove ANSI codes
    ANSIRender                  // Render colors
    ANSIPassthrough             // Pass through as-is
)

func (t *ToolUse) renderBashOutput(output string) string {
    switch t.ansiMode {
    case ANSIStrip:
        return stripANSI(output)
    case ANSIRender:
        return renderANSI(output)  // Convert to lipgloss
    default:
        return output
    }
}

func stripANSI(s string) string {
    re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
    return re.ReplaceAllString(s, "")
}
```

---

## Best Practices

1. **Always handle errors** - Don't ignore returned errors
2. **Categorize errors** - Different handling for different types
3. **Save partial state** - Allow recovery from interruptions
4. **Inform the user** - Show clear error messages
5. **Provide retry options** - Let users trigger manual retry
6. **Graceful degradation** - Continue with reduced functionality when possible
7. **Log errors** - Record for debugging without exposing to user
