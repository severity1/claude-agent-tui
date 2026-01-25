package focus_test

import (
	"testing"

	"github.com/severity1/claude-agent-tui/internal/focus"
)

func TestFocus_SetsFocusedToTrue(t *testing.T) {
	var s focus.State

	// Initially should be unfocused
	if s.Focused() {
		t.Error("State should be unfocused by default")
	}

	// Focus should set focused to true
	cmd := s.Focus()
	if !s.Focused() {
		t.Error("Focus() should set focused to true")
	}
	if cmd != nil {
		t.Error("Focus() should return nil")
	}
}

func TestBlur_SetsFocusedToFalse(t *testing.T) {
	var s focus.State

	// Set focused first
	s.Focus()
	if !s.Focused() {
		t.Error("State should be focused after Focus()")
	}

	// Blur should set focused to false
	s.Blur()
	if s.Focused() {
		t.Error("Blur() should set focused to false")
	}
}

func TestFocused_ReturnsCurrentState(t *testing.T) {
	var s focus.State

	// Default state
	if s.Focused() {
		t.Error("Focused() should return false by default")
	}

	// After Focus
	s.Focus()
	if !s.Focused() {
		t.Error("Focused() should return true after Focus()")
	}

	// After Blur
	s.Blur()
	if s.Focused() {
		t.Error("Focused() should return false after Blur()")
	}
}

func TestSetFocused_SetsStateDirectly(t *testing.T) {
	var s focus.State

	// Set to true
	s.SetFocused(true)
	if !s.Focused() {
		t.Error("SetFocused(true) should set focused to true")
	}

	// Set to false
	s.SetFocused(false)
	if s.Focused() {
		t.Error("SetFocused(false) should set focused to false")
	}
}

func TestState_ZeroValueIsUnfocused(t *testing.T) {
	// Verify zero value behavior
	var s focus.State
	if s.Focused() {
		t.Error("Zero value of State should be unfocused")
	}
}

func TestState_MultipleTransitions(t *testing.T) {
	var s focus.State

	// Multiple focus/blur cycles
	s.Focus()
	s.Blur()
	s.Focus()
	s.Focus() // double focus should still be focused
	if !s.Focused() {
		t.Error("State should be focused after multiple Focus() calls")
	}

	s.Blur()
	s.Blur() // double blur should still be blurred
	if s.Focused() {
		t.Error("State should be blurred after multiple Blur() calls")
	}
}

// Test embedding pattern
type simpleComponent struct {
	focus.State
	name string
}

func TestState_EmbeddingPattern(t *testing.T) {
	c := simpleComponent{name: "test"}

	// Should have access to Focus, Blur, Focused via embedding
	if c.Focused() {
		t.Error("Embedded State should start unfocused")
	}

	c.Focus()
	if !c.Focused() {
		t.Error("Embedded State should be focused after Focus()")
	}

	c.Blur()
	if c.Focused() {
		t.Error("Embedded State should be unfocused after Blur()")
	}
}
