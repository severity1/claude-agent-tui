# SDK Mapping

Complete mapping between `claude-agent-sdk-go` types and TUI components.

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     claude-agent-sdk-go                          │
│                                                                  │
│  Events          Callbacks         Tools           Lifecycle    │
│  ───────         ─────────         ─────           ─────────    │
│  StreamEvent     canUseTool        AskUserQuestion Connect      │
│  AssistantMsg    PreToolUse        EnterPlanMode   Query        │
│  ResultMsg       PostToolUse       ExitPlanMode    Disconnect   │
│                                    TodoWrite       Interrupt    │
│                                    TaskOutput                   │
│                                    KillShell                    │
│                                    Skill                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Adapter Layer                               │
│                                                                  │
│  stream.go       control.go        message.go      client.go    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     TUI Components                               │
│                                                                  │
│  Output          Input             System                        │
│  ──────          ─────             ──────                        │
│  streamtext      chat              theme                         │
│  message         permission        keymap                        │
│  thinking        question          animation                     │
│  tooluse         planenter         responsive                    │
│  codeblock       planexit          mouse                         │
│  status          interrupt         focus                         │
│  toast           filepicker        recovery                      │
│  diffview        todo              accessibility                 │
│  globresults     taskmgr                                         │
│  grepresults     session                                         │
│  taskstatus      search                                          │
└─────────────────────────────────────────────────────────────────┘
```

---

## Stream Events

Events received from `ReceiveMessages()` channel.

### content_block_start

Signals the beginning of a new content block.

```go
// SDK Event
type StreamEvent struct {
    Event map[string]any  // type: "content_block_start"
}

// Adapter Message
type StreamBlockStartMsg struct {
    BlockType string  // "text", "tool_use", "thinking"
    Index     int
}

// Component Response
// - streamtext: Initialize new text buffer
// - tooluse: Create new tool card (if tool_use)
// - thinking: Create new thinking block (if thinking)
```

### content_block_delta

Incremental content update during streaming.

```go
// SDK Event
type StreamEvent struct {
    Event map[string]any  // type: "content_block_delta"
    // delta.text contains the text chunk
}

// Adapter Message
type StreamDeltaMsg struct {
    Text  string
    Index int
}

// Component Response
// - streamtext: Append text, show cursor
// - Animation: Pulse on delta (optional)
```

### content_block_stop

Signals completion of a content block.

```go
// SDK Event
type StreamEvent struct {
    Event map[string]any  // type: "content_block_stop"
}

// Adapter Message
type StreamBlockStopMsg struct {
    Index int
}

// Component Response
// - streamtext: Hide cursor, finalize content
// - message: Convert stream to permanent message block
```

### message_stop

Signals the end of assistant's response.

```go
// SDK Event
type StreamEvent struct {
    Event map[string]any  // type: "message_stop"
}

// Adapter Message
type StreamDoneMsg struct{}

// Component Response
// - status: Update with final token counts
// - toast: Show completion notification (optional)
// - Layout: Return to chat input state
```

---

## Message Types

Complete messages received from the SDK.

### AssistantMessage

Full assistant response with content blocks.

```go
// SDK Type
type AssistantMessage struct {
    ID      string
    Role    string
    Content []ContentBlock
    Model   string
    Usage   Usage
}

type ContentBlock interface {
    // TextBlock, ToolUseBlock, ThinkingBlock
}

// Adapter Message
type AssistantMsg struct {
    Message *claudecode.AssistantMessage
}

// Component Response
// - message: Create MessageBubble from content blocks
// - Each block type maps to a sub-component:
//   - TextBlock -> markdown rendered text
//   - ToolUseBlock -> tooluse card
//   - ThinkingBlock -> thinking block
```

### ResultMessage

Final result with usage statistics.

```go
// SDK Type
type ResultMessage struct {
    SessionID    string
    Model        string
    InputTokens  int
    OutputTokens int
    TotalCost    float64
    Duration     time.Duration
}

// Adapter Message
type ResultMsg struct {
    Message *claudecode.ResultMessage
}

// Component Response
// - status: Update all fields
// - toast: Show cost notification (optional)
```

---

## Callbacks

User input requests from the SDK.

### canUseTool

Permission callback for tool execution.

```go
// SDK Callback
type CanUseTool func(
    toolName string,
    input map[string]any,
    context ToolPermissionContext,
) (PermissionResult, error)

// Adapter Message
type ShowPermissionPromptMsg struct {
    ToolName  string
    ToolInput map[string]any
    ToolUseID string
    Respond   func(decision string, alwaysAllow bool)
}

