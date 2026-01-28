package keyhints

import (
	"github.com/severity1/claude-agent-tui/internal/keybind"
)

// KeyBinding is an alias for keybind.Binding.
// It provides configurable keybindings with Display and Desc fields.
type KeyBinding = keybind.Binding

// NewKeyBinding creates a new KeyBinding with the given display text, description, and keys.
// This is a convenience wrapper around keybind.New().
func NewKeyBinding(display, desc string, keys ...string) KeyBinding {
	return keybind.New(display, desc, keys...)
}

// KeyMap defines keybindings for the keyhints component.
type KeyMap struct {
	ToggleHelp KeyBinding
	CloseHelp  KeyBinding
}

// DefaultKeyMap returns the default keybindings for keyhints.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		ToggleHelp: NewKeyBinding("?", "Toggle help", "?"),
		CloseHelp:  NewKeyBinding("Esc", "Close help", "esc"),
	}
}
