package keybind

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Binding wraps key.Binding with display text and description.
// The Display field is used for rendering in keyhints, while
// the key.Binding handles actual key matching via key.Matches().
type Binding struct {
	key.Binding        // Embedded for key.Matches() compatibility
	Display     string // Short display text (e.g., "Enter", "Alt+Enter")
	Desc        string // Description (e.g., "Submit message")
}

// New creates a new Binding with the given display text, description, and trigger keys.
// Multiple keys can be provided to match any of them.
//
// Examples:
//
//	submit := New("Enter", "Submit message", "enter")
//	newline := New("Alt+Enter/Ctrl+J", "Insert newline", "alt+enter", "ctrl+j")
func New(display, desc string, keys ...string) Binding {
	return Binding{
		Binding: key.NewBinding(key.WithKeys(keys...), key.WithHelp(display, desc)),
		Display: display,
		Desc:    desc,
	}
}

// Matches returns true if the given key message matches this binding.
// This is a convenience wrapper around key.Matches().
func (b Binding) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, b.Binding)
}