// Component: permission
// User selects: Allow, Deny, or Always Allow
// Response sent back via Respond callback
```

#### Tool-Specific Input Schemas

| Tool | Input Fields |
|------|--------------|
| `Bash` | `command`, `description`, `timeout`, `run_in_background` |
| `Write` | `file_path`, `content` |
| `Edit` | `file_path`, `old_string`, `new_string` |
| `Read` | `file_path`, `offset`, `limit` |
| `Glob` | `pattern`, `path` |
| `Grep` | `pattern`, `path`, `type` |
| `WebSearch` | `query` |
| `WebFetch` | `url`, `prompt` |
| `Task` | `prompt`, `subagent_type`, `run_in_background` |
| `TaskOutput` | `task_id`, `block`, `timeout` |
| `KillShell` | `shell_id` |
| `Skill` | `skill`, `args` |

### AskUserQuestion (via canUseTool)

When `toolName == "AskUserQuestion"`, the input contains questions.

```go
// Input Structure
{
    "questions": [
        {
            "question": "Which database?",
            "header": "Database",
            "options": [
                {"label": "PostgreSQL", "description": "..."},
                {"label": "MySQL", "description": "..."}
            ],
            "multiSelect": false
        }
    ]
}

// Adapter Message
type ShowQuestionPromptMsg struct {
    Question    string
    Header      string
    Options     []QuestionOption
    MultiSelect bool
    Respond     func(selections []string, customText string)
}

// Component: question
// User selects option(s) or enters custom text
// Response: { "answers": { "Which database?": "PostgreSQL" } }
```

---

## Special Tools

Tools that trigger specific UI flows.

### EnterPlanMode

Claude requests to enter planning mode.

```go
// Detected via: toolName == "EnterPlanMode" in canUseTool

// Adapter Message
type ShowPlanEnterPromptMsg struct {
    Message string
    Respond func(confirmed bool)
}

// Component: planenter
// User confirms or skips planning
// If confirmed: Layout enters planning state
```

### ExitPlanMode

Claude presents plan for approval.

```go
// Detected via: toolName == "ExitPlanMode" in canUseTool

// Input Structure
{
    "allowedPrompts": [
        {"tool": "Bash", "prompt": "run tests"},
        {"tool": "Bash", "prompt": "build project"}
    ]
}

// Adapter Message
type ShowPlanExitPromptMsg struct {
    PlanContent string           // Read from plan file
    Permissions []RequestedPermission
    Respond     func(action, feedback string, approvedPerms []string)
}

// Component: planexit
// User can: Approve, Reject, or Request Changes
// Approved permissions are pre-granted for implementation
```

### TodoWrite

Claude updates the task list.

```go
// Detected via: toolName == "TodoWrite" in canUseTool or PostToolUse

// Input Structure
{
    "todos": [
        {"content": "Fix bug", "status": "completed", "activeForm": "Fixing bug"},
        {"content": "Add tests", "status": "in_progress", "activeForm": "Adding tests"}
    ]
}

// Adapter Message
type ShowTodoMsg struct {
    Items []TodoItem
}

// Component: todo
// Displays task list with progress
// Updates in real-time as Claude works
```

### TaskOutput

Retrieves output from background tasks.

```go
// Detected via: toolName == "TaskOutput" in canUseTool

// Input Structure
{
    "task_id": "bg-task-abc123",
    "block": true,       // Wait for completion
    "timeout": 30000     // Max wait time in ms
}

// Adapter Message
type TaskOutputRequestMsg struct {
    TaskID  string
    Block   bool
    Timeout int
}

type TaskOutputResultMsg struct {
    TaskID string
    Output string
    Status string  // "running", "completed", "failed"
    Done   bool
}

// Component: taskmgr
// Shows task output in expandable panel
// Auto-scrolls to latest output
```

### KillShell

Terminates a running background shell.

```go
// Detected via: toolName == "KillShell" in canUseTool

// Input Structure
{
    "shell_id": "bg-shell-xyz789"
}

// Adapter Message
type KillShellRequestMsg struct {
    ShellID string
    Respond func(confirmed bool)
}

type KillShellResultMsg struct {
    ShellID string
    Success bool
}

// Component: taskmgr (kill confirmation)
// Shows confirmation dialog before killing
// Updates task status after kill
```

### Skill

Invokes a registered skill (slash command).

```go
// Detected via: toolName == "Skill" in canUseTool

// Input Structure
{
    "skill": "commit",
    "args": "-m 'Fix bug'"
}

// Adapter Message
type SkillInvokeMsg struct {
    SkillName string
    Args      string
}

type SkillResultMsg struct {
    SkillName string
    Output    string
    Success   bool
}

// Component: tooluse (skill card)
// Shows skill name and arguments
// Displays skill output when complete
```

---

## Background Task Tools

Tools for managing background operations.

### run_in_background Parameter

Both `Bash` and `Task` tools support background execution.

```go
// Bash with run_in_background
{
    "command": "npm run build",
    "run_in_background": true
}

// Task with run_in_background
{
    "prompt": "Run tests",
    "subagent_type": "test-runner",
    "run_in_background": true
}

// Adapter Messages
type TaskStartedMsg struct {
    TaskID      string
    Type        string   // "bash" or "agent"
    Description string
    OutputFile  string   // Path to output file
}

type TaskProgressMsg struct {
    TaskID string
    Output string
}

