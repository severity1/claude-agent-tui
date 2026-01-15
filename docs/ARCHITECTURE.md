# Architecture

## Project Structure

```
claude-agent-tui/
├── component/
│   ├── output/                    # Display components
│   │   ├── streamtext/            # Live streaming text with cursor
│   │   │   ├── streamtext.go      # Component implementation
│   │   │   └── styles.go          # Lip Gloss style definitions
│   │   ├── message/               # Message bubble (user/assistant)
│   │   ├── thinking/              # Collapsible thinking block
│   │   ├── tooluse/               # Tool invocation card
│   │   ├── codeblock/             # Syntax highlighted code
│   │   ├── status/                # Token counter, model info
│   │   ├── toast/                 # Notification toasts
│   │   ├── diffview/              # Diff visualization for Edit tool
│   │   ├── globresults/           # File list/tree for Glob results
│   │   ├── grepresults/           # Search match display for Grep
│   │   └── taskstatus/            # Background task status indicator
│   │
│   ├── input/                     # Input components
│   │   ├── chat/                  # Main chat input + mode toggles
│   │   ├── permission/            # Tool permission prompt
│   │   ├── question/              # AskUserQuestion multi-choice
│   │   ├── planenter/             # EnterPlanMode confirmation
│   │   ├── planexit/              # ExitPlanMode plan review
│   │   ├── interrupt/             # Stop confirmation modal
│   │   ├── confirm/               # Generic yes/no prompt
│   │   ├── filepicker/            # File selection dialog
│   │   ├── todo/                  # TodoWrite task list
│   │   ├── taskmgr/               # Background task manager
│   │   ├── session/               # Session picker and management
│   │   └── search/                # Message history search
│   │
│   └── shared/                    # Shared primitives
│       ├── button/                # Reusable button
│       ├── chip/                  # Toggle chip/pill
│       ├── card/                  # Card container
│       └── overlay/               # Overlay/modal wrapper
│
├── system/                        # Cross-cutting systems
│   ├── theme/                     # Theming engine
│   │   ├── theme.go               # Theme interface + registry
│   │   ├── palette.go             # Color palette definitions
│   │   └── builtin/               # Built-in themes
│   │       ├── catppuccin.go      # Catppuccin variants
│   │       ├── dracula.go         # Dracula
│   │       ├── nord.go            # Nord
│   │       ├── gruvbox.go         # Gruvbox dark/light
│   │       ├── tokyonight.go      # Tokyo Night
│   │       └── onedark.go         # One Dark
│   │
│   ├── keymap/                    # Keybinding system
│   │   ├── keymap.go              # Keymap registry
│   │   ├── defaults.go            # Default keybindings
│   │   └── help.go                # Help overlay generator
│   │
│   ├── animation/                 # Animation engine
│   │   ├── animation.go           # Animation primitives
│   │   ├── spring.go              # Spring-based motion (harmonica)
│   │   ├── easing.go              # Easing functions
│   │   └── triggers.go            # Event trigger definitions
│   │
│   ├── responsive/                # Responsive layout
│   │   ├── breakpoint.go          # Breakpoint definitions
│   │   ├── resize.go              # Resize handling
│   │   └── constraints.go         # Min/max constraints
│   │
│   ├── mouse/                     # Mouse support
│   │   ├── mouse.go               # Mouse event routing
│   │   ├── hover.go               # Hover state management
│   │   └── zones.go               # Click zone registry
│   │
│   ├── focus/                     # Focus management
│   │   └── focus.go               # Focus ring + traversal
│   │
│   ├── recovery/                  # Error recovery system
│   │   ├── recovery.go            # Recovery manager
│   │   ├── reconnect.go           # Auto-reconnect logic
│   │   └── state.go               # Connection state machine
│   │
│   └── accessibility/             # Accessibility support
│       ├── accessibility.go       # Accessibility config
│       ├── announce.go            # Screen reader announcements
│       └── contrast.go            # High contrast mode
│
├── adapter/                       # SDK integration layer
│   ├── stream.go                  # StreamEvent -> tea.Msg
│   ├── control.go                 # canUseTool -> prompts
│   ├── message.go                 # Message type adapters
│   └── client.go                  # Client lifecycle management
│
├── layout/                        # High-level composite screens
│   ├── chat/                      # Full chat interface
│   │   ├── chat.go                # Layout composition
│   │   ├── state.go               # Input state machine
│   │   └── history.go             # Message history management
│   └── agent/                     # Agent dashboard with tools view
│
├── style/                         # Default styles
│   └── defaults.go                # Sensible defaults
│
└── tui.go                         # Public API exports
```

