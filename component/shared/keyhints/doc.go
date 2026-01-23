// Package keyhints provides a reusable component for displaying keybinding hints.
//
// The component supports inline hints with expand/collapse and a modal help window.
//
// # Features
//
// The component supports:
//   - Configurable key bindings with descriptions
//   - Inline collapsed view with "? N more" indicator
//   - Inline expanded view showing all bindings
//   - Modal help window with boxed display
//   - Focus management for keyboard interaction
//   - Customizable styles for keys, descriptions, separators, and borders
//
// # API
//
// Constructor with functional options:
//   - New(bindings []Binding, opts ...Option): Create a new keyhints model
//   - WithSeparator(sep string): Set separator between hints (default " | ")
//   - WithKeyStyle(style lipgloss.Style): Set style for key names
//   - WithDescStyle(style lipgloss.Style): Set style for descriptions
//   - WithSepStyle(style lipgloss.Style): Set style for separators
//   - WithWidth(w int): Set maximum display width
//   - WithDefaultVisible(n int): Set number of hints shown when collapsed
//   - WithCollapsed(): Start in collapsed state
//   - WithFocusStyle(style lipgloss.Style): Set style for focus indicator
//   - WithHelpStyle(style lipgloss.Style): Set style for help window border
//   - WithPadding(padding string): Set left padding for inline views (default 1 space)
//
// Binding management:
//   - SetBindings(bindings []Binding): Replace all bindings
//   - AddBinding(b Binding): Add a single binding
//   - Bindings(): Get current bindings
//
// Expansion control:
//   - Toggle(): Toggle between expanded and collapsed inline view
//   - Expanded(): Check if inline view is expanded
//   - SetExpanded(bool): Set expansion state directly
//
// Focus management:
//   - Focus(): Set focus on the component
//   - Blur(): Remove focus from the component
//   - Focused(): Check if the component has focus
//
// Help window:
//   - ShowHelp(): Open the modal help window
//   - HideHelp(): Close the help window
//   - ToggleHelp(): Toggle help window visibility
//   - ShowingHelp(): Check if help window is visible
//
// # Views
//
// The component provides three view methods:
//   - View(): Inline hints with focus indicator and expand/collapse
//   - ViewCompact(): Compact view showing only key names
//   - ViewHelp(): Modal boxed window showing all bindings
//
// # Keyboard Interaction
//
// When focused, the component handles these keys:
//   - ?: Toggle the help window
//   - Esc: Close the help window (if open)
//
// # Implementation
//
// The component implements the Bubble Tea Model interface (Init, Update, View).
// When focused, it intercepts "?" and "Esc" keys for help window control.
// Parent components should call Focus()/Blur() to manage keyboard routing.
package keyhints
