# Adapter Layer

The adapter layer bridges `claude-agent-sdk-go` to Bubble Tea's message-passing architecture.

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Bubble Tea Program                        │
│                                                              │
│    Model.Update(msg tea.Msg) (tea.Model, tea.Cmd)           │
│         ▲                                                    │
│         │ tea.Msg                                            │
│         │                                                    │
├─────────┴────────────────────────────────────────────────────┤
│                     Adapter Layer                            │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │  stream.go   │  │  control.go  │  │  client.go   │       │
│  │              │  │              │  │              │       │
│  │ StreamCmd()  │  │ canUseTool   │  │ ClientCmd()  │       │
│  │ adaptMsg()   │  │ ToolAdapter  │  │ Lifecycle    │       │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘       │
│         │                 │                 │                │
└─────────┼─────────────────┼─────────────────┼────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────┐
│                   claude-agent-sdk-go                        │
│                                                              │
│  ReceiveMessages()    canUseTool()        Client API        │
│  StreamEvent          ToolUseRequest      Connect()         │
│  AssistantMessage     ToolResult          Query()           │
│  ResultMessage        Callbacks           Disconnect()      │
└─────────────────────────────────────────────────────────────┘
```

---

## stream.go - Message Stream Adapter

Converts SDK streaming messages to Bubble Tea messages.

### Output Message Types

```go
package adapter

import (
    tea "github.com/charmbracelet/bubbletea"
    claudecode "github.com/severity1/claude-agent-sdk-go"
)

// StreamDeltaMsg contains incremental text from streaming
type StreamDeltaMsg struct {
    Text      string
    BlockType string  // "text", "thinking", "tool_use"
    Index     int
}

// StreamBlockStartMsg signals a new content block
type StreamBlockStartMsg struct {
    BlockType string  // "text", "thinking", "tool_use"
    Index     int
    ToolName  string  // Only for tool_use blocks
    ToolID    string  // Only for tool_use blocks
}

// StreamBlockStopMsg signals content block completion
type StreamBlockStopMsg struct {
    Index int
}

// AssistantMsg contains a complete assistant message
type AssistantMsg struct {
    Message *claudecode.AssistantMessage
}

// ResultMsg contains the final result with usage stats
type ResultMsg struct {
    Message      *claudecode.ResultMessage
    InputTokens  int
    OutputTokens int
    TotalCost    float64
}

// StreamDoneMsg signals stream completion
type StreamDoneMsg struct{}

// StreamErrorMsg signals a stream error
type StreamErrorMsg struct {
    Err error
}

// ThinkingDeltaMsg contains incremental thinking text
type ThinkingDeltaMsg struct {
    Text  string
    Index int
}

// ToolUseDeltaMsg contains incremental tool input JSON
// Note: ToolName/ToolID come from StreamBlockStartMsg, not delta events.
// The TUI layer tracks state to associate tool info with deltas by index.
type ToolUseDeltaMsg struct {
    PartialJSON string
    Index       int
}

// UserMsg contains a user message (typically tool results)
type UserMsg struct {
    UUID            string
    Content         any    // string or content blocks
    ParentToolUseID string
}

// SystemInitMsg contains session initialization data
type SystemInitMsg struct {
    SessionID      string
    Model          string
    Cwd            string
    Tools          []string
    PermissionMode string
}

// SystemHookResponseMsg contains hook execution results
type SystemHookResponseMsg struct {
    SessionID string
    HookName  string
    HookEvent string
    Stdout    string
    Stderr    string
    ExitCode  int
}

// MessageStartMsg contains initial message metadata with model info
type MessageStartMsg struct {
    MessageID    string
    Model        string
    InputTokens  int
    OutputTokens int
}

// MessageDeltaMsg contains final usage stats and stop reason
type MessageDeltaMsg struct {
    StopReason   string
    InputTokens  int
    OutputTokens int
}

// ControlRequestMsg contains control protocol request (informational)
type ControlRequestMsg struct {
    RequestID string
    Subtype   string
}

// ControlResponseMsg contains control protocol response (informational)
type ControlResponseMsg struct {
    RequestID string
    IsSuccess bool
    Error     string
}
```

### Stream Command

```go
// StreamCmd returns a tea.Cmd that reads from SDK message channel.
// Call this in a loop to continuously receive messages.
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

// adaptMessage converts SDK message to appropriate tea.Msg
func adaptMessage(msg claudecode.Message) tea.Msg {
    switch m := msg.(type) {
    case *claudecode.StreamEvent:
        return adaptStreamEvent(m)
    case *claudecode.AssistantMessage:
        return AssistantMsg{Message: m}
    case *claudecode.ResultMessage:
        return ResultMsg{
            Message:      m,
            InputTokens:  m.Usage.InputTokens,
            OutputTokens: m.Usage.OutputTokens,
            TotalCost:    m.CalculateCost(),
        }
    default:
        return nil
    }
}

