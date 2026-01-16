# teatest Reference

**Package**: `github.com/charmbracelet/x/exp/teatest`
**Documentation**: https://github.com/charmbracelet/x
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/x/exp/teatest

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Basic Usage](#basic-usage) - Creating test models
- [TestModel API](#testmodel-api) - Send, Output, Quit, WaitFinished
- [Sending Key Messages](#sending-key-messages) - Keyboard input simulation
- [Waiting for Output](#waiting-for-output) - WaitFor with conditions
- [Testing Navigation](#testing-navigation) - Cursor and selection tests
- [Testing Forms](#testing-forms) - Input and submission tests
- [Testing Async Operations](#testing-async-operations) - Loading state tests
- [Testing Window Resize](#testing-window-resize) - Responsive layout tests
- [Golden File Testing](#golden-file-testing) - Snapshot testing
- [Testing Final State](#testing-final-state) - FinalModel assertions
- [Testing Error States](#testing-error-states) - Error handling tests
- [Testing Patterns](#testing-patterns) - Table-driven tests, helpers

## Overview

teatest provides testing utilities for Bubble Tea applications. It allows you to send messages to your model, capture output, and verify behavior without manual testing.

## Installation

```bash
go get github.com/charmbracelet/x/exp/teatest
```

## Basic Usage

```go
import (
    "testing"
    "github.com/charmbracelet/x/exp/teatest"
    tea "github.com/charmbracelet/bubbletea"
)

func TestModel(t *testing.T) {
    m := NewModel()
    tm := teatest.NewTestModel(t, m)

    // Send key messages
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Wait for output to contain expected text
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "expected")
    })

    // Clean up
    tm.Quit()
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

## TestModel API

```go
// Create test model
tm := teatest.NewTestModel(t, model)
tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

// Send messages
tm.Send(msg tea.Msg)

// Get output reader
output := tm.Output()

// Get final model (after quit)
finalModel := tm.FinalModel(t)

// Quit the program
tm.Quit()

// Wait for program to finish
tm.WaitFinished(t)
tm.WaitFinished(t, teatest.WithFinalTimeout(5*time.Second))
```

## Sending Key Messages

```go
// Single keys
tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
tm.Send(tea.KeyMsg{Type: tea.KeyTab})
tm.Send(tea.KeyMsg{Type: tea.KeyUp})
tm.Send(tea.KeyMsg{Type: tea.KeyDown})
tm.Send(tea.KeyMsg{Type: tea.KeyLeft})
tm.Send(tea.KeyMsg{Type: tea.KeyRight})

// Character input
tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

// Control keys
tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
tm.Send(tea.KeyMsg{Type: tea.KeyCtrlD})

// Special keys
tm.Send(tea.KeyMsg{Type: tea.KeySpace})
tm.Send(tea.KeyMsg{Type: tea.KeyBackspace})
tm.Send(tea.KeyMsg{Type: tea.KeyDelete})
```

## Waiting for Output

```go
// Wait for specific text
teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
    return strings.Contains(string(bts), "Loading complete")
})

// With timeout
teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
    return strings.Contains(string(bts), "result")
}, teatest.WithDuration(5*time.Second))

// Wait for multiple conditions
teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
    output := string(bts)
    return strings.Contains(output, "item 1") &&
           strings.Contains(output, "item 2")
})
```

## Testing Navigation

```go
func TestNavigation(t *testing.T) {
    m := NewListModel([]string{"item 1", "item 2", "item 3"})
    tm := teatest.NewTestModel(t, m)

    // Move down
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})

    // Verify cursor moved
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "> item 2")
    })

    // Move down again
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})

    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "> item 3")
    })

    tm.Quit()
    tm.WaitFinished(t)
}
```

## Testing Forms

```go
func TestForm(t *testing.T) {
    m := NewFormModel()
    tm := teatest.NewTestModel(t, m)

    // Type in input field
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("John Doe")})

    // Tab to next field
    tm.Send(tea.KeyMsg{Type: tea.KeyTab})

    // Type in second field
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("john@example.com")})

    // Submit form
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Verify submission
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Form submitted")
    })

    tm.Quit()
    tm.WaitFinished(t)
}
```

## Testing Async Operations

```go
func TestAsyncLoading(t *testing.T) {
    m := NewModelWithAsyncLoad()
    tm := teatest.NewTestModel(t, m)

    // Wait for loading to start
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Loading...")
    })

    // Wait for loading to complete (with longer timeout)
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Loaded 10 items")
    }, teatest.WithDuration(10*time.Second))

    tm.Quit()
    tm.WaitFinished(t)
}
```

## Testing Window Resize

```go
func TestResize(t *testing.T) {
    m := NewModel()
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Send resize message
    tm.Send(tea.WindowSizeMsg{Width: 120, Height: 40})

    // Verify layout adapted
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        // Check for wider layout indicators
        return len(strings.Split(string(bts), "\n")[0]) > 100
    })

    tm.Quit()
    tm.WaitFinished(t)
}
```

## Golden File Testing

```go
import "github.com/charmbracelet/x/exp/golden"

