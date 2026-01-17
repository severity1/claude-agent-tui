// Package chat provides a complete chat layout composing multiple components
// into a full chat interface.
//
// # Components
//
// The layout includes:
//   - Message history viewport with scrolling
//   - Chat input area (chatinput component)
//   - Status bar with token/cost info (status component)
//   - Permission prompts overlay (permission component)
//   - Question prompts (question component)
//   - Plan mode prompts (planenter, planexit components)
//
// # Features
//
// The layout:
//   - Integrates with SDK via adapter layer (StreamCmd, ToolControlAdapter)
//   - Manages input state transitions via InputState state machine
//   - Implements focus management between components
//   - Supports accessibility via WithAccessibility() option
//
// # Implementation
//
// The layout implements tea.Model and coordinates all child components.
// See docs/ARCHITECTURE.md for complete specification.
package chat
