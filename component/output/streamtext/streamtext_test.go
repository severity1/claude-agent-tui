package streamtext_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/severity1/claude-agent-tui/component/output/streamtext"
)

// ============================================================================
// Model Initialization Tests
// ============================================================================

func TestNew_DefaultValues(t *testing.T) {
	m := streamtext.New()

	if m.Content() != "" {
		t.Errorf("Content() = %q, want empty string", m.Content())
	}

	// Verify not streaming by checking View has no cursor
	view := m.View()
	if strings.Contains(view, "▌") {
		t.Error("View() contains cursor, want no cursor when not streaming")
	}
}

func TestNew_WithWidth(t *testing.T) {
	m := streamtext.New(streamtext.WithWidth(80))

	// Width is internal, but we can verify it was set by checking View behavior
	// For now, just verify the model was created
	if m.Content() != "" {
		t.Errorf("Content() = %q, want empty string", m.Content())
	}
}

// ============================================================================
// Append Functionality Tests
// ============================================================================

func TestModel_Append_SingleChunk(t *testing.T) {
	m := streamtext.New()
	m.Append("hello")

	if m.Content() != "hello" {
		t.Errorf("Content() = %q, want %q", m.Content(), "hello")
	}
}

func TestModel_Append_MultipleChunks(t *testing.T) {
	m := streamtext.New()
	m.Append("hello")
	m.Append(" ")
	m.Append("world")

	if m.Content() != "hello world" {
		t.Errorf("Content() = %q, want %q", m.Content(), "hello world")
	}
}

func TestModel_Append_EmptyString(t *testing.T) {
	m := streamtext.New()
	m.Append("hello")
	m.Append("")
	m.Append("world")

	if m.Content() != "helloworld" {
		t.Errorf("Content() = %q, want %q", m.Content(), "helloworld")
	}
}

// ============================================================================
// Clear Functionality Tests
// ============================================================================

func TestModel_Clear_ResetsContent(t *testing.T) {
	m := streamtext.New()
	m.Append("hello")
	m.Clear()

	if m.Content() != "" {
		t.Errorf("Content() = %q, want empty string after Clear()", m.Content())
	}
}

func TestModel_Clear_AfterMultipleAppends(t *testing.T) {
	m := streamtext.New()
	m.Append("hello")
	m.Append(" ")
	m.Append("world")
	m.Clear()

	if m.Content() != "" {
		t.Errorf("Content() = %q, want empty string after Clear()", m.Content())
	}
}

// ============================================================================
// Streaming State Tests
// ============================================================================

func TestModel_SetStreaming_True(t *testing.T) {
	m := streamtext.New()
	m.SetStreaming(true)

	view := m.View()
	if !strings.Contains(view, "▌") {
		t.Error("View() missing cursor when streaming=true")
	}
}

func TestModel_SetStreaming_False(t *testing.T) {
	m := streamtext.New()
	m.SetStreaming(true)
	m.SetStreaming(false)

	view := m.View()
	if strings.Contains(view, "▌") {
		t.Error("View() contains cursor when streaming=false")
	}
}

// ============================================================================
// Content Accessor Tests
// ============================================================================

func TestModel_Content_ReturnsAccumulatedText(t *testing.T) {
	m := streamtext.New()
	m.Append("line1\n")
	m.Append("line2\n")
	m.Append("line3")

	want := "line1\nline2\nline3"
	if m.Content() != want {
		t.Errorf("Content() = %q, want %q", m.Content(), want)
	}
}

// ============================================================================
// View Rendering Tests
// ============================================================================

func TestModel_View_StreamingShowsCursor(t *testing.T) {
	m := streamtext.New()
	m.Append("Hello")
	m.SetStreaming(true)

	view := m.View()
	if !strings.HasSuffix(view, "▌") {
		t.Errorf("View() = %q, want suffix %q", view, "▌")
	}
}

func TestModel_View_NotStreamingNoCursor(t *testing.T) {
	m := streamtext.New()
	m.Append("Hello")
	m.SetStreaming(false)

	view := m.View()
	if view != "Hello" {
		t.Errorf("View() = %q, want %q", view, "Hello")
	}
}

func TestModel_View_EmptyContent(t *testing.T) {
	m := streamtext.New()

	// Not streaming - empty view
	view := m.View()
	if view != "" {
		t.Errorf("View() = %q, want empty string", view)
	}

	// Streaming - just cursor
	m.SetStreaming(true)
	view = m.View()
	if view != "▌" {
		t.Errorf("View() = %q, want %q", view, "▌")
	}
}

func TestModel_View_WithWidth(t *testing.T) {
	m := streamtext.New(streamtext.WithWidth(10))
	m.Append("short")

	// Width is used for wrapping but doesn't affect short content
	view := m.View()
	if view != "short" {
		t.Errorf("View() = %q, want %q", view, "short")
	}
}

// ============================================================================
// Bubble Tea Interface Tests
// ============================================================================

func TestModel_Init_ReturnsNil(t *testing.T) {
	m := streamtext.New()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init() returned non-nil command, want nil")
	}
}

func TestModel_Update_ReturnsModelAndNil(t *testing.T) {
	m := streamtext.New()

	// Unknown message type
	type unknownMsg struct{}
	updated, cmd := m.Update(unknownMsg{})

	if cmd != nil {
		t.Error("Update() returned non-nil command for unknown msg, want nil")
	}

	// Verify we get back a Model (type assertion should work)
	if _, ok := updated.(streamtext.Model); !ok {
		t.Errorf("Update() returned %T, want streamtext.Model", updated)
	}
}

func TestModel_Update_WithKeyMsg(t *testing.T) {
	m := streamtext.New()

	// Key messages should be ignored (passed through)
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(keyMsg)

	if cmd != nil {
		t.Error("Update() returned non-nil command for key msg, want nil")
	}

	if _, ok := updated.(streamtext.Model); !ok {
		t.Errorf("Update() returned %T, want streamtext.Model", updated)
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestModel_Append_AfterClear(t *testing.T) {
	m := streamtext.New()
	m.Append("first")
	m.Clear()
	m.Append("second")

	if m.Content() != "second" {
		t.Errorf("Content() = %q, want %q", m.Content(), "second")
	}
}

func TestModel_SetStreaming_Toggle(t *testing.T) {
	m := streamtext.New()
	m.Append("text")

	// Toggle on
	m.SetStreaming(true)
	if !strings.Contains(m.View(), "▌") {
		t.Error("View() missing cursor after SetStreaming(true)")
	}

	// Toggle off
	m.SetStreaming(false)
	if strings.Contains(m.View(), "▌") {
		t.Error("View() has cursor after SetStreaming(false)")
	}

	// Toggle on again
	m.SetStreaming(true)
	if !strings.Contains(m.View(), "▌") {
		t.Error("View() missing cursor after second SetStreaming(true)")
	}
}

func TestModel_Append_SpecialCharacters(t *testing.T) {
	m := streamtext.New()
	m.Append("Hello, ")
	m.Append("世界")
	m.Append("! 🎉")

	want := "Hello, 世界! 🎉"
	if m.Content() != want {
		t.Errorf("Content() = %q, want %q", m.Content(), want)
	}
}

func TestModel_Append_Newlines(t *testing.T) {
	m := streamtext.New()
	m.Append("line1\n")
	m.Append("line2\r\n")
	m.Append("line3")

	want := "line1\nline2\r\nline3"
	if m.Content() != want {
		t.Errorf("Content() = %q, want %q", m.Content(), want)
	}
}
