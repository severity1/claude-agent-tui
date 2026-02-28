// Package focus provides a reusable focus state tracking for TUI components.
package focus

import tea "github.com/charmbracelet/bubbletea"

// State tracks focus state for a component.
// Embed this in component Models to get standard focus tracking.
//
// Components can either:
//   - Embed directly (simple components like keyhints)
//   - Delegate with override (components with nested focusable elements)
//
// Example for keyhints (simple):
//
//	type Model struct {
//	    focus.State  // embedded - gets Focus(), Blur(), Focused() for free
//	    // ... other fields
//	}
//
// Example for chatinput (delegation):
//
//	type Model struct {
//	    focus.State  // embedded for state tracking
//	    textarea textarea.Model
//	}
//
//	// Override Focus to also focus nested component
//	func (m *Model) Focus() tea.Cmd {
//	    m.State.Focus()  // track our state
//	    return m.textarea.Focus()  // delegate and return cmd
//	}
//
//	// Override Blur to also blur nested component
//	func (m *Model) Blur() {
//	    m.State.Blur()
//	    m.textarea.Blur()
//	}
//	// Focused() inherited from embedded State - no override needed
type State struct {
	focused bool
}

// Focus sets the focused state to true.
// Returns nil - override in component if you need to return a tea.Cmd.
func (s *State) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur sets the focused state to false.
func (s *State) Blur() {
	s.focused = false
}

// Focused returns the current focus state.
func (s State) Focused() bool {
	return s.focused
}

// SetFocused sets the focus state directly.
func (s *State) SetFocused(focused bool) {
	s.focused = focused
}
