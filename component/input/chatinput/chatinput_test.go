package chatinput_test

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/severity1/claude-agent-tui/component/input/chatinput"
)

// ============================================================================
// Model Initialization Tests
// ============================================================================

func TestNew_DefaultValues(t *testing.T) {
	m := chatinput.New()

	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
	if m.Focused() {
		t.Error("Focused() = true, want false for new model")
	}
}

func TestNew_WithPlaceholder(t *testing.T) {
	m := chatinput.New(chatinput.WithPlaceholder("Type here..."))

	// Placeholder appears in View() when input is empty
	view := m.View()
	if !strings.Contains(view, "Type here...") {
		t.Errorf("View() = %q, want to contain placeholder 'Type here...'", view)
	}
}

func TestNew_WithWidth(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(80), chatinput.WithShowCounter(true))

	// Width affects counter alignment - verify model renders
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string, want textarea rendering")
	}
}

func TestNew_WithHeight(t *testing.T) {
	m := chatinput.New(chatinput.WithHeight(5))

	// Height affects textarea rendering - verify model renders
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string, want textarea rendering")
	}
}

// ============================================================================
// Focus Management Tests
// ============================================================================

func TestModel_Focus_SetsState(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	if !m.Focused() {
		t.Error("Focused() = false after Focus(), want true")
	}
}

func TestModel_Blur_ClearsState(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.Blur()

	if m.Focused() {
		t.Error("Focused() = true after Blur(), want false")
	}
}

func TestModel_Focused_ReturnsState(t *testing.T) {
	m := chatinput.New()

	if m.Focused() {
		t.Error("Focused() = true for new model, want false")
	}

	m.Focus()
	if !m.Focused() {
		t.Error("Focused() = false after Focus(), want true")
	}

	m.Blur()
	if m.Focused() {
		t.Error("Focused() = true after Blur(), want false")
	}
}

// ============================================================================
// Value Management Tests
// ============================================================================

func TestModel_Value_ReturnsContent(t *testing.T) {
	m := chatinput.New()
	m.SetValue("hello world")

	if m.Value() != "hello world" {
		t.Errorf("Value() = %q, want %q", m.Value(), "hello world")
	}
}

func TestModel_SetValue_SetsContent(t *testing.T) {
	m := chatinput.New()
	m.SetValue("test content")

	if m.Value() != "test content" {
		t.Errorf("Value() = %q, want %q", m.Value(), "test content")
	}
}

func TestModel_SetValue_Empty(t *testing.T) {
	m := chatinput.New()
	m.SetValue("something")
	m.SetValue("")

	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
}

// ============================================================================
// Key Handling Tests
// ============================================================================

func TestChatInput_Submit(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("hello")

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msg)
	model := updated.(chatinput.Model)

	// Should emit SubmitMsg
	if cmd == nil {
		t.Fatal("Update() returned nil cmd for Enter, want SubmitMsg command")
	}

	// Execute the command to get the message
	resultMsg := cmd()
	submitMsg, ok := resultMsg.(chatinput.SubmitMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want chatinput.SubmitMsg", resultMsg)
	}

	if submitMsg.Text != "hello" {
		t.Errorf("SubmitMsg.Text = %q, want %q", submitMsg.Text, "hello")
	}

	// Input should be cleared after submit
	if model.Value() != "" {
		t.Errorf("Value() = %q after submit, want empty string", model.Value())
	}
}

func TestChatInput_Newline(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("line1")

	// Press Alt+Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter, Alt: true}
	updated, _ := m.Update(msg)
	model := updated.(chatinput.Model)

	// Should insert newline, not submit
	value := model.Value()
	if value != "line1\n" {
		t.Errorf("Value() = %q after Alt+Enter, want %q", value, "line1\n")
	}
}

func TestModel_Update_CtrlJ_InsertsNewline(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("line1")

	// Press Ctrl+J
	msg := tea.KeyMsg{Type: tea.KeyCtrlJ}
	updated, _ := m.Update(msg)
	model := updated.(chatinput.Model)

	// Should insert newline, not submit
	value := model.Value()
	if value != "line1\n" {
		t.Errorf("Value() = %q after Ctrl+J, want %q", value, "line1\n")
	}
}

func TestChatInput_Clear(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("some text")

	// Press Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := m.Update(msg)
	model := updated.(chatinput.Model)

	if model.Value() != "" {
		t.Errorf("Value() = %q after Esc, want empty string", model.Value())
	}
}

