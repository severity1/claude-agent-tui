// Package keybind provides a configurable keybinding system for TUI components.
//
// The package wraps the bubbles/key package to provide:
//   - Binding type that combines key.Binding with display text and description
//   - New() constructor for creating bindings with multiple trigger keys
//   - Compatibility with key.Matches() for event handling
//
// # Usage
//
// Create a binding:
//
//	submit := keybind.New("Enter", "Submit message", "enter")
//	newline := keybind.New("Alt+Enter/Ctrl+J", "Insert newline", "alt+enter", "ctrl+j")
//
// Check for matches in Update():
//
//	case key.Matches(msg, m.keyMap.Submit.Binding):
//	    return m.handleSubmit()
//
// # Integration with KeyHints
//
// Bindings can be converted to keyhints.Binding for display:
//
//	hints := []keyhints.Binding{
//	    {Key: binding.Display, Desc: binding.Desc},
//	}
package keybind
