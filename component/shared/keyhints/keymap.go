package keyhints

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/severity1/claude-agent-tui/internal/keybind"
)

// KeyBinding wraps keybind.Binding with a Matches helper method.
// This provides a cleaner API for checking key matches in Update().
type KeyBinding struct {
	keybind.Binding
}

// NewKeyBinding creates a new KeyBinding with the given display text, description, and keys.
func NewKeyBinding(display, desc string, keys ...string) KeyBinding {
	return KeyBinding{
		Binding: keybind.New(display, desc, keys...),
	}
}

// Matches returns true if the given key message matches this binding.
// This is a convenience wrapper around key.Matches().
func (kb KeyBinding) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, kb.Binding.Binding)
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
