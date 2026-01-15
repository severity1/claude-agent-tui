# Animation System

Smooth animations for component state changes using spring-based physics.

## Overview

The animation system provides:
- Spring-based physics via `harmonica`
- Event-driven animation triggers
- Configurable animation types
- Per-component animation settings
- Idle/busy state animations

---

## Animation Triggers

Events that initiate animations.

```go
package animation

type Trigger int

const (
    // User input triggers
    TriggerKeyPress Trigger = iota
    TriggerMouseClick
    TriggerMouseHover
    TriggerMouseScroll
    TriggerFocusGain
    TriggerFocusLoss

    // SDK event triggers
    TriggerStreamStart
    TriggerStreamDelta
    TriggerStreamDone
    TriggerToolStart
    TriggerToolComplete
    TriggerToolError
    TriggerPermissionRequest
    TriggerQuestionRequest

    // State triggers
    TriggerIdle
    TriggerBusy
    TriggerSuccess
    TriggerError
    TriggerTransition
)
```

---

## Animation Interface

```go
type Animation interface {
    // Start initializes and returns the first frame command
    Start() tea.Cmd

    // Update processes animation frames
    Update(msg tea.Msg) (Animation, tea.Cmd)

    // View returns the current animation frame
    View() string

    // Done returns true when animation is complete
    Done() bool

    // Reset restarts the animation
    Reset()
}
```

---

## Built-in Animation Types

### FadeIn / FadeOut

Opacity transition using ANSI color manipulation.

```go
type FadeIn struct {
    content  string
    opacity  float64
    spring   *harmonica.Spring
    style    lipgloss.Style
}

func NewFadeIn(content string, opts ...FadeOption) *FadeIn

type FadeOption func(*FadeIn)

func WithDuration(d time.Duration) FadeOption
func WithStyle(s lipgloss.Style) FadeOption
```

### SlideIn / SlideOut

Content slides from edge.

```go
type Direction int

const (
    DirectionLeft Direction = iota
    DirectionRight
    DirectionUp
    DirectionDown
)

type SlideIn struct {
    content   string
    direction Direction
    offset    float64
    spring    *harmonica.Spring
    width     int
    height    int
}

func NewSlideIn(content string, dir Direction, opts ...SlideOption) *SlideIn
```

### Pulse

Brief highlight effect.

```go
type Pulse struct {
    content    string
    intensity  float64
    spring     *harmonica.Spring
    baseStyle  lipgloss.Style
    pulseColor lipgloss.Color
}

func NewPulse(content string, opts ...PulseOption) *Pulse
```

### Shake

Horizontal shake for errors/warnings.

```go
type Shake struct {
    content   string
    offset    float64
    cycles    int
    current   int
    spring    *harmonica.Spring
}

func NewShake(content string, cycles int) *Shake
```

### Bounce

Vertical bounce effect.

```go
type Bounce struct {
    content string
    offset  float64
    spring  *harmonica.Spring
}

func NewBounce(content string) *Bounce
```

### Breathe

Subtle pulsing for idle states.

```go
type Breathe struct {
    content   string
    intensity float64
    spring    *harmonica.Spring
    inhaling  bool
}

func NewBreathe(content string) *Breathe
```

### Blink

Cursor blinking.

```go
type Blink struct {
    content  string
    visible  bool
    interval time.Duration
}

func NewBlink(content string, interval time.Duration) *Blink
```

### Spin

Loading spinner animation.

```go
type Spin struct {
    frames   []string
    current  int
    interval time.Duration
}

// Default spinner frames
var DefaultSpinFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func NewSpin(frames []string, interval time.Duration) *Spin
```

### TypeWriter

Character-by-character reveal.

```go
type TypeWriter struct {
    content  string
    revealed int
    interval time.Duration
}

func NewTypeWriter(content string, interval time.Duration) *TypeWriter
```

---

## Spring Configuration

Using `harmonica` for physics-based motion.

```go
import "github.com/charmbracelet/harmonica"

type SpringConfig struct {
    Frequency float64  // Angular frequency (speed)
    Damping   float64  // Damping ratio
}

// Preset configurations
var (
    // Fast, slight overshoot
    SpringSnappy = SpringConfig{
        Frequency: 8.0,
        Damping:   0.8,
    }

    // Smooth, no overshoot
    SpringSmooth = SpringConfig{
        Frequency: 5.0,
        Damping:   1.0,
    }

    // Bouncy, visible oscillation
    SpringBouncy = SpringConfig{
        Frequency: 6.0,
        Damping:   0.5,
    }

    // Slow, gentle motion
    SpringGentle = SpringConfig{
        Frequency: 3.0,
        Damping:   1.2,
    }
)

// Create spring from config
func (c SpringConfig) NewSpring() *harmonica.Spring {
    return harmonica.NewSpring(harmonica.FPS(60), c.Frequency, c.Damping)
}
```

