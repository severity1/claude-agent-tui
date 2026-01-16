# Harmonica Reference

**Package**: `github.com/charmbracelet/harmonica`
**Documentation**: https://github.com/charmbracelet/harmonica
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/harmonica

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Overview](#overview)
- [Spring Basics](#spring-basics)
- [Configuration Parameters](#configuration-parameters)
- [Usage Patterns](#usage-patterns)
- [Bubble Tea Integration](#bubble-tea-integration)

## Overview

Harmonica provides spring-based physics animations for smooth, natural motion in terminal applications. It's framework-agnostic and works with any rendering system including Bubble Tea.

## Spring Basics

```go
import "github.com/charmbracelet/harmonica"

// Create a spring with:
// - timeDelta: frame timing (use FPS helper)
// - angularFrequency: animation speed (6.0 is typical)
// - dampingRatio: springiness (0.5 = bouncy, 1.0 = smooth)
spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5)

// Update each frame
position, velocity = spring.Update(position, velocity, target)
```

## Configuration Parameters

### Time Delta
Frame interval for animation updates:
```go
// 60 FPS (16.67ms per frame)
delta := harmonica.FPS(60)

// 30 FPS
delta := harmonica.FPS(30)

// Custom delta (e.g., from game engine)
delta := 0.016 // 16ms
```

### Angular Frequency
Controls animation speed (higher = faster):
```go
// Slow, gentle motion
spring := harmonica.NewSpring(delta, 3.0, dampingRatio)

// Medium speed (typical)
spring := harmonica.NewSpring(delta, 6.0, dampingRatio)

// Fast, snappy motion
spring := harmonica.NewSpring(delta, 12.0, dampingRatio)
```

### Damping Ratio
Determines "springiness" behavior:

| Value | Behavior | Use Case |
|-------|----------|----------|
| < 1.0 | Under-damped (oscillates) | Bouncy UI, playful motion |
| = 1.0 | Critically damped | Fastest smooth settling |
| > 1.0 | Over-damped | Slow, gentle approach |

```go
// Bouncy (overshoots target)
spring := harmonica.NewSpring(delta, 6.0, 0.3)

// Smooth (no overshoot)
spring := harmonica.NewSpring(delta, 6.0, 1.0)

// Gentle (slower settling)
spring := harmonica.NewSpring(delta, 6.0, 1.5)
```

## Usage Patterns

### Basic Animation Loop
```go
type AnimatedValue struct {
    position float64
    velocity float64
    target   float64
    spring   harmonica.Spring
}

func NewAnimatedValue(target float64) *AnimatedValue {
    return &AnimatedValue{
        position: 0,
        velocity: 0,
        target:   target,
        spring:   harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
    }
}

func (a *AnimatedValue) Update() {
    a.position, a.velocity = a.spring.Update(a.position, a.velocity, a.target)
}

func (a *AnimatedValue) SetTarget(target float64) {
    a.target = target
}

func (a *AnimatedValue) Value() float64 {
    return a.position
}
```

### Multi-Dimensional Animation
```go
type Position struct {
    x, xVel float64
    y, yVel float64
    spring  harmonica.Spring
}

func (p *Position) Update(targetX, targetY float64) {
    p.x, p.xVel = p.spring.Update(p.x, p.xVel, targetX)
    p.y, p.yVel = p.spring.Update(p.y, p.yVel, targetY)
}
```

## Bubble Tea Integration

### Animated Component
```go
type Model struct {
    scroll       float64
    scrollVel    float64
    scrollTarget float64
    spring       harmonica.Spring
}

func New() Model {
    return Model{
        spring: harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.8),
    }
}

// Custom tick message for animation frames
type tickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m Model) Init() tea.Cmd {
    return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.scrollTarget += 10
        case "k", "up":
            m.scrollTarget -= 10
        }
    case tickMsg:
        m.scroll, m.scrollVel = m.spring.Update(m.scroll, m.scrollVel, m.scrollTarget)
        return m, tick()
    }
    return m, nil
}
```

### Respecting REDUCE_MOTION
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        if os.Getenv("REDUCE_MOTION") != "" {
            // Instant transition
            m.scroll = m.scrollTarget
            m.scrollVel = 0
        } else {
            // Animated transition
            m.scroll, m.scrollVel = m.spring.Update(m.scroll, m.scrollVel, m.scrollTarget)
        }
        return m, tick()
    }
    return m, nil
}
```

## Common Presets

| Animation Type | Frequency | Damping | Notes |
|---------------|-----------|---------|-------|
| Snappy UI | 8.0 | 1.0 | Fast, no bounce |
| Bouncy | 6.0 | 0.3 | Playful overshoot |
| Smooth scroll | 6.0 | 0.8 | Natural feeling |
| Gentle reveal | 4.0 | 1.2 | Slow fade-in |

## Best Practices

1. **Match frame rate** to your rendering loop
2. **Use tick commands** in Bubble Tea for consistent updates
3. **Respect REDUCE_MOTION** for accessibility
4. **Single spring per dimension** for predictable motion
5. **Cache spring instances** rather than recreating each frame
