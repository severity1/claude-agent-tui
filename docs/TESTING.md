# Testing Strategy

Comprehensive testing approach for `claude-agent-tui`.

## Overview

Testing a TUI library requires multiple strategies:

1. **Unit tests** - Individual component logic
2. **Integration tests** - Component interactions via teatest
3. **Visual tests** - Golden file output verification via teatest
4. **Performance tests** - Benchmarks for View() and Update()
5. **Mock tests** - SDK integration without real API
6. **Manual tests** - Real-world usage scenarios

## Dependencies

```go
import (
    "testing"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/x/exp/teatest"
    "github.com/charmbracelet/lipgloss"
    "github.com/muesli/termenv"
)
```

Install teatest:
```bash
go get github.com/charmbracelet/x/exp/teatest@latest
```

## Test Environment Setup

**Critical for CI consistency:** Force ASCII color profile to ensure golden files match across environments.

```go
// testutil/setup.go
package testutil

import (
    "github.com/charmbracelet/lipgloss"
    "github.com/muesli/termenv"
)

func init() {
    // Force consistent output across all environments
    lipgloss.SetColorProfile(termenv.Ascii)
}

// IMPORTANT: ASCII mode strips all ANSI styling.
// Golden files created this way do NOT test:
// - Color correctness
// - Bold/italic/underline rendering
// - Theme color application
//
// For color testing, consider:
// - Separate color tests with termenv.TrueColor profile
// - Visual regression testing with screenshots
// - Manual review of color themes
```

Add to `.gitattributes` to prevent line ending issues:
```
*.golden -text
```

---

## Unit Testing

### Component Logic Tests

Test component behavior without rendering:

```go
// component/streamtext/streamtext_test.go
package streamtext

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestStreamText_Append(t *testing.T) {
    m := New()

    m.Append("hello ")
    m.Append("world")

    if got := m.Content(); got != "hello world" {
        t.Errorf("Content() = %q, want %q", got, "hello world")
    }
}

func TestStreamText_Clear(t *testing.T) {
    m := New()
    m.Append("content")
    m.Clear()

    if got := m.Content(); got != "" {
        t.Errorf("Content() after Clear() = %q, want empty", got)
    }
}

func TestStreamText_StreamingCursor(t *testing.T) {
    m := New()
    m.SetStreaming(true)

    view := m.View()
    if !strings.Contains(view, CursorChar) {
        t.Error("View() should contain cursor when streaming")
    }

    m.SetStreaming(false)
    view = m.View()
    if strings.Contains(view, CursorChar) {
        t.Error("View() should not contain cursor when not streaming")
    }
}
```

### Input Component Tests

```go
// component/chatinput/chatinput_test.go
package chatinput

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestChatInput_Submit(t *testing.T) {
    m := New()

    // Simulate typing
    for _, r := range "test message" {
        m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }

    // Press enter
    var cmd tea.Cmd
    m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Check for SubmitMsg
    if cmd == nil {
        t.Fatal("expected command from submit")
    }

    msg := cmd()
    submitMsg, ok := msg.(SubmitMsg)
    if !ok {
        t.Fatalf("expected SubmitMsg, got %T", msg)
    }

    if submitMsg.Text != "test message" {
        t.Errorf("SubmitMsg.Text = %q, want %q", submitMsg.Text, "test message")
    }
}

func TestChatInput_EmptySubmit(t *testing.T) {
    m := New()

    // Press enter with empty input
    m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    if cmd != nil {
        t.Error("empty submit should not produce command")
    }
}

func TestChatInput_AltEnterNewline(t *testing.T) {
    m := New()

    // Type some text
    for _, r := range "line1" {
        m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }

    // Alt+Enter for newline
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter, Alt: true})

    // Type more
    for _, r := range "line2" {
        m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }

    value := m.Value()
    if !strings.Contains(value, "\n") {
        t.Error("Alt+Enter should insert newline")
    }
}
```

### Permission Component Tests

