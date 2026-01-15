# Theming System

Comprehensive theming for consistent styling across all components.

## Overview

The theme system provides:
- Color palette with semantic naming
- Component-specific style generation
- Built-in popular themes
- Custom theme support
- Light/dark mode adaptation
- Terminal capability detection

---

## Theme Interface

The original 18-method interface was over-engineered. Themes provide colors, components build their own styles.

```go
package theme

import "github.com/charmbracelet/lipgloss"

// Theme provides identity and colors. That's it.
type Theme interface {
    Name() string
    IsDark() bool
    Palette() Palette
}

// BaseTheme provides a default implementation
type BaseTheme struct {
    name    string
    dark    bool
    palette Palette
}

func (t *BaseTheme) Name() string      { return t.name }
func (t *BaseTheme) IsDark() bool      { return t.dark }
func (t *BaseTheme) Palette() Palette  { return t.palette }
```

### Component Style Factories

Components generate their own styles from a Palette. No circular imports.

```go
// In component/streamtext/styles.go
package streamtext

import "github.com/severity1/claude-agent-tui/system/theme"

type Styles struct {
    Container lipgloss.Style
    Text      lipgloss.Style
    Cursor    lipgloss.Style
}

// StylesFromPalette creates styles from a theme palette
func StylesFromPalette(p theme.Palette) Styles {
    return Styles{
        Container: lipgloss.NewStyle().
            Background(p.Background),
        Text: lipgloss.NewStyle().
            Foreground(p.Text),
        Cursor: lipgloss.NewStyle().
            Foreground(p.Primary).
            Bold(true),
    }
}
```

### Usage

```go
// Get theme and build component styles
t := theme.Get("catppuccin-mocha")
p := t.Palette()

stream := streamtext.New(
    streamtext.WithStyles(streamtext.StylesFromPalette(p)),
)
```

This avoids the 18-method interface, prevents circular imports, and lets components own their styling logic.

---

## Color Palette

```go
type Palette struct {
    // Base colors
    Background lipgloss.AdaptiveColor
    Foreground lipgloss.AdaptiveColor

    // Surface colors (layered backgrounds)
    Surface0 lipgloss.AdaptiveColor  // Slightly elevated
    Surface1 lipgloss.AdaptiveColor  // More elevated
    Surface2 lipgloss.AdaptiveColor  // Most elevated

    // Brand/accent colors
    Primary   lipgloss.AdaptiveColor  // Main brand color
    Secondary lipgloss.AdaptiveColor  // Secondary brand
    Accent    lipgloss.AdaptiveColor  // Highlight/focus

    // Semantic colors
    Success lipgloss.AdaptiveColor  // Green - confirmations
    Warning lipgloss.AdaptiveColor  // Yellow/orange - caution
    Error   lipgloss.AdaptiveColor  // Red - errors/danger
    Info    lipgloss.AdaptiveColor  // Blue - information

    // Text hierarchy
    Text       lipgloss.AdaptiveColor  // Primary text
    TextMuted  lipgloss.AdaptiveColor  // Secondary text
    TextSubtle lipgloss.AdaptiveColor  // Tertiary text

    // Borders and dividers
    Border      lipgloss.AdaptiveColor
    BorderFocus lipgloss.AdaptiveColor

    // Syntax highlighting
    Keyword  lipgloss.AdaptiveColor
    String   lipgloss.AdaptiveColor
    Number   lipgloss.AdaptiveColor
    Comment  lipgloss.AdaptiveColor
    Function lipgloss.AdaptiveColor
    Variable lipgloss.AdaptiveColor
    Type     lipgloss.AdaptiveColor
    Operator lipgloss.AdaptiveColor
    Constant lipgloss.AdaptiveColor
}
```

---

## Built-in Themes

### Catppuccin

Soothing pastel theme with 4 variants.

