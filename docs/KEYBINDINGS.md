# Keybinding System

Customizable keyboard shortcuts with vim-style defaults.

## Overview

The keybinding system provides:
- Default keybindings (vim + arrow keys)
- Per-component keymaps
- Runtime customization
- Conflict detection
- Help overlay generation

---

## Keymap Structure

```go
package keymap

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
    // Global
    Quit   key.Binding
    Help   key.Binding
    Cancel key.Binding

    // Navigation
    Up       key.Binding
    Down     key.Binding
    Left     key.Binding
    Right    key.Binding
    PageUp   key.Binding
    PageDown key.Binding
    Home     key.Binding
    End      key.Binding
    Tab      key.Binding
    ShiftTab key.Binding

    // Actions
    Submit   key.Binding
    Select   key.Binding
    Toggle   key.Binding
    Expand   key.Binding
    Collapse key.Binding
    Copy     key.Binding
    Paste    key.Binding

    // Chat-specific
    ToggleAutoAccept key.Binding
    TogglePlanMode   key.Binding
    Interrupt        key.Binding
    NewLine          key.Binding

    // Permission-specific
    Allow       key.Binding
    Deny        key.Binding
    AlwaysAllow key.Binding

    // Question-specific
    SelectOption key.Binding
    CustomInput  key.Binding

    // Plan-specific
    Approve        key.Binding
    Reject         key.Binding
    RequestChanges key.Binding
    TogglePerm     key.Binding

    // Scroll
    ScrollUp   key.Binding
    ScrollDown key.Binding
}
```

---

## Default Keybindings

```go
var DefaultKeymap = Keymap{
    // Global
    Quit:   key.NewBinding(
        key.WithKeys("ctrl+c", "q"),
        key.WithHelp("q/ctrl+c", "quit"),
    ),
    Help:   key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "help"),
    ),
    Cancel: key.NewBinding(
        key.WithKeys("esc"),
        key.WithHelp("esc", "cancel"),
    ),

    // Navigation (vim + arrows)
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("up/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("down/j", "down"),
    ),
    Left: key.NewBinding(
        key.WithKeys("left", "h"),
        key.WithHelp("left/h", "left"),
    ),
    Right: key.NewBinding(
        key.WithKeys("right", "l"),
        key.WithHelp("right/l", "right"),
    ),
    PageUp: key.NewBinding(
        key.WithKeys("pgup", "ctrl+u"),
        key.WithHelp("pgup", "page up"),
    ),
    PageDown: key.NewBinding(
        key.WithKeys("pgdown", "ctrl+d"),
        key.WithHelp("pgdn", "page down"),
    ),
    Home: key.NewBinding(
        key.WithKeys("home", "g"),
        key.WithHelp("home/g", "top"),
    ),
    End: key.NewBinding(
        key.WithKeys("end", "G"),
        key.WithHelp("end/G", "bottom"),
    ),
    Tab: key.NewBinding(
        key.WithKeys("tab"),
        key.WithHelp("tab", "next"),
    ),
    ShiftTab: key.NewBinding(
        key.WithKeys("shift+tab"),
        key.WithHelp("shift+tab", "prev"),
    ),

    // Actions
    Submit: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "submit"),
    ),
    Select: key.NewBinding(
        key.WithKeys(" ", "space"),
        key.WithHelp("space", "select"),
    ),
    Toggle: key.NewBinding(
        key.WithKeys(" ", "space"),
        key.WithHelp("space", "toggle"),
    ),
    Expand: key.NewBinding(
        key.WithKeys("enter", "l", "right"),
        key.WithHelp("enter", "expand"),
    ),
    Collapse: key.NewBinding(
        key.WithKeys("esc", "h", "left"),
        key.WithHelp("esc", "collapse"),
    ),
    Copy: key.NewBinding(
        key.WithKeys("c", "ctrl+c"),
        key.WithHelp("c", "copy"),
    ),

    // Chat
    ToggleAutoAccept: key.NewBinding(
        key.WithKeys("ctrl+a"),
        key.WithHelp("ctrl+a", "auto-accept"),
    ),
    TogglePlanMode: key.NewBinding(
        key.WithKeys("ctrl+p"),
        key.WithHelp("ctrl+p", "plan mode"),
    ),
    Interrupt: key.NewBinding(
        key.WithKeys("ctrl+c"),
        key.WithHelp("ctrl+c", "interrupt"),
    ),
    NewLine: key.NewBinding(
        key.WithKeys("alt+enter"),
        key.WithHelp("alt+enter", "new line"),
    ),

    // Permission
    Allow: key.NewBinding(
        key.WithKeys("a", "y"),
        key.WithHelp("a/y", "allow"),
    ),
    Deny: key.NewBinding(
        key.WithKeys("d", "n"),
        key.WithHelp("d/n", "deny"),
    ),
    AlwaysAllow: key.NewBinding(
        key.WithKeys("A"),
        key.WithHelp("A", "always allow"),
    ),

    // Question
    SelectOption: key.NewBinding(
        key.WithKeys("1", "2", "3", "4"),
        key.WithHelp("1-4", "select option"),
    ),
    CustomInput: key.NewBinding(
        key.WithKeys("o"),
        key.WithHelp("o", "other"),
    ),

    // Plan
    Approve: key.NewBinding(
        key.WithKeys("a"),
        key.WithHelp("a", "approve"),
    ),
    Reject: key.NewBinding(
        key.WithKeys("r"),
        key.WithHelp("r", "reject"),
    ),
    RequestChanges: key.NewBinding(
        key.WithKeys("c"),
        key.WithHelp("c", "changes"),
    ),
    TogglePerm: key.NewBinding(
        key.WithKeys("space"),
        key.WithHelp("space", "toggle perm"),
    ),

    // Scroll
    ScrollUp: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("up/k", "scroll up"),
    ),
    ScrollDown: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("down/j", "scroll down"),
    ),
}
```

