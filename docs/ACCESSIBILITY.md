# Accessibility

Support for users with accessibility needs.

## Overview

The accessibility system provides:
- Reduced motion support
- High contrast themes
- Screen reader announcements
- Keyboard-only navigation
- Visible focus indicators
- Color-blindness safe design

---

## Configuration

### Environment Variables

```bash
# Disable animations
export REDUCE_MOTION=1

# Enable high contrast mode
export HIGH_CONTRAST=1

# Enable screen reader announcements
export SCREEN_READER=1
```

### Programmatic Configuration

```go
package accessibility

type Config struct {
    ReducedMotion    bool  // Disable all animations
    HighContrast     bool  // Use high contrast theme
    ScreenReader     bool  // Enable screen reader announcements
    FocusIndicators  bool  // Show visible focus ring
    LargeText        bool  // Increase text sizes
}

// LoadConfig creates config from environment
func LoadConfig() Config {
    return Config{
        ReducedMotion:   os.Getenv("REDUCE_MOTION") == "1",
        HighContrast:    os.Getenv("HIGH_CONTRAST") == "1",
        ScreenReader:    os.Getenv("SCREEN_READER") == "1",
        FocusIndicators: true,  // Always enabled by default
        LargeText:       os.Getenv("LARGE_TEXT") == "1",
    }
}

// Apply config to a layout
func WithAccessibility(cfg Config) layout.Option {
    return func(l *layout.Layout) {
        l.accessibility = cfg
        if cfg.ReducedMotion {
            l.animation.Disable()
        }
        if cfg.HighContrast {
            l.theme = theme.Get("high-contrast")
        }
    }
}
```

---

## Reduced Motion

When `REDUCE_MOTION=1` is set:

### Disabled Animations

| Animation | Replacement |
|-----------|-------------|
| Fade in/out | Instant show/hide |
| Slide transitions | Instant position change |
| Cursor blink | Solid cursor |
| Spinner rotation | Static indicator |
| Toast slide | Instant appear |
| Pulse effects | Static highlight |

### Implementation

```go
func (m *Model) Animate(trigger animation.Trigger) tea.Cmd {
    if m.accessibility.ReducedMotion {
        // Skip animation, apply final state immediately
        return m.applyFinalState(trigger)
    }
    return m.animation.Start(trigger)
}

// Check in animation system
func (a *AnimationEngine) Start(anim Animation) tea.Cmd {
    if a.reducedMotion {
        return nil  // No animation
    }
    return anim.Start()
}
```

---

## High Contrast Theme

Enhanced contrast for better visibility.

### Color Requirements

| Element | Minimum Contrast Ratio |
|---------|----------------------|
| Body text | 7:1 (AAA) |
| Large text | 4.5:1 (AA) |
| Interactive elements | 3:1 |
| Focus indicators | 3:1 |

### High Contrast Palette

```go
var HighContrastPalette = Palette{
    // Maximum contrast backgrounds
    Background:    lipgloss.Color("#000000"),
    Foreground:    lipgloss.Color("#FFFFFF"),

    // High contrast surfaces
    Surface0:      lipgloss.Color("#1A1A1A"),
    Surface1:      lipgloss.Color("#333333"),
    Surface2:      lipgloss.Color("#4D4D4D"),

    // Bright, saturated accents
    Primary:       lipgloss.Color("#00FFFF"),  // Cyan
    Secondary:     lipgloss.Color("#FFFF00"),  // Yellow
    Accent:        lipgloss.Color("#FF00FF"),  // Magenta

    // Semantic colors (bright)
    Success:       lipgloss.Color("#00FF00"),  // Green
    Warning:       lipgloss.Color("#FFFF00"),  // Yellow
    Error:         lipgloss.Color("#FF0000"),  // Red
    Info:          lipgloss.Color("#00FFFF"),  // Cyan

    // Text colors
    Text:          lipgloss.Color("#FFFFFF"),
    TextMuted:     lipgloss.Color("#CCCCCC"),
    TextSubtle:    lipgloss.Color("#999999"),
}
```

### Usage

