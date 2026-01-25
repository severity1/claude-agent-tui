package chatinput_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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

// extractCopiedMsg executes a batch command and extracts the CopiedMsg.
// CopyCmd returns a tea.BatchMsg containing both the clipboard command and flash timer.
func extractCopiedMsg(t *testing.T, cmd tea.Cmd) chatinput.CopiedMsg {
	t.Helper()
	if cmd == nil {
		t.Fatal("cmd is nil")
	}

	// Execute the batch command
	msg := cmd()
	batchMsg, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want tea.BatchMsg", msg)
	}

	// Execute each command in the batch and find CopiedMsg
	for _, batchCmd := range batchMsg {
		if batchCmd == nil {
			continue
		}
		result := batchCmd()
		if copiedMsg, ok := result.(chatinput.CopiedMsg); ok {
			return copiedMsg
		}
	}

	t.Fatal("batch did not contain CopiedMsg")
	return chatinput.CopiedMsg{}
}

func TestModel_CopyCmd_ReturnsCopiedMsg(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("text to copy")

	// Call CopyCmd directly (exported method)
	cmd := m.CopyCmd()
	if cmd == nil {
		t.Fatal("CopyCmd() returned nil, want non-nil command")
	}

	// Extract CopiedMsg from batch
	copiedMsg := extractCopiedMsg(t, cmd)

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

	copiedMsg := extractCopiedMsg(t, cmd)

	if copiedMsg.Text != "" {
		t.Errorf("CopiedMsg.Text = %q, want empty string", copiedMsg.Text)
	}
}

func TestModel_CopyCmd_MultilineText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("line1\nline2\nline3")

	cmd := m.CopyCmd()
	copiedMsg := extractCopiedMsg(t, cmd)

	if copiedMsg.Text != "line1\nline2\nline3" {
		t.Errorf("CopiedMsg.Text = %q, want multiline text", copiedMsg.Text)
	}
}

