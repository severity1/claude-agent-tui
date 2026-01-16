# Lip Gloss Reference

**Package**: `github.com/charmbracelet/lipgloss`
**Documentation**: https://github.com/charmbracelet/lipgloss
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/lipgloss

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Creating Styles](#creating-styles) - NewStyle, chaining methods
- [Colors](#colors) - Color, AdaptiveColor, CompleteColor
- [Spacing](#spacing) - Padding, Margin
- [Borders](#borders) - Border types and styling
- [Dimensions](#dimensions) - Width, Height, MaxWidth, MaxHeight
- [Text Formatting](#text-formatting) - Bold, Italic, Underline, etc.
- [Layout Composition](#layout-composition) - JoinHorizontal, JoinVertical, Place
- [Alignment](#alignment) - Horizontal and vertical alignment
- [Style Inheritance](#style-inheritance) - Inherit, Copy
- [Measuring Content](#measuring-content) - Width, Height functions

## Overview

Lip Gloss provides declarative, composable styling for terminal UIs. Styles are immutable - methods return new styles.

## Creating Styles

```go
// Create and chain style properties
style := lipgloss.NewStyle().
    Bold(true).
    Italic(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(1, 2).
    Margin(1).
    Width(40)

// Render text with style
output := style.Render("Hello, World!")

// Copy and modify (styles are immutable)
errorStyle := style.Foreground(lipgloss.Color("#FF0000"))
```

## Colors

```go
// Hex colors (true color)
lipgloss.Color("#FF6AD5")

// ANSI colors (0-255)
lipgloss.Color("205")

// ANSI basic colors (0-15)
lipgloss.Color("9")  // Bright red

// Adaptive colors (light/dark mode)
lipgloss.AdaptiveColor{
    Light: "#000000",
    Dark:  "#FFFFFF",
}

// Complete color (with fallbacks)
lipgloss.CompleteColor{
    TrueColor: "#FF6AD5",
    ANSI256:   "219",
    ANSI:      "13",
}

// No color (transparent)
lipgloss.NoColor{}
```

## Spacing

```go
// Padding (inside border)
style.Padding(1)           // All sides
style.Padding(1, 2)        // Vertical, horizontal
style.Padding(1, 2, 3, 4)  // Top, right, bottom, left
style.PaddingTop(1)
style.PaddingRight(2)
style.PaddingBottom(1)
style.PaddingLeft(2)

// Margin (outside border)
style.Margin(1)            // All sides
style.Margin(1, 2)         // Vertical, horizontal
style.MarginTop(1)
style.MarginRight(2)
style.MarginBottom(1)
style.MarginLeft(2)
```

## Borders

```go
// Border types
lipgloss.NormalBorder()
lipgloss.RoundedBorder()
lipgloss.DoubleBorder()
lipgloss.ThickBorder()
lipgloss.HiddenBorder()
lipgloss.BlockBorder()

// Apply border
style.Border(lipgloss.RoundedBorder())
style.BorderForeground(lipgloss.Color("#874BFD"))
style.BorderBackground(lipgloss.Color("#1a1a2e"))

// Individual borders
style.BorderTop(true)
style.BorderRight(true)
style.BorderBottom(true)
style.BorderLeft(true)

// Border colors per side
style.BorderTopForeground(lipgloss.Color("#FF0000"))
style.BorderBottomForeground(lipgloss.Color("#00FF00"))
```

## Dimensions

```go
style.Width(40)           // Exact width
style.Height(10)          // Exact height
style.MaxWidth(80)        // Maximum width
style.MaxHeight(20)       // Maximum height

// Inline makes element not take full width
style.Inline(true)
```

## Text Formatting

```go
style.Bold(true)
style.Italic(true)
style.Underline(true)
style.Strikethrough(true)
style.Reverse(true)       // Swap fg/bg
style.Blink(true)
style.Faint(true)
```

## Layout Composition

```go
// Horizontal layout
row := lipgloss.JoinHorizontal(
    lipgloss.Top,     // Alignment: Top, Center, Bottom
    left,
    center,
    right,
)

// Vertical layout
column := lipgloss.JoinVertical(
    lipgloss.Left,    // Alignment: Left, Center, Right
    header,
    body,
    footer,
)

// Centering/positioning content
centered := lipgloss.Place(
    width, height,
    lipgloss.Center, lipgloss.Center,  // Horizontal, vertical
    content,
)

// Place with custom fill
placed := lipgloss.PlaceHorizontal(width, lipgloss.Center, content)
placed := lipgloss.PlaceVertical(height, lipgloss.Center, content)
```

## Alignment

```go
// Text alignment within width
style.Align(lipgloss.Center)
style.Align(lipgloss.Left)
style.Align(lipgloss.Right)

// Vertical alignment
style.AlignVertical(lipgloss.Center)
style.AlignVertical(lipgloss.Top)
style.AlignVertical(lipgloss.Bottom)

// Both
style.AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center)
```

## Style Inheritance

```go
// Copy a style
baseStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF"))
derivedStyle := baseStyle.Background(lipgloss.Color("#000"))

// Inherit from another style
style := lipgloss.NewStyle().Inherit(baseStyle)

// Unset properties
style.UnsetBold()
style.UnsetForeground()
style.UnsetPadding()
```

## Measuring Content

```go
// Get rendered width/height
width := lipgloss.Width(renderedString)
height := lipgloss.Height(renderedString)

// Get size
w, h := lipgloss.Size(renderedString)
```

## Common Patterns

### Card Component
```go
cardStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    Padding(1, 2).
    Width(40)

selectedCardStyle := cardStyle.
    BorderForeground(lipgloss.Color("#FF6AD5"))
```

### Status Indicators
```go
statusStyles := map[string]lipgloss.Style{
    "working": lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
    "blocked": lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")),
    "error":   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
}
```

### Responsive Width
```go
func (m Model) View() string {
    maxWidth := m.width - 4  // Account for margins
    style := lipgloss.NewStyle().Width(maxWidth)
    return style.Render(content)
}
```

## Best Practices

1. **Define styles as package variables** - Avoid recreating in View()
2. **Use composition** - Build complex layouts from simple parts
3. **Store dimensions in model** - Update on WindowSizeMsg
4. **Use adaptive colors** - Support light/dark terminals
5. **Styles are immutable** - Methods return new styles

## Common Pitfalls

- Mutating styles (they're immutable, chain methods)
- Hardcoding widths (use model dimensions)
- Creating styles in View() (performance impact)
- Forgetting border takes space (accounts for 2 chars)
- Not handling zero dimensions gracefully
