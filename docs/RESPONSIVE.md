# Responsive System

Components adapt to terminal dimensions automatically.

## Overview

The responsive system provides:
- Breakpoint-based layouts
- Min/max size constraints
- Automatic variant switching
- Resize handling
- Graceful degradation

---

## Breakpoints

```go
package responsive

type Breakpoint int

const (
    BreakpointXS Breakpoint = iota  // < 40 columns
    BreakpointSM                     // 40-79 columns
    BreakpointMD                     // 80-119 columns
    BreakpointLG                     // 120-159 columns
    BreakpointXL                     // 160+ columns
)

// Column thresholds
var BreakpointWidths = map[Breakpoint]int{
    BreakpointXS: 0,
    BreakpointSM: 40,
    BreakpointMD: 80,
    BreakpointLG: 120,
    BreakpointXL: 160,
}

// Height breakpoints
type HeightBreakpoint int

const (
    HeightXS HeightBreakpoint = iota  // < 15 rows
    HeightSM                           // 15-24 rows
    HeightMD                           // 25-39 rows
    HeightLG                           // 40+ rows
)

var BreakpointHeights = map[HeightBreakpoint]int{
    HeightXS: 0,
    HeightSM: 15,
    HeightMD: 25,
    HeightLG: 40,
}
```

---

## Current Breakpoint

```go
// CurrentBreakpoint returns the breakpoint for given width
func CurrentBreakpoint(width int) Breakpoint {
    switch {
    case width >= 160:
        return BreakpointXL
    case width >= 120:
        return BreakpointLG
    case width >= 80:
        return BreakpointMD
    case width >= 40:
        return BreakpointSM
    default:
        return BreakpointXS
    }
}

// CurrentHeightBreakpoint returns the breakpoint for given height
func CurrentHeightBreakpoint(height int) HeightBreakpoint {
    switch {
    case height >= 40:
        return HeightLG
    case height >= 25:
        return HeightMD
    case height >= 15:
        return HeightSM
    default:
        return HeightXS
    }
}
```

---

## Responsive Config

Per-component responsive settings.

```go
type ResponsiveConfig struct {
    // Size constraints
    MinWidth  int
    MaxWidth  int
    MinHeight int
    MaxHeight int

    // Breakpoint-specific variants
    VariantXS Variant
    VariantSM Variant
    VariantMD Variant
    VariantLG Variant
    VariantXL Variant

    // Breakpoint-specific visibility
    HiddenXS bool
    HiddenSM bool

    // Padding per breakpoint
    PaddingXS int
    PaddingSM int
    PaddingMD int
    PaddingLG int
    PaddingXL int
}

// VariantForWidth returns the appropriate variant
func (r ResponsiveConfig) VariantForWidth(width int) Variant {
    bp := CurrentBreakpoint(width)
    switch bp {
    case BreakpointXL:
        return r.VariantXL
    case BreakpointLG:
        return r.VariantLG
    case BreakpointMD:
        return r.VariantMD
    case BreakpointSM:
        return r.VariantSM
    default:
        return r.VariantXS
    }
}

// PaddingForWidth returns the appropriate padding
func (r ResponsiveConfig) PaddingForWidth(width int) int {
    bp := CurrentBreakpoint(width)
    switch bp {
    case BreakpointXL:
        return r.PaddingXL
    case BreakpointLG:
        return r.PaddingLG
    case BreakpointMD:
        return r.PaddingMD
    case BreakpointSM:
        return r.PaddingSM
    default:
        return r.PaddingXS
    }
}
```

---

## Default Responsive Configs

```go
var DefaultResponsiveConfigs = map[string]ResponsiveConfig{
    "chat": {
        MinWidth:  30,
        MinHeight: 10,
        VariantXS: VariantCompact,
        VariantSM: VariantView,
        VariantMD: VariantView,
        VariantLG: VariantWindow,
        VariantXL: VariantFullscreen,
        PaddingXS: 0,
        PaddingSM: 1,
        PaddingMD: 2,
        PaddingLG: 2,
        PaddingXL: 3,
    },
    "permission": {
        MinWidth:  40,
        MinHeight: 8,
        VariantXS: VariantInline,
        VariantSM: VariantOverlay,
        VariantMD: VariantModal,
        VariantLG: VariantModal,
        VariantXL: VariantModal,
    },
    "status": {
        VariantXS: VariantCompact,
        VariantSM: VariantCompact,
        VariantMD: VariantBar,
        VariantLG: VariantBar,
        VariantXL: VariantBar,
        HiddenXS:  true,  // Hide on very small terminals
    },
    "todo": {
        VariantXS: VariantCollapsible,
        VariantSM: VariantCollapsible,
        VariantMD: VariantSidebar,
        VariantLG: VariantSidebar,
        VariantXL: VariantSidebar,
    },
}
```

---

## Size Constraints