func TestModel_Update_Enter_EmptyInput_NoSubmit(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	// Leave input empty

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(msg)

	// Should NOT emit SubmitMsg for empty input
	if cmd != nil {
		t.Error("Update() returned non-nil cmd for Enter with empty input, want nil")
	}
}

// ============================================================================
// Bubble Tea Interface Tests
// ============================================================================

func TestModel_Init_ReturnsNil(t *testing.T) {
	m := chatinput.New()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init() returned non-nil command, want nil")
	}
}

func TestModel_Update_PassthroughOtherKeys(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Regular character key should be passed to textarea
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updated, _ := m.Update(msg)
	model := updated.(chatinput.Model)

	// Character should be added to textarea
	if model.Value() != "a" {
		t.Errorf("Value() = %q after typing 'a', want %q", model.Value(), "a")
	}
}

func TestModel_View_RendersTextarea(t *testing.T) {
	m := chatinput.New()
	m.SetValue("test content")

	view := m.View()

	// View should contain the content
	if view == "" {
		t.Error("View() returned empty string, want textarea rendering")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestModel_Submit_MultilineText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("line1\nline2\nline3")

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("Update() returned nil cmd for Enter with multiline text")
	}

	resultMsg := cmd()
	submitMsg, ok := resultMsg.(chatinput.SubmitMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want chatinput.SubmitMsg", resultMsg)
	}

	if submitMsg.Text != "line1\nline2\nline3" {
		t.Errorf("SubmitMsg.Text = %q, want %q", submitMsg.Text, "line1\nline2\nline3")
	}
}

func TestModel_Update_WhenNotFocused(t *testing.T) {
	m := chatinput.New()
	// Not focused

	// Try to type
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updated, _ := m.Update(msg)
	model := updated.(chatinput.Model)

	// Should not add character when not focused
	if model.Value() != "" {
		t.Errorf("Value() = %q when not focused, want empty string", model.Value())
	}
}

func TestModel_Focus_ReturnsCmd(t *testing.T) {
	m := chatinput.New()
	cmd := m.Focus()

	// Focus returns a command for textarea cursor blink
	if cmd == nil {
		t.Error("Focus() returned nil, want cursor blink command")
	}
}

// ============================================================================
// History Buffer Tests
// ============================================================================

func TestNew_WithHistorySize(t *testing.T) {
	m := chatinput.New(chatinput.WithHistorySize(10))

	// History size is internal, verify model created correctly
	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
}

func TestModel_History_UpArrow_RecallsPrevious(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit some messages to build history
	m.SetValue("first")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	m.SetValue("second")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Press Up to recall previous message
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	if m.Value() != "second" {
		t.Errorf("Value() = %q after Up, want %q", m.Value(), "second")
	}

	// Press Up again to get older message
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	if m.Value() != "first" {
		t.Errorf("Value() = %q after second Up, want %q", m.Value(), "first")
	}
}

func TestModel_History_DownArrow_NavigatesForward(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit messages
	m.SetValue("first")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	m.SetValue("second")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Go back in history
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	// Now at "first", go forward
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(chatinput.Model)

	if m.Value() != "second" {
		t.Errorf("Value() = %q after Down, want %q", m.Value(), "second")
	}
}

func TestModel_History_DownArrow_RestoresDraftAtEnd(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit a message
	m.SetValue("message")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Start with empty draft, go back then forward past the end
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(chatinput.Model)

	// Should restore to empty draft (was empty when we started)
	if m.Value() != "" {
		t.Errorf("Value() = %q after Down past end, want empty string (draft)", m.Value())
	}
}

func TestModel_History_UpArrow_StopsAtOldest(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit one message
	m.SetValue("only")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Press Up multiple times
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	// Should stay at oldest message
	if m.Value() != "only" {
		t.Errorf("Value() = %q after multiple Ups, want %q", m.Value(), "only")
	}
}

func TestModel_History_EmptyHistory_UpDoesNothing(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("current")

	// Press Up with no history
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	// Should keep current value
	if m.Value() != "current" {
		t.Errorf("Value() = %q after Up with no history, want %q", m.Value(), "current")
	}
}