### Damping Ratios

| Ratio | Behavior |
|-------|----------|
| < 1.0 | Under-damped: oscillates before settling |
| = 1.0 | Critically damped: fastest without overshoot |
| > 1.0 | Over-damped: slow, no overshoot |

---

## Component Animation Config

Per-component animation settings.

```go
type AnimationConfig struct {
    // Enable/disable animations
    Enabled bool

    // Per-trigger animations
    OnKeyPress    Animation
    OnClick       Animation
    OnHover       Animation
    OnFocus       Animation
    OnBlur        Animation
    OnStreamDelta Animation
    OnToolStart   Animation
    OnToolComplete Animation
    OnError       Animation
    OnIdle        Animation

    // Transition animations
    EnterAnimation Animation
    ExitAnimation  Animation

    // Timing
    IdleDelay time.Duration  // Delay before idle animation starts
}

// Default config
var DefaultAnimationConfig = AnimationConfig{
    Enabled:       true,
    OnKeyPress:    nil,  // No animation
    OnClick:       NewPulse("", WithDuration(100*time.Millisecond)),
    OnHover:       nil,
    OnFocus:       NewPulse("", WithDuration(150*time.Millisecond)),
    OnBlur:        nil,
    OnStreamDelta: nil,
    OnToolStart:   NewSpin(DefaultSpinFrames, 80*time.Millisecond),
    OnToolComplete: NewPulse("", WithDuration(200*time.Millisecond)),
    OnError:       NewShake("", 3),
    OnIdle:        NewBreathe(""),
    IdleDelay:     5 * time.Second,
}
```

---

## AnimatedComponent Interface

Components that support animations.

```go
type AnimatedComponent interface {
    tea.Model

    // Trigger an animation
    Animate(trigger Trigger) tea.Cmd

    // Get/set animation config
    AnimationConfig() AnimationConfig
    SetAnimationConfig(AnimationConfig)
}
```

---

## Usage Examples

### Basic Animation

```go
// Create component with animation
input := chatinput.New(
    chatinput.WithAnimation(AnimationConfig{
        Enabled: true,
        OnFocus: NewPulse(""),
        OnError: NewShake("", 3),
    }),
)

// Trigger animation in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter && m.input.Value() == "" {
            // Shake on empty submit
            return m, m.input.Animate(TriggerError)
        }
    }
    return m, nil
}
```

### Toast Slide Animation

```go
// Toast slides in from right
func (t *ToastManager) Show(level Level, title, message string) tea.Cmd {
    toast := Toast{
        ID:        uuid.New().String(),
        Level:     level,
        Title:     title,
        Message:   message,
        Animation: NewSlideIn(t.renderToast(toast), DirectionRight),
    }
    t.toasts = append(t.toasts, toast)
    return toast.Animation.Start()
}
```

### Streaming Cursor Blink

```go
// StreamText with blinking cursor
func (m *StreamText) View() string {
    content := m.content.String()
    if m.streaming {
        cursor := m.cursorAnim.View()
        return content + cursor
    }
    return content
}

func (m *StreamText) Init() tea.Cmd {
    m.cursorAnim = NewBlink("█", 530*time.Millisecond)
    return m.cursorAnim.Start()
}
```

### Loading Spinner

```go
// ToolUse with loading spinner
func (m *ToolUse) View() string {
    if m.isLoading {
        spinner := m.spinner.View()
        return m.styles.Container.Render(
            m.styles.ToolName.Render(m.toolName) + " " + spinner,
        )
    }
    // ... render complete state
}
```

---

## Animation Messages

```go
// Frame tick message
type AnimationTickMsg struct {
    ID string  // Animation instance ID
}

// Animation complete message
type AnimationDoneMsg struct {
    ID      string
    Trigger Trigger
}
```

---

## Performance Considerations

1. **Frame Rate**: Animations run at 60 FPS by default
2. **Batching**: Multiple animations can run concurrently
3. **Cleanup**: Animations auto-cleanup when Done() returns true
4. **Disable**: Set `Enabled: false` to skip all animations
5. **Reduce Motion**: Respect user preferences for reduced motion

```go
// Check for reduced motion preference
func ShouldReduceMotion() bool {
    // Check REDUCE_MOTION env var
    if os.Getenv("REDUCE_MOTION") == "1" {
        return true
    }
    return false
}

// Apply reduced motion
func GetAnimationConfig() AnimationConfig {
    if ShouldReduceMotion() {
        return AnimationConfig{Enabled: false}
    }
    return DefaultAnimationConfig
}
```
