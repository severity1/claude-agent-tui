# Component Specifications

## Component Overview

| Component | Category | SDK Mapping | Charmbracelet Deps |
|-----------|----------|-------------|-------------------|
| streamtext | output | StreamEvent (delta) | lipgloss |
| message | output | AssistantMessage | lipgloss, glamour |
| thinking | output | ThinkingBlock | lipgloss, bubbles/viewport |
| tooluse | output | ToolUseBlock, ToolResultBlock | lipgloss, bubbles/viewport |
| codeblock | output | TextBlock (code) | lipgloss, chroma |
| status | output | ResultMessage | lipgloss |
| toast | output | Errors, completions | lipgloss, harmonica |
| diffview | output | Edit tool (old/new strings) | lipgloss, bubbles/viewport, chroma |
| globresults | output | Glob tool results | lipgloss, bubbles/list |
| grepresults | output | Grep tool results | lipgloss, bubbles/list |
| taskstatus | output | Background task indicator | lipgloss, bubbles/spinner |
| chat | input | User prompt | bubbles/textarea, lipgloss |
| permission | input | canUseTool callback | bubbles/list, lipgloss |
| question | input | AskUserQuestion tool | bubbles/list, bubbles/textarea, lipgloss |
| planenter | input | EnterPlanMode tool | lipgloss |
| planexit | input | ExitPlanMode tool | bubbles/viewport, bubbles/textarea, lipgloss |
| interrupt | input | Context cancellation | lipgloss |
| confirm | input | Generic confirmation | lipgloss |
| filepicker | input | File path inputs | bubbles/filepicker, lipgloss |
| todo | input | TodoWrite tool | bubbles/list, lipgloss |
| taskmgr | input | KillShell, TaskOutput tools | bubbles/list, bubbles/viewport, lipgloss |
| session | input | Session management | bubbles/list, lipgloss |
| search | input | Message history search | bubbles/textinput, bubbles/list, lipgloss |

---

## Output Components

### StreamText

Live-updating text display with optional typing cursor for streaming responses.

**SDK Mapping:** `StreamEvent` with `content_block_delta` type
**Charmbracelet Deps:** `lipgloss`
**Variants:** View, Inline

```go
package streamtext

type Model struct {
    content   strings.Builder
    cursor    bool
    blinking  bool
    variant   Variant
    styles    Styles
    animation AnimationConfig
}

type Styles struct {
    Text      lipgloss.Style
    Cursor    lipgloss.Style
}

// Messages
type DeltaMsg struct { Text string }
type DoneMsg struct{}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithCursor(show bool) Option
func WithBlinking(blink bool) Option
func WithTheme(t Theme) Option
func WithAnimation(a AnimationConfig) Option
```

**Behavior:**
- Appends text on `DeltaMsg`
- Shows blinking cursor while streaming
- Hides cursor on `DoneMsg`
- Supports copy-to-clipboard
- Pulse animation on new delta (optional)

---

### Message

Styled conversation message supporting multiple content block types.

**SDK Mapping:** `AssistantMessage`, `UserMessage`
**Charmbracelet Deps:** `lipgloss`, `glamour`
**Variants:** View, Inline, Card

```go
package message

type Model struct {
    role    Role
    blocks  []Block
    variant Variant
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
func WithVariant(v Variant) Option
func WithBlocks(blocks []Block) Option
func WithTheme(t Theme) Option
```

**Behavior:**
- Renders user messages with right-aligned style
- Renders assistant messages with left-aligned style
- Iterates through content blocks and renders each
- Detects code fences and renders with syntax highlighting
- Markdown rendering via glamour

---

### Thinking

Collapsible extended thinking display.

**SDK Mapping:** `ThinkingBlock` content type
**Charmbracelet Deps:** `lipgloss`, `bubbles/viewport`
**Variants:** View, Collapsible, Modal