```go
// Mocha (dark)
var CatppuccinMocha = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#1e1e2e", Light: "#eff1f5"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#cdd6f4", Light: "#4c4f69"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#313244", Light: "#ccd0da"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#45475a", Light: "#bcc0cc"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#585b70", Light: "#acb0be"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#89b4fa", Light: "#1e66f5"},  // Blue
    Secondary:  lipgloss.AdaptiveColor{Dark: "#cba6f7", Light: "#8839ef"},  // Mauve
    Accent:     lipgloss.AdaptiveColor{Dark: "#f5c2e7", Light: "#ea76cb"},  // Pink
    Success:    lipgloss.AdaptiveColor{Dark: "#a6e3a1", Light: "#40a02b"},  // Green
    Warning:    lipgloss.AdaptiveColor{Dark: "#f9e2af", Light: "#df8e1d"},  // Yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#f38ba8", Light: "#d20f39"},  // Red
    Info:       lipgloss.AdaptiveColor{Dark: "#89dceb", Light: "#209fb5"},  // Sky
    Text:       lipgloss.AdaptiveColor{Dark: "#cdd6f4", Light: "#4c4f69"},
    TextMuted:  lipgloss.AdaptiveColor{Dark: "#a6adc8", Light: "#6c6f85"},
    TextSubtle: lipgloss.AdaptiveColor{Dark: "#6c7086", Light: "#9ca0b0"},
    // ... syntax colors
}

// Other variants: Latte (light), Frappe, Macchiato
```

### Dracula

Classic dark theme with vibrant colors.

```go
var Dracula = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#282a36", Light: "#f8f8f2"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#f8f8f2", Light: "#282a36"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#343746", Light: "#e8e8e2"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#44475a", Light: "#d8d8d2"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#545770", Light: "#c8c8c2"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#bd93f9", Light: "#6d53c9"},  // Purple
    Secondary:  lipgloss.AdaptiveColor{Dark: "#ff79c6", Light: "#cf4996"},  // Pink
    Accent:     lipgloss.AdaptiveColor{Dark: "#8be9fd", Light: "#5bc9cd"},  // Cyan
    Success:    lipgloss.AdaptiveColor{Dark: "#50fa7b", Light: "#20ca4b"},  // Green
    Warning:    lipgloss.AdaptiveColor{Dark: "#f1fa8c", Light: "#c1ca5c"},  // Yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#ff5555", Light: "#cf2525"},  // Red
    Info:       lipgloss.AdaptiveColor{Dark: "#8be9fd", Light: "#5bc9cd"},  // Cyan
    // ...
}
```

### Nord

Arctic, north-bluish palette.

```go
var Nord = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#2e3440", Light: "#eceff4"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#eceff4", Light: "#2e3440"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#3b4252", Light: "#e5e9f0"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#434c5e", Light: "#d8dee9"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#4c566a", Light: "#c8ced9"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#88c0d0", Light: "#5e81ac"},  // Frost
    Secondary:  lipgloss.AdaptiveColor{Dark: "#81a1c1", Light: "#4c6e8f"},
    Accent:     lipgloss.AdaptiveColor{Dark: "#8fbcbb", Light: "#5f9c9b"},
    Success:    lipgloss.AdaptiveColor{Dark: "#a3be8c", Light: "#738c5c"},  // Aurora green
    Warning:    lipgloss.AdaptiveColor{Dark: "#ebcb8b", Light: "#bb9b5b"},  // Aurora yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#bf616a", Light: "#8f313a"},  // Aurora red
    Info:       lipgloss.AdaptiveColor{Dark: "#88c0d0", Light: "#5e81ac"},
    // ...
}
```

### Gruvbox

Retro groove with warm colors.