// adaptStreamEvent extracts typed content from StreamEvent
func adaptStreamEvent(event *claudecode.StreamEvent) tea.Msg {
    eventType, _ := event.Event["type"].(string)

    switch eventType {
    case "content_block_start":
        return adaptBlockStart(event)

    case "content_block_delta":
        return adaptBlockDelta(event)

    case "content_block_stop":
        index, _ := event.Event["index"].(float64)
        return StreamBlockStopMsg{Index: int(index)}

    case "message_stop":
        return StreamDoneMsg{}

    default:
        return nil
    }
}

func adaptBlockStart(event *claudecode.StreamEvent) tea.Msg {
    index, _ := event.Event["index"].(float64)
    cb, _ := event.Event["content_block"].(map[string]any)
    blockType, _ := cb["type"].(string)

    msg := StreamBlockStartMsg{
        BlockType: blockType,
        Index:     int(index),
    }

    // Extract tool info for tool_use blocks
    if blockType == "tool_use" {
        msg.ToolName, _ = cb["name"].(string)
        msg.ToolID, _ = cb["id"].(string)
    }

    return msg
}

func adaptBlockDelta(event *claudecode.StreamEvent) tea.Msg {
    index, _ := event.Event["index"].(float64)
    delta, _ := event.Event["delta"].(map[string]any)
    deltaType, _ := delta["type"].(string)

    switch deltaType {
    case "text_delta":
        text, _ := delta["text"].(string)
        return StreamDeltaMsg{
            Text:      text,
            BlockType: "text",
            Index:     int(index),
        }

    case "thinking_delta":
        text, _ := delta["thinking"].(string)
        return ThinkingDeltaMsg{
            Text:  text,
            Index: int(index),
        }

    case "input_json_delta":
        partialJSON, _ := delta["partial_json"].(string)
        return ToolUseDeltaMsg{
            PartialJSON: partialJSON,
            Index:       int(index),
        }

    default:
        return nil
    }
}
```

### Usage Example

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.StreamDeltaMsg:
        // Route to appropriate component based on block type
        m.streamText.AppendText(msg.Text)
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.ThinkingDeltaMsg:
        m.thinking.AppendText(msg.Text)
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.ToolUseDeltaMsg:
        m.toolUse.AppendJSON(msg.PartialJSON)
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.StreamBlockStartMsg:
        switch msg.BlockType {
        case "thinking":
            m.thinking.Start()
        case "tool_use":
            m.toolUse.Start(msg.ToolName, msg.ToolID)
        }
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.ResultMsg:
        m.status.Update(msg.InputTokens, msg.OutputTokens, msg.TotalCost)
        m.streaming = false
        return m, nil

    case adapter.StreamErrorMsg:
        return m, m.toast.Show(toast.Error, "Stream Error", msg.Err.Error())
    }
    return m, nil
}
```

---

## control.go - Tool Control Adapter

The `canUseTool` callback is the primary mechanism for tool permission handling in the SDK. It intercepts tool execution requests and allows the TUI to prompt the user.

### Tool Request Types

```go
package adapter

// ToolUseRequest represents an intercepted tool execution
type ToolUseRequest struct {
    ToolName  string
    ToolInput map[string]any
    ToolUseID string
}

// ToolDecision represents the user's response
type ToolDecision struct {
    Allow       bool
    AlwaysAllow bool
    Reason      string
}
```

### TUI Message Types

