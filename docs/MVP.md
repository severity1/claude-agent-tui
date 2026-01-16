# MVP Implementation Plan

Actionable steps to ship a working chat TUI.

## Philosophy

**Ship something that works, then iterate.**

The full spec has 27 components, 12 variants, 11 themes, and 10 animation types. None of that matters until you have a working chat loop.

---

## MVP Scope

A minimal chat TUI that:

1. Connects to Claude via `claude-agent-sdk-go`
2. Streams text responses with a cursor
3. Shows permission prompts for Bash tool
4. Has basic keyboard navigation
5. Handles Ctrl+C to interrupt/quit

**Not in MVP:** Themes, animations, sessions, mouse support, responsive layouts, 12 variants, accessibility (add immediately after MVP works).

---

## Directory Structure

```
claude-agent-tui/
├── go.mod
├── go.sum
├── example/
│   └── chat/
│       └── main.go                    # Working example
├── component/
│   ├── output/
│   │   └── streamtext/
│   │       └── streamtext.go          # Streaming text display
│   ├── input/
│   │   ├── chatinput/
│   │   │   └── chatinput.go           # Text input with submit
│   │   └── permission/
│   │       └── permission.go          # Tool approval prompt
│   └── shared/                        # Shared primitives (added post-MVP)
├── adapter/
│   ├── stream.go                      # SDK events to tea.Msg
│   └── control.go                     # canUseTool callback handler
└── layout/
    └── chat/
        └── chat.go                    # Composes components
```

No `system/`, no `style/`, no theme files. Add when needed.

---

## Implementation Order

### Phase 1: Foundation (Day 1-2)

1. **Project setup**
   ```bash
   go mod init github.com/severity1/claude-agent-tui
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/lipgloss
   go get github.com/charmbracelet/bubbles/textarea
   go get github.com/charmbracelet/bubbles/viewport
   ```

2. **adapter/stream.go** - SDK channel to tea.Msg
   - `StreamTextMsg{Text string}` - text delta
   - `StreamDoneMsg{}` - stream complete
   - `StreamErrorMsg{Err error}` - error occurred
   - `StreamCmd()` - reads from channel, returns next message

3. **component/output/streamtext/streamtext.go** - displays streaming text
   - `Append(text string)` - add text
   - `Clear()` - reset content
   - `SetStreaming(bool)` - show/hide cursor
   - Simple `View()` - content + optional cursor

### Phase 2: Input (Day 2-3)

4. **component/input/chatinput/chatinput.go** - text input
   - Wraps `bubbles/textarea`
   - Enter submits, Alt+Enter newline
   - Emits `SubmitMsg{Text string}`
   - `Focus()`, `Blur()` methods

5. **layout/chat/chat.go** - composes stream + input
   - Stream at top, input at bottom
   - Routes messages to correct component
   - Manages focus state

6. **example/chat/main.go** - working example
   - Test stream + input without SDK
   - Mock data for development

### Phase 3: SDK Integration (Day 3-4)

7. **adapter/control.go** - permission handling
   - `ShowPermissionMsg` with response callback
   - Channel-based blocking for SDK callback
   - Auto-approve tracking per tool

8. **component/input/permission/permission.go** - approval UI
   - Shows tool name and input
   - Allow/Deny/Always options
   - Keyboard shortcuts (a/y, d/n, A)

9. **Update layout/chat/chat.go** - full integration
   - Connect to SDK on init
   - Handle permission overlay
   - Stream queries to display

### Phase 4: Polish (Day 4-5)

10. **Error handling**
    - Connection errors show message
    - Stream errors allow retry
    - Ctrl+C interrupts or quits appropriately

11. **Basic styling**
    - Hardcoded colors (extract to theme later)
    - Borders on permission dialog
    - Visual feedback for streaming state

12. **Test with real SDK**
    - Verify all event types handled
    - Test permission flow
    - Test interruption

---

## Key Interfaces

### StreamText Component

```go
type Model struct {
    content   strings.Builder
    streaming bool
    width     int
}

// Methods
func (m *Model) Append(text string)
func (m *Model) Clear()
func (m *Model) SetStreaming(streaming bool)
func (m Model) View() string
```

