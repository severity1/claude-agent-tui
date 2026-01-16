# Bubbles Reference

**Package**: `github.com/charmbracelet/bubbles`
**Documentation**: https://github.com/charmbracelet/bubbles
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/bubbles

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [List Component](#list-component) - Filterable, navigable item lists
- [Viewport](#viewport-scrollable-content) - Scrollable content areas
- [Text Input](#text-input-single-line) - Single-line text entry
- [Text Area](#text-area-multi-line) - Multi-line text entry
- [Spinner](#spinner) - Loading indicators
- [Progress Bar](#progress-bar) - Progress indicators
- [Table](#table) - Tabular data display
- [Help](#help) - Key binding help display
- [Key Bindings](#key-bindings) - Input handling
- [Timer](#timer) - Countdown timer
- [Stopwatch](#stopwatch) - Elapsed time counter
- [Paginator](#paginator) - Page navigation
- [File Picker](#file-picker) - File system navigation
- [Cursor](#cursor) - Text cursor management

## Overview

Bubbles provides pre-built, customizable TUI components for Bubble Tea applications.

## List Component

```go
import "github.com/charmbracelet/bubbles/list"

// Implement list.Item interface
type item struct {
    title, desc string
}
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Create list
items := []list.Item{
    item{"First", "Description 1"},
    item{"Second", "Description 2"},
}
l := list.New(items, list.NewDefaultDelegate(), 30, 20)
l.Title = "My List"
l.SetShowStatusBar(true)
l.SetShowFilter(true)
l.SetShowHelp(true)
l.SetFilteringEnabled(true)

// In Update
newList, cmd := l.Update(msg)
l = newList

// Get selected item
selected := l.SelectedItem().(item)
```

## Viewport (Scrollable Content)

```go
import "github.com/charmbracelet/bubbles/viewport"

vp := viewport.New(80, 20)
vp.SetContent(longContent)

// Navigation
vp.GotoTop()
vp.GotoBottom()
vp.LineDown(5)
vp.LineUp(5)
vp.HalfViewDown()
vp.HalfViewUp()
vp.ViewDown()
vp.ViewUp()

// Check position
atTop := vp.AtTop()
atBottom := vp.AtBottom()
percent := vp.ScrollPercent()

// In Update
vp, cmd = vp.Update(msg)
```

## Text Input (Single Line)

```go
import "github.com/charmbracelet/bubbles/textinput"

ti := textinput.New()
ti.Placeholder = "Enter text..."
ti.Focus()
ti.CharLimit = 156
ti.Width = 40

// For passwords
ti.EchoMode = textinput.EchoPassword
ti.EchoCharacter = '*'

// Get value
value := ti.Value()

// In Update
ti, cmd = ti.Update(msg)
```

## Text Area (Multi-line)

```go
import "github.com/charmbracelet/bubbles/textarea"

ta := textarea.New()
ta.Placeholder = "Type here..."
ta.SetWidth(60)
ta.SetHeight(10)
ta.Focus()
ta.CharLimit = 1000
ta.ShowLineNumbers = true

// Get value
value := ta.Value()

// In Update
ta, cmd = ta.Update(msg)
```

## Spinner

```go
import "github.com/charmbracelet/bubbles/spinner"

s := spinner.New()
s.Spinner = spinner.Dot  // Or: Line, MiniDot, Jump, Pulse, Points, Globe, Moon, Monkey, Meter, Hamburger
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

// In Update (must forward tick messages)
s, cmd = s.Update(msg)

// Init returns the tick command
func (m Model) Init() tea.Cmd {
    return m.spinner.Tick
}
```

## Progress Bar

```go
import "github.com/charmbracelet/bubbles/progress"

// With gradient
p := progress.New(progress.WithDefaultGradient())
p := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

// Solid color
p := progress.New(progress.WithSolidFill("#FF0000"))

// Set progress (0.0 to 1.0)
p.SetPercent(0.5)

// Animated progress
cmd := p.SetPercent(0.75)  // Returns animation command

// In Update
progressModel, cmd := p.Update(msg)
p = progressModel.(progress.Model)
```

## Table

```go
import "github.com/charmbracelet/bubbles/table"

columns := []table.Column{
    {Title: "Name", Width: 20},
    {Title: "Age", Width: 10},
    {Title: "City", Width: 15},
}

rows := []table.Row{
    {"Alice", "30", "NYC"},
    {"Bob", "25", "LA"},
}

t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(10),
)

// Styling
s := table.DefaultStyles()
s.Header = s.Header.Bold(true)
s.Selected = s.Selected.Background(lipgloss.Color("57"))
t.SetStyles(s)

// Get selected row
selected := t.SelectedRow()

// In Update
t, cmd = t.Update(msg)
```

## Help

```go
import "github.com/charmbracelet/bubbles/help"

h := help.New()
h.ShowAll = true  // Show full help vs short

// Implement help.KeyMap interface
type keyMap struct {
    Up   key.Binding
    Down key.Binding
    Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down},
        {k.Quit},
    }
}

// Render help
helpView := h.View(keys)
```

## Key Bindings

```go
import "github.com/charmbracelet/bubbles/key"

// Define bindings
var keys = struct {
    Up   key.Binding
    Down key.Binding
    Quit key.Binding
}{
    Up: key.NewBinding(
        key.WithKeys("k", "up"),
        key.WithHelp("k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("j", "down"),
        key.WithHelp("j", "move down"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}

// Check if key matches
if key.Matches(msg, keys.Up) {
    m.cursor--
}

// Enable/disable bindings
keys.Up.SetEnabled(false)
```

## Timer

```go
import "github.com/charmbracelet/bubbles/timer"

t := timer.NewWithInterval(5*time.Minute, time.Second)

// In Update
switch msg := msg.(type) {
case timer.TickMsg:
    var cmd tea.Cmd
    t, cmd = t.Update(msg)
    return m, cmd
case timer.TimeoutMsg:
    // Timer finished
}

// Start/stop
t.Toggle()
t.Start()
t.Stop()

// Check state
t.Running()
t.Timedout()
```

## Stopwatch

```go
import "github.com/charmbracelet/bubbles/stopwatch"

sw := stopwatch.NewWithInterval(time.Second)

// Start in Init
func (m Model) Init() tea.Cmd {
    return m.stopwatch.Start()
}

// Control
sw.Toggle()
sw.Start()
sw.Stop()
sw.Reset()

// Get elapsed
elapsed := sw.Elapsed()
```

## Paginator

```go
import "github.com/charmbracelet/bubbles/paginator"

p := paginator.New()
p.Type = paginator.Dots  // Or: Arabic
p.PerPage = 10
p.SetTotalPages(len(items) / p.PerPage)

// Navigation
p.NextPage()
p.PrevPage()
p.Page = 0

// Get current page items
start, end := p.GetSliceBounds(len(items))
pageItems := items[start:end]
```

## File Picker

```go
import "github.com/charmbracelet/bubbles/filepicker"

fp := filepicker.New()
fp.CurrentDirectory, _ = os.UserHomeDir()
fp.AllowedTypes = []string{".go", ".md", ".txt"}
fp.ShowHidden = false
fp.DirAllowed = false
fp.FileAllowed = true

// Check if file selected
if didSelect, path := fp.DidSelectFile(msg); didSelect {
    // User selected: path
}

// Check if disabled file selected
if didSelect, path := fp.DidSelectDisabledFile(msg); didSelect {
    // User tried to select disabled file
}
```

## Cursor

```go
import "github.com/charmbracelet/bubbles/cursor"

c := cursor.New()
c.SetMode(cursor.CursorBlink)  // Or: CursorStatic, CursorHide
c.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

// In Update
c, cmd = c.Update(msg)
```

## Best Practices

1. **Forward messages to components** - They need tick/key messages
2. **Initialize with Init()** - Return component's Init cmd
3. **Store updated components** - They return new instances
4. **Customize with styles** - Most have SetStyles methods
5. **Check component state** - Use accessor methods