```go
// ShowPermissionPromptMsg triggers permission input component
type ShowPermissionPromptMsg struct {
    Request   ToolUseRequest
    Respond   func(decision ToolDecision)
}

// ShowQuestionPromptMsg triggers question input component
type ShowQuestionPromptMsg struct {
    ToolUseID   string
    Questions   []Question
    Respond     func(answers map[string]any)
}

type Question struct {
    Question    string
    Header      string
    Options     []QuestionOption
    MultiSelect bool
}

type QuestionOption struct {
    Label       string
    Description string
}

// ShowPlanEnterPromptMsg triggers plan entry confirmation
type ShowPlanEnterPromptMsg struct {
    ToolUseID string
    Respond   func(confirmed bool)
}

// ShowPlanExitPromptMsg triggers plan review component
type ShowPlanExitPromptMsg struct {
    ToolUseID   string
    PlanFile    string
    Permissions []RequestedPermission
    Respond     func(action PlanAction, feedback string, approvedPerms []string)
}

type RequestedPermission struct {
    Tool   string
    Prompt string
}

type PlanAction int

const (
    PlanApprove PlanAction = iota
    PlanReject
    PlanRequestChanges
)

// ShowTodoUpdateMsg notifies of TodoWrite tool updates
type ShowTodoUpdateMsg struct {
    ToolUseID string
    Todos     []TodoItem
}

type TodoItem struct {
    Content    string
    Status     string  // "pending", "in_progress", "completed"
    ActiveForm string
}

// Background Task Messages

// TaskStartedMsg signals a new background task has started
type TaskStartedMsg struct {
    TaskID      string
    Type        string       // "bash" or "agent"
    Description string
    OutputFile  string       // Path to output file for incremental reading
}

// TaskProgressMsg contains incremental output from a background task
type TaskProgressMsg struct {
    TaskID string
    Output string
}

// TaskCompletedMsg signals a background task has finished
type TaskCompletedMsg struct {
    TaskID   string
    Status   string    // "completed", "failed", "killed"
    Output   string
    ExitCode *int      // nil if killed or agent type
}

// ShowTaskManagerMsg triggers the task manager panel
type ShowTaskManagerMsg struct {
    Tasks []BackgroundTask
}

type BackgroundTask struct {
    ID          string
    Type        string       // "bash" or "agent"
    Description string
    Status      string       // "running", "completed", "failed", "killed"
    StartedAt   time.Time
    Output      string
    OutputFile  string
    ExitCode    *int
}

// TaskOutputRequestMsg triggers a task output fetch
type ShowTaskOutputRequestMsg struct {
    ToolUseID string
    TaskID    string
    Block     bool
    Timeout   int
    Respond   func(output string, done bool)
}

// KillShellRequestMsg triggers kill confirmation
type ShowKillShellRequestMsg struct {
    ToolUseID string
    ShellID   string
    Respond   func(confirmed bool)
}

// Skill Messages

// ShowSkillInvokeMsg signals a skill is being invoked
type ShowSkillInvokeMsg struct {
    ToolUseID string
    SkillName string
    Args      string
}

// Session Messages

// ShowSessionPickerMsg triggers the session picker
type ShowSessionPickerMsg struct {
    Sessions []SessionInfo
    Current  string
    Respond  func(action SessionAction, sessionID string)
}

type SessionInfo struct {
    ID           string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    MessageCount int
    TokensUsed   int
    Model        string
    Summary      string
    IsActive     bool
}

type SessionAction int

const (
    SessionActionSelect SessionAction = iota
    SessionActionResume
    SessionActionDelete
    SessionActionExport
)

// SessionSelectedMsg signals a session selection
type SessionSelectedMsg struct {
    SessionID string
    Action    SessionAction
}

// SessionResumedMsg signals a session has been resumed
type SessionResumedMsg struct {
    SessionID string
}

// Search Messages

// ShowSearchMsg triggers the search overlay
type ShowSearchMsg struct {
    Respond func(result SearchResult)
}

// SearchResultsMsg contains search results
type SearchResultsMsg struct {
    Query   string
    Results []SearchResult
}

type SearchResult struct {
    MessageIndex int
    Role         string
    Snippet      string
    Highlights   []HighlightRange
    Timestamp    time.Time
}

type HighlightRange struct {
    Start int
    End   int
}
```

### canUseTool Callback Adapter

The `canUseTool` callback provides fine-grained control over tool execution.

**Critical:** The callback runs in the SDK's goroutine, not the Bubble Tea goroutine.
This requires careful handling of:
1. Context cancellation - user closes app while waiting
2. Program shutdown - `tea.Program` may be nil or closed
3. Deadlock prevention - `program.Send()` might block if program is busy