func TestModel_CopyCmd_UnicodeText(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("Hello 世界 🌍")

	cmd := m.CopyCmd()
	copiedMsg := extractCopiedMsg(t, cmd)

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

	// Extract and verify CopiedMsg from batch
	copiedMsg := extractCopiedMsg(t, cmd)

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
// Flash Animation Tests
// ============================================================================

func TestModel_CopyCmd_SetsFlashing(t *testing.T) {
	m := chatinput.New()
	m.Focus()
	m.SetValue("text")

	// Before copy, flashing should be false
	if m.Flashing() {
		t.Error("Flashing() = true before CopyCmd, want false")
	}

	// Call CopyCmd
	m.CopyCmd()

	// After copy, flashing should be true
	if !m.Flashing() {
		t.Error("Flashing() = false after CopyCmd, want true")
	}
}

func TestModel_Flashing_DefaultFalse(t *testing.T) {
	m := chatinput.New()

	if m.Flashing() {
		t.Error("Flashing() = true for new model, want false")
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

// ============================================================================
// Styling Tests
// ============================================================================

func TestNew_DefaultStyles(t *testing.T) {
	m := chatinput.New()

	// Default prompt should be empty string
	if m.Prompt() != "" {
		t.Errorf("Prompt() = %q, want empty string", m.Prompt())
	}
}

func TestNew_WithPrompt(t *testing.T) {
	m := chatinput.New(chatinput.WithPrompt("$ "))

	if m.Prompt() != "$ " {
		t.Errorf("Prompt() = %q, want %q", m.Prompt(), "$ ")
	}
}

func TestNew_WithPrompt_Empty(t *testing.T) {
	m := chatinput.New(chatinput.WithPrompt(""))

	// Empty prompt should be allowed
	if m.Prompt() != "" {
		t.Errorf("Prompt() = %q, want empty string", m.Prompt())
	}
}

func TestModel_View_DefaultStyleHasNoVisibleBorder(t *testing.T) {
	m := chatinput.New()

	view := m.View()

	// Default style now uses background instead of border
	// Should NOT contain rounded border characters
	if strings.Contains(view, "\u256d") || strings.Contains(view, "╭") {
		t.Errorf("View() should not contain rounded border with background style, got: %s", view)
	}
}

func TestModel_View_ContainsPrompt(t *testing.T) {
	m := chatinput.New(chatinput.WithPrompt(">>> "))

	view := m.View()

	// View should contain the prompt
	if !strings.Contains(view, ">>>") {
		t.Errorf("View() should contain prompt '>>>', got: %s", view)
	}
}

func TestModel_View_FocusUseCursor(t *testing.T) {
	// Focus is indicated by cursor visibility, not background change
	// The textarea styles are the same for both focus states
	m := chatinput.New()

	// Get unfocused view
	unfocusedView := m.View()

	// Focus and get focused view
	m.Focus()
	focusedView := m.View()

	// Both views render without error
	if unfocusedView == "" {
		t.Error("unfocused View() returned empty string")
	}
	if focusedView == "" {
		t.Error("focused View() returned empty string")
	}
}

func TestNew_WithPromptStyle(t *testing.T) {
	// Custom prompt style
	customStyle := lipgloss.NewStyle().Bold(true)
	m := chatinput.New(chatinput.WithPromptStyle(customStyle))

	// Model should be created without error
	if m.Prompt() != "" {
		t.Errorf("WithPromptStyle should not affect prompt text, got %q", m.Prompt())
	}
}

func TestNew_WithFocusedPromptStyle(t *testing.T) {
	// Custom focused prompt style
	customStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000"))
	m := chatinput.New(chatinput.WithFocusedPromptStyle(customStyle))

	// Model should be created without error
	if m.Prompt() != "" {
		t.Errorf("WithFocusedPromptStyle should not affect prompt text, got %q", m.Prompt())
	}
}

func TestNew_WithFocusStyle(t *testing.T) {
	// Custom focus style for border
	customStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	m := chatinput.New(chatinput.WithFocusStyle(customStyle))

	// Model should be created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom focus style")
	}
}

func TestModel_View_FocusedPromptStyleApplied(t *testing.T) {
	// Verify that the component can render with both focused and unfocused states
	// Note: lipgloss disables colors when not connected to a TTY, so we test
	// that both views render without error rather than comparing ANSI codes

	unfocusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	m := chatinput.New(
		chatinput.WithPrompt(">"),
		chatinput.WithPromptStyle(unfocusedStyle),
		chatinput.WithFocusedPromptStyle(focusedStyle),
	)

	// Get unfocused view
	unfocusedView := m.View()
	if unfocusedView == "" {
		t.Error("unfocused View() returned empty string")
	}
	// Verify prompt is present
	if !strings.Contains(unfocusedView, ">") {
		t.Error("unfocused View() should contain prompt character")
	}

	// Focus and get focused view
	m.Focus()
	focusedView := m.View()
	if focusedView == "" {
		t.Error("focused View() returned empty string")
	}
	// Verify prompt is present in focused view too
	if !strings.Contains(focusedView, ">") {
		t.Error("focused View() should contain prompt character")
	}
}

// ============================================================================
// Runtime Width/Height Tests
// ============================================================================

func TestModel_SetWidth(t *testing.T) {
	m := chatinput.New()

	m.SetWidth(100)

	if m.Width() != 100 {
		t.Errorf("Width() = %d after SetWidth(100), want 100", m.Width())
	}
}

func TestModel_Width_ReturnsConfiguredWidth(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(80))

	if m.Width() != 80 {
		t.Errorf("Width() = %d, want 80", m.Width())
	}
}

func TestModel_Width_ReturnsZeroIfNotSet(t *testing.T) {
	m := chatinput.New()

	if m.Width() != 0 {
		t.Errorf("Width() = %d for new model without width, want 0", m.Width())
	}
}

func TestModel_SetWidth_OverridesWithWidth(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(60))

	// Override at runtime
	m.SetWidth(120)

	if m.Width() != 120 {
		t.Errorf("Width() = %d after SetWidth override, want 120", m.Width())
	}
}

func TestModel_SetWidth_UpdatesCounter(t *testing.T) {
	m := chatinput.New(chatinput.WithShowCounter(true))
	m.SetValue("hello")

	// Set width and verify counter aligns
	m.SetWidth(80)

	view := m.View()
	// Counter should be present
	if !strings.Contains(view, "5 chars") {
		t.Errorf("View() missing char count after SetWidth, got: %s", view)
	}
}