---

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        Bubble Tea Program                        │
│                                                                  │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│  │  Model   │───►│  Update  │───►│  View    │───►│ Terminal │  │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘  │
│       ▲               │                                         │
│       │               │                                         │
│       │          tea.Msg                                        │
│       │               │                                         │
│  ┌────┴───────────────┴────────────────────────────────────┐   │
│  │                    Adapter Layer                         │   │
│  │                                                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │ StreamCmd() │  │CanUseTool() │  │ ClientCmd() │      │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘      │   │
│  └─────────┼────────────────┼────────────────┼─────────────┘   │
│            │                │                │                  │
└────────────┼────────────────┼────────────────┼──────────────────┘
             │                │                │
             ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    claude-agent-sdk-go                           │
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ StreamEvent │  │ canUseTool  │  │   Client    │              │
│  │   Channel   │  │  Callback   │  │  Lifecycle  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Variants

Components support a small set of display variants. Keep it simple - most components only need 1-2.

```go
type Variant int

const (
    VariantDefault Variant = iota  // Standard display (embeddable, borderless)
    VariantInline                   // Flows within text, minimal chrome
    VariantModal                    // Centered dialog with backdrop
    VariantCompact                  // Minimal footprint for small terminals
    VariantFullscreen               // Takes entire terminal
)
```

**Design principle:** If you need "Card" or "Window" style, wrap the component in a styled container instead of adding more variants. If you need "Collapsible", that's a behavior (Expandable interface), not a variant.

### Component Variant Usage

| Component | Primary Variant | When to use Compact | When to use Modal |
|-----------|----------------|---------------------|-------------------|
| streamtext | Default | Small terminals | Never |
| message | Default | Small terminals | Never |
| thinking | Default (collapsed by default) | - | Expand on click |
| tooluse | Default | Small terminals | View full output |
| codeblock | Default | - | Full-screen view |
| status | Compact | Always use Compact | Never |
| toast | Compact | Always | Never |
| permission | Modal | Inline in small terminal | Default |
| question | Modal | - | Default |
| planexit | Fullscreen | - | - |
| filepicker | Modal | - | Default |
| todo | Default | - | - |

### Expandable Interface

For components that can expand/collapse (thinking, tooluse, todo):

```go
type Expandable interface {
    Expanded() bool
    SetExpanded(bool)
    Toggle()
}
```

This is separate from Variant. A component in `VariantDefault` can still be collapsible.

---

## System Components

### Theming System

The theme system provides consistent styling across all components:

```go
type Theme interface {
    Name() string
    Palette() Palette

    // Component-specific styles
    StreamText() streamtext.Styles
    Message() message.Styles
    // ... etc for each component
}

type Palette struct {
    // Base colors
    Background    lipgloss.AdaptiveColor
    Foreground    lipgloss.AdaptiveColor

    // Surface colors (layered backgrounds)
    Surface0, Surface1, Surface2 lipgloss.AdaptiveColor

    // Semantic colors
    Primary, Secondary, Accent   lipgloss.AdaptiveColor
    Success, Warning, Error, Info lipgloss.AdaptiveColor

    // Text colors
    Text, TextMuted, TextSubtle  lipgloss.AdaptiveColor

    // Syntax highlighting
    Keyword, String, Number, Comment lipgloss.AdaptiveColor
    Function, Variable, Type, Operator lipgloss.AdaptiveColor
}
```