```go
// ToolControlAdapter creates canUseTool callback for SDK
type ToolControlAdapter struct {
    program      *tea.Program
    autoApproved map[string]bool  // Tools auto-approved by user
    shutdown     chan struct{}    // Signals shutdown in progress
    mu           sync.Mutex
}

func NewToolControlAdapter(p *tea.Program) *ToolControlAdapter {
    return &ToolControlAdapter{
        program:      p,
        autoApproved: make(map[string]bool),
        shutdown:     make(chan struct{}),
    }
}

// Shutdown signals all pending callbacks to abort
// Call this before tea.Program.Quit() to prevent deadlocks
func (t *ToolControlAdapter) Shutdown() {
    close(t.shutdown)
}

// CanUseTool returns the callback function for SDK configuration
func (t *ToolControlAdapter) CanUseTool() claudecode.CanUseToolFunc {
    return func(
        ctx context.Context,
        toolName string,
        toolInput map[string]any,
        toolUseID string,
    ) (bool, error) {
        // Handle special tools that need UI interaction
        switch toolName {
        case "AskUserQuestion":
            return t.handleAskUserQuestion(ctx, toolInput, toolUseID)

        case "EnterPlanMode":
            return t.handleEnterPlanMode(ctx, toolUseID)

        case "ExitPlanMode":
            return t.handleExitPlanMode(ctx, toolInput, toolUseID)

        case "TodoWrite":
            return t.handleTodoWrite(ctx, toolInput, toolUseID)

        case "TaskOutput":
            return t.handleTaskOutput(ctx, toolInput, toolUseID)

        case "KillShell":
            return t.handleKillShell(ctx, toolInput, toolUseID)

        case "Skill":
            return t.handleSkill(ctx, toolInput, toolUseID)
        }

        // Check if tool is auto-approved
        t.mu.Lock()
        if t.autoApproved[toolName] {
            t.mu.Unlock()
            return true, nil
        }
        t.mu.Unlock()

        // Show permission prompt for other tools
        return t.showPermissionPrompt(ctx, toolName, toolInput, toolUseID)
    }
}

// showPermissionPrompt displays permission UI and waits for response
func (t *ToolControlAdapter) showPermissionPrompt(
    ctx context.Context,
    toolName string,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    // Check for shutdown before sending
    select {
    case <-t.shutdown:
        return false, fmt.Errorf("adapter shutdown")
    default:
    }

    resultCh := make(chan ToolDecision, 1)

    // Non-blocking send to prevent deadlock if program is busy
    // The buffered channel in tea.Program handles this, but be defensive
    t.program.Send(ShowPermissionPromptMsg{
        Request: ToolUseRequest{
            ToolName:  toolName,
            ToolInput: toolInput,
            ToolUseID: toolUseID,
        },
        Respond: func(decision ToolDecision) {
            // Respond callback may be called after shutdown
            select {
            case resultCh <- decision:
            default:
                // Channel closed or full, ignore
            }
        },
    })

    select {
    case decision := <-resultCh:
        if decision.AlwaysAllow {
            t.mu.Lock()
            t.autoApproved[toolName] = true
            t.mu.Unlock()
        }
        return decision.Allow, nil
    case <-ctx.Done():
        return false, ctx.Err()
    case <-t.shutdown:
        return false, fmt.Errorf("adapter shutdown")
    }
}
```

### Special Tool Handlers

```go
// handleAskUserQuestion shows question UI
func (t *ToolControlAdapter) handleAskUserQuestion(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    questionsRaw, _ := toolInput["questions"].([]any)
    if len(questionsRaw) == 0 {
        return true, nil
    }

    questions := make([]Question, 0, len(questionsRaw))
    for _, q := range questionsRaw {
        qMap := q.(map[string]any)
        question := Question{
            Question:    qMap["question"].(string),
            Header:      qMap["header"].(string),
            MultiSelect: qMap["multiSelect"].(bool),
        }

        if opts, ok := qMap["options"].([]any); ok {
            for _, opt := range opts {
                optMap := opt.(map[string]any)
                question.Options = append(question.Options, QuestionOption{
                    Label:       optMap["label"].(string),
                    Description: optMap["description"].(string),
                })
            }
        }
        questions = append(questions, question)
    }

    resultCh := make(chan map[string]any, 1)

    t.program.Send(ShowQuestionPromptMsg{
        ToolUseID: toolUseID,
        Questions: questions,
        Respond: func(answers map[string]any) {
            resultCh <- answers
        },
    })

    select {
    case <-resultCh:
        // Answers are handled by the SDK through tool result
        return true, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// handleEnterPlanMode shows plan entry confirmation
func (t *ToolControlAdapter) handleEnterPlanMode(
    ctx context.Context,
    toolUseID string,
) (bool, error) {
    resultCh := make(chan bool, 1)

    t.program.Send(ShowPlanEnterPromptMsg{
        ToolUseID: toolUseID,
        Respond: func(confirmed bool) {
            resultCh <- confirmed
        },
    })

    select {
    case confirmed := <-resultCh:
        return confirmed, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// handleExitPlanMode shows plan review UI
func (t *ToolControlAdapter) handleExitPlanMode(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    // Extract plan file path
    planFile, _ := toolInput["planFile"].(string)

    // Extract requested permissions
    var permissions []RequestedPermission
    if perms, ok := toolInput["allowedPrompts"].([]any); ok {
        for _, p := range perms {
            perm := p.(map[string]any)
            permissions = append(permissions, RequestedPermission{
                Tool:   perm["tool"].(string),
                Prompt: perm["prompt"].(string),
            })
        }
    }

    resultCh := make(chan struct {
        action        PlanAction
        feedback      string
        approvedPerms []string
    }, 1)

    t.program.Send(ShowPlanExitPromptMsg{
        ToolUseID:   toolUseID,
        PlanFile:    planFile,
        Permissions: permissions,
        Respond: func(action PlanAction, feedback string, approvedPerms []string) {
            resultCh <- struct {
                action        PlanAction
                feedback      string
                approvedPerms []string
            }{action, feedback, approvedPerms}
        },
    })

    select {
    case result := <-resultCh:
        return result.action == PlanApprove, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// handleTodoWrite updates todo display (always allowed)
func (t *ToolControlAdapter) handleTodoWrite(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    todosRaw, _ := toolInput["todos"].([]any)
    todos := make([]TodoItem, 0, len(todosRaw))

    for _, todo := range todosRaw {
        todoMap := todo.(map[string]any)
        todos = append(todos, TodoItem{
            Content:    todoMap["content"].(string),
            Status:     todoMap["status"].(string),
            ActiveForm: todoMap["activeForm"].(string),
        })
    }

    // Send update to UI (non-blocking)
    t.program.Send(ShowTodoUpdateMsg{
        ToolUseID: toolUseID,
        Todos:     todos,
    })

    // TodoWrite is always allowed
    return true, nil
}

// handleTaskOutput retrieves output from background task
func (t *ToolControlAdapter) handleTaskOutput(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    taskID, _ := toolInput["task_id"].(string)
    block, _ := toolInput["block"].(bool)
    timeout := 30000 // default
    if to, ok := toolInput["timeout"].(float64); ok {
        timeout = int(to)
    }

    resultCh := make(chan struct {
        output string
        done   bool
    }, 1)

    t.program.Send(ShowTaskOutputRequestMsg{
        ToolUseID: toolUseID,
        TaskID:    taskID,
        Block:     block,
        Timeout:   timeout,
        Respond: func(output string, done bool) {
            resultCh <- struct {
                output string
                done   bool
            }{output, done}
        },
    })

    select {
    case <-resultCh:
        // Output retrieved successfully
        return true, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// handleKillShell shows kill confirmation before terminating shell
func (t *ToolControlAdapter) handleKillShell(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    shellID, _ := toolInput["shell_id"].(string)

    resultCh := make(chan bool, 1)

    t.program.Send(ShowKillShellRequestMsg{
        ToolUseID: toolUseID,
        ShellID:   shellID,
        Respond: func(confirmed bool) {
            resultCh <- confirmed
        },
    })

    select {
    case confirmed := <-resultCh:
        return confirmed, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// handleSkill notifies UI of skill invocation (always allowed)
func (t *ToolControlAdapter) handleSkill(
    ctx context.Context,
    toolInput map[string]any,
    toolUseID string,
) (bool, error) {
    skillName, _ := toolInput["skill"].(string)
    args, _ := toolInput["args"].(string)

    // Notify UI of skill invocation (non-blocking)
    t.program.Send(ShowSkillInvokeMsg{
        ToolUseID: toolUseID,
        SkillName: skillName,
        Args:      args,
    })

    // Skill invocations are always allowed
    return true, nil
}
```