```go
package thinking

type Model struct {
    content   string
    collapsed bool
    viewport  viewport.Model
    variant   Variant
    styles    Styles
    hovered   bool
}

type Styles struct {
    Container  lipgloss.Style
    Header     lipgloss.Style
    Content    lipgloss.Style
    Collapsed  lipgloss.Style
    Toggle     lipgloss.Style
}

// Constructor
func New(content string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithCollapsed(collapsed bool) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Enter` or `Space`: Toggle expand/collapse
- `Up/Down` or `j/k`: Scroll when expanded

**Mouse Support:**
- Click header to toggle
- Scroll wheel when expanded
- Hover highlight

**Behavior:**
- Shows "Thinking..." header with toggle indicator
- Expands/collapses on Enter or click
- Scrollable when expanded (uses viewport)
- Shows signature if available

---

### ToolUse

Tool invocation display with input and optional output.

**SDK Mapping:** `ToolUseBlock`, `ToolResultBlock`, `PreToolUseHookInput`, `PostToolUseHookInput`
**Charmbracelet Deps:** `lipgloss`, `bubbles/viewport`, `bubbles/spinner`
**Variants:** View, Card, Inline, Modal

```go
package tooluse

type Model struct {
    toolName   string
    toolInput  map[string]any
    toolOutput any
    isError    bool
    isLoading  bool
    variant    Variant
    spinner    spinner.Model
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
    Spinner      lipgloss.Style
}

// Constructor
func New(name string, input map[string]any, opts ...Option) Model
func FromToolUseBlock(block *claudecode.ToolUseBlock, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithOutput(output any, isError bool) Option
func WithTheme(t Theme) Option
```

**Tool-Specific Display:**
- `Read`: File path with icon
- `Write`: File path with content preview
- `Edit`: Diff view with old/new
- `Bash`: Command with description
- `Glob`: File list
- `Grep`: Match highlights
- `WebSearch`: Results list
- `WebFetch`: URL with status

**Behavior:**
- Displays tool name as header
- Shows JSON-formatted input (collapsible for large inputs)
- Loading spinner while executing
- Updates with output when available
- Error styling for failed tools

---

### CodeBlock

Syntax-highlighted code rendering.

**SDK Mapping:** Code content in `TextBlock`
**Charmbracelet Deps:** `lipgloss`, `chroma`
**Variants:** View, Inline, Modal, Window

```go
package codeblock

type Model struct {
    code       string
    language   string
    showLines  bool
    viewport   viewport.Model
    variant    Variant
    styles     Styles
}

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    Code        lipgloss.Style
    LineNumbers lipgloss.Style
    CopyButton  lipgloss.Style
}

// Constructor
func New(code, language string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithLineNumbers(show bool) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `c`: Copy code to clipboard
- `Up/Down`: Scroll (when large)

**Mouse Support:**
- Click copy button
- Scroll wheel
- Select text (if supported)

**Behavior:**
- Detects language from fence or content
- Applies syntax highlighting (using chroma)
- Optional line numbers
- Copy button/shortcut
- Theme-aware syntax colors

---

### Status

Token count, cost, and model information display.

**SDK Mapping:** `ResultMessage`
**Charmbracelet Deps:** `lipgloss`, `bubbles/spinner`
**Variants:** Bar, Inline, Compact

```go
package status

type Model struct {
    model        string
    inputTokens  int
    outputTokens int
    totalCost    float64
    duration     time.Duration
    isStreaming  bool
    variant      Variant
    spinner      spinner.Model
    styles       Styles
}

