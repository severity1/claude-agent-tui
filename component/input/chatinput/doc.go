// Package chatinput provides the main chat input component with mode toggles.
//
// SDK Mapping: User prompt submission, QueryParams
//
// # Features
//
// The component supports:
//   - Multi-line text input via bubbles textarea
//   - Mode toggles (plan mode, auto-accept) rendered as chips
//   - Keyboard shortcuts: Enter/Ctrl+Enter (submit), Ctrl+A (auto-accept),
//     Ctrl+P (plan mode), Esc (clear)
//   - Focus management via Focus(), Blur(), Focused()
//   - Theming and style customization via WithStyles(), WithTheme()
//   - Multiple display variants via WithVariant(): Fullscreen, Window, View, Overlay
//
// # Messages
//
// The component emits:
//   - SubmitMsg: Contains text, auto-accept state, and plan mode state
//
// # Implementation
//
// The component implements tea.Model and the FocusableComponent interface.
// See docs/COMPONENTS.md for complete specification.
package chatinput