### Tool-Specific Input Parsing

Helper functions for parsing tool inputs by tool type.

```go
// ParseBashInput extracts Bash tool parameters
func ParseBashInput(input map[string]any) (command, description string, timeout int) {
    command, _ = input["command"].(string)
    description, _ = input["description"].(string)
    if t, ok := input["timeout"].(float64); ok {
        timeout = int(t)
    }
    return
}

// ParseWriteInput extracts Write tool parameters
func ParseWriteInput(input map[string]any) (filePath, content string) {
    filePath, _ = input["file_path"].(string)
    content, _ = input["content"].(string)
    return
}

// ParseEditInput extracts Edit tool parameters
func ParseEditInput(input map[string]any) (filePath, oldString, newString string, replaceAll bool) {
    filePath, _ = input["file_path"].(string)
    oldString, _ = input["old_string"].(string)
    newString, _ = input["new_string"].(string)
    replaceAll, _ = input["replace_all"].(bool)
    return
}

// ParseReadInput extracts Read tool parameters
func ParseReadInput(input map[string]any) (filePath string, offset, limit int) {
    filePath, _ = input["file_path"].(string)
    if o, ok := input["offset"].(float64); ok {
        offset = int(o)
    }
    if l, ok := input["limit"].(float64); ok {
        limit = int(l)
    }
    return
}

// ParseGlobInput extracts Glob tool parameters
func ParseGlobInput(input map[string]any) (pattern, path string) {
    pattern, _ = input["pattern"].(string)
    path, _ = input["path"].(string)
    return
}

// ParseGrepInput extracts Grep tool parameters
func ParseGrepInput(input map[string]any) (pattern, path, glob, outputMode string) {
    pattern, _ = input["pattern"].(string)
    path, _ = input["path"].(string)
    glob, _ = input["glob"].(string)
    outputMode, _ = input["output_mode"].(string)
    return
}

// ParseWebFetchInput extracts WebFetch tool parameters
func ParseWebFetchInput(input map[string]any) (url, prompt string) {
    url, _ = input["url"].(string)
    prompt, _ = input["prompt"].(string)
    return
}

// ParseWebSearchInput extracts WebSearch tool parameters
func ParseWebSearchInput(input map[string]any) (query string, allowedDomains, blockedDomains []string) {
    query, _ = input["query"].(string)
    if domains, ok := input["allowed_domains"].([]any); ok {
        for _, d := range domains {
            allowedDomains = append(allowedDomains, d.(string))
        }
    }
    if domains, ok := input["blocked_domains"].([]any); ok {
        for _, d := range domains {
            blockedDomains = append(blockedDomains, d.(string))
        }
    }
    return
}

// ParseTaskInput extracts Task (subagent) tool parameters
func ParseTaskInput(input map[string]any) (description, prompt, subagentType string) {
    description, _ = input["description"].(string)
    prompt, _ = input["prompt"].(string)
    subagentType, _ = input["subagent_type"].(string)
    return
}

// ParseNotebookEditInput extracts NotebookEdit tool parameters
func ParseNotebookEditInput(input map[string]any) (notebookPath, newSource, cellType, editMode string) {
    notebookPath, _ = input["notebook_path"].(string)
    newSource, _ = input["new_source"].(string)
    cellType, _ = input["cell_type"].(string)
    editMode, _ = input["edit_mode"].(string)
    return
}

// ParseTaskOutputInput extracts TaskOutput tool parameters
func ParseTaskOutputInput(input map[string]any) (taskID string, block bool, timeout int) {
    taskID, _ = input["task_id"].(string)
    block, _ = input["block"].(bool)
    if t, ok := input["timeout"].(float64); ok {
        timeout = int(t)
    } else {
        timeout = 30000 // default
    }
    return
}

// ParseKillShellInput extracts KillShell tool parameters
func ParseKillShellInput(input map[string]any) (shellID string) {
    shellID, _ = input["shell_id"].(string)
    return
}

// ParseSkillInput extracts Skill tool parameters
func ParseSkillInput(input map[string]any) (skill, args string) {
    skill, _ = input["skill"].(string)
    args, _ = input["args"].(string)
    return
}

// ParseBashBackgroundInput extracts Bash tool parameters including background flag
func ParseBashBackgroundInput(input map[string]any) (command, description string, timeout int, runInBackground bool) {
    command, _ = input["command"].(string)
    description, _ = input["description"].(string)
    if t, ok := input["timeout"].(float64); ok {
        timeout = int(t)
    }
    runInBackground, _ = input["run_in_background"].(bool)
    return
}

// ParseTaskBackgroundInput extracts Task (subagent) tool parameters including background flag
func ParseTaskBackgroundInput(input map[string]any) (description, prompt, subagentType string, runInBackground bool) {
    description, _ = input["description"].(string)
    prompt, _ = input["prompt"].(string)
    subagentType, _ = input["subagent_type"].(string)
    runInBackground, _ = input["run_in_background"].(bool)
    return
}
```

