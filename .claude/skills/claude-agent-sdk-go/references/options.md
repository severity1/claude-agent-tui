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