### ChatInput Component

```go
type SubmitMsg struct {
    Text string
}

type Model struct {
    textarea textarea.Model
    focused  bool
}

// Methods
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Value() string
```

### Permission Component

```go
type ResultMsg struct {
    Allow       bool
    AlwaysAllow bool
}

type Model struct {
    toolName  string
    toolInput map[string]any
    visible   bool
    respond   func(allow, always bool)
}

// Methods
func (m *Model) Show(toolName string, input map[string]any, respond func(bool, bool))
func (m Model) Visible() bool
```

### Stream Adapter

```go
// Messages
type StreamTextMsg struct{ Text string }
type StreamDoneMsg struct{}
type StreamErrorMsg struct{ Err error }

// Command
func StreamCmd(ctx context.Context, ch <-chan Message) tea.Cmd
```

### Control Adapter

```go
type ShowPermissionMsg struct {
    ToolName  string
    ToolInput map[string]any
    Respond   func(allow, always bool)
}

type ToolControl struct {
    program     *tea.Program
    autoApprove map[string]bool
}

func (t *ToolControl) CanUseTool() func(context.Context, string, map[string]any, string) (bool, error)
```

---

## Message Flow

```
User types -> ChatInput -> SubmitMsg -> Layout
                                          |
                                          v
                                    SDK.Query()
                                          |
                                          v
SDK streams <- StreamCmd() <- msgChan <- SDK
     |
     v
StreamTextMsg -> StreamText.Append()
     |
     v
[Tool request] -> ShowPermissionMsg -> Permission.Show()
                                            |
                                            v
                                     User decision
                                            |
                                            v
                                     Respond callback
                                            |
                                            v
                                     SDK continues
```

---

## Testing Strategy

### Unit Tests

Test files live alongside implementation (e.g., `component/output/streamtext/streamtext_test.go`).

```go
func TestStreamText_Append(t *testing.T) {
    m := streamtext.New()
    m.Append("hello ")
    m.Append("world")

    if m.Content() != "hello world" {
        t.Errorf("expected 'hello world', got %q", m.Content())
    }
}

func TestChatInput_Submit(t *testing.T) {
    m := chatinput.New()
    // Simulate typing and enter
    // Assert SubmitMsg emitted
}

func TestPermission_Shortcuts(t *testing.T) {
    m := permission.New()
    var result permission.ResultMsg

    m.Show("Bash", map[string]any{"command": "ls"}, func(a, aa bool) {
        result = permission.ResultMsg{Allow: a, AlwaysAllow: aa}
    })

    // Press 'y'
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

    if !result.Allow {
        t.Error("expected allow=true after 'y'")
    }
}
```

### Integration Test (layout/chat/chat_test.go)

```go
func TestChatLayout_FullFlow(t *testing.T) {
    // Mock SDK that returns canned responses
    // Verify message routing
    // Verify permission flow
}
```

### Manual Testing Checklist

- [ ] Start app, see input prompt
- [ ] Type message, press Enter, see streaming
- [ ] Ctrl+C during stream stops it
- [ ] Tool request shows permission dialog
- [ ] 'y' allows, continues streaming
- [ ] 'A' allows and remembers
- [ ] 'd' denies
- [ ] Ctrl+C with no stream quits app

---

## Post-MVP Priorities

After MVP works:

1. **Accessibility (immediate)** - REDUCE_MOTION, keyboard-only works
2. **Basic theme** - Extract hardcoded colors to Palette struct
3. **Message history** - Viewport with scroll
4. **More tools** - Write, Edit, Read displays
5. **Error toasts** - Show errors without breaking flow

NOT priorities:
- 12 variants (stick with 2-3)
- 10 animation types (add one when requested)
- 11 themes (one default theme is fine)
- Session persistence (complex, defer)

---

## Success Criteria

MVP is done when:

1. `go run example/chat/main.go` starts without error
2. User can type a message and see Claude respond
3. Permission prompts appear and work
4. Ctrl+C interrupts streams or quits cleanly
5. No panics under normal usage

Everything else is iteration.