func TestModel_SetHeight(t *testing.T) {
	m := chatinput.New()

	m.SetHeight(10)

	// Height is internal to textarea, verify model works after change
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string after SetHeight")
	}
}

func TestModel_SetSize(t *testing.T) {
	m := chatinput.New()

	m.SetSize(100, 5)

	// Width should be updated
	if m.Width() != 100 {
		t.Errorf("Width() = %d after SetSize(100, 5), want 100", m.Width())
	}

	// Verify model renders correctly
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string after SetSize")
	}
}

func TestModel_SetSize_UpdatesBoth(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(50), chatinput.WithHeight(3))

	// Update both dimensions
	m.SetSize(120, 8)

	if m.Width() != 120 {
		t.Errorf("Width() = %d after SetSize, want 120", m.Width())
	}
}

// ============================================================================
// Auto-Expanding Height Tests
// ============================================================================

func TestNew_WithMinHeight(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(3))

	if m.MinHeight() != 3 {
		t.Errorf("MinHeight() = %d, want 3", m.MinHeight())
	}
}

func TestNew_WithMaxHeight(t *testing.T) {
	m := chatinput.New(chatinput.WithMaxHeight(8))

	if m.MaxHeight() != 8 {
		t.Errorf("MaxHeight() = %d, want 8", m.MaxHeight())
	}
}

func TestNew_DefaultMinMaxHeight(t *testing.T) {
	m := chatinput.New()

	// Default minHeight is 1
	if m.MinHeight() != 1 {
		t.Errorf("MinHeight() = %d, want 1 (default)", m.MinHeight())
	}

	// Default maxHeight is 10
	if m.MaxHeight() != 10 {
		t.Errorf("MaxHeight() = %d, want 10 (default)", m.MaxHeight())
	}
}

func TestModel_Height_SingleLine(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(1), chatinput.WithMaxHeight(5))
	m.SetValue("hello")

	// Single line content should have height of 1
	if m.Height() != 1 {
		t.Errorf("Height() = %d for single line, want 1", m.Height())
	}
}

func TestModel_Height_MultiLine(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(1), chatinput.WithMaxHeight(5))
	m.SetValue("line1\nline2\nline3")

	// Three lines should have height of 3
	if m.Height() != 3 {
		t.Errorf("Height() = %d for 3 lines, want 3", m.Height())
	}
}

func TestModel_Height_ClampedToMin(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(3), chatinput.WithMaxHeight(10))
	m.SetValue("single line")

	// Even with single line, height should be clamped to minHeight
	if m.Height() != 3 {
		t.Errorf("Height() = %d, want 3 (minHeight)", m.Height())
	}
}

func TestModel_Height_ClampedToMax(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(1), chatinput.WithMaxHeight(3))
	m.SetValue("line1\nline2\nline3\nline4\nline5")

	// Five lines should be clamped to maxHeight of 3
	if m.Height() != 3 {
		t.Errorf("Height() = %d, want 3 (maxHeight)", m.Height())
	}
}

func TestModel_Height_UnlimitedMax(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(1), chatinput.WithMaxHeight(0))
	m.SetValue("line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11\nline12")

	// With maxHeight=0 (unlimited), height equals line count
	if m.Height() != 12 {
		t.Errorf("Height() = %d with unlimited max, want 12", m.Height())
	}
}

func TestModel_SetMinHeight(t *testing.T) {
	m := chatinput.New()
	m.SetMinHeight(5)

	if m.MinHeight() != 5 {
		t.Errorf("MinHeight() = %d after SetMinHeight(5), want 5", m.MinHeight())
	}
}

func TestModel_SetMaxHeight(t *testing.T) {
	m := chatinput.New()
	m.SetMaxHeight(15)

	if m.MaxHeight() != 15 {
		t.Errorf("MaxHeight() = %d after SetMaxHeight(15), want 15", m.MaxHeight())
	}
}

func TestModel_SetMinHeight_InvalidValue(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(3))
	m.SetMinHeight(0) // Invalid, should be ignored

	if m.MinHeight() != 3 {
		t.Errorf("MinHeight() = %d after SetMinHeight(0), want 3 (unchanged)", m.MinHeight())
	}
}