type Styles struct {
    Container   lipgloss.Style
    Model       lipgloss.Style
    Tokens      lipgloss.Style
    Cost        lipgloss.Style
    Duration    lipgloss.Style
    Spinner     lipgloss.Style
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
func WithVariant(v Variant) Option
func WithTheme(t Theme) Option
```

**Behavior:**
- Updates incrementally via `UpdateMsg`
- Formats cost as USD
- Shows tokens as "in/out"
- Spinner while streaming
- Compact mode for small terminals

---

### Toast

Notification toast for events and errors.

**SDK Mapping:** Errors, tool completions, stream events
**Charmbracelet Deps:** `lipgloss`, `harmonica`
**Variants:** Toast

```go
package toast

type Model struct {
    toasts     []ToastItem
    position   Position
    maxVisible int
    styles     Styles
}

type ToastItem struct {
    ID        string
    Level     Level
    Title     string
    Message   string
    Duration  time.Duration
    CreatedAt time.Time
}

type Level int

const (
    LevelInfo Level = iota
    LevelSuccess
    LevelWarning
    LevelError
)

type Position int

const (
    PositionTopRight Position = iota
    PositionTopLeft
    PositionBottomRight
    PositionBottomLeft
)

type Styles struct {
    Container lipgloss.Style
    Info      lipgloss.Style
    Success   lipgloss.Style
    Warning   lipgloss.Style
    Error     lipgloss.Style
    Title     lipgloss.Style
    Message   lipgloss.Style
}

// Constructor
func New(opts ...Option) Model

// Methods
func (m *Model) Show(level Level, title, message string)
func (m *Model) Dismiss(id string)
func (m *Model) DismissAll()

// Options
func WithStyles(s Styles) Option
func WithPosition(p Position) Option
func WithMaxVisible(n int) Option
func WithDuration(d time.Duration) Option
func WithTheme(t Theme) Option
```

**Animations:**
- Slide in from edge
- Fade out on dismiss
- Stack management

**Behavior:**
- Auto-dismiss after duration
- Manual dismiss with click or key
- Stacks multiple toasts
- Different styles per level

---

### DiffView

Side-by-side or unified diff visualization for Edit tool operations.

**SDK Mapping:** `Edit` tool `old_string`/`new_string` parameters
**Charmbracelet Deps:** `lipgloss`, `bubbles/viewport`, `chroma`
**Variants:** View, Modal, Window

```go
package diffview

type Model struct {
    filePath   string
    oldContent string
    newContent string
    mode       DiffMode
    viewport   viewport.Model
    lineSync   bool
    variant    Variant
    styles     Styles
}

type DiffMode int

const (
    DiffUnified DiffMode = iota  // Unified diff format (default)
    DiffSplit                     // Side-by-side comparison
    DiffInline                    // Inline with highlights
)

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    FilePath    lipgloss.Style
    Added       lipgloss.Style
    Removed     lipgloss.Style
    Changed     lipgloss.Style
    Context     lipgloss.Style
    LineNumber  lipgloss.Style
    Gutter      lipgloss.Style
}

// Messages
type ContentMsg struct {
    FilePath   string
    OldContent string
    NewContent string
}

