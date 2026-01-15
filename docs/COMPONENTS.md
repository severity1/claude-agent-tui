# Component Specifications

## Output Components

### StreamText

Live-updating text display with optional typing cursor for streaming responses.

```go
package streamtext

type Model struct {
    content   strings.Builder
    cursor    bool
    blinking  bool
    styles    Styles
}

type Styles struct {
    Text      lipgloss.Style
    Cursor    lipgloss.Style
}

// Messages
type DeltaMsg struct {
    Text string
}

type DoneMsg struct{}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithCursor(show bool) Option
func WithBlinking(blink bool) Option
```

**Behavior:**
- Appends text on `DeltaMsg`
- Shows blinking cursor while streaming
- Hides cursor on `DoneMsg`
- Supports copy-to-clipboard

---

### MessageBubble

Styled conversation message supporting multiple content block types.

```go
package message

type Model struct {
    role    Role        // User or Assistant
    blocks  []Block     // Content blocks
    styles  Styles
}

type Role string

const (
    RoleUser      Role = "user"
    RoleAssistant Role = "assistant"
)

type Block interface {
    View(styles Styles) string
}

type Styles struct {
    UserContainer      lipgloss.Style
    AssistantContainer lipgloss.Style
    TextBlock          lipgloss.Style
    CodeBlock          lipgloss.Style
    ToolUseBlock       lipgloss.Style
    ThinkingBlock      lipgloss.Style
}

// Constructor
func New(role Role, opts ...Option) Model
func FromAssistantMessage(msg *claudecode.AssistantMessage, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithBlocks(blocks []Block) Option
```

**Behavior:**
- Renders user messages with right-aligned style
- Renders assistant messages with left-aligned style
- Iterates through content blocks and renders each
- Detects code fences and renders with syntax highlighting

---

### ThinkingBlock

Collapsible extended thinking display.

```go
package thinking

type Model struct {
    content   string
    collapsed bool
    viewport  viewport.Model
    styles    Styles
}

type Styles struct {
    Container  lipgloss.Style
    Header     lipgloss.Style
    Content    lipgloss.Style
    Collapsed  lipgloss.Style
}

// Constructor
func New(content string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithCollapsed(collapsed bool) Option
```

**Behavior:**
- Shows "Thinking..." header with toggle indicator
- Expands/collapses on Enter or click
- Scrollable when expanded (uses viewport)
- Shows signature if available

---

### ToolUseCard

Tool invocation display with input and optional output.

```go
package tooluse

type Model struct {
    toolName   string
    toolInput  map[string]any
    toolOutput any
    isError    bool
    styles     Styles
}

type Styles struct {
    Container    lipgloss.Style
    ToolName     lipgloss.Style
    InputLabel   lipgloss.Style
    InputJSON    lipgloss.Style
    OutputLabel  lipgloss.Style
    OutputJSON   lipgloss.Style
    ErrorOutput  lipgloss.Style
}

// Constructor
func New(name string, input map[string]any, opts ...Option) Model
func FromToolUseBlock(block *claudecode.ToolUseBlock, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithOutput(output any, isError bool) Option
```

**Behavior:**
- Displays tool name as header
- Shows JSON-formatted input
- Updates with output when available
- Error styling for failed tools

---

### CodeBlock

Syntax-highlighted code rendering.

```go
package codeblock

type Model struct {
    code     string
    language string
    styles   Styles
}

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    Code        lipgloss.Style
    LineNumbers lipgloss.Style
}

// Constructor
func New(code, language string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithLineNumbers(show bool) Option
```

**Behavior:**
- Detects language from fence or content
- Applies syntax highlighting (using chroma or similar)
- Optional line numbers
- Copy button/shortcut

---

### StatusBar

Token count, cost, and model information display.

```go
package status

type Model struct {
    model       string
    inputTokens  int
    outputTokens int
    totalCost   float64
    duration    time.Duration
    styles      Styles
}

type Styles struct {
    Container   lipgloss.Style
    Model       lipgloss.Style
    Tokens      lipgloss.Style
    Cost        lipgloss.Style
    Duration    lipgloss.Style
}

// Messages
type UpdateMsg struct {
    Model        *string
    InputTokens  *int
    OutputTokens *int
    TotalCost    *float64
    Duration     *time.Duration
}

// Constructor
func New(opts ...Option) Model
func FromResultMessage(msg *claudecode.ResultMessage, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
```

**Behavior:**
- Updates incrementally via `UpdateMsg`
- Formats cost as USD
- Shows tokens as "in/out"

---

## Input Components

### ChatInput

Main conversation input with mode toggles.

```go
package chatinput

type Model struct {
    textarea   textarea.Model
    autoAccept bool
    planMode   bool
    styles     Styles
}

type Styles struct {
    Container   lipgloss.Style
    TextArea    lipgloss.Style
    ToggleOn    lipgloss.Style
    ToggleOff   lipgloss.Style
    Label       lipgloss.Style
    Hint        lipgloss.Style
}

// Messages
type SubmitMsg struct {
    Text       string
    AutoAccept bool
    PlanMode   bool
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithPlaceholder(text string) Option
func WithAutoAccept(enabled bool) Option
func WithPlanMode(enabled bool) Option
```