func TestModel_SetMaxHeight_ZeroAllowed(t *testing.T) {
	m := chatinput.New()
	m.SetMaxHeight(0) // 0 means unlimited

	if m.MaxHeight() != 0 {
		t.Errorf("MaxHeight() = %d after SetMaxHeight(0), want 0", m.MaxHeight())
	}
}

func TestWithMinHeight_InvalidValue(t *testing.T) {
	m := chatinput.New(chatinput.WithMinHeight(-1))

	// Invalid value should keep default
	if m.MinHeight() != 1 {
		t.Errorf("MinHeight() = %d with invalid option, want 1 (default)", m.MinHeight())
	}
}

// ============================================================================
// Info Bar Tests
// ============================================================================

func TestNew_WithInfoBar(t *testing.T) {
	m := chatinput.New(chatinput.WithInfoBar("Claude Sonnet 4.5"))

	if m.InfoBar() != "Claude Sonnet 4.5" {
		t.Errorf("InfoBar() = %q, want %q", m.InfoBar(), "Claude Sonnet 4.5")
	}
}

func TestNew_DefaultInfoBar(t *testing.T) {
	m := chatinput.New()

	// Default info bar is empty
	if m.InfoBar() != "" {
		t.Errorf("InfoBar() = %q, want empty string (default)", m.InfoBar())
	}
}

func TestModel_SetInfoBar(t *testing.T) {
	m := chatinput.New()
	m.SetInfoBar("GPT-4")

	if m.InfoBar() != "GPT-4" {
		t.Errorf("InfoBar() = %q after SetInfoBar, want %q", m.InfoBar(), "GPT-4")
	}
}

func TestModel_SetInfoBar_Empty(t *testing.T) {
	m := chatinput.New(chatinput.WithInfoBar("initial"))
	m.SetInfoBar("")

	if m.InfoBar() != "" {
		t.Errorf("InfoBar() = %q after SetInfoBar(''), want empty", m.InfoBar())
	}
}

func TestModel_View_WithInfoBar(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithInfoBar("Claude Opus 4"),
	)

	view := m.View()

	// View should contain the info bar text
	if !strings.Contains(view, "Claude Opus 4") {
		t.Errorf("View() should contain info bar text, got: %s", view)
	}
}

func TestModel_View_WithoutInfoBar(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(60))

	view := m.View()

	// View should render without error when no info bar is set
	if view == "" {
		t.Error("View() returned empty string without info bar")
	}
}

func TestNew_WithInfoStyle(t *testing.T) {
	customStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	m := chatinput.New(chatinput.WithInfoStyle(customStyle))

	// Model should be created without error
	if m.InfoBar() != "" {
		t.Errorf("WithInfoStyle should not affect info bar text, got %q", m.InfoBar())
	}
}

// ============================================================================
// View Structure Tests (Auto-expand + Info Bar)
// ============================================================================

func TestModel_View_AutoExpands(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithMinHeight(1),
		chatinput.WithMaxHeight(5),
	)

	// Start with single line
	m.SetValue("line1")
	view1 := m.View()

	// Add more lines
	m.SetValue("line1\nline2\nline3")
	view2 := m.View()

	// View should be taller with more content (more lines in output)
	lines1 := strings.Count(view1, "\n")
	lines2 := strings.Count(view2, "\n")

	if lines2 <= lines1 {
		t.Errorf("View with 3 lines of content (%d newlines) should be taller than 1 line (%d newlines)",
			lines2, lines1)
	}
}

func TestModel_View_ContainerPadding(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithInfoBar("Test Model"),
	)

	view := m.View()

	// View should contain the thick border characters
	if !strings.Contains(view, "┃") {
		t.Errorf("View() should contain thick border character, got: %s", view)
	}
}

// ============================================================================
// Visual Styling Options Tests
// ============================================================================

func TestNew_WithContentBackground(t *testing.T) {
	// Should accept and apply without error
	m := chatinput.New(chatinput.WithContentBackground(lipgloss.Color("#FF0000")))

	// Verify model is created (no panic or error)
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom content background")
	}
}

