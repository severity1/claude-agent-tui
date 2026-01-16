# Functional Options Reference

The SDK uses the functional options pattern for configuration. All options follow the `With*` naming convention.

## Model Configuration

```go
// Set primary model
claudecode.WithModel("sonnet")      // sonnet, haiku, opus
claudecode.WithModel("opus")

// Set fallback model
claudecode.WithFallbackModel("haiku")
```

## Tool Control

```go
// Allow specific tools
claudecode.WithAllowedTools([]string{"Read", "Write", "Bash"})

// Disallow specific tools
claudecode.WithDisallowedTools([]string{"Write", "Bash"})

// Provide custom MCP tools
claudecode.WithTools(myTools)

// Use tool presets
claudecode.WithToolsPreset("coding")    // Common coding tools
claudecode.WithToolsPreset("readonly")  // Read-only tools
```

## Behavior Control

```go
// System prompt
claudecode.WithSystemPrompt("You are a Go expert")

// Limit conversation turns
claudecode.WithMaxTurns(10)

// Permission mode
claudecode.WithPermissionMode("default")     // Ask for permissions
claudecode.WithPermissionMode("acceptEdits") // Auto-accept edits
claudecode.WithPermissionMode("bypassPermissions") // Accept all
```

## Environment

```go
// Working directory
claudecode.WithCwd("/path/to/project")

// Additional directories
claudecode.WithAddDirs([]string{"/other/dir"})

// Environment variables
claudecode.WithEnv(map[string]string{"KEY": "value"})
claudecode.WithEnvVar("KEY", "value")
```

## Sandboxing

```go
// Enable sandbox mode
claudecode.WithSandbox(true)
```

## MCP Servers

```go
// External MCP servers
claudecode.WithMcpServers(map[string]McpServer{
    "my-server": {Command: "node", Args: []string{"server.js"}},
})

// In-process SDK MCP server
server := claudecode.CreateSDKMcpServer("my-tools", tools)
claudecode.WithSdkMcpServer(server)
```

## Hooks and Callbacks

```go
// Permission callback
claudecode.WithPermissionCallback(func(input CanUseToolInput) PermissionResult {
    if input.Tool == "Write" {
        return PermissionResultDeny("Write not allowed")
    }
    return PermissionResultAllow(nil)
})

// Lifecycle hooks
claudecode.WithHooks(Hooks{
    PreToolUse:  preToolUseFunc,
    PostToolUse: postToolUseFunc,
})
```

## Agents

```go
// Define custom agent
agent := claudecode.AgentDefinition{
    Name:         "my-agent",
    SystemPrompt: "You are a specialized agent",
    Tools:        []string{"Read", "Grep"},
    Model:        claudecode.ModelSonnet,
}

claudecode.WithAgent(agent)
claudecode.WithAgents([]AgentDefinition{agent1, agent2})
```

## File Checkpointing

```go
// Enable file change tracking
claudecode.WithFileCheckpointing(true)
```

## Plugins

```go
// Load plugins
claudecode.WithPlugins([]string{"plugin1", "plugin2"})
```

## TUI-Relevant Options

For TUI integration, commonly used options:

```go
claudecode.WithModel("sonnet"),
claudecode.WithMaxTurns(100),
claudecode.WithPermissionCallback(tuiPermissionHandler),
claudecode.WithCwd(projectDir),
claudecode.WithFileCheckpointing(true),
```

## Complete TUI Setup Example

```go
func createTUIClient(ctx context.Context, projectDir string, adapter *ControlAdapter) (claudecode.Client, error) {
    return claudecode.NewClient(ctx,
        // Model selection
        claudecode.WithModel("sonnet"),
        claudecode.WithFallbackModel("haiku"),

        // Behavior
        claudecode.WithMaxTurns(100),
        claudecode.WithSystemPrompt("You are a helpful assistant working in a TUI."),

        // Environment
        claudecode.WithCwd(projectDir),
        claudecode.WithEnv(map[string]string{
            "EDITOR": "vim",
        }),

        // Tool control
        claudecode.WithAllowedTools([]string{
            "Read", "Write", "Edit", "Bash", "Glob", "Grep",
        }),

        // Permission handling (required for TUI)
        claudecode.WithPermissionCallback(adapter.CanUseTool),

        // File tracking
        claudecode.WithFileCheckpointing(true),

        // Sandbox for safety
        claudecode.WithSandbox(true),
    )
}
```

## Permission Callback Pattern

```go
type ControlAdapter struct {
    promptCh   chan ShowPermissionMsg
    responseCh chan PermissionResponse
    ctx        context.Context
}

func (a *ControlAdapter) CanUseTool(input claudecode.CanUseToolInput) claudecode.PermissionResult {
    // Auto-approve safe tools
    if input.Tool == "Read" || input.Tool == "Glob" || input.Tool == "Grep" {
        return claudecode.PermissionResultAllow(nil)
    }

    // Send to TUI for approval
    select {
    case a.promptCh <- ShowPermissionMsg{
        Tool:      input.Tool,
        Input:     input.Input,
        ToolUseID: input.ToolUseID,
    }:
    case <-a.ctx.Done():
        return claudecode.PermissionResultDeny("Operation canceled")
    }

    // Wait for user response
    select {
    case response := <-a.responseCh:
        if response.Allow {
            return claudecode.PermissionResultAllow(nil)
        }
        return claudecode.PermissionResultDeny(response.Reason)
    case <-time.After(5 * time.Minute):
        return claudecode.PermissionResultDeny("Timed out waiting for permission")
    case <-a.ctx.Done():
        return claudecode.PermissionResultDeny("Operation canceled")
    }
}
```

## Option Ordering and Precedence

Options are applied in order. Later options override earlier ones:

```go
// Second WithModel overrides first
claudecode.WithModel("haiku"),
claudecode.WithModel("sonnet"),  // sonnet wins
```

## Common Option Combinations

### Read-Only Mode

```go
claudecode.WithAllowedTools([]string{"Read", "Glob", "Grep"}),
claudecode.WithPermissionMode("default"),  // Still ask for permission
```

### Full Auto-Accept (Dangerous)

```go
claudecode.WithPermissionMode("bypassPermissions"),
// Warning: Bypasses all permission prompts
```

### Plan Mode Only

```go
claudecode.WithPermissionMode("plan"),
// Only allow read operations and planning
```

### Constrained Environment

```go
claudecode.WithSandbox(true),
claudecode.WithCwd("/safe/directory"),
claudecode.WithAllowedTools([]string{"Read", "Write"}),
claudecode.WithDisallowedTools([]string{"Bash"}),
```

## Edge Cases

### Empty or Nil Options

```go
// Safe - nil options are ignored
claudecode.NewClient(ctx, nil)

// Empty array is valid
claudecode.WithAllowedTools([]string{})
```

### Conflicting Options

```go
// Both specified - DisallowedTools takes precedence
claudecode.WithAllowedTools([]string{"Bash", "Write"}),
claudecode.WithDisallowedTools([]string{"Bash"}),
// Result: Only Write is allowed
```

### Invalid Values

```go
// Invalid model - SDK will return error
claudecode.WithModel("invalid-model")

// Invalid path - SDK will return error at runtime
claudecode.WithCwd("/nonexistent/path")
```

## Live Documentation

For latest options, check:
- SDK source: github.com/severity1/claude-agent-sdk-go
- Option godoc: pkg.go.dev/github.com/severity1/claude-agent-sdk-go#Option