// Constructor
func New(filePath, oldContent, newContent string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithMode(m DiffMode) Option
func WithContextLines(n int) Option
func WithSyntaxHighlight(lang string) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Scroll diff
- `m`: Cycle diff mode (unified/split/inline)
- `c`: Copy diff to clipboard

**Mouse Support:**
- Scroll wheel navigation
- Click to select lines

**Behavior:**
- Parses old/new content and generates diff
- Syntax highlighting for both sides
- Line number alignment in split mode
- Scroll synchronization between panes
- Context lines configuration
- Highlights added/removed/changed lines

---

### GlobResults

File list or tree display for Glob tool results.

**SDK Mapping:** `Glob` tool result (file paths)
**Charmbracelet Deps:** `lipgloss`, `bubbles/list`
**Variants:** View, Card, Collapsible

```go
package globresults

type Model struct {
    pattern   string
    path      string
    files     []FileMatch
    treeView  bool
    cursor    int
    collapsed bool
    variant   Variant
    styles    Styles
}

type FileMatch struct {
    Path    string
    ModTime time.Time
    Size    int64
    IsDir   bool
}

type Styles struct {
    Container  lipgloss.Style
    Header     lipgloss.Style
    Pattern    lipgloss.Style
    File       lipgloss.Style
    Directory  lipgloss.Style
    Selected   lipgloss.Style
    Size       lipgloss.Style
    ModTime    lipgloss.Style
    TreeBranch lipgloss.Style
}

// Messages
type ResultsMsg struct {
    Pattern string
    Path    string
    Files   []FileMatch
}

// Constructor
func New(pattern string, files []FileMatch, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithTreeView(enabled bool) Option
func WithPath(path string) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate files
- `Enter`: Select/open file
- `t`: Toggle tree/flat view
- `c`: Copy file path

**Mouse Support:**
- Click to select file
- Double-click to open
- Scroll wheel navigation

**Behavior:**
- Displays file matches from Glob tool
- Tree view shows directory hierarchy
- Flat view shows sorted file list
- Shows file metadata (size, mod time)
- Icons for files vs directories

---

### GrepResults

Search match display for Grep tool results with context.

**SDK Mapping:** `Grep` tool result (matches with context)
**Charmbracelet Deps:** `lipgloss`, `bubbles/list`
**Variants:** View, Card, Collapsible

```go
package grepresults

type Model struct {
    pattern      string
    matches      []GrepMatch
    contextLines int
    cursor       int
    expanded     map[int]bool
    variant      Variant
    styles       Styles
}

type GrepMatch struct {
    File    string
    Line    int
    Content string
    Before  []string  // Context lines before match
    After   []string  // Context lines after match
}

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    Pattern     lipgloss.Style
    FileName    lipgloss.Style
    LineNumber  lipgloss.Style
    Match       lipgloss.Style
    Highlight   lipgloss.Style
    Context     lipgloss.Style
    Selected    lipgloss.Style
}

// Messages
type ResultsMsg struct {
    Pattern string
    Matches []GrepMatch
}

