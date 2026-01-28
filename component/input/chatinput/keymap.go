package chatinput

import (
	"github.com/severity1/claude-agent-tui/component/shared/keyhints"
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

// KeyMap defines all keybindings for the chatinput component.
type KeyMap struct {
	Submit        KeyBinding
	InsertNewline KeyBinding
	Clear         KeyBinding
	HistoryPrev   KeyBinding
	HistoryNext   KeyBinding
	Copy          KeyBinding
}

// DefaultKeyMap returns the default keybindings for chatinput.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Submit:        NewKeyBinding("Enter", "Submit message", "enter"),
		InsertNewline: NewKeyBinding("Alt+Enter", "Insert newline", "alt+enter", "ctrl+j"),
		Clear:         NewKeyBinding("Esc", "Clear input", "esc"),
		HistoryPrev:   NewKeyBinding("Up", "Previous history", "up"),
		HistoryNext:   NewKeyBinding("Down", "Next history", "down"),
		Copy:          NewKeyBinding("Ctrl+Y", "Copy to clipboard", "ctrl+y", "ctrl+shift+c"),
	}
}

// ToHintBindings converts the KeyMap to a slice of keyhints.Binding for display.
// This allows keyhints to display the current keybindings.
func (km KeyMap) ToHintBindings() []keyhints.Binding {
	return []keyhints.Binding{
		{Key: km.Submit.Display, Desc: km.Submit.Desc},
		{Key: km.InsertNewline.Display, Desc: km.InsertNewline.Desc},
		{Key: km.Clear.Display, Desc: km.Clear.Desc},
		{Key: km.HistoryPrev.Display, Desc: km.HistoryPrev.Desc},
		{Key: km.HistoryNext.Display, Desc: km.HistoryNext.Desc},
		{Key: km.Copy.Display, Desc: km.Copy.Desc},
	}
}
