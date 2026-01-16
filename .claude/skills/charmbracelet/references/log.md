# Charmbracelet Log Reference

**Package**: `github.com/charmbracelet/log`
**Documentation**: https://github.com/charmbracelet/log
**Go Doc**: https://pkg.go.dev/github.com/charmbracelet/log

> For the latest API changes, use WebFetch on the GitHub URL above.

## Table of Contents
- [Overview](#overview)
- [Logger Creation](#logger-creation)
- [Log Levels](#log-levels)
- [Structured Logging](#structured-logging)
- [Formatters](#formatters)
- [Styling](#styling)
- [TUI Integration](#tui-integration)

## Overview

Charmbracelet log is a structured logging library with beautiful terminal output. It uses Lip Gloss for styling and supports multiple output formats.

## Logger Creation

### Global Logger
```go
import "github.com/charmbracelet/log"

// Use default global logger
log.Info("Starting application")
log.Debug("Debug info", "key", "value")
```

### Custom Logger
```go
import "github.com/charmbracelet/log"

// Basic custom logger
logger := log.New(os.Stderr)

// Logger with options
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportCaller:    true,
    ReportTimestamp: true,
    TimeFormat:      time.Kitchen,
    Prefix:          "myapp",
    Level:           log.DebugLevel,
})
```

### File Logging (for TUI apps)
```go
// Log to file to avoid interfering with TUI
f, _ := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
logger := log.New(f)
```

## Log Levels

### Available Levels (low to high priority)
```go
log.DebugLevel  // Detailed debugging information
log.InfoLevel   // General operational information
log.WarnLevel   // Warning conditions
log.ErrorLevel  // Error conditions
log.FatalLevel  // Fatal errors (calls os.Exit)
```

### Logging Functions
```go
// Standard functions
log.Debug("message", "key", value)
log.Info("message", "key", value)
log.Warn("message", "key", value)
log.Error("message", "key", value)
log.Fatal("message", "key", value)  // Exits with code 1

// Format functions
log.Debugf("user %s logged in", username)
log.Infof("processing %d items", count)
log.Warnf("retry attempt %d", attempt)
log.Errorf("failed to connect: %v", err)

// Print (always outputs regardless of level)
log.Print("always shown")
```

### Setting Log Level
```go
// Global
log.SetLevel(log.DebugLevel)

// Per logger
logger.SetLevel(log.WarnLevel)
```

## Structured Logging

### Key-Value Pairs
```go
log.Info("user action",
    "user_id", 123,
    "action", "login",
    "ip", "192.168.1.1",
)
```

### Sub-loggers with Context
```go
// Create sub-logger with persistent fields
userLogger := logger.With("user_id", 123, "session", "abc")

// All subsequent logs include user_id and session
userLogger.Info("page viewed", "page", "/home")
userLogger.Info("clicked button", "button", "submit")
```

### Helper for Caller Skipping
```go
// For wrapper functions
func myLog(msg string) {
    log.Helper()  // Skip this frame in caller reporting
    log.Info(msg)
}
```

## Formatters

### Text Formatter (Default)
```go
logger := log.New(os.Stderr)
// Output: INFO Starting server port=8080
```

### JSON Formatter
```go
logger := log.NewWithOptions(os.Stderr, log.Options{
    Formatter: log.JSONFormatter,
})
// Output: {"level":"info","msg":"Starting server","port":8080}
```

### Logfmt Formatter
```go
logger := log.NewWithOptions(os.Stderr, log.Options{
    Formatter: log.LogfmtFormatter,
})
// Output: level=info msg="Starting server" port=8080
```

## Styling

### Custom Styles
```go
styles := log.DefaultStyles()

// Customize level colors
styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
    SetString("ERROR").
    Bold(true).
    Foreground(lipgloss.Color("204"))

// Customize key style
styles.Key = lipgloss.NewStyle().
    Foreground(lipgloss.Color("99"))

// Customize value style
styles.Value = lipgloss.NewStyle().
    Foreground(lipgloss.Color("248"))

logger.SetStyles(styles)
```

### Disable Styling
```go
// For non-TTY output (pipes, files)
// Styling is automatically disabled when output is not a TTY
```

## TUI Integration

### Logging to File
```go
// In TUI apps, log to a file to avoid corrupting the display
type Model struct {
    logger *log.Logger
}

func New() Model {
    f, _ := os.CreateTemp("", "app-*.log")
    return Model{
        logger: log.New(f),
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.logger.Debug("received message", "type", fmt.Sprintf("%T", msg))
    // ...
}
```

### Debug Mode Toggle
```go
type Model struct {
    logger *log.Logger
    debug  bool
}

func (m *Model) SetDebug(enabled bool) {
    m.debug = enabled
    if enabled {
        m.logger.SetLevel(log.DebugLevel)
    } else {
        m.logger.SetLevel(log.InfoLevel)
    }
}
```

### Slog Integration
```go
import (
    "log/slog"
    "github.com/charmbracelet/log"
)

// Use charmbracelet/log as slog handler
handler := log.New(os.Stderr)
slog.SetDefault(slog.New(handler))

// Now use standard slog API
slog.Info("message", "key", "value")
```

## Best Practices

1. **Log to file in TUI apps** to avoid display corruption
2. **Use structured logging** with key-value pairs for queryability
3. **Create sub-loggers** for components with persistent context
4. **Set appropriate levels** (Debug in dev, Info/Warn in prod)
5. **Disable timestamps** in TUI output for cleaner display
6. **Use JSON formatter** for log aggregation systems