Built-in themes: Catppuccin (4 variants), Dracula, Nord, Gruvbox (2 variants), Tokyo Night, One Dark, System.

### Animation System

Smooth animations for all component state changes:

```go
type Trigger int

const (
    // User input triggers
    TriggerKeyPress Trigger = iota
    TriggerMouseClick
    TriggerMouseHover
    TriggerFocusGain
    TriggerFocusLoss

    // SDK event triggers
    TriggerStreamStart
    TriggerStreamDelta
    TriggerStreamDone
    TriggerToolStart
    TriggerToolComplete
    TriggerToolError

    // State triggers
    TriggerIdle
    TriggerBusy
    TriggerTransition
)
```

Uses `harmonica` for spring-based physics animations.

### Responsive System

Components adapt to terminal dimensions:

```go
type Breakpoint int

const (
    BreakpointXS Breakpoint = iota  // < 40 cols
    BreakpointSM                     // 40-79 cols
    BreakpointMD                     // 80-119 cols
    BreakpointLG                     // 120-159 cols
    BreakpointXL                     // 160+ cols
)
```

### Keybinding System

Customizable keybindings with vim-style defaults:

```go
type Keymap struct {
    // Navigation
    Up, Down, Left, Right key.Binding
    PageUp, PageDown      key.Binding

    // Actions
    Submit, Select, Toggle key.Binding
    Cancel, Help, Quit     key.Binding

    // Component-specific
    Allow, Deny, AlwaysAllow key.Binding
    ToggleAutoAccept, TogglePlanMode key.Binding
}
```

### Mouse Support

Full mouse interaction:

```go
type Zone struct {
    ID       string
    Bounds   Rect
    OnClick  func(x, y int)
    OnHover  func(entered bool)
    OnScroll func(delta int)
}

type MouseManager struct {
    zones   []Zone
    hovered string
}
```

### Recovery System

Automatic error recovery and reconnection:

```go
type RecoveryManager struct {
    maxRetries    int
    retryDelay    time.Duration
    backoffFactor float64
    onDisconnect  func() tea.Cmd
    onReconnect   func() tea.Cmd
}

type ConnectionState int

const (
    ConnectionHealthy ConnectionState = iota
    ConnectionUnstable
    ConnectionLost
    ConnectionReconnecting
)

// Recovery messages
type ConnectionLostMsg struct { Err error }
type ReconnectingMsg struct { Attempt, MaxAttempts int }
type ReconnectedMsg struct{}
type ReconnectFailedMsg struct { Err error }
```

Features:
- Automatic reconnection with exponential backoff
- Partial message recovery
- User notification of connection state
- Manual retry option

### Accessibility System

Support for users with accessibility needs:

```go
type AccessibilityConfig struct {
    ReducedMotion    bool  // Disable animations
    HighContrast     bool  // High contrast theme
    ScreenReader     bool  // Enable announcements
    FocusIndicators  bool  // Visible focus ring
    LargeText        bool  // Larger font sizes
}

// Load from environment
func LoadAccessibilityConfig() AccessibilityConfig {
    return AccessibilityConfig{
        ReducedMotion:   os.Getenv("REDUCE_MOTION") == "1",
        HighContrast:    os.Getenv("HIGH_CONTRAST") == "1",
        ScreenReader:    os.Getenv("SCREEN_READER") == "1",
        FocusIndicators: true,  // Always on by default
    }
}

// Screen reader announcements
type Announcer struct {
    program *tea.Program
}

func (a *Announcer) Announce(message string) {
    // Send announcement to screen reader
    fmt.Fprintf(os.Stderr, "\x1b]777;%s\x07", message)
}
```

