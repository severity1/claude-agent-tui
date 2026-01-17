// Package permission provides a tool permission prompt component.
//
// SDK Mapping: canUseTool callback
//
// # Features
//
// The component supports:
//   - Displaying tool invocation details with JSON-formatted input
//   - Accept/Reject/Always buttons
//   - Danger detection: highlights destructive commands (rm -rf, sudo, etc.)
//   - Keyboard shortcuts: y/a (allow), n/d (deny), A (always allow), j/k (navigate)
//   - Focus management via Focus(), Blur(), Focused()
//   - Theming and style customization via WithStyles(), WithTheme()
//   - Multiple display variants via WithVariant(): Modal, Overlay, Inline
//   - Mouse support: click options, hover highlight
//
// # Messages
//
// The component emits:
//   - ResultMsg: Contains decision ("allow"/"deny") and always-allow flag
//
// # Implementation
//
// The component implements tea.Model and the FocusableComponent interface.
// See docs/COMPONENTS.md for complete specification.
package permission