```go
type Constraints struct {
    Width  Range
    Height Range
}

type Range struct {
    Min int
    Max int  // 0 = no max
}

// Clamp constrains a value to the range
func (r Range) Clamp(value int) int {
    if value < r.Min {
        return r.Min
    }
    if r.Max > 0 && value > r.Max {
        return r.Max
    }
    return value
}

// Apply constraints to dimensions
func (c Constraints) Apply(width, height int) (int, int) {
    return c.Width.Clamp(width), c.Height.Clamp(height)
}
```

---

## Resize Handling

```go
import "github.com/charmbracelet/x/term"

// WindowSizeMsg is sent when terminal resizes
// (built into Bubble Tea as tea.WindowSizeMsg)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.breakpoint = CurrentBreakpoint(msg.Width)
        m.heightBreakpoint = CurrentHeightBreakpoint(msg.Height)

        // Update all component sizes
        return m, m.resize()
    }
    return m, nil
}

func (m *Model) resize() tea.Cmd {
    // Recalculate layout
    m.layout.SetSize(m.width, m.height)

    // Update component variants based on breakpoint
    for _, c := range m.components {
        if rc, ok := c.(ResponsiveComponent); ok {
            rc.SetBreakpoint(m.breakpoint)
        }
    }

    return nil
}
```

---

## ResponsiveComponent Interface

```go
type ResponsiveComponent interface {
    tea.Model

    // SetSize updates component dimensions
    SetSize(width, height int)

    // SetBreakpoint changes the current breakpoint
    SetBreakpoint(bp Breakpoint)

    // ResponsiveConfig returns the component's config
    ResponsiveConfig() ResponsiveConfig
}
```

---

## Layout Calculations

### Available Space

```go
func (l *Layout) calculateAvailableSpace() (int, int) {
    width := l.width
    height := l.height

    // Subtract fixed elements
    if l.showStatus {
        height -= 1  // Status bar
    }
    if l.showInput {
        height -= l.inputHeight
    }

    // Apply padding
    padding := l.responsive.PaddingForWidth(width)
    width -= padding * 2
    height -= padding * 2

    return width, height
}
```

### Split Layouts

```go
type SplitLayout struct {
    ratio      float64  // 0.0 to 1.0
    minPrimary int
    minSecondary int
    direction  Direction  // Horizontal or Vertical
}

func (s *SplitLayout) Calculate(available int) (primary, secondary int) {
    primary = int(float64(available) * s.ratio)
    secondary = available - primary

    // Apply minimums
    if primary < s.minPrimary {
        primary = s.minPrimary
        secondary = available - primary
    }
    if secondary < s.minSecondary {
        secondary = s.minSecondary
        primary = available - secondary
    }

    return primary, secondary
}
```

---

## Graceful Degradation

Behavior when terminal is too small.

```go
func (m Model) View() string {
    // Check minimum size
    if m.width < 30 || m.height < 10 {
        return m.renderTooSmall()
    }

    // Normal rendering
    return m.renderNormal()
}

func (m Model) renderTooSmall() string {
    msg := "Terminal too small"
    if m.width >= 20 {
        msg = fmt.Sprintf("Resize to at least %dx%d", 30, 10)
    }
    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        msg,
    )
}
```

---

## Breakpoint-Specific Rendering

```go
func (m *Message) View() string {
    switch m.breakpoint {
    case BreakpointXS:
        return m.renderCompact()
    case BreakpointSM:
        return m.renderNarrow()
    default:
        return m.renderFull()
    }
}

func (m *Message) renderCompact() string {
    // No avatars, minimal padding
    return m.styles.Compact.Render(m.text)
}

func (m *Message) renderNarrow() string {
    // Small avatar, reduced padding
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        m.renderSmallAvatar(),
        m.styles.Narrow.Render(m.text),
    )
}

func (m *Message) renderFull() string {
    // Full avatar, timestamps, etc.
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        m.renderAvatar(),
        m.renderContent(),
    )
}
```

---

## Responsive Text Wrapping

```go
import "github.com/muesli/reflow/wordwrap"

func (m *Message) wrapText(text string) string {
    maxWidth := m.width - m.padding*2

    // Adjust wrap width based on breakpoint
    if m.breakpoint == BreakpointXS {
        maxWidth = m.width  // No padding
    }

    return wordwrap.String(text, maxWidth)
}
```

---

## Usage Example

```go
func NewChatLayout(opts ...Option) *ChatLayout {
    l := &ChatLayout{
        responsive: DefaultResponsiveConfigs["chat"],
    }
    for _, opt := range opts {
        opt(l)
    }
    return l
}

func WithResponsiveConfig(rc ResponsiveConfig) Option {
    return func(l *ChatLayout) {
        l.responsive = rc
    }
}

func (l *ChatLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        l.SetSize(msg.Width, msg.Height)

        // Get variant for current size
        variant := l.responsive.VariantForWidth(msg.Width)
        l.setVariant(variant)

        return l, nil
    }
    return l, nil
}
```