func TestNew_WithBorderStyle(t *testing.T) {
	// Test with rounded border (different from default thick border)
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderStyle(lipgloss.RoundedBorder()),
	)

	view := m.View()
	// Rounded border uses different characters than thick border
	// Thick: ┃ │ Rounded: │
	if strings.Contains(view, "┃") {
		t.Errorf("View() with RoundedBorder should not contain thick border character ┃")
	}
}

func TestNew_WithBorderColor(t *testing.T) {
	// Should accept custom border color without error
	m := chatinput.New(chatinput.WithBorderColor(lipgloss.Color("#00FF00")))

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom border color")
	}
}

func TestNew_WithFlashBorderColor(t *testing.T) {
	// Should accept custom flash border color without error
	m := chatinput.New(chatinput.WithFlashBorderColor(lipgloss.Color("#0000FF")))

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom flash border color")
	}
}

func TestNew_WithFlashDuration(t *testing.T) {
	// Should accept custom flash duration
	m := chatinput.New(chatinput.WithFlashDuration(200 * time.Millisecond))

	// Verify model is created correctly
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom flash duration")
	}
}

func TestNew_WithFlashDuration_Zero(t *testing.T) {
	// Zero duration should be ignored, keeping default
	m := chatinput.New(chatinput.WithFlashDuration(0))

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string when flash duration was zero")
	}
}

func TestNew_WithBorderPadding(t *testing.T) {
	// Default has 1,1,0,0 padding. Setting to 0,0,0,0 should result in shorter view
	mDefault := chatinput.New(chatinput.WithWidth(60))
	mNoPadding := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderPadding(0, 0, 0, 0),
	)

	viewDefault := mDefault.View()
	viewNoPadding := mNoPadding.View()

	// View with no padding should have fewer lines
	linesDefault := strings.Count(viewDefault, "\n")
	linesNoPadding := strings.Count(viewNoPadding, "\n")

	if linesNoPadding >= linesDefault {
		t.Errorf("View with no padding (%d lines) should be shorter than default (%d lines)",
			linesNoPadding, linesDefault)
	}
}

func TestNew_WithBorderPadding_NegativeIgnored(t *testing.T) {
	// Negative values should be ignored
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderPadding(-1, -1, -1, -1),
	)

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string when negative padding was provided")
	}
}

func TestNew_VisualOptions_Combined(t *testing.T) {
	// Test all visual options together
	m := chatinput.New(
		chatinput.WithWidth(80),
		chatinput.WithContentBackground(lipgloss.Color("#1a1a2e")),
		chatinput.WithBorderStyle(lipgloss.RoundedBorder()),
		chatinput.WithBorderColor(lipgloss.Color("#4a4a6a")),
		chatinput.WithFlashBorderColor(lipgloss.Color("#e94560")),
		chatinput.WithFlashDuration(300*time.Millisecond),
		chatinput.WithBorderPadding(2, 2, 0, 0),
	)

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with all visual options combined")
	}

	// Should use rounded border (not thick)
	if strings.Contains(view, "┃") {
		t.Error("View() should use rounded border, but contains thick border character")
	}
}

// ============================================================================
// Boundary Value Tests
// ============================================================================

func TestModel_SetWidth_ZeroValue(t *testing.T) {
	m := chatinput.New()
	m.SetWidth(0)

	// SetWidth(0) should be accepted and stored
	if m.Width() != 0 {
		t.Errorf("Width() = %d after SetWidth(0), want 0", m.Width())
	}

	// View should still render (using minContentWidth internally)
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string after SetWidth(0)")
	}
}

func TestModel_SetWidth_NegativeValue(t *testing.T) {
	m := chatinput.New()
	m.SetWidth(-1)

	// SetWidth(-1) should be accepted and stored
	if m.Width() != -1 {
		t.Errorf("Width() = %d after SetWidth(-1), want -1", m.Width())
	}

	// View should still render (using minContentWidth internally)
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string after SetWidth(-1)")
	}
}