```go
// Dark variant
var GruvboxDark = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#282828", Light: "#fbf1c7"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#ebdbb2", Light: "#3c3836"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#f2e5bc"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#504945", Light: "#ebdbb2"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#665c54", Light: "#d5c4a1"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#458588", Light: "#076678"},  // Blue
    Secondary:  lipgloss.AdaptiveColor{Dark: "#b16286", Light: "#8f3f71"},  // Purple
    Accent:     lipgloss.AdaptiveColor{Dark: "#d79921", Light: "#b57614"},  // Yellow
    Success:    lipgloss.AdaptiveColor{Dark: "#98971a", Light: "#79740e"},  // Green
    Warning:    lipgloss.AdaptiveColor{Dark: "#d79921", Light: "#b57614"},  // Yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#cc241d", Light: "#9d0006"},  // Red
    Info:       lipgloss.AdaptiveColor{Dark: "#83a598", Light: "#427b58"},  // Aqua
    // ...
}

// Light variant: GruvboxLight
```

### Tokyo Night

Inspired by Tokyo city lights.

```go
var TokyoNight = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#1a1b26", Light: "#d5d6db"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#c0caf5", Light: "#343b58"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#24283b", Light: "#cbccd1"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#414868", Light: "#b4b5b9"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#545c7e", Light: "#9699a3"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#7aa2f7", Light: "#2e7de9"},  // Blue
    Secondary:  lipgloss.AdaptiveColor{Dark: "#bb9af7", Light: "#9854f1"},  // Purple
    Accent:     lipgloss.AdaptiveColor{Dark: "#7dcfff", Light: "#007197"},  // Cyan
    Success:    lipgloss.AdaptiveColor{Dark: "#9ece6a", Light: "#587539"},  // Green
    Warning:    lipgloss.AdaptiveColor{Dark: "#e0af68", Light: "#8c6c3e"},  // Yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#f7768e", Light: "#f52a65"},  // Red
    Info:       lipgloss.AdaptiveColor{Dark: "#7dcfff", Light: "#007197"},  // Cyan
    // ...
}
```

### One Dark

Atom-inspired dark theme.

```go
var OneDark = Palette{
    Background: lipgloss.AdaptiveColor{Dark: "#282c34", Light: "#fafafa"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#abb2bf", Light: "#383a42"},
    Surface0:   lipgloss.AdaptiveColor{Dark: "#31353f", Light: "#f0f0f0"},
    Surface1:   lipgloss.AdaptiveColor{Dark: "#3e4451", Light: "#e5e5e6"},
    Surface2:   lipgloss.AdaptiveColor{Dark: "#4b5263", Light: "#d4d4d5"},
    Primary:    lipgloss.AdaptiveColor{Dark: "#61afef", Light: "#4078f2"},  // Blue
    Secondary:  lipgloss.AdaptiveColor{Dark: "#c678dd", Light: "#a626a4"},  // Purple
    Accent:     lipgloss.AdaptiveColor{Dark: "#56b6c2", Light: "#0184bc"},  // Cyan
    Success:    lipgloss.AdaptiveColor{Dark: "#98c379", Light: "#50a14f"},  // Green
    Warning:    lipgloss.AdaptiveColor{Dark: "#e5c07b", Light: "#c18401"},  // Yellow
    Error:      lipgloss.AdaptiveColor{Dark: "#e06c75", Light: "#e45649"},  // Red
    Info:       lipgloss.AdaptiveColor{Dark: "#56b6c2", Light: "#0184bc"},  // Cyan
    // ...
}
```

### System

Adapts to terminal's native colors.

```go
var System = Palette{
    // Uses terminal's default colors
    Background: lipgloss.AdaptiveColor{Dark: "", Light: ""},  // Inherit
    Foreground: lipgloss.AdaptiveColor{Dark: "", Light: ""},  // Inherit
    // Uses ANSI color numbers for compatibility
    Primary:   lipgloss.AdaptiveColor{Dark: "4", Light: "4"},   // Blue
    Secondary: lipgloss.AdaptiveColor{Dark: "5", Light: "5"},   // Magenta
    Accent:    lipgloss.AdaptiveColor{Dark: "6", Light: "6"},   // Cyan
    Success:   lipgloss.AdaptiveColor{Dark: "2", Light: "2"},   // Green
    Warning:   lipgloss.AdaptiveColor{Dark: "3", Light: "3"},   // Yellow
    Error:     lipgloss.AdaptiveColor{Dark: "1", Light: "1"},   // Red
    Info:      lipgloss.AdaptiveColor{Dark: "6", Light: "6"},   // Cyan
    // ...
}
```

