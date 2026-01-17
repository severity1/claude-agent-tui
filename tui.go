// Package tui provides reusable terminal UI components for building Claude Code
// interfaces using the Charmbracelet ecosystem.
//
// This library offers 20+ Bubble Tea components that integrate with
// claude-agent-sdk-go, providing a complete toolkit for building
// interactive Claude agent terminals.
//
// Components are organized into three categories:
//   - output: Display components (streamtext, message, thinking, tooluse, etc.)
//   - input: Input components (chatinput, permission, question, confirm, etc.)
//   - shared: Shared primitives (button, chip, card, overlay)
//
// The adapter package bridges the SDK's streaming events to Bubble Tea messages,
// enabling seamless integration between the SDK and TUI components.
package tui