---

## Component-Specific Keymaps

Each component uses a subset of the global keymap.

### ChatInput Keymap

```go
type ChatInputKeymap struct {
    Submit           key.Binding
    NewLine          key.Binding
    Cancel           key.Binding
    ToggleAutoAccept key.Binding
    TogglePlanMode   key.Binding
}

func (k Keymap) ChatInput() ChatInputKeymap {
    return ChatInputKeymap{
        Submit:           k.Submit,
        NewLine:          k.NewLine,
        Cancel:           k.Cancel,
        ToggleAutoAccept: k.ToggleAutoAccept,
        TogglePlanMode:   k.TogglePlanMode,
    }
}
```

### Permission Keymap

```go
type PermissionKeymap struct {
    Up          key.Binding
    Down        key.Binding
    Select      key.Binding
    Allow       key.Binding
    Deny        key.Binding
    AlwaysAllow key.Binding
}

func (k Keymap) Permission() PermissionKeymap {
    return PermissionKeymap{
        Up:          k.Up,
        Down:        k.Down,
        Select:      k.Submit,
        Allow:       k.Allow,
        Deny:        k.Deny,
        AlwaysAllow: k.AlwaysAllow,
    }
}
```

### Question Keymap

```go
type QuestionKeymap struct {
    Up           key.Binding
    Down         key.Binding
    Select       key.Binding
    Toggle       key.Binding
    Submit       key.Binding
    CustomInput  key.Binding
    Cancel       key.Binding
}

func (k Keymap) Question() QuestionKeymap {
    return QuestionKeymap{
        Up:          k.Up,
        Down:        k.Down,
        Select:      k.Select,
        Toggle:      k.Toggle,
        Submit:      k.Submit,
        CustomInput: k.CustomInput,
        Cancel:      k.Cancel,
    }
}
```

---

## Customization API