### High Contrast

Accessibility-focused theme with maximum contrast ratios (WCAG AAA).

```go
var HighContrast = Palette{
    // Maximum contrast backgrounds
    Background: lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"},
    Foreground: lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#000000"},

    // High contrast surfaces
    Surface0: lipgloss.AdaptiveColor{Dark: "#1A1A1A", Light: "#F5F5F5"},
    Surface1: lipgloss.AdaptiveColor{Dark: "#333333", Light: "#E5E5E5"},
    Surface2: lipgloss.AdaptiveColor{Dark: "#4D4D4D", Light: "#D5D5D5"},

    // Bright, saturated primary colors
    Primary:   lipgloss.AdaptiveColor{Dark: "#00FFFF", Light: "#0000CC"},  // Cyan/Blue
    Secondary: lipgloss.AdaptiveColor{Dark: "#FFFF00", Light: "#CC6600"},  // Yellow/Orange
    Accent:    lipgloss.AdaptiveColor{Dark: "#FF00FF", Light: "#9900CC"},  // Magenta/Purple

    // Maximum contrast semantic colors
    Success: lipgloss.AdaptiveColor{Dark: "#00FF00", Light: "#006600"},  // Bright green
    Warning: lipgloss.AdaptiveColor{Dark: "#FFFF00", Light: "#996600"},  // Bright yellow
    Error:   lipgloss.AdaptiveColor{Dark: "#FF0000", Light: "#CC0000"},  // Bright red
    Info:    lipgloss.AdaptiveColor{Dark: "#00FFFF", Light: "#0066CC"},  // Bright cyan

    // High contrast text
    Text:       lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#000000"},
    TextMuted:  lipgloss.AdaptiveColor{Dark: "#CCCCCC", Light: "#333333"},
    TextSubtle: lipgloss.AdaptiveColor{Dark: "#999999", Light: "#666666"},

    // Strong borders
    Border:      lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#000000"},
    BorderFocus: lipgloss.AdaptiveColor{Dark: "#00FFFF", Light: "#0000CC"},

    // High contrast syntax highlighting
    Keyword:  lipgloss.AdaptiveColor{Dark: "#FF00FF", Light: "#9900CC"},
    String:   lipgloss.AdaptiveColor{Dark: "#00FF00", Light: "#006600"},
    Number:   lipgloss.AdaptiveColor{Dark: "#00FFFF", Light: "#0066CC"},
    Comment:  lipgloss.AdaptiveColor{Dark: "#999999", Light: "#666666"},
    Function: lipgloss.AdaptiveColor{Dark: "#FFFF00", Light: "#996600"},
    Variable: lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#000000"},
    Type:     lipgloss.AdaptiveColor{Dark: "#00FFFF", Light: "#0066CC"},
    Operator: lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#000000"},
    Constant: lipgloss.AdaptiveColor{Dark: "#FF0000", Light: "#CC0000"},
}
```

**Accessibility Features:**
- Minimum 7:1 contrast ratio for all text (WCAG AAA)
- Minimum 3:1 contrast ratio for interactive elements
- Pure black/white backgrounds for maximum contrast
- Bright, saturated colors that are distinguishable
- Works well for users with low vision

**Enable via environment:**
```bash
export HIGH_CONTRAST=1
```

**Or programmatically:**
```go
// Load based on environment
if os.Getenv("HIGH_CONTRAST") == "1" {
    theme := theme.Get("high-contrast")
}

// Or explicitly
theme := theme.Get("high-contrast")
```

---

## Theme Registry