```go
// Automatically enabled with HIGH_CONTRAST=1
theme := theme.Get("high-contrast")

// Or explicitly
layout := chat.New(
    chat.WithTheme(theme.HighContrast()),
)
```

---

## Screen Reader Support

Announcements for screen reader users.

### Announcement Types

```go
type AnnouncementPriority int

const (
    PriorityPolite     AnnouncementPriority = iota  // Wait for idle
    PriorityAssertive                                 // Interrupt current
)

type Announcement struct {
    Message  string
    Priority AnnouncementPriority
}
```

### Announcement Triggers

| Event | Announcement |
|-------|--------------|
| Message received | "Assistant: [first 100 chars]" |
| Tool started | "Running [tool name]" |
| Tool completed | "[tool name] complete" |
| Permission request | "Permission required for [tool]" |
| Question asked | "Question: [question text]" |
| Error occurred | "Error: [error message]" |
| Plan ready | "Plan ready for review" |
| Todo updated | "[count] tasks, [completed] complete" |

### Implementation

```go
type Announcer struct {
    enabled bool
    output  io.Writer  // Usually os.Stderr
}

func NewAnnouncer(cfg Config) *Announcer {
    return &Announcer{
        enabled: cfg.ScreenReader,
        output:  os.Stderr,
    }
}

// Announce sends message to screen reader
func (a *Announcer) Announce(msg Announcement) {
    if !a.enabled {
        return
    }

    // OSC 777 for notifications (works in many terminals)
    fmt.Fprintf(a.output, "\x1b]777;notify;%s\x07", msg.Message)
}

// AnnouncePolite waits for idle before announcing
func (a *Announcer) AnnouncePolite(message string) {
    a.Announce(Announcement{
        Message:  message,
        Priority: PriorityPolite,
    })
}

// AnnounceAssertive interrupts current speech
func (a *Announcer) AnnounceAssertive(message string) {
    a.Announce(Announcement{
        Message:  message,
        Priority: PriorityAssertive,
    })
}
```

### Component Integration

```go
func (m *Message) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case AssistantMsg:
        m.content = msg.Content

        // Announce for screen readers
        if m.announcer != nil {
            preview := truncate(msg.Content, 100)
            m.announcer.AnnouncePolite("Assistant: " + preview)
        }

        return m, nil
    }
    return m, nil
}
```

---

## Keyboard Navigation

All functionality accessible via keyboard.

### Global Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Next focusable element |
| `Shift+Tab` | Previous focusable element |
| `Enter` | Activate focused element |
| `Esc` | Cancel/close modal |
| `?` | Show help overlay |

### Component-Specific

| Component | Keys |
|-----------|------|
| Lists | `j/k` or arrows to navigate |
| Toggles | `Space` to toggle |
| Buttons | `Enter` to activate |
| Text areas | Standard text editing |
| Modals | `Esc` to close |

### Focus Traversal Order

```go
// Default focus order in chat layout
var DefaultFocusOrder = []string{
    "chat-input",
    "auto-accept-toggle",
    "plan-mode-toggle",
    "message-list",
    "status-bar",
}

// Modal focus trapping
func (m *Modal) trapFocus() {
    // Focus stays within modal until dismissed
    m.focusManager.SetTrap(m.focusableChildren)
}
```

---

## Focus Indicators

Visible focus ring for all interactive elements.

### Focus Ring Styles

```go
var DefaultFocusRing = FocusRingStyle{
    // Bright, visible ring
    Color:     lipgloss.Color("#FFFFFF"),
    Width:     2,
    Style:     FocusRingSolid,  // or FocusRingDashed

    // High contrast variant
    HighContrast: FocusRingStyle{
        Color: lipgloss.Color("#00FFFF"),
        Width: 3,
        Style: FocusRingDouble,
    },
}

func (c *Component) View() string {
    style := c.baseStyle

    if c.focused {
        // Add focus ring
        style = style.
            BorderForeground(c.focusRing.Color).
            BorderStyle(lipgloss.DoubleBorder())
    }

    return style.Render(c.content)
}
```

### Focus Ring Configuration

