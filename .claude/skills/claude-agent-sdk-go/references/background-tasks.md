# Background Tasks Reference

Handle SDK background tasks, TaskOutput, and KillShell integrations.

## Background Task Concept

Claude can run tasks in the background using the `run_in_background` parameter on the Task tool. Background tasks continue executing while the main conversation proceeds.

## Task Messages

### Task Lifecycle Messages

```go
// Task started in background
type TaskStartedMsg struct {
    TaskID      string
    Description string
    StartTime   time.Time
}

// Task progress update
type TaskProgressMsg struct {
    TaskID   string
    Progress float64  // 0.0 to 1.0
    Status   string   // Current status text
}

// Task completed
type TaskCompletedMsg struct {
    TaskID    string
    Success   bool
    Output    string
    Error     error
    Duration  time.Duration
}

// Task output available (for TaskOutput tool)
type TaskOutputMsg struct {
    TaskID  string
    Content string
    IsError bool
}

// Task killed (for KillShell tool)
type TaskKilledMsg struct {
    TaskID string
    Reason string
}
```

### Detecting Background Tasks

```go
func convertToTeaMsg(msg claudecode.Message) tea.Msg {
    switch m := msg.(type) {
    case *claudecode.ToolUseBlock:
        // Check if Task tool with run_in_background
        if m.Name == "Task" {
            if input, ok := m.Input.(map[string]interface{}); ok {
                if bg, _ := input["run_in_background"].(bool); bg {
                    return TaskStartedMsg{
                        TaskID:      m.ID,
                        Description: input["description"].(string),
                        StartTime:   time.Now(),
                    }
                }
            }
        }

        // Check for TaskOutput tool
        if m.Name == "TaskOutput" {
            return TaskOutputRequestMsg{
                TaskID: m.Input.(map[string]interface{})["task_id"].(string),
            }
        }

        // Check for KillShell tool
        if m.Name == "KillShell" {
            return KillShellRequestMsg{
                ShellID: m.Input.(map[string]interface{})["shell_id"].(string),
            }
        }

    case *claudecode.ToolResultBlock:
        // Task completed - check if was background
        if isBackgroundTask(m.ToolUseID) {
            return TaskCompletedMsg{
                TaskID:  m.ToolUseID,
                Success: !m.IsError,
                Output:  m.Content,
            }
        }
    }

    return nil
}
```

## Task Manager Integration

Track and display background tasks.

```go
type TaskManager struct {
    tasks   map[string]*BackgroundTask
    mu      sync.RWMutex
    maxShow int  // Max visible in compact view
}

type BackgroundTask struct {
    ID          string
    Description string
    StartTime   time.Time
    Status      TaskStatus
    Progress    float64
    Output      strings.Builder
    Error       error
}

type TaskStatus int

const (
    TaskStatusPending TaskStatus = iota
    TaskStatusRunning
    TaskStatusCompleted
    TaskStatusFailed
    TaskStatusKilled
)

func (m *TaskManager) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case TaskStartedMsg:
        m.mu.Lock()
        m.tasks[msg.TaskID] = &BackgroundTask{
            ID:          msg.TaskID,
            Description: msg.Description,
            StartTime:   msg.StartTime,
            Status:      TaskStatusRunning,
        }
        m.mu.Unlock()

    case TaskProgressMsg:
        m.mu.Lock()
        if task, ok := m.tasks[msg.TaskID]; ok {
            task.Progress = msg.Progress
        }
        m.mu.Unlock()

    case TaskCompletedMsg:
        m.mu.Lock()
        if task, ok := m.tasks[msg.TaskID]; ok {
            if msg.Success {
                task.Status = TaskStatusCompleted
            } else {
                task.Status = TaskStatusFailed
                task.Error = msg.Error
            }
            task.Output.WriteString(msg.Output)
        }
        m.mu.Unlock()

    case TaskKilledMsg:
        m.mu.Lock()
        if task, ok := m.tasks[msg.TaskID]; ok {
            task.Status = TaskStatusKilled
        }
        m.mu.Unlock()
    }

    return nil
}
```

## TaskOutput Tool Handling

Display output from background tasks.