```go
// component/permission/permission_test.go
package permission

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestPermission_Shortcuts(t *testing.T) {
    tests := []struct {
        name        string
        key         string
        wantAllow   bool
        wantAlways  bool
    }{
        {"y allows", "y", true, false},
        {"a allows", "a", true, false},
        {"n denies", "n", false, false},
        {"d denies", "d", false, false},
        {"A always allows", "A", true, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := New()
            var result ResultMsg

            m.Show("Bash", map[string]any{"command": "ls"}, func(a, aa bool) {
                result = ResultMsg{Allow: a, AlwaysAllow: aa}
            })

            // Press the key
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
            m, _ = m.Update(msg)

            if result.Allow != tt.wantAllow {
                t.Errorf("Allow = %v, want %v", result.Allow, tt.wantAllow)
            }
            if result.AlwaysAllow != tt.wantAlways {
                t.Errorf("AlwaysAllow = %v, want %v", result.AlwaysAllow, tt.wantAlways)
            }
        })
    }
}

func TestPermission_Navigation(t *testing.T) {
    m := New()
    m.Show("Bash", map[string]any{}, func(bool, bool) {})

    // Default selection is 0 (Allow)
    if m.selected != 0 {
        t.Errorf("initial selection = %d, want 0", m.selected)
    }

    // Navigate down
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
    if m.selected != 1 {
        t.Errorf("after down selection = %d, want 1", m.selected)
    }

    // Navigate down again
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
    if m.selected != 2 {
        t.Errorf("after second down selection = %d, want 2", m.selected)
    }

    // Can't go past last option
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
    if m.selected != 2 {
        t.Errorf("selection should stay at 2, got %d", m.selected)
    }
}
```

---

## teatest Integration Testing

The official way to test Bubble Tea programs. Use `github.com/charmbracelet/x/exp/teatest`.

### Basic teatest Pattern

```go
// layout/chat_teatest_test.go
package layout

import (
    "bytes"
    "io"
    "testing"
    "time"

    "github.com/charmbracelet/x/exp/teatest"
)

func TestChatLayout_FullFlow(t *testing.T) {
    m := New()

    // Create test model with fixed terminal size
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(80, 24),
    )

    // Type a message
    tm.Type("Hello Claude")

    // Wait for input to appear in output
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Hello Claude"))
    }, teatest.WithDuration(2*time.Second))

    // Send enter key to submit
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Wait for streaming indicator
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("█")) // Streaming cursor
    }, teatest.WithDuration(2*time.Second))

    // Get final output for golden file comparison
    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    tm.WaitFinished(t, teatest.WithFinalTimeout(5*time.Second))
}
```

### Golden File Testing with teatest

```go
func TestChatLayout_Output_Golden(t *testing.T) {
    m := New()
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(80, 24),
    )

    // Simulate some interactions
    tm.Type("test input")
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Quit and get final output
    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

    out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(5*time.Second)))
    if err != nil {
        t.Fatal(err)
    }

    // Compare against golden file (creates testdata/TestChatLayout_Output_Golden.golden)
    teatest.RequireEqualOutput(t, out)
}
```

Update golden files:
```bash
go test ./... -update
```

### Testing Final Model State

```go
func TestChatLayout_FinalState(t *testing.T) {
    m := New()
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(80, 24),
    )

    // Type and submit
    tm.Type("Hello")
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Quit
    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

    // Get final model state
    fm := tm.FinalModel(t, teatest.WithFinalTimeout(5*time.Second))
    finalModel, ok := fm.(Model)
    if !ok {
        t.Fatalf("expected Model, got %T", fm)
    }

    // Assert internal state
    if finalModel.streaming {
        t.Error("should not be streaming after quit")
    }
}
```

### WaitFor with Custom Intervals

```go
func TestPermission_Appears(t *testing.T) {
    m := New()
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Trigger permission prompt
    tm.Send(adapter.ShowPermissionMsg{
        ToolName:  "Bash",
        ToolInput: map[string]any{"command": "ls -la"},
        Respond:   func(bool, bool) {},
    })

    // Wait for permission dialog with custom timing
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Allow Bash?")) &&
               bytes.Contains(bts, []byte("ls -la"))
    },
        teatest.WithDuration(3*time.Second),
        teatest.WithCheckInterval(50*time.Millisecond),
    )

    // Press 'y' to allow
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

    // Verify dialog dismissed
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return !bytes.Contains(bts, []byte("Allow Bash?"))
    }, teatest.WithDuration(time.Second))

    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

### Testing Streaming Responses

```go
func TestStreamText_StreamingFlow(t *testing.T) {
    m := New()
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Simulate stream start
    tm.Send(adapter.StreamBlockStartMsg{BlockType: "text", Index: 0})

    // Send deltas
    tm.Send(adapter.StreamTextMsg{Text: "Hello "})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Hello"))
    }, teatest.WithDuration(time.Second))

    tm.Send(adapter.StreamTextMsg{Text: "world!"})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Hello world!"))
    }, teatest.WithDuration(time.Second))

    // Stream complete
    tm.Send(adapter.StreamDoneMsg{})

    // Verify cursor gone
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return !bytes.Contains(bts, []byte("█"))
    }, teatest.WithDuration(time.Second))

    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