func TestModel_History_MaxSize(t *testing.T) {
	m := chatinput.New(chatinput.WithHistorySize(3))
	m.Focus()

	// Submit more messages than max size
	for i := 1; i <= 5; i++ {
		m.SetValue(fmt.Sprintf("msg%d", i))
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = updated.(chatinput.Model)
	}

	// Navigate back - should only have last 3 messages
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	if m.Value() != "msg5" {
		t.Errorf("Value() = %q, want %q", m.Value(), "msg5")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	if m.Value() != "msg4" {
		t.Errorf("Value() = %q, want %q", m.Value(), "msg4")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	if m.Value() != "msg3" {
		t.Errorf("Value() = %q, want %q", m.Value(), "msg3")
	}

	// Should not go further back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)
	if m.Value() != "msg3" {
		t.Errorf("Value() = %q after extra Up, want %q", m.Value(), "msg3")
	}
}

func TestModel_History_DefaultSize(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Default should allow reasonable history (at least 100)
	for i := 1; i <= 50; i++ {
		m.SetValue(fmt.Sprintf("msg%d", i))
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = updated.(chatinput.Model)
	}

	// Should be able to recall all 50
	for i := 50; i >= 1; i-- {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
		m = updated.(chatinput.Model)
		want := fmt.Sprintf("msg%d", i)
		if m.Value() != want {
			t.Errorf("Value() = %q, want %q", m.Value(), want)
			break
		}
	}
}

// ============================================================================
// Draft Save Tests
// ============================================================================

func TestModel_History_SavesDraft(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit a message to have history
	m.SetValue("history")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Type a draft
	m.SetValue("my draft")

	// Navigate to history
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	if m.Value() != "history" {
		t.Errorf("Value() = %q after Up, want %q", m.Value(), "history")
	}

	// Navigate back - should restore draft
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(chatinput.Model)

	if m.Value() != "my draft" {
		t.Errorf("Value() = %q after Down, want %q (draft)", m.Value(), "my draft")
	}
}

func TestModel_History_DraftClearedOnSubmit(t *testing.T) {
	m := chatinput.New()
	m.Focus()

	// Submit to have history
	m.SetValue("first")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Type draft, go to history
	m.SetValue("draft")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(chatinput.Model)

	// Submit from history
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(chatinput.Model)

	// Should be empty, not the old draft
	if m.Value() != "" {
		t.Errorf("Value() = %q after submit from history, want empty", m.Value())
	}
}

// ============================================================================
// Counter Tests
// ============================================================================

func TestNew_WithShowCounter(t *testing.T) {
	m := chatinput.New(chatinput.WithShowCounter(true))

	// Verify model created correctly
	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
}

func TestModel_CharCount_Empty(t *testing.T) {
	m := chatinput.New()

	if m.CharCount() != 0 {
		t.Errorf("CharCount() = %d, want 0", m.CharCount())
	}
}

func TestModel_CharCount_WithContent(t *testing.T) {
	m := chatinput.New()
	m.SetValue("hello")

	if m.CharCount() != 5 {
		t.Errorf("CharCount() = %d, want 5", m.CharCount())
	}
}

func TestModel_LineCount_Empty(t *testing.T) {
	m := chatinput.New()

	if m.LineCount() != 1 {
		t.Errorf("LineCount() = %d, want 1", m.LineCount())
	}
}

func TestModel_LineCount_SingleLine(t *testing.T) {
	m := chatinput.New()
	m.SetValue("hello world")

	if m.LineCount() != 1 {
		t.Errorf("LineCount() = %d, want 1", m.LineCount())
	}
}

func TestModel_LineCount_MultiLine(t *testing.T) {
	m := chatinput.New()
	m.SetValue("line1\nline2\nline3")

	if m.LineCount() != 3 {
		t.Errorf("LineCount() = %d, want 3", m.LineCount())
	}
}

func TestModel_View_WithCounter(t *testing.T) {
	m := chatinput.New(chatinput.WithShowCounter(true))
	m.SetValue("hello")

	view := m.View()

	// View should contain counter info
	if !strings.Contains(view, "5 chars") {
		t.Errorf("View() missing char count, got: %s", view)
	}
	if !strings.Contains(view, "1 lines") {
		t.Errorf("View() missing line count, got: %s", view)
	}
}

func TestModel_View_WithoutCounter(t *testing.T) {
	m := chatinput.New() // counter off by default
	m.SetValue("hello")

	view := m.View()

	// View should NOT contain counter info
	if strings.Contains(view, "chars") {
		t.Errorf("View() should not contain counter when disabled, got: %s", view)
	}
}

// ============================================================================
// Line Numbers Tests
// ============================================================================

func TestNew_WithLineNumbers(t *testing.T) {
	m := chatinput.New(chatinput.WithLineNumbers(true))

	// Line numbers is internal to textarea, verify model created
	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
}

func TestNew_WithLineNumbers_False(t *testing.T) {
	m := chatinput.New(chatinput.WithLineNumbers(false))

	// Verify model created correctly
	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty string", m.Value())
	}
}

// ============================================================================
// Copy to Clipboard Tests (OSC 52)
// ============================================================================

func TestModel_CopyCmd_ReturnsCopiedMsg(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("text to copy")

	// Call CopyCmd directly (exported method)
	cmd := m.CopyCmd()
	if cmd == nil {
		t.Fatal("CopyCmd() returned nil, want non-nil command")
	}

	// Execute the command
	msg := cmd()
	copiedMsg, ok := msg.(chatinput.CopiedMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want chatinput.CopiedMsg", msg)
	}

	// Verify the copied text
	if copiedMsg.Text != "text to copy" {
		t.Errorf("CopiedMsg.Text = %q, want %q", copiedMsg.Text, "text to copy")
	}

	// OSC 52 always returns nil error
	if copiedMsg.Err != nil {
		t.Errorf("CopiedMsg.Err = %v, want nil", copiedMsg.Err)
	}
}

func TestModel_CopyCmd_EmptyText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	// Leave empty

	cmd := m.CopyCmd()
	if cmd == nil {
		t.Fatal("CopyCmd() returned nil for empty text, want non-nil command")
	}

	msg := cmd()
	copiedMsg, ok := msg.(chatinput.CopiedMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want chatinput.CopiedMsg", msg)
	}

	if copiedMsg.Text != "" {
		t.Errorf("CopiedMsg.Text = %q, want empty string", copiedMsg.Text)
	}
}