func TestModel_CopyCmd_WhenUnfocused(t *testing.T) {
	m := chatinput.New()
	// Intentionally NOT calling Focus()
	m.SetValue("text to copy unfocused")

	// CopyCmd should work even when unfocused
	cmd := m.CopyCmd()
	if cmd == nil {
		t.Fatal("CopyCmd() returned nil when unfocused, want non-nil command")
	}

	// Extract and verify CopiedMsg from batch
	copiedMsg := extractCopiedMsg(t, cmd)

	if copiedMsg.Text != "text to copy unfocused" {
		t.Errorf("CopiedMsg.Text = %q, want %q", copiedMsg.Text, "text to copy unfocused")
	}

	// Flash state should be set even when unfocused
	if !m.Flashing() {
		t.Error("Flashing() = false after CopyCmd when unfocused, want true")
	}
}

// ============================================================================
// Border Left/Right Toggle Tests
// ============================================================================

func TestNew_DefaultBorderLeft(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(60))

	view := m.View()
	// Default borderLeft is true, so left border should be visible (thick border character)
	if !strings.Contains(view, "┃") && !strings.Contains(view, "\u2503") {
		t.Errorf("View() with default borderLeft=true should show thick left border, got: %s", view)
	}
}

func TestNew_DefaultBorderRight(t *testing.T) {
	m := chatinput.New(chatinput.WithWidth(60))

	view := m.View()
	// Default borderRight is true, so right border should be visible
	if !strings.Contains(view, "┃") && !strings.Contains(view, "\u2503") {
		t.Errorf("View() with default borderRight=true should show thick border, got: %s", view)
	}
}

func TestNew_WithBorderLeft_False(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderLeft(false),
	)

	view := m.View()
	// borderLeft=false should hide left border (use space)
	// View should start with a space character instead of thick border
	if strings.HasPrefix(view, "┃") || strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with borderLeft=false should not start with thick left border, got first char: %q", view[:10])
	}
}

func TestNew_WithBorderRight_False(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderRight(false),
	)

	view := m.View()
	lines := strings.Split(view, "\n")
	// At least one line should NOT end with thick border character
	hasInvisibleRight := false
	for _, line := range lines {
		if len(line) > 0 && !strings.HasSuffix(line, "┃") && !strings.HasSuffix(line, "\u2503") {
			hasInvisibleRight = true
			break
		}
	}
	if !hasInvisibleRight {
		t.Errorf("View() with borderRight=false should have invisible right border")
	}
}

func TestNew_WithBorderLeft_True(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderLeft(true),
	)

	view := m.View()
	// borderLeft=true should show visible thick border
	if !strings.HasPrefix(view, "┃") && !strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with borderLeft=true should start with thick left border, got: %q", view[:20])
	}
}

func TestNew_WithBorderRight_True(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderRight(true),
	)

	view := m.View()
	// borderRight=true should show visible thick border on the right
	lines := strings.Split(view, "\n")
	hasVisibleRight := false
	for _, line := range lines {
		if len(line) > 0 && (strings.HasSuffix(line, "┃") || strings.HasSuffix(line, "\u2503")) {
			hasVisibleRight = true
			break
		}
	}
	if !hasVisibleRight {
		t.Errorf("View() with borderRight=true should have visible right border")
	}
}

func TestModel_View_BothBordersHidden(t *testing.T) {
	m := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderLeft(false),
		chatinput.WithBorderRight(false),
	)

	view := m.View()
	// Both borders hidden but space reserved
	// Should not start with thick border
	if strings.HasPrefix(view, "┃") || strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with both borders hidden should not start with thick border")
	}
	// View should still render correctly
	if view == "" {
		t.Error("View() with both borders hidden returned empty string")
	}
}

func TestModel_View_BordersPreserveWidth(t *testing.T) {
	// Test that hidden borders still reserve space (width remains consistent)
	m1 := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderLeft(true),
		chatinput.WithBorderRight(true),
	)
	m2 := chatinput.New(
		chatinput.WithWidth(60),
		chatinput.WithBorderLeft(false),
		chatinput.WithBorderRight(false),
	)

	view1 := m1.View()
	view2 := m2.View()

	// Both should have the same number of lines (height)
	lines1 := strings.Count(view1, "\n")
	lines2 := strings.Count(view2, "\n")
	if lines1 != lines2 {
		t.Errorf("View() line count differs: visible borders=%d, hidden borders=%d", lines1, lines2)
	}
}