---

## Performance Testing

Benchmark critical paths to catch performance regressions.

### View() Benchmarks

```go
// component/streamtext/streamtext_bench_test.go
package streamtext

import "testing"

func BenchmarkStreamText_View_Empty(b *testing.B) {
    m := New()
    m.SetWidth(80)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

func BenchmarkStreamText_View_Small(b *testing.B) {
    m := New()
    m.SetWidth(80)
    m.Append("Hello, world!")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

func BenchmarkStreamText_View_Large(b *testing.B) {
    m := New()
    m.SetWidth(80)
    // Simulate large response (10KB)
    content := strings.Repeat("Lorem ipsum dolor sit amet. ", 400)
    m.Append(content)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

func BenchmarkStreamText_View_VeryLarge(b *testing.B) {
    m := New()
    m.SetWidth(80)
    // Simulate very large response (100KB)
    content := strings.Repeat("Lorem ipsum dolor sit amet. ", 4000)
    m.Append(content)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}
```

### Update() Benchmarks

```go
func BenchmarkStreamText_Update_Delta(b *testing.B) {
    m := New()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Append("x")
    }
}

func BenchmarkChatInput_Update_Typing(b *testing.B) {
    m := chatinput.New()
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m, _ = m.Update(msg)
    }
}

func BenchmarkLayout_Update_StreamDelta(b *testing.B) {
    m := layout.New()
    msg := adapter.StreamTextMsg{Text: "x"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m, _ = m.Update(msg)
    }
}
```

### Memory Allocation Benchmarks

```go
func BenchmarkStreamText_Append_Allocs(b *testing.B) {
    m := New()

    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Append("chunk of text ")
    }
}

func BenchmarkLayout_View_Allocs(b *testing.B) {
    m := layout.New()
    // Pre-populate with some content
    m.stream.Append(strings.Repeat("text ", 100))

    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run with memory stats
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkStreamText ./component/streamtext/

# Compare benchmarks (requires benchstat)
go test -bench=. -count=10 ./... > old.txt
# Make changes
go test -bench=. -count=10 ./... > new.txt
benchstat old.txt new.txt
```

### Performance Targets

| Operation | Target | Notes |
|-----------|--------|-------|
| StreamText.View() (small) | < 100μs | Normal message |
| StreamText.View() (large) | < 1ms | 10KB content |
| StreamText.Append() | < 1μs | Single delta |
| Layout.View() | < 5ms | Full render |
| Layout.Update() | < 100μs | Message routing |

### CI Performance Regression

```yaml
# .github/workflows/bench.yml
name: Benchmarks

on:
  pull_request:
    branches: [main]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run benchmarks on PR
        run: go test -bench=. -benchmem -count=5 ./... > new.txt

      - name: Checkout main
        run: git checkout main

      - name: Run benchmarks on main
        run: go test -bench=. -benchmem -count=5 ./... > old.txt

      - name: Compare benchmarks
        run: |
          go install golang.org/x/perf/cmd/benchstat@latest
          benchstat old.txt new.txt
```

---

## Integration Testing

Test component interactions in layout.

```go
// layout/chat_test.go
package layout

func TestChatLayout_StateTransitions(t *testing.T) {
    m := New()

    // Initial state: chat input
    if m.state != InputStateChat {
        t.Errorf("initial state = %v, want InputStateChat", m.state)
    }

    // Trigger permission prompt
    m, _ = m.Update(adapter.ShowPermissionMsg{
        ToolName:  "Bash",
        ToolInput: map[string]any{"command": "ls"},
        Respond:   func(bool, bool) {},
    })

    if m.state != InputStatePermission {
        t.Errorf("state after permission = %v, want InputStatePermission", m.state)
    }

    // Resolve permission
    m, _ = m.Update(permission.ResultMsg{Allow: true})

    if m.state != InputStateChat {
        t.Errorf("state after permission resolved = %v, want InputStateChat", m.state)
    }
}

func TestChatLayout_FocusManagement(t *testing.T) {
    m := New()

    // Input should be focused initially
    if !m.input.Focused() {
        t.Error("input should be focused initially")
    }

    // During streaming, input should blur
    m, _ = m.Update(chatinput.SubmitMsg{Text: "test"})
    if m.input.Focused() {
        t.Error("input should blur during streaming")
    }

    // After stream done, input should refocus
    m, _ = m.Update(adapter.StreamDoneMsg{})
    if !m.input.Focused() {
        t.Error("input should refocus after streaming")
    }
}
```