```go
// Request to view task output
type TaskOutputRequestMsg struct {
    TaskID string
}

// Response with task output
type TaskOutputResponseMsg struct {
    TaskID  string
    Content string
    IsError bool
}

func (a *ControlAdapter) HandleTaskOutput(taskID string) tea.Cmd {
    return func() tea.Msg {
        // Get task output from manager
        task := a.taskManager.Get(taskID)
        if task == nil {
            return TaskOutputResponseMsg{
                TaskID:  taskID,
                Content: "Task not found",
                IsError: true,
            }
        }

        return TaskOutputResponseMsg{
            TaskID:  taskID,
            Content: task.Output.String(),
            IsError: task.Status == TaskStatusFailed,
        }
    }
}
```

## KillShell Tool Handling

Kill running background shells.

```go
// Request to kill a shell
type KillShellRequestMsg struct {
    ShellID string
}

// Confirmation dialog
type KillShellConfirmMsg struct {
    ShellID     string
    Description string
    Respond     func(confirmed bool)
}

func (a *ControlAdapter) HandleKillShell(input claudecode.CanUseToolInput) tea.Cmd {
    shellID := input.Input.(map[string]interface{})["shell_id"].(string)

    // Get shell info for confirmation
    task := a.taskManager.Get(shellID)
    description := "unknown task"
    if task != nil {
        description = task.Description
    }

    // Send confirmation prompt to TUI
    return func() tea.Msg {
        return KillShellConfirmMsg{
            ShellID:     shellID,
            Description: description,
            Respond: func(confirmed bool) {
                a.responseCh <- PermissionResponse{Allow: confirmed}
            },
        }
    }
}
```

## Task Status Indicator

Show running task count.

```go
type TaskStatusIndicator struct {
    manager *TaskManager
    styles  TaskStatusStyles
}

func (t *TaskStatusIndicator) View() string {
    running := t.manager.RunningCount()
    if running == 0 {
        return ""
    }

    // Show running count
    return t.styles.Running.Render(
        fmt.Sprintf("[%d tasks running]", running),
    )
}

func (m *TaskManager) RunningCount() int {
    m.mu.RLock()
    defer m.mu.RUnlock()

    count := 0
    for _, task := range m.tasks {
        if task.Status == TaskStatusRunning {
            count++
        }
    }
    return count
}
```

## Task List View

Full task manager view.

```go
func (m *TaskManager) View() string {
    m.mu.RLock()
    defer m.mu.RUnlock()

    var lines []string

    for _, task := range m.tasks {
        status := statusSymbol(task.Status)
        elapsed := time.Since(task.StartTime).Round(time.Second)

        line := fmt.Sprintf("%s %s (%s)", status, task.Description, elapsed)

        if task.Status == TaskStatusRunning && task.Progress > 0 {
            line += fmt.Sprintf(" %.0f%%", task.Progress*100)
        }

        lines = append(lines, line)
    }

    return strings.Join(lines, "\n")
}

func statusSymbol(status TaskStatus) string {
    switch status {
    case TaskStatusPending:
        return "[.]"
    case TaskStatusRunning:
        return "[>]"
    case TaskStatusCompleted:
        return "[+]"
    case TaskStatusFailed:
        return "[!]"
    case TaskStatusKilled:
        return "[x]"
    default:
        return "[-]"
    }
}
```

## SDK Options for Background Tasks

```go
// Configure task behavior
claudecode.WithMaxBackgroundTasks(5)     // Limit concurrent tasks
claudecode.WithTaskTimeout(10*time.Minute) // Task timeout
```

## TUI Integration Pattern

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case TaskStartedMsg:
        // Add to task manager
        m.taskManager.Update(msg)
        // Show toast notification
        return m, m.toast.Show(toast.Info, "Task Started", msg.Description)

    case TaskCompletedMsg:
        m.taskManager.Update(msg)
        if msg.Success {
            return m, m.toast.Show(toast.Success, "Task Completed", msg.TaskID)
        }
        return m, m.toast.Show(toast.Error, "Task Failed", msg.Error.Error())

    case TaskOutputRequestMsg:
        // Show task output in modal
        return m, m.showTaskOutput(msg.TaskID)

    case KillShellConfirmMsg:
        // Show confirmation dialog
        m.showKillConfirm(msg)
        return m, nil
    }

    return m, nil
}
```

## Related Documentation

- [messages.md](messages.md) - Message types
- [options.md](options.md) - SDK configuration options
- docs/COMPONENTS.md - TaskManager component specification