```go
type FocusRingStyle struct {
    Color lipgloss.Color
    Width int
    Style FocusRingType
}

type FocusRingType int

const (
    FocusRingSolid FocusRingType = iota
    FocusRingDashed
    FocusRingDouble
    FocusRingThick
)
```

---

## Color-Blindness Safe Design

Design that doesn't rely solely on color.

### Principles

1. **Use symbols with colors** - Never color alone
2. **Sufficient contrast** - Works in grayscale
3. **Patterns for differentiation** - Different borders/fills
4. **Labels for status** - Text in addition to color

### Status Indicators

```go
type StatusIndicator struct {
    Status  Status
    Label   string
    Symbol  string
    Color   lipgloss.Color
}

var StatusIndicators = map[Status]StatusIndicator{
    StatusPending:    {Label: "Pending",    Symbol: "[ ]", Color: "#808080"},
    StatusInProgress: {Label: "In Progress", Symbol: "[~]", Color: "#FFFF00"},
    StatusCompleted:  {Label: "Completed",  Symbol: "[x]", Color: "#00FF00"},
    StatusFailed:     {Label: "Failed",     Symbol: "[!]", Color: "#FF0000"},
}

func (s StatusIndicator) Render() string {
    // Symbol always present, color optional enhancement
    return lipgloss.NewStyle().
        Foreground(s.Color).
        Render(s.Symbol + " " + s.Label)
}
```

### Tool Result Indicators

| Status | Symbol | Color | Label |
|--------|--------|-------|-------|
| Success | `[OK]` | Green | "Success" |
| Error | `[!!]` | Red | "Error" |
| Warning | `[!]` | Yellow | "Warning" |
| Info | `[i]` | Blue | "Info" |
| Loading | `[~]` | Cyan | "Loading..." |

---

## Testing Accessibility

### Manual Testing Checklist

- [ ] All features work with keyboard only
- [ ] Focus ring visible on all interactive elements
- [ ] Screen reader announces important events
- [ ] High contrast mode provides sufficient contrast
- [ ] Animations disabled with `REDUCE_MOTION=1`
- [ ] Status indicators have symbols, not just colors

### Automated Testing

```go
func TestAccessibility(t *testing.T) {
    t.Run("keyboard navigation", func(t *testing.T) {
        m := NewModel()

        // Tab cycles through focusable elements
        m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
        assert.True(t, m.components[0].Focused())

        m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
        assert.True(t, m.components[1].Focused())
    })

    t.Run("reduced motion", func(t *testing.T) {
        os.Setenv("REDUCE_MOTION", "1")
        defer os.Unsetenv("REDUCE_MOTION")

        cfg := LoadConfig()
        assert.True(t, cfg.ReducedMotion)
    })

    t.Run("focus indicators visible", func(t *testing.T) {
        m := NewModel()
        m.Focus()

        view := m.View()
        // Check that focus ring is rendered
        assert.Contains(t, view, FocusRingChar)
    })
}
```

---

## Best Practices

### For Component Authors

1. **Always include keyboard shortcuts** for all mouse interactions
2. **Emit announcements** for significant state changes
3. **Support focus management** with `Focus()`, `Blur()`, `Focused()`
4. **Use semantic colors** from theme palette
5. **Include symbols** alongside color-coded status
6. **Test with reduced motion** enabled

### Usage Example

```go
func main() {
    // Load accessibility config from environment
    cfg := accessibility.LoadConfig()

    // Create layout with accessibility support
    layout := chat.New(
        chat.WithAccessibility(cfg),
    )

    // Announcer for screen readers
    announcer := accessibility.NewAnnouncer(cfg)
    layout.SetAnnouncer(announcer)

    p := tea.NewProgram(layout)
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## Environment Variable Reference

| Variable | Values | Effect |
|----------|--------|--------|
| `REDUCE_MOTION` | `1` | Disable animations |
| `HIGH_CONTRAST` | `1` | Enable high contrast theme |
| `SCREEN_READER` | `1` | Enable announcements |
| `LARGE_TEXT` | `1` | Increase text sizes |
| `NO_COLOR` | `1` | Disable colors (standard) |