---

## Keyboard Navigation Tests

```go
func TestKeyboardNavigation(t *testing.T) {
    m := New()

    keys := []struct {
        key      tea.KeyMsg
        expected string
    }{
        {tea.KeyMsg{Type: tea.KeyTab}, "next element"},
        {tea.KeyMsg{Type: tea.KeyShiftTab}, "prev element"},
        {tea.KeyMsg{Type: tea.KeyUp}, "navigate up"},
        {tea.KeyMsg{Type: tea.KeyDown}, "navigate down"},
    }

    for _, k := range keys {
        t.Run(k.expected, func(t *testing.T) {
            _, cmd := m.Update(k.key)
            // Verify appropriate action taken
        })
    }
}
```

---

## Manual Testing Checklist

### MVP Checklist

- [ ] `go run example/chat/main.go` starts without error
- [ ] Typing in input shows characters
- [ ] Enter submits message
- [ ] Alt+Enter inserts newline
- [ ] Streaming response shows cursor
- [ ] Streaming completes, cursor disappears
- [ ] Ctrl+C during stream stops it
- [ ] Ctrl+C without stream quits app
- [ ] Permission prompt appears for Bash tool
- [ ] 'y' allows and continues
- [ ] 'n' denies
- [ ] 'A' always allows
- [ ] Arrow keys navigate permission options
- [ ] Enter selects current option

### Theme Testing

- [ ] Default theme renders correctly
- [ ] Theme switch preserves input text
- [ ] Colors visible in 256-color terminal
- [ ] Colors visible in true-color terminal
- [ ] Fallback works in 16-color terminal

### Accessibility Testing

- [ ] `REDUCE_MOTION=1` disables animations
- [ ] `HIGH_CONTRAST=1` uses high contrast theme
- [ ] All features work keyboard-only
- [ ] Focus ring visible on all interactive elements
- [ ] Status indicators have symbols, not just colors

### Responsive Testing

- [ ] Renders in 40x15 terminal
- [ ] Renders in 80x24 terminal
- [ ] Renders in 160x50 terminal
- [ ] Resize during operation doesn't crash
- [ ] "Too small" message at <30x10

---

## CI Configuration

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test -v -race ./...

      - name: Run tests with coverage
        run: go test -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out

  teatest:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run teatest integration tests
        run: go test -v -tags=integration ./...

      - name: Verify golden files are up to date
        run: |
          go test ./... -update
          git diff --exit-code -- '*.golden' || (echo "Golden files out of date. Run 'go test ./... -update'" && exit 1)

  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run benchmarks
        run: go test -bench=. -benchmem -run=^$ ./... | tee bench.txt

      - name: Check for performance regression
        run: |
          # Fail if any benchmark is significantly slower than target
          # This is a simple check - consider benchstat for more accurate comparison
          if grep -E "ns/op|allocs/op" bench.txt | grep -E "[0-9]{8,}"; then
            echo "Warning: Some benchmarks may be too slow"
          fi

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build
        run: go build ./...

      - name: Build example
        run: go build ./example/chat
```

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

set -e

echo "Running fmt..."
go fmt ./...

echo "Running vet..."
go vet ./...

echo "Running short tests..."
go test -short ./...

echo "Checking golden files..."
go test ./... -update
if ! git diff --quiet -- '*.golden'; then
    echo "Golden files changed. Please review and commit."
    git diff --stat -- '*.golden'
    exit 1
fi

echo "All checks passed!"
```

---

## Test Organization

```
claude-agent-tui/
├── component/
│   ├── streamtext/
│   │   ├── streamtext.go
│   │   ├── streamtext_test.go
│   │   └── testdata/
│   │       └── golden/
│   ├── chatinput/
│   │   ├── chatinput.go
│   │   └── chatinput_test.go
│   └── permission/
│       ├── permission.go
│       └── permission_test.go
├── adapter/
│   ├── stream.go
│   ├── stream_test.go
│   ├── control.go
│   ├── control_test.go
│   └── mock_test.go
├── layout/
│   ├── chat.go
│   └── chat_test.go
└── testutil/
    ├── mock.go          # Shared mock helpers
    └── assert.go        # Custom assertions
```

---

## Coverage Goals