Features:
- `REDUCE_MOTION` environment variable support
- High contrast theme variant
- Screen reader announcements
- Visible focus indicators on all interactive elements
- Color-blindness safe indicators (symbols + colors)

---

## Adapter Pattern

The adapter layer bridges the SDK's channel-based streaming to Bubble Tea's message-passing architecture.

### Stream Adapter

Converts SDK `StreamEvent` messages to Bubble Tea messages:

```go
// StreamCmd returns a tea.Cmd that reads from SDK channel
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

// Message types emitted by adapter
type StreamDeltaMsg struct { Text string }
type StreamBlockStartMsg struct { BlockType string; Index int }
type StreamBlockStopMsg struct { Index int }
type StreamDoneMsg struct{}
type StreamErrorMsg struct { Err error }
type AssistantMsg struct { Message *claudecode.AssistantMessage }
type ResultMsg struct { Message *claudecode.ResultMessage }
```

### Control Adapter

Handles the SDK's `canUseTool` callback to show input prompts:

```go
// Permission request from canUseTool callback
type ShowPermissionPromptMsg struct {
    ToolName  string
    ToolInput map[string]any
    ToolUseID string
    Respond   func(decision string, alwaysAllow bool)
}

// Question request from AskUserQuestion tool
type ShowQuestionPromptMsg struct {
    Question    string
    Header      string
    Options     []QuestionOption
    MultiSelect bool
    Respond     func(selections []string, customText string)
}

// Plan mode requests
type ShowPlanEnterPromptMsg struct {
    Message string
    Respond func(confirmed bool)
}

type ShowPlanExitPromptMsg struct {
    PlanContent string
    Permissions []RequestedPermission
    Respond     func(action, feedback string, approvedPerms []string)
}
```

---

## Input State Machine

The chat layout manages which input component has focus:

```
                         ┌─────────────────┐
                         │   ChatInput     │◄── Default state
                         └────────┬────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│PermissionPrompt │    │ QuestionPrompt  │    │ PlanEnterConfirm│
└────────┬────────┘    └────────┬────────┘    └────────┬────────┘
         │                      │                      │
         │                      │                      │ (confirmed)
         │                      │                      ▼
         │                      │             ┌─────────────────┐
         │                      │             │  Planning...    │
         │                      │             │  (streaming)    │
         │                      │             └────────┬────────┘
         │                      │                      │
         │                      │                      ▼
         │                      │             ┌─────────────────┐
         │                      │             │ PlanExitPrompt  │
         │                      │             └────────┬────────┘
         │                      │                      │
         └──────────────────────┴──────────────────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │   ChatInput     │
                         └─────────────────┘
```

### State Enum

```go
type InputState int

const (
    InputStateChat InputState = iota
    InputStatePermission
    InputStateQuestion
    InputStatePlanEnter
    InputStatePlanning
    InputStatePlanExit
    InputStateInterrupt
    InputStateFilePicker
    InputStateTodo
)
```

---

## Component Interface

All components follow the Bubble Tea Model interface:

```go
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

### With Focus Management

Input components include focus tracking:

```go
type FocusableComponent interface {
    Component
    Focus() tea.Cmd
    Blur()
    Focused() bool
}
```

### With Mouse Support

Interactive components implement mouse handling:

```go
type MouseSupport interface {
    Zones() []Zone
    HandleMouse(msg tea.MouseMsg) tea.Cmd
    SetHovered(hovered bool)
    IsHovered() bool
}
```

### With Animation

Animated components implement animation lifecycle:

```go
type AnimatedComponent interface {
    Component
    Animate(trigger Trigger) tea.Cmd
    AnimationConfig() AnimationConfig
}
```

---

## Error Boundaries

Component `View()` methods can panic. The layout should catch and handle these gracefully.

```go
// SafeView wraps a component's View() with panic recovery
func SafeView(c Component, fallback string) string {
    defer func() {
        if r := recover(); r != nil {
            // Log error, return fallback
            log.Printf("panic in View(): %v", r)
        }
    }()
    return c.View()
}

