# Huh Reference

**Package**: `github.com/charmbracelet/huh`
**Documentation**: https://github.com/charmbracelet/huh
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/huh

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Overview](#overview)
- [Field Types](#field-types)
- [Form Creation](#form-creation)
- [Validation](#validation)
- [Bubble Tea Integration](#bubble-tea-integration)
- [Dynamic Fields](#dynamic-fields)

## Overview

Huh is a forms library for terminal applications. It provides pre-built field types for user input and integrates seamlessly with Bubble Tea as a `tea.Model`.

## Field Types

### Input (Single Line)
```go
huh.NewInput().
    Title("Name").
    Placeholder("Enter your name").
    Value(&name)
```

### Text (Multi-Line)
```go
huh.NewText().
    Title("Description").
    Placeholder("Enter description...").
    Value(&description)
```

### Select (Single Choice)
```go
huh.NewSelect[string]().
    Title("Color").
    Options(
        huh.NewOption("Red", "red"),
        huh.NewOption("Blue", "blue"),
        huh.NewOption("Green", "green"),
    ).
    Value(&color)
```

### MultiSelect (Multiple Choices)
```go
huh.NewMultiSelect[string]().
    Title("Features").
    Options(
        huh.NewOption("Auth", "auth"),
        huh.NewOption("Logging", "logging"),
        huh.NewOption("Metrics", "metrics"),
    ).
    Value(&features)
```

### Confirm (Yes/No)
```go
huh.NewConfirm().
    Title("Proceed?").
    Affirmative("Yes").
    Negative("No").
    Value(&proceed)
```

## Form Creation

### Basic Form
```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewInput().Title("Name").Value(&name),
        huh.NewInput().Title("Email").Value(&email),
    ),
)

err := form.Run()
if err != nil {
    // Handle error
}
```

### Multi-Page Form (Groups)
```go
form := huh.NewForm(
    // Page 1
    huh.NewGroup(
        huh.NewInput().Title("Username").Value(&username),
        huh.NewInput().Title("Password").EchoMode(huh.EchoModePassword).Value(&password),
    ),
    // Page 2
    huh.NewGroup(
        huh.NewSelect[string]().Title("Role").Options(...).Value(&role),
    ),
)
```

## Validation

```go
huh.NewInput().
    Title("Username").
    Value(&username).
    Validate(func(s string) error {
        if len(s) < 3 {
            return errors.New("username must be at least 3 characters")
        }
        if s == "admin" {
            return errors.New("username 'admin' is reserved")
        }
        return nil
    })
```

## Bubble Tea Integration

Huh forms implement `tea.Model`, making them composable:

```go
type Model struct {
    form *huh.Form
    done bool
}

func New() Model {
    return Model{
        form: huh.NewForm(
            huh.NewGroup(
                huh.NewInput().Title("Name").Value(&name),
            ),
        ),
    }
}

func (m Model) Init() tea.Cmd {
    return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }

    if m.form.State == huh.StateCompleted {
        m.done = true
    }

    return m, cmd
}

func (m Model) View() string {
    if m.done {
        return "Form completed!"
    }
    return m.form.View()
}
```

### Accessing Values
```go
// After form completes
name := m.form.GetString("name")
count := m.form.GetInt("count")
```

## Dynamic Fields

Create dependent fields that update based on other values:

```go
var country string
var state string

huh.NewSelect[string]().
    Title("Country").
    Options(
        huh.NewOption("USA", "usa"),
        huh.NewOption("Canada", "canada"),
    ).
    Value(&country),

huh.NewSelect[string]().
    Title("State/Province").
    OptionsFunc(func() []huh.Option[string] {
        if country == "usa" {
            return []huh.Option[string]{
                huh.NewOption("California", "ca"),
                huh.NewOption("Texas", "tx"),
            }
        }
        return []huh.Option[string]{
            huh.NewOption("Ontario", "on"),
            huh.NewOption("Quebec", "qc"),
        }
    }, &country). // Re-evaluate when country changes
    Value(&state)
```

## Standalone Field Usage

Run individual fields without a form:

```go
var name string
err := huh.NewInput().
    Title("What's your name?").
    Value(&name).
    Run()
```

## Best Practices

1. **Use Groups** for multi-step forms with logical sections
2. **Validate early** to provide immediate feedback
3. **Use Value pointers** to capture results
4. **Handle StateCompleted** in Bubble Tea integration
5. **Use OptionsFunc** for dynamic option lists