---

## client.go - Client Lifecycle Adapter

Manages SDK client lifecycle within Bubble Tea.

### Client State

```go
package adapter

// ClientState represents the SDK client connection state
type ClientState int

const (
    ClientStateDisconnected ClientState = iota
    ClientStateConnecting
    ClientStateConnected
    ClientStateStreaming
    ClientStateWaitingInput    // Waiting for user input (permission, question, etc.)
    ClientStatePlanMode        // In plan mode
    ClientStateError
)

func (s ClientState) String() string {
    switch s {
    case ClientStateDisconnected:
        return "disconnected"
    case ClientStateConnecting:
        return "connecting"
    case ClientStateConnected:
        return "connected"
    case ClientStateStreaming:
        return "streaming"
    case ClientStateWaitingInput:
        return "waiting"
    case ClientStatePlanMode:
        return "planning"
    case ClientStateError:
        return "error"
    default:
        return "unknown"
    }
}

// ClientStateMsg reports client state changes
type ClientStateMsg struct {
    State ClientState
    Err   error
}
```

### Client Adapter

```go
// ClientAdapter manages SDK client lifecycle
type ClientAdapter struct {
    client       claudecode.Client
    toolControl  *ToolControlAdapter
    ctx          context.Context
    cancel       context.CancelFunc
    state        ClientState
    options      []claudecode.Option

    // Session info
    sessionID    string
    model        string
    planMode     bool
    autoAccept   bool

    mu           sync.RWMutex
}

// NewClientAdapter creates a new client adapter
func NewClientAdapter(program *tea.Program, opts ...claudecode.Option) *ClientAdapter {
    toolControl := NewToolControlAdapter(program)

    return &ClientAdapter{
        toolControl: toolControl,
        options:     opts,
        state:       ClientStateDisconnected,
    }
}

// ConnectCmd returns a tea.Cmd that establishes SDK connection
func (c *ClientAdapter) ConnectCmd() tea.Cmd {
    return func() tea.Msg {
        c.mu.Lock()
        c.ctx, c.cancel = context.WithCancel(context.Background())

        // Inject canUseTool callback
        allOpts := append(c.options, claudecode.WithCanUseTool(c.toolControl.CanUseTool()))

        c.client = claudecode.NewClient(allOpts...)
        c.mu.Unlock()

        if err := c.client.Connect(c.ctx); err != nil {
            c.setState(ClientStateError)
            return ClientStateMsg{State: ClientStateError, Err: err}
        }

        c.setState(ClientStateConnected)
        return ClientStateMsg{State: ClientStateConnected}
    }
}

// QueryCmd returns a tea.Cmd that sends a query
func (c *ClientAdapter) QueryCmd(prompt string) tea.Cmd {
    return func() tea.Msg {
        c.mu.RLock()
        client := c.client
        ctx := c.ctx
        c.mu.RUnlock()

        if client == nil {
            return ClientStateMsg{
                State: ClientStateError,
                Err:   fmt.Errorf("client not connected"),
            }
        }

        if err := client.Query(ctx, prompt); err != nil {
            return ClientStateMsg{State: ClientStateError, Err: err}
        }

        c.setState(ClientStateStreaming)
        return ClientStateMsg{State: ClientStateStreaming}
    }
}

// MessageChannel returns the SDK message channel for streaming
func (c *ClientAdapter) MessageChannel() <-chan claudecode.Message {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if c.client == nil {
        return nil
    }
    return c.client.ReceiveMessages(c.ctx)
}

// DisconnectCmd returns a tea.Cmd that closes the connection
func (c *ClientAdapter) DisconnectCmd() tea.Cmd {
    return func() tea.Msg {
        c.mu.Lock()
        defer c.mu.Unlock()

        if c.cancel != nil {
            c.cancel()
        }
        if c.client != nil {
            c.client.Disconnect(c.ctx)
        }
        c.state = ClientStateDisconnected
        return ClientStateMsg{State: ClientStateDisconnected}
    }
}

// InterruptCmd returns a tea.Cmd that interrupts current operation
func (c *ClientAdapter) InterruptCmd() tea.Cmd {
    return func() tea.Msg {
        c.mu.RLock()
        client := c.client
        ctx := c.ctx
        c.mu.RUnlock()

        if client != nil {
            client.Interrupt(ctx)
        }
        return nil
    }
}

// State returns current client state
func (c *ClientAdapter) State() ClientState {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.state
}

func (c *ClientAdapter) setState(state ClientState) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.state = state
}

// SetPlanMode updates plan mode state
func (c *ClientAdapter) SetPlanMode(enabled bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.planMode = enabled
    if enabled {
        c.state = ClientStatePlanMode
    }
}

// SetAutoAccept updates auto-accept state
func (c *ClientAdapter) SetAutoAccept(enabled bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.autoAccept = enabled
}

// AutoApprove adds a tool to auto-approved list
func (c *ClientAdapter) AutoApprove(toolName string) {
    c.toolControl.mu.Lock()
    defer c.toolControl.mu.Unlock()
    c.toolControl.autoApproved[toolName] = true
}

// ClearAutoApprovals removes all auto-approvals
func (c *ClientAdapter) ClearAutoApprovals() {
    c.toolControl.mu.Lock()
    defer c.toolControl.mu.Unlock()
    c.toolControl.autoApproved = make(map[string]bool)
}
```

