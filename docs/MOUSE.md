# Mouse Support

Full mouse interaction with click, hover, and scroll support.

## Overview

The mouse system provides:
- Click handling with zones
- Hover states and feedback
- Scroll wheel support
- Drag support (where applicable)
- Cursor style hints

---

## Mouse Zones

Clickable regions within components.

```go
package mouse

type Zone struct {
    ID       string
    Bounds   Rect
    OnClick  func(x, y int) tea.Cmd
    OnHover  func(entered bool) tea.Cmd
    OnScroll func(delta int) tea.Cmd
    Cursor   CursorStyle
}

type Rect struct {
    X      int
    Y      int
    Width  int
    Height int
}

// Check if point is within bounds
func (r Rect) Contains(x, y int) bool {
    return x >= r.X && x < r.X+r.Width &&
           y >= r.Y && y < r.Y+r.Height
}

type CursorStyle int

const (
    CursorDefault CursorStyle = iota
    CursorPointer
    CursorText
    CursorNotAllowed
    CursorGrab
    CursorGrabbing
)
```

---

## Mouse Manager

Manages zones and routes mouse events.

```go
type MouseManager struct {
    zones      []Zone
    hovered    string
    pressed    string
    scrollable map[string]bool
}

func NewMouseManager() *MouseManager {
    return &MouseManager{
        zones:      make([]Zone, 0),
        scrollable: make(map[string]bool),
    }
}

// Register adds a clickable zone
func (m *MouseManager) Register(zone Zone) {
    m.zones = append(m.zones, zone)
}

// Unregister removes a zone by ID
func (m *MouseManager) Unregister(id string) {
    for i, z := range m.zones {
        if z.ID == id {
            m.zones = append(m.zones[:i], m.zones[i+1:]...)
            return
        }
    }
}

// Clear removes all zones
func (m *MouseManager) Clear() {
    m.zones = m.zones[:0]
    m.hovered = ""
    m.pressed = ""
}

// Handle processes a mouse event
func (m *MouseManager) Handle(msg tea.MouseMsg) tea.Cmd {
    var cmds []tea.Cmd

    for _, zone := range m.zones {
        if !zone.Bounds.Contains(msg.X, msg.Y) {
            // Mouse left this zone
            if m.hovered == zone.ID && zone.OnHover != nil {
                cmds = append(cmds, zone.OnHover(false))
                m.hovered = ""
            }
            continue
        }

        switch msg.Action {
        case tea.MouseActionPress:
            if msg.Button == tea.MouseButtonLeft {
                m.pressed = zone.ID
            }

        case tea.MouseActionRelease:
            if m.pressed == zone.ID && zone.OnClick != nil {
                cmds = append(cmds, zone.OnClick(msg.X, msg.Y))
            }
            m.pressed = ""

        case tea.MouseActionMotion:
            if m.hovered != zone.ID && zone.OnHover != nil {
                if m.hovered != "" {
                    // Leave previous zone
                    for _, z := range m.zones {
                        if z.ID == m.hovered && z.OnHover != nil {
                            cmds = append(cmds, z.OnHover(false))
                        }
                    }
                }
                cmds = append(cmds, zone.OnHover(true))
                m.hovered = zone.ID
            }

        case tea.MouseActionWheel:
            if zone.OnScroll != nil {
                delta := 1
                if msg.Button == tea.MouseButtonWheelUp {
                    delta = -1
                }
                cmds = append(cmds, zone.OnScroll(delta))
            }
        }
    }

    return tea.Batch(cmds...)
}

// HoveredZone returns the currently hovered zone ID
func (m *MouseManager) HoveredZone() string {
    return m.hovered
}

// SetScrollable marks a zone as scrollable
func (m *MouseManager) SetScrollable(id string, scrollable bool) {
    m.scrollable[id] = scrollable
}
```

---

## MouseSupport Interface

Components with mouse support implement this interface.

```go
type MouseSupport interface {
    // Zones returns all clickable zones
    Zones() []Zone

    // HandleMouse processes mouse events
    HandleMouse(msg tea.MouseMsg) tea.Cmd

    // Hover state
    SetHovered(hovered bool)
    IsHovered() bool
}
```

---

## Component Implementation

### Button with Mouse

```go
type Button struct {
    label   string
    bounds  Rect
    hovered bool
    pressed bool
    onClick func() tea.Cmd
    styles  ButtonStyles
}

func (b *Button) Zones() []Zone {
    return []Zone{{
        ID:      "button-" + b.label,
        Bounds:  b.bounds,
        OnClick: func(x, y int) tea.Cmd { return b.onClick() },
        OnHover: func(entered bool) tea.Cmd {
            b.hovered = entered
            return nil
        },
        Cursor: CursorPointer,
    }}
}

func (b *Button) View() string {
    style := b.styles.Normal
    if b.pressed {
        style = b.styles.Pressed
    } else if b.hovered {
        style = b.styles.Hover
    }
    return style.Render(b.label)
}
```