type TaskCompletedMsg struct {
    TaskID string
    Status string   // "completed", "failed", "killed"
    Output string
}

// Component: taskstatus, taskmgr
// taskstatus: Shows running task count in status bar
// taskmgr: Full task management panel
```

### Task Flow

```
1. Background task started
   └─> TaskStartedMsg{TaskID: "bg-123", Type: "bash"}
   └─> taskstatus: Show spinner, increment count

2. Task produces output
   └─> TaskProgressMsg{TaskID: "bg-123", Output: "..."}
   └─> taskmgr: Append to output view

3. Task completes
   └─> TaskCompletedMsg{TaskID: "bg-123", Status: "completed"}
   └─> taskstatus: Decrement count
   └─> taskmgr: Update status indicator
```

---

## Built-in Tools Display

How each SDK tool is displayed in the TUI.

| Tool | Display Component | Visual |
|------|-------------------|--------|
| **Read** | `tooluse` | File icon + path, line range if specified |
| **Write** | `tooluse` | File icon + path, content preview (collapsed) |
| **Edit** | `tooluse`, `diffview` | Diff view: old (red) -> new (green), side-by-side option |
| **Bash** | `tooluse` | Terminal icon, command, description |
| **Glob** | `tooluse`, `globresults` | Folder icon, pattern, matched files tree/list |
| **Grep** | `tooluse`, `grepresults` | Search icon, pattern, match highlights with context |
| **WebSearch** | `tooluse` | Globe icon, query, results list |
| **WebFetch** | `tooluse` | Link icon, URL, content summary |
| **Task** | `message`, `taskmgr` | Nested agent indicator; background tasks in task manager |
| **TaskOutput** | `taskmgr` | Task output viewer, scrollable content |
| **KillShell** | `taskmgr` | Kill confirmation dialog, task status update |
| **Skill** | `tooluse` | Slash command icon, skill name, arguments |
| **TodoWrite** | `todo` | Task list with checkboxes |
| **NotebookEdit** | `tooluse` | Notebook icon, cell preview |

---

## Error Handling

SDK errors mapped to TUI responses.

```go
// SDK Errors
type CLINotFoundError struct { ... }
type ProcessError struct { ... }
type CLIJSONDecodeError struct { ... }

// Adapter Message
type StreamErrorMsg struct {
    Err error
}

// Component Response
// - toast: Show error notification (Level: Error)
// - status: Clear streaming state
// - Layout: Enable chat input, show error state
```

---

## Lifecycle Events

Client lifecycle mapped to TUI states.

```go
// Client States
type ClientState int

const (
    ClientStateDisconnected ClientState = iota
    ClientStateConnecting
    ClientStateConnected
    ClientStateStreaming
    ClientStateError
)

// Adapter Message
type ClientStateMsg struct {
    State ClientState
    Err   error
}

// Layout Response
// - Disconnected: Show reconnect option
// - Connecting: Show connection spinner
// - Connected: Enable chat input
// - Streaming: Show streaming UI, disable input
// - Error: Show error, offer retry
```

---

## Permission Modes

SDK permission modes and their TUI behavior.

| Mode | SDK Constant | TUI Behavior |
|------|--------------|--------------|
| Default | `PermissionModeDefault` | Show permission prompt for all tools |
| Accept Edits | `PermissionModeAcceptEdits` | Auto-approve Read, Write, Edit |
| Plan | `PermissionModePlan` | Read-only streaming, wait for plan approval |
| Bypass | `PermissionModeBypass` | No prompts (dangerous, use carefully) |

```go
// ChatInput toggle -> Permission mode
if autoAccept && planMode {
    // Not a valid combination, prefer plan mode
    mode = PermissionModePlan
} else if autoAccept {
    mode = PermissionModeAcceptEdits
} else if planMode {
    mode = PermissionModePlan
} else {
    mode = PermissionModeDefault
}
```

---

## Message Flow Example

Complete flow from user input to tool execution:

```
1. User types in ChatInput, presses Enter
   └─> chatinput.SubmitMsg{Text: "fix the bug", AutoAccept: false}

2. Layout calls SDK Query
   └─> client.Query(ctx, "fix the bug")

3. SDK streams response
   └─> StreamDeltaMsg{Text: "I'll read the file..."}
   └─> StreamText appends, shows cursor

4. SDK requests tool use (Read)
   └─> canUseTool("Read", {file_path: "main.go"})
   └─> ShowPermissionPromptMsg displayed

5. User approves in Permission component
   └─> permission.ResultMsg{Decision: "allow"}
   └─> Respond callback sends approval

6. SDK executes tool, streams result
   └─> ToolUse card shows file content

7. SDK streams more text
   └─> StreamDeltaMsg{Text: "I found the issue..."}

8. SDK requests Edit
   └─> canUseTool("Edit", {file_path: "main.go", ...})
   └─> User approves

9. SDK completes
   └─> ResultMsg with token counts
   └─> Status bar updates
   └─> Layout returns to ChatInput state
```