func TestViewGolden(t *testing.T) {
    m := NewModel()
    m.width = 80
    m.height = 24
    m.items = []string{"item 1", "item 2", "item 3"}
    m.cursor = 1

    output := m.View()

    // Compare against golden file
    // Creates/updates file with: go test -update
    golden.RequireEqual(t, []byte(output))
}

func TestMultipleViews(t *testing.T) {
    tests := []struct {
        name   string
        setup  func(*Model)
    }{
        {"empty", func(m *Model) { m.items = nil }},
        {"loaded", func(m *Model) { m.items = []string{"a", "b"} }},
        {"selected", func(m *Model) { m.items = []string{"a"}; m.cursor = 0 }},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := NewModel()
            tt.setup(&m)
            golden.RequireEqual(t, []byte(m.View()))
        })
    }
}
```

## Testing Final State

```go
func TestFinalState(t *testing.T) {
    m := NewModel()
    tm := teatest.NewTestModel(t, m)

    // Perform actions
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")})
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Quit
    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    tm.WaitFinished(t)

    // Get final model and assert state
    final := tm.FinalModel(t)
    model := final.(Model)

    if model.input != "test" {
        t.Errorf("expected input 'test', got '%s'", model.input)
    }
    if !model.submitted {
        t.Error("expected form to be submitted")
    }
}
```

## Testing Error States

```go
func TestErrorHandling(t *testing.T) {
    m := NewModel()
    tm := teatest.NewTestModel(t, m)

    // Trigger error
    tm.Send(triggerErrorMsg{})

    // Verify error displayed
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Error:")
    })

    // Verify error can be dismissed
    tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return !strings.Contains(string(bts), "Error:")
    })

    tm.Quit()
    tm.WaitFinished(t)
}
```

## Testing Patterns

### Table-Driven Tests
```go
func TestKeyBindings(t *testing.T) {
    tests := []struct {
        name     string
        key      tea.KeyMsg
        expected string
    }{
        {"quit", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}, "Goodbye"},
        {"help", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}, "Help:"},
        {"down", tea.KeyMsg{Type: tea.KeyDown}, "> item 2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := NewModel()
            tm := teatest.NewTestModel(t, m)

            tm.Send(tt.key)

            teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
                return strings.Contains(string(bts), tt.expected)
            })

            tm.Quit()
            tm.WaitFinished(t)
        })
    }
}
```

### Test Helpers
```go
func sendKeys(tm *teatest.TestModel, keys string) {
    for _, r := range keys {
        tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }
}

func waitForText(t *testing.T, tm *teatest.TestModel, text string) {
    t.Helper()
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), text)
    }, teatest.WithDuration(2*time.Second))
}
```

## Best Practices

1. **Always call WaitFinished** - Prevents test resource leaks
2. **Use timeouts** - Avoid hanging tests with WithDuration
3. **Test edge cases** - Empty states, errors, boundaries
4. **Use golden files** - For complex view assertions
5. **Clean test names** - Describe what's being tested
6. **Parallel when safe** - Use t.Parallel() for independent tests

## Common Pitfalls

- Forgetting to call `tm.Quit()` and `tm.WaitFinished()`
- WaitFor conditions that never become true (test hangs)
- Not setting initial terminal size for size-dependent views
- Testing implementation details instead of behavior
- Golden files not updated after intentional changes (use `-update`)