```go
package theme

var registry = map[string]Theme{
    "catppuccin-mocha":     &CatppuccinMocha{},
    "catppuccin-latte":     &CatppuccinLatte{},
    "catppuccin-frappe":    &CatppuccinFrappe{},
    "catppuccin-macchiato": &CatppuccinMacchiato{},
    "dracula":              &Dracula{},
    "nord":                 &Nord{},
    "gruvbox-dark":         &GruvboxDark{},
    "gruvbox-light":        &GruvboxLight{},
    "tokyonight":           &TokyoNight{},
    "onedark":              &OneDark{},
    "system":               &System{},
}

// Get retrieves a theme by name
func Get(name string) Theme {
    if t, ok := registry[name]; ok {
        return t
    }
    return registry["catppuccin-mocha"]  // Default
}

// Register adds a custom theme
func Register(name string, theme Theme) {
    registry[name] = theme
}

// List returns all available theme names
func List() []string {
    names := make([]string, 0, len(registry))
    for name := range registry {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}
```

---

## Creating Custom Themes

### Using BaseTheme

```go
type MyTheme struct {
    BaseTheme  // Embeds default implementations
    palette Palette
}

func NewMyTheme() *MyTheme {
    return &MyTheme{
        palette: Palette{
            Background: lipgloss.AdaptiveColor{Dark: "#1a1a2e", Light: "#f5f5f5"},
            Primary:    lipgloss.AdaptiveColor{Dark: "#e94560", Light: "#c73550"},
            // ... define all colors
        },
    }
}

func (t *MyTheme) Name() string { return "my-theme" }
func (t *MyTheme) IsDark() bool { return true }
func (t *MyTheme) Palette() Palette { return t.palette }

// Register it
theme.Register("my-theme", NewMyTheme())
```

### From JSON (Runtime)

```go
// Load theme from JSON file
func LoadThemeFromFile(path string) (Theme, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config ThemeConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return NewThemeFromConfig(config), nil
}

// JSON structure
type ThemeConfig struct {
    Name   string            `json:"name"`
    Dark   bool              `json:"dark"`
    Colors map[string]string `json:"colors"`
}
```

Example JSON theme:
```json
{
    "name": "my-custom-theme",
    "dark": true,
    "colors": {
        "background": "#1a1a2e",
        "foreground": "#eaeaea",
        "primary": "#e94560",
        "secondary": "#0f3460",
        "accent": "#16213e",
        "success": "#00bf63",
        "warning": "#ffc107",
        "error": "#dc3545",
        "info": "#17a2b8"
    }
}
```

---

## Using Themes

### With Components

```go
// Get theme
theme := theme.Get("catppuccin-mocha")

// Create component with theme
input := chatinput.New(
    chatinput.WithTheme(theme),
)

// Or use theme's pre-built styles
input := chatinput.New(
    chatinput.WithStyles(theme.ChatInput()),
)
```

### With Layouts

```go
// Apply theme to entire layout
chat := layout.NewChat(
    layout.WithTheme(theme.Get("dracula")),
)
```

### Runtime Switching

```go
// Message to switch themes
type SwitchThemeMsg struct {
    Name string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case SwitchThemeMsg:
        m.theme = theme.Get(msg.Name)
        m.applyTheme()
        return m, nil
    }
    // ...
}

func (m *Model) applyTheme() {
    m.chatInput = chatinput.New(chatinput.WithTheme(m.theme))
    m.status = status.New(status.WithTheme(m.theme))
    // ... update all components
}
```

---

## Terminal Capability Detection

```go
import "github.com/charmbracelet/x/ansi"

// Detect color profile
profile := lipgloss.DefaultRenderer().ColorProfile()

switch profile {
case lipgloss.TrueColor:
    // Full 24-bit color support
case lipgloss.ANSI256:
    // 256 color support
case lipgloss.ANSI:
    // 16 color support
case lipgloss.Ascii:
    // No color support
}
```

Lipgloss automatically degrades colors to match terminal capabilities.