```go
// Set changes a binding
func (k *Keymap) Set(action string, keys ...string) error {
    binding := key.NewBinding(key.WithKeys(keys...))

    switch action {
    case "quit":
        k.Quit = binding
    case "help":
        k.Help = binding
    case "up":
        k.Up = binding
    // ... other actions
    default:
        return fmt.Errorf("unknown action: %s", action)
    }
    return nil
}

// Get retrieves a binding by action name
func (k *Keymap) Get(action string) (key.Binding, bool) {
    switch action {
    case "quit":
        return k.Quit, true
    case "help":
        return k.Help, true
    // ... other actions
    default:
        return key.Binding{}, false
    }
}

// Merge overlays another keymap on top
func (k *Keymap) Merge(overlay Keymap) {
    // Only merge non-empty bindings
    if len(overlay.Quit.Keys()) > 0 {
        k.Quit = overlay.Quit
    }
    // ... for each binding
}

// Conflicts checks for key conflicts
func (k *Keymap) Conflicts() []string {
    keyMap := make(map[string][]string)

    // Collect all keys and their actions
    for _, key := range k.Quit.Keys() {
        keyMap[key] = append(keyMap[key], "quit")
    }
    // ... for each binding

    // Find conflicts
    var conflicts []string
    for key, actions := range keyMap {
        if len(actions) > 1 {
            conflicts = append(conflicts,
                fmt.Sprintf("%s: %s", key, strings.Join(actions, ", ")))
        }
    }
    return conflicts
}
```

---

## Loading Custom Keymaps

### From JSON

```go
func LoadKeymapFromFile(path string) (*Keymap, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config map[string][]string
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    km := DefaultKeymap
    for action, keys := range config {
        if err := km.Set(action, keys...); err != nil {
            return nil, err
        }
    }

    return &km, nil
}
```

Example JSON keymap:
```json
{
    "quit": ["ctrl+q"],
    "up": ["up", "ctrl+p"],
    "down": ["down", "ctrl+n"],
    "submit": ["enter", "ctrl+m"],
    "allow": ["a", "y", "enter"],
    "deny": ["d", "n", "backspace"]
}
```

---

## Help Overlay

Auto-generated help from keybindings.

```go
import "github.com/charmbracelet/bubbles/help"

type HelpOverlay struct {
    help   help.Model
    keymap Keymap
    width  int
}

func NewHelpOverlay(keymap Keymap) *HelpOverlay {
    return &HelpOverlay{
        help:   help.New(),
        keymap: keymap,
    }
}

// ShortHelp returns abbreviated help
func (h *HelpOverlay) ShortHelp() []key.Binding {
    return []key.Binding{
        h.keymap.Up,
        h.keymap.Down,
        h.keymap.Submit,
        h.keymap.Cancel,
        h.keymap.Help,
    }
}

// FullHelp returns complete help grouped by category
func (h *HelpOverlay) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        // Navigation
        {h.keymap.Up, h.keymap.Down, h.keymap.Left, h.keymap.Right},
        {h.keymap.PageUp, h.keymap.PageDown, h.keymap.Home, h.keymap.End},
        // Actions
        {h.keymap.Submit, h.keymap.Select, h.keymap.Cancel},
        // Mode toggles
        {h.keymap.ToggleAutoAccept, h.keymap.TogglePlanMode},
    }
}

func (h *HelpOverlay) View() string {
    return h.help.View(h)
}
```

---

## Context-Aware Keybindings

Keybindings change based on current state.

```go
func (m Model) activeKeymap() Keymap {
    km := m.keymap

    switch m.state {
    case StateChat:
        // Disable navigation during text input
        km.Up = key.NewBinding(key.WithDisabled())
        km.Down = key.NewBinding(key.WithDisabled())

    case StatePermission:
        // Enable quick keys
        km.Allow.SetEnabled(true)
        km.Deny.SetEnabled(true)

    case StateStreaming:
        // Only interrupt available
        km.Submit = key.NewBinding(key.WithDisabled())
        km.Interrupt.SetEnabled(true)
    }

    return km
}
```

---

## Usage in Components

```go
// Component with custom keymap
func NewChatInput(opts ...Option) *ChatInput {
    c := &ChatInput{
        keymap: DefaultKeymap.ChatInput(),
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

func WithKeymap(km ChatInputKeymap) Option {
    return func(c *ChatInput) {
        c.keymap = km
    }
}

func (c *ChatInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, c.keymap.Submit):
            return c, c.submit()
        case key.Matches(msg, c.keymap.ToggleAutoAccept):
            c.autoAccept = !c.autoAccept
            return c, nil
        case key.Matches(msg, c.keymap.TogglePlanMode):
            c.planMode = !c.planMode
            return c, nil
        }
    }
    // ...
}
```