// Usage in layout
func (m *Model) View() string {
    streamView := SafeView(m.stream, "[Error rendering stream]")
    inputView := SafeView(m.input, "[Error rendering input]")
    return lipgloss.JoinVertical(lipgloss.Left, streamView, inputView)
}
```

For Update panics, wrap at the layout level:

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("panic in Update(): %v", r)
            // Return unchanged model, show error toast
        }
    }()
    // ... normal update logic
}
```

---

## State Persistence

When rebuilding components (e.g., theme change), preserve user state.

### Stateful Interface

Components with user-editable state implement:

```go
type Stateful interface {
    // State returns serializable state
    State() any

    // RestoreState applies previously saved state
    RestoreState(any)
}
```

### ChatInput Example

```go
func (c *ChatInput) State() any {
    return chatInputState{
        Value:   c.textarea.Value(),
        Cursor:  c.textarea.CursorPos(),
        Focused: c.focused,
    }
}

func (c *ChatInput) RestoreState(s any) {
    if state, ok := s.(chatInputState); ok {
        c.textarea.SetValue(state.Value)
        c.textarea.SetCursor(state.Cursor)
        if state.Focused {
            c.textarea.Focus()
        }
    }
}
```

### Theme Switch Pattern

```go
func (m *Model) switchTheme(themeName string) {
    // Save state from old components
    inputState := m.input.State()

    // Create new components with new theme
    theme := theme.Get(themeName)
    m.input = chatinput.New(chatinput.WithTheme(theme))

    // Restore state
    m.input.RestoreState(inputState)
}
```

---

## Styling System

Components accept styles via functional options and themes:

```go
// Using functional options
component := streamtext.New(
    streamtext.WithStyles(styles),
    streamtext.WithVariant(VariantInline),
)

// Using theme
theme := theme.Get("catppuccin-mocha")
component := streamtext.New(
    streamtext.WithTheme(theme),
)
```

---

## SDK Integration Points

### Message Types Handled

| SDK Type | tea.Msg | Component |
|----------|---------|-----------|
| `StreamEvent` (content_block_delta) | `StreamDeltaMsg` | StreamText |
| `AssistantMessage` | `AssistantMsg` | Message |
| `ResultMessage` | `ResultMsg` | Status |
| `canUseTool` callback | `ShowPermissionPromptMsg` | Permission |
| `AskUserQuestion` tool | `ShowQuestionPromptMsg` | Question |
| `TodoWrite` tool | `ShowTodoMsg` | Todo |

### Tool Coverage

| Tool | Display Component | Input Component |
|------|-------------------|-----------------|
| Read | ToolUse (file card) | - |
| Write | ToolUse (file card) | Permission, FilePicker |
| Edit | ToolUse, DiffView | Permission |
| Bash | ToolUse (command card) | Permission |
| Glob | ToolUse, GlobResults | FilePicker |
| Grep | ToolUse, GrepResults | - |
| WebSearch | ToolUse (results card) | - |
| WebFetch | ToolUse (content card) | - |
| AskUserQuestion | - | Question |
| Task | Message, TaskMgr (if background) | Permission |
| TaskOutput | TaskMgr | - |
| KillShell | - | TaskMgr (confirmation) |
| Skill | ToolUse (skill card) | - |
| TodoWrite | Todo (task list) | - |
| NotebookEdit | ToolUse (cell view) | Permission |

### Permission Modes

`ChatInput` mode toggles map to SDK permission modes:

| Toggle | SDK Mode |
|--------|----------|
| Auto-accept ON | `PermissionModeAcceptEdits` |
| Plan mode ON | `PermissionModePlan` |
| Both OFF | `PermissionModeDefault` |
