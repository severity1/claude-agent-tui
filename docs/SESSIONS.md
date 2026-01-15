# Session Management

Multi-session support with persistence and resume capabilities.

## Overview

The session system provides:
- Multiple concurrent sessions
- Session persistence and restore
- Resume from previous conversations
- Session metadata and search
- Export and import

---

## Session Data Model

### Session Info

```go
package session

type Session struct {
    ID           string       `json:"id"`
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`
    Model        string       `json:"model"`
    Messages     []Message    `json:"messages"`
    TokensUsed   TokenUsage   `json:"tokens_used"`
    Summary      string       `json:"summary"`
    Tags         []string     `json:"tags"`
    Metadata     Metadata     `json:"metadata"`
}

type TokenUsage struct {
    Input  int `json:"input"`
    Output int `json:"output"`
    Total  int `json:"total"`
}

type Metadata struct {
    WorkingDir  string            `json:"working_dir"`
    ProjectName string            `json:"project_name"`
    GitBranch   string            `json:"git_branch,omitempty"`
    Custom      map[string]string `json:"custom,omitempty"`
}

type Message struct {
    Role      string    `json:"role"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
    TokenCount int      `json:"token_count,omitempty"`
}
```

### Session Info (for picker display)

```go
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

func (s *Session) ToInfo() SessionInfo {
    summary := s.Summary
    if summary == "" && len(s.Messages) > 0 {
        // Use first user message as summary
        for _, m := range s.Messages {
            if m.Role == "user" {
                summary = truncate(m.Content, 50)
                break
            }
        }
    }

    return SessionInfo{
        ID:           s.ID,
        CreatedAt:    s.CreatedAt,
        UpdatedAt:    s.UpdatedAt,
        MessageCount: len(s.Messages),
        TokensUsed:   s.TokensUsed.Total,
        Model:        s.Model,
        Summary:      summary,
    }
}
```

---

## Session Store

Persistent session storage.

```go
type Store interface {
    // Create new session
    Create(session *Session) error

    // Get session by ID
    Get(id string) (*Session, error)

    // Update existing session
    Update(session *Session) error

    // Delete session
    Delete(id string) error

    // List all sessions
    List() ([]SessionInfo, error)

    // Search sessions by query
    Search(query string) ([]SessionInfo, error)
}
```

### File-based Store

```go
type FileStore struct {
    baseDir string
    mu      sync.RWMutex
}

func NewFileStore(baseDir string) (*FileStore, error) {
    // Create directory if not exists
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        return nil, err
    }
    return &FileStore{baseDir: baseDir}, nil
}

func (s *FileStore) sessionPath(id string) string {
    return filepath.Join(s.baseDir, id+".json")
}

func (s *FileStore) Create(session *Session) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    data, err := json.MarshalIndent(session, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(s.sessionPath(session.ID), data, 0644)
}

func (s *FileStore) Get(id string) (*Session, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    data, err := os.ReadFile(s.sessionPath(id))
    if err != nil {
        return nil, err
    }

    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }

    return &session, nil
}

func (s *FileStore) List() ([]SessionInfo, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    entries, err := os.ReadDir(s.baseDir)
    if err != nil {
        return nil, err
    }

    var sessions []SessionInfo
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
            id := strings.TrimSuffix(entry.Name(), ".json")
            session, err := s.Get(id)
            if err != nil {
                continue
            }
            sessions = append(sessions, session.ToInfo())
        }
    }

    // Sort by updated time, newest first
    sort.Slice(sessions, func(i, j int) bool {
        return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
    })

    return sessions, nil
}

func (s *FileStore) Search(query string) ([]SessionInfo, error) {
    all, err := s.List()
    if err != nil {
        return nil, err
    }

    query = strings.ToLower(query)
    var results []SessionInfo

    for _, info := range all {
        // Search in summary and ID
        if strings.Contains(strings.ToLower(info.Summary), query) ||
           strings.Contains(strings.ToLower(info.ID), query) {
            results = append(results, info)
        }
    }

    return results, nil
}
```

### Default Storage Location

```go
func DefaultStoreDir() string {
    // ~/.claude-tui/sessions/
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".claude-tui", "sessions")
}
```

---

## Session Manager

High-level session management.

```go
type Manager struct {
    store       Store
    current     *Session
    currentID   string
    autoSave    bool
    saveDebounce time.Duration
    saveTimer   *time.Timer
}

func NewManager(store Store, opts ...Option) *Manager {
    m := &Manager{
        store:        store,
        autoSave:     true,
        saveDebounce: 5 * time.Second,
    }
    for _, opt := range opts {
        opt(m)
    }
    return m
}
```

### Session Lifecycle

```go
// NewSession creates and starts a new session
func (m *Manager) NewSession(model string) (*Session, error) {
    session := &Session{
        ID:        generateID(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Model:     model,
        Metadata: Metadata{
            WorkingDir: mustGetwd(),
        },
    }

    if err := m.store.Create(session); err != nil {
        return nil, err
    }

    m.current = session
    m.currentID = session.ID
    return session, nil
}

// Resume loads and continues an existing session
func (m *Manager) Resume(id string) (*Session, error) {
    session, err := m.store.Get(id)
    if err != nil {
        return nil, err
    }

    m.current = session
    m.currentID = session.ID
    return session, nil
}

// Current returns the active session
func (m *Manager) Current() *Session {
    return m.current
}

// Close saves and closes the current session
func (m *Manager) Close() error {
    if m.current == nil {
        return nil
    }

    m.current.UpdatedAt = time.Now()
    if err := m.store.Update(m.current); err != nil {
        return err
    }

    m.current = nil
    m.currentID = ""
    return nil
}
```

### Message Management

```go
// AddMessage appends a message to the current session
func (m *Manager) AddMessage(role, content string) error {
    if m.current == nil {
        return errors.New("no active session")
    }

    msg := Message{
        Role:      role,
        Content:   content,
        Timestamp: time.Now(),
    }

    m.current.Messages = append(m.current.Messages, msg)
    m.current.UpdatedAt = time.Now()

    // Auto-save with debounce
    if m.autoSave {
        m.scheduleSave()
    }

    return nil
}

// UpdateTokens updates token usage
func (m *Manager) UpdateTokens(input, output int) {
    if m.current == nil {
        return
    }

    m.current.TokensUsed.Input += input
    m.current.TokensUsed.Output += output
    m.current.TokensUsed.Total = m.current.TokensUsed.Input + m.current.TokensUsed.Output
}

func (m *Manager) scheduleSave() {
    if m.saveTimer != nil {
        m.saveTimer.Stop()
    }

    m.saveTimer = time.AfterFunc(m.saveDebounce, func() {
        m.store.Update(m.current)
    })
}
```

---

## TUI Integration

### Session Messages

```go
// ShowSessionPickerMsg triggers the session picker
type ShowSessionPickerMsg struct {
    Sessions []SessionInfo
    Current  string
}

// SessionSelectedMsg signals a session was selected
type SessionSelectedMsg struct {
    SessionID string
    Action    SessionAction
}

type SessionAction int

const (
    SessionActionSelect SessionAction = iota  // Switch to session
    SessionActionResume                        // Resume session
    SessionActionDelete                        // Delete session
    SessionActionExport                        // Export session
)

// SessionResumedMsg signals session context loaded
type SessionResumedMsg struct {
    Session *Session
}
```

### Session Picker Component

```go
type SessionPicker struct {
    sessions  []SessionInfo
    current   string
    selected  int
    filter    string
    styles    SessionPickerStyles
    respond   func(action SessionAction, sessionID string)
}

func (s *SessionPicker) View() string {
    var items []string

    for i, sess := range s.sessions {
        item := s.renderSessionItem(sess, i == s.selected)
        items = append(items, item)
    }

    return s.styles.Container.Render(
        lipgloss.JoinVertical(lipgloss.Left,
            s.styles.Header.Render("Sessions"),
            strings.Join(items, "\n"),
        ),
    )
}

func (s *SessionPicker) renderSessionItem(sess SessionInfo, selected bool) string {
    style := s.styles.Item
    if selected {
        style = s.styles.Selected
    }
    if sess.ID == s.current {
        style = s.styles.Active
    }

    // Format: [date] summary (X messages, Y tokens)
    date := sess.UpdatedAt.Format("Jan 02")
    meta := fmt.Sprintf("(%d msgs, %d tokens)",
        sess.MessageCount, sess.TokensUsed)

    return style.Render(
        fmt.Sprintf("%s %s %s", date, sess.Summary, meta),
    )
}
```

### Keybindings

```go
func (s *SessionPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            s.selected = max(0, s.selected-1)
        case "down", "j":
            s.selected = min(len(s.sessions)-1, s.selected+1)
        case "enter":
            return s, s.selectCurrent(SessionActionResume)
        case "n":
            return s, s.newSession()
        case "d":
            return s, s.selectCurrent(SessionActionDelete)
        case "e":
            return s, s.selectCurrent(SessionActionExport)
        case "esc":
            return s, s.close()
        }
    }
    return s, nil
}
```

---

## Export/Import

### Export Formats

```go
type ExportFormat int

const (
    ExportJSON ExportFormat = iota
    ExportMarkdown
    ExportPlainText
)

func (m *Manager) Export(id string, format ExportFormat) ([]byte, error) {
    session, err := m.store.Get(id)
    if err != nil {
        return nil, err
    }

    switch format {
    case ExportJSON:
        return json.MarshalIndent(session, "", "  ")
    case ExportMarkdown:
        return exportMarkdown(session)
    case ExportPlainText:
        return exportPlainText(session)
    default:
        return nil, errors.New("unknown format")
    }
}
```

### Markdown Export

```go
func exportMarkdown(session *Session) ([]byte, error) {
    var buf bytes.Buffer

    // Header
    fmt.Fprintf(&buf, "# Session: %s\n\n", session.ID)
    fmt.Fprintf(&buf, "**Created:** %s\n", session.CreatedAt.Format(time.RFC3339))
    fmt.Fprintf(&buf, "**Model:** %s\n", session.Model)
    fmt.Fprintf(&buf, "**Messages:** %d\n", len(session.Messages))
    fmt.Fprintf(&buf, "**Tokens:** %d\n\n", session.TokensUsed.Total)
    fmt.Fprintf(&buf, "---\n\n")

    // Messages
    for _, msg := range session.Messages {
        role := strings.Title(msg.Role)
        fmt.Fprintf(&buf, "## %s\n\n", role)
        fmt.Fprintf(&buf, "%s\n\n", msg.Content)
    }

    return buf.Bytes(), nil
}
```

### Import

```go
func (m *Manager) Import(data []byte) (*Session, error) {
    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }

    // Generate new ID for imported session
    session.ID = generateID()
    session.CreatedAt = time.Now()
    session.UpdatedAt = time.Now()

    if err := m.store.Create(&session); err != nil {
        return nil, err
    }

    return &session, nil
}
```

---

## Context Window Management

Track and display context window usage.

```go
type ContextWindow struct {
    Used     int
    Max      int
    Messages []MessageContext
}

type MessageContext struct {
    Index  int
    Tokens int
    Role   string
}

func (c *ContextWindow) Percentage() float64 {
    if c.Max == 0 {
        return 0
    }
    return float64(c.Used) / float64(c.Max) * 100
}

func (c *ContextWindow) View() string {
    pct := c.Percentage()
    color := "#00FF00"  // Green
    if pct > 75 {
        color = "#FFFF00"  // Yellow
    }
    if pct > 90 {
        color = "#FF0000"  // Red
    }

    bar := renderProgressBar(pct, 20, color)
    return fmt.Sprintf("Context: %s %.1f%% (%d/%d)",
        bar, pct, c.Used, c.Max)
}
```

---

## Layout Integration

### Chat Layout with Sessions

```go
type ChatLayout struct {
    sessionMgr  *session.Manager
    currentSess *session.Session

    // ... other fields
}

func (c *ChatLayout) Init() tea.Cmd {
    // Load or create session
    sessions, _ := c.sessionMgr.store.List()
    if len(sessions) > 0 {
        // Option to resume or start new
        return func() tea.Msg {
            return ShowSessionPickerMsg{
                Sessions: sessions,
            }
        }
    }

    // No existing sessions, create new
    session, _ := c.sessionMgr.NewSession(c.model)
    c.currentSess = session
    return nil
}

func (c *ChatLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case SessionSelectedMsg:
        switch msg.Action {
        case SessionActionResume:
            session, _ := c.sessionMgr.Resume(msg.SessionID)
            c.currentSess = session
            return c, c.restoreMessages(session)

        case SessionActionDelete:
            c.sessionMgr.store.Delete(msg.SessionID)
            return c, c.refreshSessionList()
        }

    case adapter.AssistantMsg:
        // Save assistant message
        c.sessionMgr.AddMessage("assistant", msg.Content)

    case chatinput.SubmitMsg:
        // Save user message
        c.sessionMgr.AddMessage("user", msg.Text)
    }

    return c, nil
}
```

### Status Bar with Session Info

```go
func (c *ChatLayout) renderStatusBar() string {
    session := c.currentSess

    parts := []string{
        c.status.View(),
    }

    if session != nil {
        sessionInfo := fmt.Sprintf("Session: %s (%d msgs)",
            truncate(session.ID, 8),
            len(session.Messages))
        parts = append(parts, sessionInfo)
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
```

---

## Best Practices

1. **Auto-save with debounce** - Don't save on every keystroke
2. **Show context usage** - Help users understand token budget
3. **Provide export options** - Let users save conversations
4. **Quick session switching** - Easy to switch between projects
5. **Session search** - Find old conversations quickly
6. **Resume capability** - Continue where you left off
7. **Clean up old sessions** - Option to delete old data
