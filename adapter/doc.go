// Package adapter provides integration between claude-agent-sdk-go and
// Bubble Tea components.
//
// The adapter package bridges the SDK's streaming events to Bubble Tea messages,
// enabling seamless integration between the SDK and TUI components.
//
// Key adapters:
//   - StreamAdapter: Converts SDK StreamEvents to tea.Msg
//   - ControlAdapter: Handles canUseTool callbacks via prompts
//   - ClientAdapter: Manages SDK client lifecycle
package adapter
