package keybind_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/severity1/claude-agent-tui/internal/keybind"
)

// ============================================================================
// Binding Creation Tests
// ============================================================================

func TestNew_BasicBinding(t *testing.T) {
	b := keybind.New("Enter", "Submit message", "enter")

	if b.Display != "Enter" {
		t.Errorf("Display = %q, want %q", b.Display, "Enter")
	}
	if b.Desc != "Submit message" {
		t.Errorf("Desc = %q, want %q", b.Desc, "Submit message")
	}
}

func TestNew_MultipleKeys(t *testing.T) {
	b := keybind.New("Alt+Enter/Ctrl+J", "Insert newline", "alt+enter", "ctrl+j")

	// Verify both keys are registered via key.Matches
	altEnterMsg := tea.KeyMsg{Type: tea.KeyEnter, Alt: true}
	ctrlJMsg := tea.KeyMsg{Type: tea.KeyCtrlJ}

	if !key.Matches(altEnterMsg, b.Binding) {
		t.Error("Binding should match Alt+Enter")
	}
	if !key.Matches(ctrlJMsg, b.Binding) {
		t.Error("Binding should match Ctrl+J")
	}
}

func TestNew_EmptyKeys(t *testing.T) {
	// Should not panic with empty keys
	b := keybind.New("Empty", "No keys")

	if b.Display != "Empty" {
		t.Errorf("Display = %q, want %q", b.Display, "Empty")
	}
}

// ============================================================================
// Key Matching Tests
// ============================================================================

func TestBinding_MatchesKeyMsg_SingleKey(t *testing.T) {
	b := keybind.New("Enter", "Submit", "enter")

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	if !key.Matches(enterMsg, b.Binding) {
		t.Error("Binding should match Enter key")
	}
}

func TestBinding_MatchesKeyMsg_MultipleKeys(t *testing.T) {
	b := keybind.New("Up/k", "History previous", "up", "k")

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}

	if !key.Matches(upMsg, b.Binding) {
		t.Error("Binding should match Up key")
	}
	if !key.Matches(kMsg, b.Binding) {
		t.Error("Binding should match 'k' key")
	}
}

func TestBinding_MatchesKeyMsg_Modifiers(t *testing.T) {
	tests := []struct {
		name    string
		keys    []string
		msg     tea.KeyMsg
		matches bool
	}{
		{
			name:    "Alt modifier",
			keys:    []string{"alt+enter"},
			msg:     tea.KeyMsg{Type: tea.KeyEnter, Alt: true},
			matches: true,
		},
		{
			name:    "Ctrl modifier via key type",
			keys:    []string{"ctrl+y"},
			msg:     tea.KeyMsg{Type: tea.KeyCtrlY},
			matches: true,
		},
		{
			name:    "Wrong modifier",
			keys:    []string{"alt+enter"},
			msg:     tea.KeyMsg{Type: tea.KeyEnter, Alt: false},
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := keybind.New("Test", "test", tt.keys...)
			got := key.Matches(tt.msg, b.Binding)
			if got != tt.matches {
				t.Errorf("key.Matches() = %v, want %v", got, tt.matches)
			}
		})
	}
}

func TestBinding_MatchesKeyMsg_NoMatch(t *testing.T) {
	b := keybind.New("Enter", "Submit", "enter")

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	if key.Matches(escMsg, b.Binding) {
		t.Error("Binding should NOT match Esc key")
	}
}

// ============================================================================
// Help Text Tests
// ============================================================================

func TestBinding_HelpText(t *testing.T) {
	b := keybind.New("Enter", "Submit message", "enter")

	help := b.Binding.Help()
	if help.Key != "Enter" {
		t.Errorf("Help().Key = %q, want %q", help.Key, "Enter")
	}
	if help.Desc != "Submit message" {
		t.Errorf("Help().Desc = %q, want %q", help.Desc, "Submit message")
	}
}

// ============================================================================
// Enable/Disable Tests
// ============================================================================

func TestBinding_SetEnabled(t *testing.T) {
	b := keybind.New("Enter", "Submit", "enter")

	// Initially enabled
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	if !key.Matches(enterMsg, b.Binding) {
		t.Error("Binding should match when enabled by default")
	}

	// Disable
	b.Binding.SetEnabled(false)
	if key.Matches(enterMsg, b.Binding) {
		t.Error("Binding should NOT match when disabled")
	}

	// Re-enable
	b.Binding.SetEnabled(true)
	if !key.Matches(enterMsg, b.Binding) {
		t.Error("Binding should match when re-enabled")
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestBinding_UsedInSwitch(t *testing.T) {
	// Test that bindings work correctly in a switch statement pattern
	submit := keybind.New("Enter", "Submit", "enter")
	clear := keybind.New("Esc", "Clear", "esc")
	newline := keybind.New("Alt+Enter", "Newline", "alt+enter", "ctrl+j")

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want string
	}{
		{"Enter submits", tea.KeyMsg{Type: tea.KeyEnter}, "submit"},
		{"Esc clears", tea.KeyMsg{Type: tea.KeyEsc}, "clear"},
		{"Alt+Enter newline", tea.KeyMsg{Type: tea.KeyEnter, Alt: true}, "newline"},
		{"Ctrl+J newline", tea.KeyMsg{Type: tea.KeyCtrlJ}, "newline"},
		{"Other key unhandled", tea.KeyMsg{Type: tea.KeyUp}, "unhandled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			switch {
			case key.Matches(tt.msg, submit.Binding):
				got = "submit"
			case key.Matches(tt.msg, newline.Binding):
				got = "newline"
			case key.Matches(tt.msg, clear.Binding):
				got = "clear"
			default:
				got = "unhandled"
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