### List with Mouse

```go
type List struct {
    items    []Item
    cursor   int
    viewport Rect
    hovered  int
}

func (l *List) Zones() []Zone {
    zones := make([]Zone, len(l.items))
    for i, item := range l.items {
        idx := i
        zones[i] = Zone{
            ID: fmt.Sprintf("list-item-%d", i),
            Bounds: Rect{
                X:      l.viewport.X,
                Y:      l.viewport.Y + i,
                Width:  l.viewport.Width,
                Height: 1,
            },
            OnClick: func(x, y int) tea.Cmd {
                l.cursor = idx
                return l.selectItem()
            },
            OnHover: func(entered bool) tea.Cmd {
                if entered {
                    l.hovered = idx
                } else if l.hovered == idx {
                    l.hovered = -1
                }
                return nil
            },
            Cursor: CursorPointer,
        }
    }
    return zones
}
```

### Scrollable Content

```go
type ScrollView struct {
    content  string
    viewport viewport.Model
    bounds   Rect
}

func (s *ScrollView) Zones() []Zone {
    return []Zone{{
        ID:     "scrollview",
        Bounds: s.bounds,
        OnScroll: func(delta int) tea.Cmd {
            if delta < 0 {
                s.viewport.LineUp(1)
            } else {
                s.viewport.LineDown(1)
            }
            return nil
        },
        Cursor: CursorDefault,
    }}
}
```

---

## Enabling Mouse in Bubble Tea

```go
func main() {
    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),  // Enable mouse motion events
    )
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Mouse Options

| Option | Description |
|--------|-------------|
| `tea.WithMouseCellMotion()` | Track mouse movement (hover) |
| `tea.WithMouseAllMotion()` | Track all mouse movement |
| `tea.WithMouseButtonTracking()` | Track button clicks only |

---

## Hover Styles

Apply different styles on hover.

```go
type HoverStyles struct {
    Normal  lipgloss.Style
    Hovered lipgloss.Style
}

func (c *Component) View() string {
    style := c.styles.Normal
    if c.hovered {
        style = c.styles.Hovered
    }
    return style.Render(c.content)
}
```

Example hover effects:
```go
var DefaultHoverStyles = HoverStyles{
    Normal: lipgloss.NewStyle().
        Background(lipgloss.Color("#282828")),
    Hovered: lipgloss.NewStyle().
        Background(lipgloss.Color("#3c3836")).
        Bold(true),
}
```

---

## Click Feedback

Provide visual feedback on click.

```go
func (b *Button) HandleMouse(msg tea.MouseMsg) tea.Cmd {
    if !b.bounds.Contains(msg.X, msg.Y) {
        return nil
    }

    switch msg.Action {
    case tea.MouseActionPress:
        if msg.Button == tea.MouseButtonLeft {
            b.pressed = true
            return nil
        }

    case tea.MouseActionRelease:
        if b.pressed {
            b.pressed = false
            // Trigger click animation
            return tea.Batch(
                b.onClick(),
                b.Animate(TriggerMouseClick),
            )
        }
    }

    return nil
}
```

---

## Zone Registration Pattern

```go
type Layout struct {
    mouse      *MouseManager
    components []MouseSupport
}

func (l *Layout) Init() tea.Cmd {
    l.mouse = NewMouseManager()
    return nil
}

func (l *Layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Re-register zones on each frame
        l.mouse.Clear()
        for _, c := range l.components {
            for _, zone := range c.Zones() {
                l.mouse.Register(zone)
            }
        }
        return l, l.mouse.Handle(msg)
    }
    return l, nil
}
```

---

## Coordinate Mapping

Components need to track their position for accurate zones.

```go
type Component struct {
    x, y   int
    width  int
    height int
}

func (c *Component) SetPosition(x, y int) {
    c.x = x
    c.y = y
}

func (c *Component) Bounds() Rect {
    return Rect{
        X:      c.x,
        Y:      c.y,
        Width:  c.width,
        Height: c.height,
    }
}
```

---

## Double-Click Detection

```go
type DoubleClickDetector struct {
    lastClick time.Time
    lastZone  string
    threshold time.Duration
}

func NewDoubleClickDetector() *DoubleClickDetector {
    return &DoubleClickDetector{
        threshold: 300 * time.Millisecond,
    }
}

func (d *DoubleClickDetector) IsDoubleClick(zoneID string) bool {
    now := time.Now()
    isDouble := zoneID == d.lastZone &&
        now.Sub(d.lastClick) < d.threshold

    d.lastClick = now
    d.lastZone = zoneID

    return isDouble
}
```

Usage:
```go
func (c *Component) HandleMouse(msg tea.MouseMsg) tea.Cmd {
    if msg.Action == tea.MouseActionRelease {
        if c.doubleClick.IsDoubleClick(c.id) {
            return c.onDoubleClick()
        }
        return c.onClick()
    }
    return nil
}
```