### Usage Example

```go
type Model struct {
    client   *adapter.ClientAdapter
    msgChan  <-chan claudecode.Message
    ctx      context.Context
    // components...
}

func NewModel() Model {
    return Model{
        ctx: context.Background(),
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m *Model) Start(p *tea.Program) tea.Cmd {
    m.client = adapter.NewClientAdapter(p,
        claudecode.WithPartialStreaming(),
        claudecode.WithMaxTurns(30),
    )
    return m.client.ConnectCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case adapter.ClientStateMsg:
        switch msg.State {
        case adapter.ClientStateConnected:
            m.connected = true

        case adapter.ClientStateStreaming:
            m.msgChan = m.client.MessageChannel()
            return m, adapter.StreamCmd(m.ctx, m.msgChan)

        case adapter.ClientStateError:
            return m, m.toast.Show(toast.Error, "Error", msg.Err.Error())
        }

    case adapter.StreamDeltaMsg:
        m.streamText.AppendText(msg.Text)
        return m, adapter.StreamCmd(m.ctx, m.msgChan)

    case adapter.StreamDoneMsg:
        m.streaming = false
        m.client.setState(adapter.ClientStateConnected)
        return m, nil

    // Handle permission prompts
    case adapter.ShowPermissionPromptMsg:
        m.permission.Show(msg.Request, msg.Respond)
        return m, nil

    // Handle question prompts
    case adapter.ShowQuestionPromptMsg:
        m.question.Show(msg.Questions, msg.Respond)
        return m, nil

    // Handle plan mode
    case adapter.ShowPlanEnterPromptMsg:
        m.planEnter.Show(msg.Respond)
        return m, nil

    case adapter.ShowPlanExitPromptMsg:
        m.planExit.Show(msg.PlanFile, msg.Permissions, msg.Respond)
        return m, nil

    // Handle todo updates
    case adapter.ShowTodoUpdateMsg:
        m.todo.Update(msg.Todos)
        return m, nil

    // Handle chat input submission
    case chatinput.SubmitMsg:
        if msg.AutoAccept {
            m.client.SetAutoAccept(true)
        }
        if msg.PlanMode {
            m.client.SetPlanMode(true)
        }
        return m, m.client.QueryCmd(msg.Text)
    }

    return m, nil
}
```

