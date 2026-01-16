# Glamour Reference

**Package**: `github.com/charmbracelet/glamour`
**Documentation**: https://github.com/charmbracelet/glamour
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/glamour

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Overview](#overview)
- [Basic Rendering](#basic-rendering)
- [Custom Renderer](#custom-renderer)
- [Styles](#styles)
- [Auto Style Detection](#auto-style-detection)
- [Word Wrapping](#word-wrapping)
- [Integration with Bubble Tea](#integration-with-bubble-tea)

## Overview

Glamour is a markdown renderer for the terminal. It converts markdown to ANSI-styled terminal output, supporting themes and customization.

## Basic Rendering

```go
import "github.com/charmbracelet/glamour"

// Simple rendering with style
out, err := glamour.Render(markdownContent, "dark")
if err != nil {
    return err
}
fmt.Print(out)

// Available built-in styles: "dark", "light", "dracula", "pink", "notty", "ascii"
```

## Custom Renderer

```go
import "github.com/charmbracelet/glamour"

// Create renderer with options
r, err := glamour.NewTermRenderer(
    glamour.WithAutoStyle(),
    glamour.WithWordWrap(80),
)
if err != nil {
    return err
}

out, err := r.Render(markdownContent)
if err != nil {
    return err
}
fmt.Print(out)
```

## Styles

### Built-in Styles

| Style | Description |
|-------|-------------|
| `dark` | Dark terminal background |
| `light` | Light terminal background |
| `dracula` | Dracula color theme |
| `pink` | Pink accent theme |
| `notty` | No TTY (minimal formatting) |
| `ascii` | ASCII-only (no unicode) |

### Custom Style via Environment

```bash
# Use built-in style
export GLAMOUR_STYLE=dark

# Use custom stylesheet file
export GLAMOUR_STYLE=/path/to/style.json
```

### Environment Config

```go
// Use GLAMOUR_STYLE environment variable
out, err := glamour.RenderWithEnvironmentConfig(markdownContent)
```

## Auto Style Detection

```go
// Automatically detect terminal background and choose style
r, _ := glamour.NewTermRenderer(
    glamour.WithAutoStyle(),
)
```

This detects whether the terminal has a dark or light background and selects the appropriate style.

## Word Wrapping

```go
// Set word wrap width (default: 80)
r, _ := glamour.NewTermRenderer(
    glamour.WithWordWrap(100),
)

// Disable word wrapping
r, _ := glamour.NewTermRenderer(
    glamour.WithWordWrap(0),
)
```

## Integration with Bubble Tea

### In a Component View

```go
import (
    "github.com/charmbracelet/glamour"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    content  string
    rendered string
    width    int
    renderer *glamour.TermRenderer
}

func New() Model {
    r, _ := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(80),
    )
    return Model{renderer: r}
}

func (m Model) View() string {
    return m.rendered
}

// Update rendered content when markdown or width changes
func (m *Model) SetContent(markdown string) {
    r, _ := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(m.width),
    )
    m.rendered, _ = r.Render(markdown)
}
```

### Handling Window Resize

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Re-render with new width
        m.width = msg.Width
        r, _ := glamour.NewTermRenderer(
            glamour.WithAutoStyle(),
            glamour.WithWordWrap(m.width),
        )
        m.rendered, _ = r.Render(m.content)
    }
    return m, nil
}
```

## Common Options

| Option | Description |
|--------|-------------|
| `WithAutoStyle()` | Auto-detect terminal background |
| `WithWordWrap(n)` | Set word wrap width (0 disables) |
| `WithEnvironmentConfig()` | Use GLAMOUR_STYLE env var |
| `WithStylePath(path)` | Load custom style from file |
| `WithStyles(styles)` | Use programmatic styles |

## Style JSON Format

Custom styles are JSON files with element-specific formatting:

```json
{
  "document": {
    "block_prefix": "",
    "block_suffix": "",
    "margin": 2
  },
  "heading": {
    "bold": true,
    "color": "99"
  },
  "code_block": {
    "color": "244",
    "margin": 2
  }
}
```

## Best Practices

1. **Cache the renderer** - Creating renderers is relatively expensive
2. **Handle errors** - Render can fail on malformed markdown
3. **Re-render on resize** - Word wrap needs updated width
4. **Use AutoStyle** - Better user experience across terminals
5. **Consider notty** - For non-TTY output (pipes, redirects)