// Constructor
func New(pattern string, matches []GrepMatch, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithContextLines(n int) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate matches
- `Enter`: Jump to file location
- `Space`: Toggle context expansion
- `c`: Copy match line

**Mouse Support:**
- Click to select match
- Double-click to open file
- Scroll wheel navigation

**Behavior:**
- Groups matches by file
- Highlights matching text
- Shows configurable context lines
- Expandable/collapsible context
- Navigate to file:line on selection

---

### TaskStatus

Background task status indicator showing active tasks.

**SDK Mapping:** `run_in_background` parameter, `TaskOutput` tool
**Charmbracelet Deps:** `lipgloss`, `bubbles/spinner`
**Variants:** Bar, Compact, Inline

```go
package taskstatus

type Model struct {
    tasks     []BackgroundTask
    spinner   spinner.Model
    variant   Variant
    styles    Styles
}

type BackgroundTask struct {
    ID          string
    Type        string       // "bash" or "agent"
    Description string
    Status      TaskState
    StartedAt   time.Time
    Progress    float64      // 0.0 to 1.0 if known
}

type TaskState int

const (
    TaskRunning TaskState = iota
    TaskCompleted
    TaskFailed
    TaskKilled
)

type Styles struct {
    Container   lipgloss.Style
    TaskCount   lipgloss.Style
    TaskName    lipgloss.Style
    Running     lipgloss.Style
    Completed   lipgloss.Style
    Failed      lipgloss.Style
    Spinner     lipgloss.Style
}

// Messages
type UpdateMsg struct {
    Tasks []BackgroundTask
}

type TaskStartedMsg struct {
    Task BackgroundTask
}

type TaskCompletedMsg struct {
    TaskID string
    Status TaskState
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithTheme(t Theme) Option
```

**Behavior:**
- Shows count of running background tasks
- Spinner animation while tasks active
- Click to open task manager
- Compact mode shows just count
- Bar mode shows task list

---

## Input Components

### ChatInput

Main conversation input with mode toggles.

**SDK Mapping:** User prompt submission, `QueryParams`
**Charmbracelet Deps:** `bubbles/textarea`, `lipgloss`
**Variants:** Fullscreen, Window, View, Overlay

```go
package chatinput

type Model struct {
    textarea   textarea.Model
    autoAccept bool
    planMode   bool
    variant    Variant
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
func WithVariant(v Variant) Option
func WithPlaceholder(text string) Option
func WithAutoAccept(enabled bool) Option
func WithPlanMode(enabled bool) Option
func WithTheme(t Theme) Option
func WithKeymap(k Keymap) Option
```

**Keybindings:**
- `Enter` or `Ctrl+Enter`: Submit
- `Ctrl+A`: Toggle auto-accept
- `Ctrl+P`: Toggle plan mode
- `Esc`: Clear input

**Mouse Support:**
- Click to focus
- Click toggle chips

**Behavior:**
- Multi-line text input
- Shows toggle states as pills/chips
- Emits `SubmitMsg` with current mode states

---

### Permission

Tool approval dialog.

**SDK Mapping:** `canUseTool` callback
**Charmbracelet Deps:** `bubbles/list`, `lipgloss`
**Variants:** Modal, Overlay, Inline

```go
package permission

type Model struct {
    toolName   string
    toolInput  map[string]any
    toolUseID  string
    selected   int
    options    []string
    variant    Variant
    styles     Styles
    respond    func(decision string, alwaysAllow bool)
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
func WithVariant(v Variant) Option
func WithOptions(options []string) Option
func WithResponder(fn func(string, bool)) Option
func WithTheme(t Theme) Option
func WithKeymap(k Keymap) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate options
- `Enter`: Select option
- `a` or `y`: Quick allow
- `d` or `n`: Quick deny
- `A`: Always allow

**Mouse Support:**
- Click options
- Hover highlight

**Danger Detection:**
- Highlights dangerous commands (rm -rf, sudo, etc.)
- Warning styling for destructive operations

**Behavior:**
- Shows tool name and formatted input
- Highlights dangerous tools
- Returns decision via respond callback

---

### Question

Multi-choice question response (AskUserQuestion).

**SDK Mapping:** `AskUserQuestion` tool
**Charmbracelet Deps:** `bubbles/list`, `bubbles/textarea`, `lipgloss`
**Variants:** Modal, Overlay, View

```go
package question

type Model struct {
    question    string
    header      string
    options     []QuestionOption
    multiSelect bool
    cursor      int
    selected    map[int]bool
    customInput textarea.Model
    showCustom  bool
    variant     Variant
    styles      Styles
    respond     func(selections []string, customText string)
}

type QuestionOption struct {
    Label       string
    Description string
}

type Styles struct {
    Container   lipgloss.Style
    Question    lipgloss.Style
    Header      lipgloss.Style
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
func New(question string, options []QuestionOption, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithMultiSelect(enabled bool) Option
func WithHeader(h string) Option
func WithResponder(fn func([]string, string)) Option
func WithTheme(t Theme) Option
func WithKeymap(k Keymap) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate options
- `Space`: Toggle selection (multi-select)
- `Enter`: Submit selection
- `o`: Open custom input ("Other")
- `Esc`: Cancel custom input

**Mouse Support:**
- Click options to select
- Click "Other" for custom input

**Behavior:**
- Single or multi-select mode
- Shows option descriptions
- "Other" option opens text input
- Returns selections via respond callback

---

### PlanEnter

Plan mode entry confirmation.

**SDK Mapping:** `EnterPlanMode` tool
**Charmbracelet Deps:** `lipgloss`
**Variants:** Modal, Overlay

```go
package planenter

type Model struct {
    message  string
    focused  int
    variant  Variant
    styles   Styles
    respond  func(confirmed bool)
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
func WithVariant(v Variant) Option
func WithResponder(fn func(bool)) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Left/Right` or `Tab`: Navigate buttons
- `Enter`: Select button
- `y`: Quick confirm
- `n`: Quick skip

**Behavior:**
- Shows Claude's reason for planning
- Two buttons: "Confirm" and "Skip (just implement)"
- Returns decision via respond callback

---

### PlanExit

Plan review and action selection.

**SDK Mapping:** `ExitPlanMode` tool, `allowedPrompts`
**Charmbracelet Deps:** `bubbles/viewport`, `bubbles/textarea`, `lipgloss`, `glamour`
**Variants:** Modal, Fullscreen, Window

```go
package planexit

type Model struct {
    planContent  string
    permissions  []Permission
    viewport     viewport.Model
    focused      int
    options      []Action
    feedback     textarea.Model
    showFeedback bool
    variant      Variant
    styles       Styles
    respond      func(action, feedback string, approvedPerms []string)
}

type Permission struct {
    Tool     string
    Prompt   string
    Approved bool
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
    Permissions []string
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
func WithVariant(v Variant) Option
func WithActions(actions []Action) Option
func WithResponder(fn func(string, string, []string)) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Scroll plan
- `Left/Right` or `Tab`: Navigate actions
- `Enter`: Select action
- `a`: Quick approve
- `r`: Quick reject
- `c`: Request changes (opens feedback)
- `Space`: Toggle permission checkbox

**Mouse Support:**
- Scroll plan content
- Click action buttons
- Click permission checkboxes

**Behavior:**
- Scrollable plan view (markdown rendered)
- Permission checkboxes
- Feedback input for "Request Changes"
- Returns decision via respond callback

---

### Interrupt

Stop/cancel confirmation modal.

**SDK Mapping:** Context cancellation, `Interrupt()`
**Charmbracelet Deps:** `lipgloss`
**Variants:** Modal, Overlay

```go
package interrupt

type Model struct {
    focused int
    variant Variant
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
func WithVariant(v Variant) Option
func WithTheme(t Theme) Option
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

---

### Confirm

Generic yes/no confirmation dialog.

**SDK Mapping:** Generic confirmations
**Charmbracelet Deps:** `lipgloss`
**Variants:** Modal, Inline

```go
package confirm

type Model struct {
    title    string
    message  string
    focused  int
    variant  Variant
    styles   Styles
}

type Styles struct {
    Container   lipgloss.Style
    Title       lipgloss.Style
    Message     lipgloss.Style
    Button      lipgloss.Style
    ButtonFocus lipgloss.Style
}

// Messages
type ResultMsg struct {
    Confirmed bool
}

// Constructor
func New(title, message string, opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Left/Right` or `Tab`: Navigate buttons
- `Enter`: Select button
- `y`: Quick confirm
- `n`: Quick cancel

**Behavior:**
- Displays title and message
- Two buttons: "Yes" and "No"
- Returns decision via `ResultMsg`

---

### FilePicker

File selection dialog.

**SDK Mapping:** File path inputs in tools (Write, Glob)
**Charmbracelet Deps:** `bubbles/filepicker`, `lipgloss`
**Variants:** Modal, Overlay, Window

```go
package filepicker

type Model struct {
    picker   filepicker.Model
    variant  Variant
    styles   Styles
}

type Styles struct {
    Container lipgloss.Style
    Selected  lipgloss.Style
    Directory lipgloss.Style
    File      lipgloss.Style
}

// Messages
type ResultMsg struct {
    Path string
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithCurrentDir(dir string) Option
func WithAllowedTypes(types []string) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate
- `Enter`: Select file/enter directory
- `Backspace`: Go up directory
- `Esc`: Cancel

**Mouse Support:**
- Click to select
- Double-click to enter/select

**Behavior:**
- Navigate filesystem
- Filter by file type
- Returns selected path

---

### Todo

TodoWrite task list display.

**SDK Mapping:** `TodoWrite` tool
**Charmbracelet Deps:** `bubbles/list`, `lipgloss`
**Variants:** View, Sidebar, Overlay, Collapsible

```go
package todo

type Model struct {
    items     []TodoItem
    collapsed bool
    variant   Variant
    styles    Styles
}

type TodoItem struct {
    Content    string
    Status     Status
    ActiveForm string
}

type Status int

const (
    StatusPending Status = iota
    StatusInProgress
    StatusCompleted
)

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    Item        lipgloss.Style
    Pending     lipgloss.Style
    InProgress  lipgloss.Style
    Completed   lipgloss.Style
    Checkbox    lipgloss.Style
}

// Messages
type UpdateMsg struct {
    Items []TodoItem
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithCollapsed(collapsed bool) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Enter` or `Space`: Toggle collapse
- `Up/Down`: Scroll items

**Mouse Support:**
- Click header to collapse
- Scroll items

**Behavior:**
- Displays task list with status indicators
- Collapsible for space saving
- Updates from SDK TodoWrite events
- Shows progress (X/Y completed)

---

### TaskManager

Background task management panel for viewing and controlling background tasks.

**SDK Mapping:** `KillShell` tool, `TaskOutput` tool, `run_in_background` parameter
**Charmbracelet Deps:** `bubbles/list`, `bubbles/viewport`, `lipgloss`
**Variants:** Sidebar, Overlay, Modal, Collapsible

```go
package taskmgr

type Model struct {
    tasks      []BackgroundTask
    selected   int
    expanded   map[string]bool
    outputView viewport.Model
    variant    Variant
    styles     Styles
    respond    func(action TaskAction, taskID string)
}

type BackgroundTask struct {
    ID          string
    Type        string       // "bash" or "agent"
    Description string
    Status      TaskState
    StartedAt   time.Time
    Output      string
    OutputFile  string       // For incremental output reading
    ExitCode    *int         // nil while running
}

type TaskState int

const (
    TaskRunning TaskState = iota
    TaskCompleted
    TaskFailed
    TaskKilled
)

type TaskAction int

const (
    ActionViewOutput TaskAction = iota
    ActionKill
    ActionRefresh
)

type Styles struct {
    Container    lipgloss.Style
    Header       lipgloss.Style
    TaskItem     lipgloss.Style
    Selected     lipgloss.Style
    Running      lipgloss.Style
    Completed    lipgloss.Style
    Failed       lipgloss.Style
    OutputView   lipgloss.Style
    ActionButton lipgloss.Style
}

// Messages
type ShowMsg struct {
    Tasks []BackgroundTask
}

type UpdateMsg struct {
    Tasks []BackgroundTask
}

type TaskOutputMsg struct {
    TaskID string
    Output string
    Done   bool
}

type KillConfirmMsg struct {
    TaskID  string
    Confirm bool
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithResponder(fn func(TaskAction, string)) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate tasks
- `Enter`: Expand/view task output
- `k` or `Delete`: Kill selected task (with confirmation)
- `r`: Refresh task status
- `Esc`: Close panel/collapse output

**Mouse Support:**
- Click to select task
- Click expand button to view output
- Click kill button to terminate
- Scroll wheel in output view

**Behavior:**
- Lists all background tasks with status
- Expandable output view per task
- Kill confirmation before terminating
- Auto-refresh task status
- Shows task duration and exit code

---

### Session

Session picker and management for multi-session workflows.

**SDK Mapping:** Session management, conversation resume
**Charmbracelet Deps:** `bubbles/list`, `lipgloss`
**Variants:** Modal, Overlay, Sidebar

```go
package session

type Model struct {
    sessions   []SessionInfo
    current    string
    selected   int
    variant    Variant
    styles     Styles
    respond    func(action SessionAction, sessionID string)
}

type SessionInfo struct {
    ID           string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    MessageCount int
    TokensUsed   int
    Model        string
    Summary      string       // First user message or auto-summary
    IsActive     bool
}

type SessionAction int

const (
    ActionSelect SessionAction = iota
    ActionResume
    ActionDelete
    ActionExport
)

type Styles struct {
    Container   lipgloss.Style
    Header      lipgloss.Style
    SessionItem lipgloss.Style
    Selected    lipgloss.Style
    Active      lipgloss.Style
    Summary     lipgloss.Style
    Metadata    lipgloss.Style
    ActionBtn   lipgloss.Style
}

// Messages
type ShowMsg struct {
    Sessions []SessionInfo
    Current  string
}

type SelectedMsg struct {
    SessionID string
    Action    SessionAction
}

type ResumedMsg struct {
    SessionID string
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithResponder(fn func(SessionAction, string)) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- `Up/Down` or `j/k`: Navigate sessions
- `Enter`: Select/resume session
- `n`: New session
- `d`: Delete session (with confirmation)
- `e`: Export session
- `Esc`: Close picker

**Mouse Support:**
- Click to select session
- Double-click to resume
- Click action buttons

**Behavior:**
- Lists previous sessions with metadata
- Shows current active session
- Session summary from first message
- Displays token usage and message count
- Supports session resume and delete

---

### Search

Message history search with result highlighting.

**SDK Mapping:** Local conversation history search
**Charmbracelet Deps:** `bubbles/textinput`, `bubbles/list`, `lipgloss`
**Variants:** Overlay, Modal, Inline

```go
package search

type Model struct {
    query      textinput.Model
    results    []SearchResult
    selected   int
    searching  bool
    variant    Variant
    styles     Styles
    respond    func(result SearchResult)
}

type SearchResult struct {
    MessageIndex int
    Role         string        // "user" or "assistant"
    Snippet      string        // Matched text with context
    Highlights   []HighlightRange
    Timestamp    time.Time
}

type HighlightRange struct {
    Start int
    End   int
}

type Styles struct {
    Container   lipgloss.Style
    Input       lipgloss.Style
    ResultItem  lipgloss.Style
    Selected    lipgloss.Style
    Highlight   lipgloss.Style
    Snippet     lipgloss.Style
    Role        lipgloss.Style
    Timestamp   lipgloss.Style
    NoResults   lipgloss.Style
}

// Messages
type QueryMsg struct {
    Query string
}

type ResultsMsg struct {
    Results []SearchResult
}

type SelectMsg struct {
    Result SearchResult
}

// Constructor
func New(opts ...Option) Model

// Options
func WithStyles(s Styles) Option
func WithVariant(v Variant) Option
func WithResponder(fn func(SearchResult)) Option
func WithTheme(t Theme) Option
```

**Keybindings:**
- Type to search (live filtering)
- `Up/Down` or `j/k`: Navigate results
- `Enter`: Jump to message
- `Esc`: Close search
- `Ctrl+L`: Clear query

**Mouse Support:**
- Click input to focus
- Click result to select
- Double-click to jump to message

**Behavior:**
- Live search as you type
- Highlights matching text in results
- Shows message context around match
- Navigate to message on selection
- Supports regex patterns

---

## Shared Components

### Button

Reusable button primitive.

```go
package button

type Model struct {
    label   string
    focused bool
    styles  Styles
}

type Styles struct {
    Normal  lipgloss.Style
    Focused lipgloss.Style
    Pressed lipgloss.Style
}
```

### Chip

Toggle chip/pill for mode toggles.

```go
package chip

type Model struct {
    label   string
    active  bool
    styles  Styles
}

type Styles struct {
    Active   lipgloss.Style
    Inactive lipgloss.Style
}
```

### Card

Card container with border and padding.

```go
package card

type Model struct {
    content string
    title   string
    styles  Styles
}

type Styles struct {
    Container lipgloss.Style
    Title     lipgloss.Style
    Content   lipgloss.Style
}
```

### Overlay

Overlay/modal wrapper with backdrop.

```go
package overlay

type Model struct {
    content tea.Model
    styles  Styles
}

type Styles struct {
    Backdrop lipgloss.Style
    Content  lipgloss.Style
}
```
