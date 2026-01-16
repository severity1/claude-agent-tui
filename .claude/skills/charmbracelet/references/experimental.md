# Charmbracelet Experimental (x/) Reference

**Repository**: `github.com/charmbracelet/x`
**Documentation**: https://github.com/charmbracelet/x
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/x

> Experimental packages may have unstable APIs. Check the repository for current status.

## Table of Contents
- [Overview](#overview)
- [ANSI Utilities](#ansi-utilities)
- [Term Utilities](#term-utilities)
- [Other Utilities](#other-utilities)

## Overview

The `x/` repository contains experimental packages that may graduate to stable releases. These provide low-level terminal utilities and helpers for TUI development.

## ANSI Utilities

**Package**: `github.com/charmbracelet/x/ansi`

### Sequence Parsing
```go
import "github.com/charmbracelet/x/ansi"

// Parse ANSI escape sequences
seq := ansi.Parse(input)

// Check sequence type
if seq.IsCSI() {
    // Control Sequence Introducer
}
if seq.IsOSC() {
    // Operating System Command
}
```

### Sequence Building
```go
// Build ANSI sequences programmatically
seq := ansi.CSI("m", "1", "31")  // Bold red: \x1b[1;31m

// Common sequences
ansi.CursorUp(n)
ansi.CursorDown(n)
ansi.CursorForward(n)
ansi.CursorBack(n)
ansi.CursorPosition(row, col)
ansi.EraseDisplay(mode)
ansi.EraseLine(mode)
```

### String Width
```go
// Calculate display width (handles wide chars, emoji)
width := ansi.StringWidth(str)

// Truncate to width
truncated := ansi.Truncate(str, maxWidth, "...")
```

## Term Utilities

**Package**: `github.com/charmbracelet/x/term`

### Terminal Detection
```go
import "github.com/charmbracelet/x/term"

// Check if running in a terminal
if term.IsTerminal(os.Stdout.Fd()) {
    // Enable rich output
}

// Get terminal size
width, height, err := term.GetSize(os.Stdout.Fd())
```

### Raw Mode
```go
// Save current state and enable raw mode
oldState, err := term.MakeRaw(os.Stdin.Fd())
if err != nil {
    return err
}
defer term.Restore(os.Stdin.Fd(), oldState)

// Terminal is now in raw mode
// Read input character by character
```

### Color Support Detection
```go
// Check color profile
profile := term.ColorProfile()
switch profile {
case term.Ascii:
    // No color support
case term.ANSI:
    // 16 colors
case term.ANSI256:
    // 256 colors
case term.TrueColor:
    // 16 million colors
}
```

## Other Utilities

### Input Reading
```go
import "github.com/charmbracelet/x/input"

// Read terminal input with proper escape sequence handling
reader := input.NewReader(os.Stdin)
for {
    ev, err := reader.ReadEvent()
    if err != nil {
        break
    }
    // Handle event
}
```

### Window Title
```go
import "github.com/charmbracelet/x/windows"

// Set terminal window title
windows.SetTitle("My Application")
```

## Compatibility Notes

These packages are experimental and may:
- Have breaking API changes
- Be moved to stable packages
- Be deprecated

Check the x/ repository for current status and migration guides when packages graduate to stable releases.

## Best Practices

1. **Pin versions** carefully in go.mod for experimental packages
2. **Check for stable alternatives** before using x/ packages
3. **Follow repository updates** for breaking changes
4. **Use teatest from x/exp** for testing (more stable than other x/ packages)
5. **Prefer bubbletea/lipgloss** over raw ANSI when possible