**Keybindings:**
- `Enter` or `Ctrl+Enter`: Submit
- `Ctrl+A`: Toggle auto-accept
- `Ctrl+P`: Toggle plan mode
- `Esc`: Clear input

**Behavior:**
- Multi-line text input
- Shows toggle states as pills/chips
- Emits `SubmitMsg` with current mode states

---

### PermissionPrompt

Tool approval dialog.

```go
package permission

type Model struct {
    toolName   string
    toolInput  map[string]any
    selected   int
    options    []string
    styles     Styles
}

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    ToolName    lipgloss.Style
    InputJSON   lipgloss.Style
    Option      lipgloss.Style
    Selected    lipgloss.Style
    Warning     lipgloss.Style
}

// Messages
type ResultMsg struct {
    Decision    string  // "allow", "deny"
    AlwaysAllow bool
}

// Constructor
func New(toolName string, toolInput map[string]any, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithOptions(options []string) Option  // Default: Allow, Deny, Always Allow
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate options
- `Enter`: Select option
- `a`: Quick allow
- `d`: Quick deny
- `A`: Always allow

**Behavior:**
- Shows tool name and formatted input
- Highlights dangerous tools (rm, sudo, etc.)
- Returns decision via `ResultMsg`

---

### QuestionPrompt

Multi-choice question response (AskUserQuestion).

```go
package question

type Model struct {
    question    string
    options     []Option
    multiSelect bool
    cursor      int
    selected    map[int]bool
    customInput textarea.Model
    showCustom  bool
    styles      Styles
}

type Option struct {
    Label       string
    Description string
}

type Styles struct {
    Container   lipgloss.Style
    Question    lipgloss.Style
    Option      lipgloss.Style
    Selected    lipgloss.Style
    Description lipgloss.Style
    CustomInput lipgloss.Style
}

// Messages
type ResultMsg struct {
    Selections []string
    CustomText string
}

// Constructor
func New(question string, options []Option, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithMultiSelect(enabled bool) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate options
- `Space`: Toggle selection (multi-select)
- `Enter`: Submit selection
- `o`: Open custom input ("Other")
- `Esc`: Cancel custom input

**Behavior:**
- Single or multi-select mode
- Shows option descriptions
- "Other" option opens text input
- Returns selections via `ResultMsg`

---

### PlanEnterConfirm

Plan mode entry confirmation.

```go
package planenter

type Model struct {
    message  string
    focused  int
    styles   Styles
}

type Styles struct {
    Container   lipgloss.Style
    Message     lipgloss.Style
    Button      lipgloss.Style
    ButtonFocus lipgloss.Style
}

// Messages
type ResultMsg struct {
    Confirmed bool
}

// Constructor
func New(message string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
```

**Keybindings:**
- `Left/Right` or `Tab`: Navigate buttons
- `Enter`: Select button
- `y`: Quick confirm
- `n`: Quick skip

**Behavior:**
- Shows Claude's reason for planning
- Two buttons: "Confirm" and "Skip (just implement)"
- Returns decision via `ResultMsg`

---

### PlanExitPrompt

Plan review and action selection.

```go
package planexit

type Model struct {
    planContent string
    permissions []Permission
    viewport    viewport.Model
    focused     int
    options     []Action
    feedback    textarea.Model
    showFeedback bool
    styles      Styles
}

type Permission struct {
    Tool   string
    Prompt string
}

type Action struct {
    Label string
    Key   string
}

type Styles struct {
    Container       lipgloss.Style
    Header          lipgloss.Style
    PlanContent     lipgloss.Style
    PermissionList  lipgloss.Style
    PermissionItem  lipgloss.Style
    Option          lipgloss.Style
    OptionFocus     lipgloss.Style
    FeedbackInput   lipgloss.Style
}

// Messages
type ResultMsg struct {
    Action      string  // "approve", "reject", "changes"
    Feedback    string
    Permissions []string  // Approved permission prompts
}

// Default actions
var DefaultActions = []Action{
    {Label: "Approve & Implement", Key: "a"},
    {Label: "Reject Plan", Key: "r"},
    {Label: "Request Changes", Key: "c"},
}

// Constructor
func New(plan string, permissions []Permission, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithActions(actions []Action) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Scroll plan
- `Left/Right` or `Tab`: Navigate actions
- `Enter`: Select action
- `a`: Quick approve
- `r`: Quick reject
- `c`: Request changes (opens feedback)
- `Space`: Toggle permission checkbox

**Behavior:**
- Scrollable plan view (markdown rendered)
- Permission checkboxes
- Feedback input for "Request Changes"
- Returns decision via `ResultMsg`

---

### InterruptConfirm

Stop/cancel confirmation modal.

```go
package interrupt

type Model struct {
    focused int
    styles  Styles
}

type Styles struct {
    Container   lipgloss.Style
    Message     lipgloss.Style
    Button      lipgloss.Style
    ButtonFocus lipgloss.Style
}

// Messages
type ResultMsg struct {
    Confirmed bool
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
```

**Keybindings:**
- `Left/Right` or `Tab`: Navigate buttons
- `Enter`: Select button
- `y`: Quick confirm
- `n`: Quick cancel
- `Esc`: Cancel

**Behavior:**
- Modal overlay
- "Stop generation?" message
- Returns decision via `ResultMsg`