| Package | Unit Test | teatest Integration | Benchmark |
|---------|-----------|---------------------|-----------|
| component/* | 80% | Yes (golden files) | View(), Update() |
| adapter/* | 90% | Yes (mock flow) | StreamCmd() |
| layout/* | 70% | Yes (full flow) | Full render |
| system/* | 60% | No | Theme apply |

### Coverage Priorities

1. **adapter/** - Highest priority. This is the SDK integration layer. If it breaks, nothing works.
2. **component/** - High priority. Each component needs unit tests for logic + golden files for output.
3. **layout/** - Medium priority. Integration tests via teatest cover most of this.
4. **system/** - Lower priority. Theme/keymap are configuration, less likely to regress.

### What to Test with teatest

| Scenario | Test Type | Example |
|----------|-----------|---------|
| User types and submits | Full flow | `TestChatLayout_SubmitMessage` |
| Permission prompt flow | Interaction | `TestPermission_AllowDeny` |
| Streaming response | Timed | `TestStreaming_DeltaAccumulation` |
| Keyboard navigation | Input | `TestKeyboardNavigation_TabCycle` |
| Component rendering | Golden | `TestChatLayout_Output_Golden` |

### What to Benchmark

| Component | Benchmark | Why |
|-----------|-----------|-----|
| StreamText.View() | Large content | Renders on every frame |
| StreamText.Append() | Frequent | Called per delta |
| Layout.View() | Full render | Most expensive operation |
| Layout.Update() | Routing | Called on every message |

Focus on adapter layer - it's the glue between SDK and UI. Get that right.

---

## SDK Integration Notes

**Assumption:** Claude Code is installed and working locally. No mocking the SDK.

### Real SDK Testing

All integration tests use the real `severity1/claude-agent-sdk-go` with actual Claude responses:

```go
// layout/chat_integration_test.go
func TestChatLayout_RealConversation(t *testing.T) {
    m := layout.New()
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

    // Real query to Claude
    tm.Type("Respond with exactly: TEST_OK")
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // Wait for real response (longer timeout for API)
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("TEST_OK"))
    }, teatest.WithDuration(30*time.Second))

    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    tm.WaitFinished(t, teatest.WithFinalTimeout(5*time.Second))
}
```

### When to Use Real SDK

| Test | Real SDK? | Timeout |
|------|-----------|---------|
| Component unit tests | No (no SDK) | 1s |
| Adapter message parsing | Yes | 30s |
| Layout integration | Yes | 60s |
| Permission flow | Yes | 30s |
| Streaming behavior | Yes | 60s |
| Golden file generation | Yes | 60s |

### Test Timeouts

Real API calls need longer timeouts:

```go
// Short timeout for unit tests
teatest.WithDuration(time.Second)

// Long timeout for real SDK
teatest.WithDuration(30 * time.Second)

// Very long for complex flows
teatest.WithFinalTimeout(60 * time.Second)
```

### CI Requirements

CI runners need Claude Code installed.

**IMPORTANT:** This is a gap in our setup. Options to address:

1. **Self-hosted runner** - Pre-install Claude Code on your own runner
2. **Mock-only CI** - Skip SDK tests in GitHub Actions, run locally
3. **Setup script** - Create a CI setup script (not yet implemented)

For now, SDK integration tests should be run locally before PR merge.

```yaml
jobs:
  unit-tests:
    runs-on: ubuntu-latest  # Unit tests don't need Claude Code
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test -v -short ./...

  integration:
    runs-on: self-hosted  # Requires Claude Code pre-installed
    steps:
      - name: Verify Claude Code
        run: |
          claude --version || (echo "Claude Code not installed" && exit 1)

      - name: Run integration tests
        run: go test -v -run Integration ./...
        timeout-minutes: 15
```

**TODO:** Create a setup script or Docker image with Claude Code for CI.

### Capturing Real Message Formats

Use real SDK to capture actual message structures for documentation:

```go
func TestCapture_StreamEventFormat(t *testing.T) {
    // Connect to real SDK
    client := claudecode.NewClient()
    ctx := context.Background()
    client.Connect(ctx)
    defer client.Disconnect(ctx)

    // Query and capture
    client.Query(ctx, "Say hi")
    for msg := range client.ReceiveMessages(ctx) {
        // Log actual message structure
        t.Logf("Message type: %T", msg)
        t.Logf("Message: %+v", msg)
    }
}
```

### Benefits

1. **No mock drift** - Tests fail when SDK changes (catches bugs early)
2. **Real golden files** - Capture actual Claude output formatting
3. **True integration** - Permission prompts, streaming, all real
4. **Documentation accuracy** - Message formats match reality