---

## Integration Summary

### Stream Events to Messages

| SDK Event | Adapter Message | Component |
|-----------|-----------------|-----------|
| `content_block_start` (text) | `StreamBlockStartMsg` | streamtext |
| `content_block_start` (thinking) | `StreamBlockStartMsg` | thinking |
| `content_block_start` (tool_use) | `StreamBlockStartMsg` | tooluse |
| `content_block_delta` (text) | `StreamDeltaMsg` | streamtext |
| `content_block_delta` (thinking) | `ThinkingDeltaMsg` | thinking |
| `content_block_delta` (input_json) | `ToolUseDeltaMsg` | tooluse |
| `content_block_stop` | `StreamBlockStopMsg` | all |
| `message_start` | `MessageStartMsg` | status |
| `message_delta` | `MessageDeltaMsg` | status |
| `message_stop` | `StreamDoneMsg` | layout |
| `UserMessage` | `UserMsg` | message |
| `SystemMessage` (init) | `SystemInitMsg` | status |
| `SystemMessage` (hook_response) | `SystemHookResponseMsg` | toast/log |
| `AssistantMessage` | `AssistantMsg` | message |
| `ResultMessage` | `ResultMsg` | status |
| `control_request` | `ControlRequestMsg` | debug |
| `control_response` | `ControlResponseMsg` | debug |
| Error | `StreamErrorMsg` | toast |

### canUseTool to Components

| Tool | Adapter Message | Component |
|------|-----------------|-----------|
| `Bash` | `ShowPermissionPromptMsg` | permission |
| `Write` | `ShowPermissionPromptMsg` | permission |
| `Edit` | `ShowPermissionPromptMsg` | permission, diffview |
| `Read` | `ShowPermissionPromptMsg` | permission |
| `Glob` | `ShowPermissionPromptMsg` | permission, globresults |
| `Grep` | `ShowPermissionPromptMsg` | permission, grepresults |
| `WebFetch` | `ShowPermissionPromptMsg` | permission |
| `WebSearch` | `ShowPermissionPromptMsg` | permission |
| `NotebookEdit` | `ShowPermissionPromptMsg` | permission |
| `Task` | `ShowPermissionPromptMsg` | permission, taskmgr (if background) |
| `AskUserQuestion` | `ShowQuestionPromptMsg` | question |
| `EnterPlanMode` | `ShowPlanEnterPromptMsg` | planenter |
| `ExitPlanMode` | `ShowPlanExitPromptMsg` | planexit |
| `TodoWrite` | `ShowTodoUpdateMsg` | todo |
| `TaskOutput` | `ShowTaskOutputRequestMsg` | taskmgr |
| `KillShell` | `ShowKillShellRequestMsg` | taskmgr |
| `Skill` | `ShowSkillInvokeMsg` | tooluse |

### Background Task Events

| Event | Adapter Message | Component |
|-------|-----------------|-----------|
| Task started | `TaskStartedMsg` | taskstatus, taskmgr |
| Task progress | `TaskProgressMsg` | taskmgr |
| Task completed | `TaskCompletedMsg` | taskstatus, taskmgr |

### Session Events

| Event | Adapter Message | Component |
|-------|-----------------|-----------|
| Show picker | `ShowSessionPickerMsg` | session |
| Session selected | `SessionSelectedMsg` | session |
| Session resumed | `SessionResumedMsg` | layout |

### Search Events

| Event | Adapter Message | Component |
|-------|-----------------|-----------|
| Show search | `ShowSearchMsg` | search |
| Search results | `SearchResultsMsg` | search |

### Client States to UI

| Client State | UI Behavior |
|--------------|-------------|
| `Disconnected` | Show connect button |
| `Connecting` | Show spinner |
| `Connected` | Enable input |
| `Streaming` | Show streaming indicator, disable input |
| `WaitingInput` | Show appropriate prompt component |
| `PlanMode` | Show plan mode indicator |
| `Error` | Show error toast, enable retry |