func TestModel_CopyCmd_MultilineText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("line1\nline2\nline3")

	cmd := m.CopyCmd()
	msg := cmd()
	copiedMsg := msg.(chatinput.CopiedMsg)

	if copiedMsg.Text != "line1\nline2\nline3" {
		t.Errorf("CopiedMsg.Text = %q, want multiline text", copiedMsg.Text)
	}
}

func TestModel_CopyCmd_UnicodeText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("Hello 世界 🌍")

	cmd := m.CopyCmd()
	msg := cmd()
	copiedMsg := msg.(chatinput.CopiedMsg)

	if copiedMsg.Text != "Hello 世界 🌍" {
		t.Errorf("CopiedMsg.Text = %q, want unicode text", copiedMsg.Text)
	}
}

func TestModel_Update_CtrlY_TriggersCopy(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("text to copy")

	// Ctrl+Y is reliably detected via tea.KeyCtrlY
	msg := tea.KeyMsg{Type: tea.KeyCtrlY}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("Update() returned nil cmd for Ctrl+Y, want CopyCmd")
	}

	// Execute and verify it's a CopiedMsg
	result := cmd()
	copiedMsg, ok := result.(chatinput.CopiedMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want chatinput.CopiedMsg", result)
	}

	if copiedMsg.Text != "text to copy" {
		t.Errorf("CopiedMsg.Text = %q, want %q", copiedMsg.Text, "text to copy")
	}
}

func TestCopiedMsg_Type(t *testing.T) {
	// Verify CopiedMsg struct exists and has expected fields
	msg := chatinput.CopiedMsg{Text: "test", Err: nil}

	if msg.Text != "test" {
		t.Errorf("CopiedMsg.Text = %q, want %q", msg.Text, "test")
	}
	if msg.Err != nil {
		t.Errorf("CopiedMsg.Err = %v, want nil", msg.Err)
	}
}

// ============================================================================
// Unicode Character Counting Tests
// ============================================================================

func TestCharCount_Unicode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "ascii only",
			input: "hello",
			want:  5,
		},
		{
			name:  "with emoji",
			input: "hello 👋",
			want:  7, // 5 letters + 1 space + 1 emoji (not 10 bytes)
		},
		{
			name:  "emoji sequence",
			input: "🎉🎊🎁",
			want:  3,
		},
		{
			name:  "japanese characters",
			input: "こんにちは",
			want:  5,
		},
		{
			name:  "mixed unicode",
			input: "Hello 世界 🌍",
			want:  10, // 5 + 1 + 2 + 1 + 1
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := chatinput.New()
			m.Focus()
			m.SetValue(tt.input)

			got := m.CharCount()
			if got != tt.want {
				t.Errorf("CharCount() = %d, want %d (input: %q, bytes: %d)",
					got, tt.want, tt.input, len(tt.input))
			}
		})
	}
}
