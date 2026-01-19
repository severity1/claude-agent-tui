package adapter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	claudecode "github.com/severity1/claude-agent-sdk-go"
	"github.com/severity1/claude-agent-tui/adapter"
)

// ============================================================================
// Message Type Tests
// ============================================================================

func TestStreamDeltaMsg_Fields(t *testing.T) {
	msg := adapter.StreamDeltaMsg{
		Text:      "Hello",
		BlockType: "text",
		Index:     0,
	}

	if msg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", msg.Text, "Hello")
	}
	if msg.BlockType != "text" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "text")
	}
	if msg.Index != 0 {
		t.Errorf("Index = %d, want %d", msg.Index, 0)
	}
}

func TestStreamBlockStartMsg_Fields(t *testing.T) {
	msg := adapter.StreamBlockStartMsg{
		BlockType: "tool_use",
		Index:     1,
		ToolName:  "Bash",
		ToolID:    "tool_123",
	}

	if msg.BlockType != "tool_use" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "tool_use")
	}
	if msg.Index != 1 {
		t.Errorf("Index = %d, want %d", msg.Index, 1)
	}
	if msg.ToolName != "Bash" {
		t.Errorf("ToolName = %q, want %q", msg.ToolName, "Bash")
	}
	if msg.ToolID != "tool_123" {
		t.Errorf("ToolID = %q, want %q", msg.ToolID, "tool_123")
	}
}

func TestStreamBlockStopMsg_Fields(t *testing.T) {
	msg := adapter.StreamBlockStopMsg{
		Index: 2,
	}

	if msg.Index != 2 {
		t.Errorf("Index = %d, want %d", msg.Index, 2)
	}
}

func TestThinkingDeltaMsg_Fields(t *testing.T) {
	msg := adapter.ThinkingDeltaMsg{
		Text:  "Let me think...",
		Index: 0,
	}

	if msg.Text != "Let me think..." {
		t.Errorf("Text = %q, want %q", msg.Text, "Let me think...")
	}
	if msg.Index != 0 {
		t.Errorf("Index = %d, want %d", msg.Index, 0)
	}
}

func TestToolUseDeltaMsg_Fields(t *testing.T) {
	msg := adapter.ToolUseDeltaMsg{
		PartialJSON: `{"command":`,
		ToolName:    "Bash",
		ToolID:      "tool_456",
		Index:       1,
	}

	if msg.PartialJSON != `{"command":` {
		t.Errorf("PartialJSON = %q, want %q", msg.PartialJSON, `{"command":`)
	}
	if msg.ToolName != "Bash" {
		t.Errorf("ToolName = %q, want %q", msg.ToolName, "Bash")
	}
	if msg.ToolID != "tool_456" {
		t.Errorf("ToolID = %q, want %q", msg.ToolID, "tool_456")
	}
	if msg.Index != 1 {
		t.Errorf("Index = %d, want %d", msg.Index, 1)
	}
}

func TestAssistantMsg_Fields(t *testing.T) {
	assistantMsg := &claudecode.AssistantMessage{}
	msg := adapter.AssistantMsg{
		Message: assistantMsg,
	}

	if msg.Message != assistantMsg {
		t.Errorf("Message = %v, want %v", msg.Message, assistantMsg)
	}
}

func TestResultMsg_Fields(t *testing.T) {
	resultMsg := &claudecode.ResultMessage{}
	msg := adapter.ResultMsg{
		Message:      resultMsg,
		InputTokens:  100,
		OutputTokens: 50,
		TotalCost:    0.0015,
	}

	if msg.Message != resultMsg {
		t.Errorf("Message = %v, want %v", msg.Message, resultMsg)
	}
	if msg.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want %d", msg.InputTokens, 100)
	}
	if msg.OutputTokens != 50 {
		t.Errorf("OutputTokens = %d, want %d", msg.OutputTokens, 50)
	}
	if msg.TotalCost != 0.0015 {
		t.Errorf("TotalCost = %f, want %f", msg.TotalCost, 0.0015)
	}
}

func TestStreamDoneMsg_IsEmpty(t *testing.T) {
	msg := adapter.StreamDoneMsg{}
	// Verify it's an empty struct by checking its zero value behavior
	_ = msg // should compile - no fields to access
}

func TestStreamErrorMsg_Fields(t *testing.T) {
	err := context.Canceled
	msg := adapter.StreamErrorMsg{
		Err: err,
	}

	if msg.Err != err {
		t.Errorf("Err = %v, want %v", msg.Err, err)
	}
}

func TestStreamErrorMsg_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantStr string
	}{
		{
			name:    "with error",
			err:     context.Canceled,
			wantStr: "context canceled",
		},
		{
			name:    "nil error",
			err:     nil,
			wantStr: "unknown stream error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := adapter.StreamErrorMsg{Err: tt.err}
			if got := msg.Error(); got != tt.wantStr {
				t.Errorf("Error() = %q, want %q", got, tt.wantStr)
			}
		})
	}
}

// ============================================================================
// StreamCmd Tests
// ============================================================================

func TestStreamCmd_ReturnsTeaCmd(t *testing.T) {
	ctx := context.Background()
	ch := make(chan claudecode.Message)
	defer close(ch)

	cmd := adapter.StreamCmd(ctx, ch)

	if cmd == nil {
		t.Fatal("StreamCmd returned nil")
	}

	// Verify it returns a function
	_, ok := interface{}(cmd).(tea.Cmd)
	if !ok {
		t.Errorf("StreamCmd should return tea.Cmd, got %T", cmd)
	}
}

func TestStreamCmd_ChannelReceive(t *testing.T) {
	ctx := context.Background()
	ch := make(chan claudecode.Message, 1)

	// Send an AssistantMessage to the channel
	assistantMsg := &claudecode.AssistantMessage{}
	ch <- assistantMsg

	cmd := adapter.StreamCmd(ctx, ch)
	result := cmd()

	msg, ok := result.(adapter.AssistantMsg)
	if !ok {
		t.Fatalf("expected AssistantMsg, got %T", result)
	}

	if msg.Message != assistantMsg {
		t.Errorf("Message = %v, want %v", msg.Message, assistantMsg)
	}
}

func TestStreamCmd_ChannelClose_ReturnsStreamDoneMsg(t *testing.T) {
	ctx := context.Background()
	ch := make(chan claudecode.Message)
	close(ch)

	cmd := adapter.StreamCmd(ctx, ch)
	result := cmd()

	_, ok := result.(adapter.StreamDoneMsg)
	if !ok {
		t.Fatalf("expected StreamDoneMsg, got %T", result)
	}
}

func TestStreamCmd_ContextCancellation_ReturnsStreamErrorMsg(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan claudecode.Message) // unbuffered, will block

	cancel() // cancel immediately

	cmd := adapter.StreamCmd(ctx, ch)
	result := cmd()

	errMsg, ok := result.(adapter.StreamErrorMsg)
	if !ok {
		t.Fatalf("expected StreamErrorMsg, got %T", result)
	}

	if errMsg.Err != context.Canceled {
		t.Errorf("Err = %v, want %v", errMsg.Err, context.Canceled)
	}
}

func TestStreamCmd_NilChannel_ReturnsStreamDoneMsg(t *testing.T) {
	ctx := context.Background()

	cmd := adapter.StreamCmd(ctx, nil)
	result := cmd()

	_, ok := result.(adapter.StreamDoneMsg)
	if !ok {
		t.Fatalf("expected StreamDoneMsg for nil channel, got %T", result)
	}
}

// ============================================================================
// adaptMessage Tests - Table Driven
// ============================================================================

func TestAdaptMessage_StreamEvent_ContentBlockDelta_Text(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(0), // JSON numbers decode as float64
			"delta": map[string]any{
				"type": "text_delta",
				"text": "Hello",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamDeltaMsg)
	if !ok {
		t.Fatalf("expected StreamDeltaMsg, got %T", result)
	}

	if msg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", msg.Text, "Hello")
	}
	if msg.BlockType != "text" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "text")
	}
	if msg.Index != 0 {
		t.Errorf("Index = %d, want %d", msg.Index, 0)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockDelta_Thinking(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(0),
			"delta": map[string]any{
				"type":     "thinking_delta",
				"thinking": "Let me consider...",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.ThinkingDeltaMsg)
	if !ok {
		t.Fatalf("expected ThinkingDeltaMsg, got %T", result)
	}

	if msg.Text != "Let me consider..." {
		t.Errorf("Text = %q, want %q", msg.Text, "Let me consider...")
	}
	if msg.Index != 0 {
		t.Errorf("Index = %d, want %d", msg.Index, 0)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockDelta_InputJSON(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(1),
			"delta": map[string]any{
				"type":         "input_json_delta",
				"partial_json": `{"command":"ls"}`,
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.ToolUseDeltaMsg)
	if !ok {
		t.Fatalf("expected ToolUseDeltaMsg, got %T", result)
	}

	if msg.PartialJSON != `{"command":"ls"}` {
		t.Errorf("PartialJSON = %q, want %q", msg.PartialJSON, `{"command":"ls"}`)
	}
	if msg.Index != 1 {
		t.Errorf("Index = %d, want %d", msg.Index, 1)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockStart_Text(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStart,
			"index": float64(0),
			"content_block": map[string]any{
				"type": "text",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStartMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStartMsg, got %T", result)
	}

	if msg.BlockType != "text" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "text")
	}
	if msg.Index != 0 {
		t.Errorf("Index = %d, want %d", msg.Index, 0)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockStart_Thinking(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStart,
			"index": float64(0),
			"content_block": map[string]any{
				"type": "thinking",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStartMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStartMsg, got %T", result)
	}

	if msg.BlockType != "thinking" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "thinking")
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockStart_ToolUse(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStart,
			"index": float64(1),
			"content_block": map[string]any{
				"type": "tool_use",
				"name": "Read",
				"id":   "tool_abc123",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStartMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStartMsg, got %T", result)
	}

	if msg.BlockType != "tool_use" {
		t.Errorf("BlockType = %q, want %q", msg.BlockType, "tool_use")
	}
	if msg.ToolName != "Read" {
		t.Errorf("ToolName = %q, want %q", msg.ToolName, "Read")
	}
	if msg.ToolID != "tool_abc123" {
		t.Errorf("ToolID = %q, want %q", msg.ToolID, "tool_abc123")
	}
	if msg.Index != 1 {
		t.Errorf("Index = %d, want %d", msg.Index, 1)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockStop(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStop,
			"index": float64(2),
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStopMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStopMsg, got %T", result)
	}

	if msg.Index != 2 {
		t.Errorf("Index = %d, want %d", msg.Index, 2)
	}
}

func TestAdaptMessage_StreamEvent_MessageStop(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageStop,
		},
	}

	result := adapter.AdaptMessage(event)

	_, ok := result.(adapter.StreamDoneMsg)
	if !ok {
		t.Fatalf("expected StreamDoneMsg, got %T", result)
	}
}

func TestAdaptMessage_AssistantMessage(t *testing.T) {
	assistantMsg := &claudecode.AssistantMessage{}

	result := adapter.AdaptMessage(assistantMsg)

	msg, ok := result.(adapter.AssistantMsg)
	if !ok {
		t.Fatalf("expected AssistantMsg, got %T", result)
	}

	if msg.Message != assistantMsg {
		t.Errorf("Message = %v, want %v", msg.Message, assistantMsg)
	}
}

func TestAdaptMessage_ResultMessage(t *testing.T) {
	cost := 0.0025
	resultMsg := &claudecode.ResultMessage{
		TotalCostUSD: &cost,
		Usage: &map[string]any{
			"input_tokens":  float64(150),
			"output_tokens": float64(75),
		},
	}

	result := adapter.AdaptMessage(resultMsg)

	msg, ok := result.(adapter.ResultMsg)
	if !ok {
		t.Fatalf("expected ResultMsg, got %T", result)
	}

	if msg.Message != resultMsg {
		t.Errorf("Message = %v, want %v", msg.Message, resultMsg)
	}
	if msg.InputTokens != 150 {
		t.Errorf("InputTokens = %d, want %d", msg.InputTokens, 150)
	}
	if msg.OutputTokens != 75 {
		t.Errorf("OutputTokens = %d, want %d", msg.OutputTokens, 75)
	}
	if msg.TotalCost != 0.0025 {
		t.Errorf("TotalCost = %f, want %f", msg.TotalCost, 0.0025)
	}
}

func TestAdaptMessage_UnknownMessageType(t *testing.T) {
	// Create a mock type that implements claudecode.Message
	result := adapter.AdaptMessage(nil)

	if result != nil {
		t.Errorf("expected nil for nil input, got %T", result)
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestAdaptMessage_NilMessage(t *testing.T) {
	result := adapter.AdaptMessage(nil)

	if result != nil {
		t.Errorf("expected nil for nil input, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_EmptyEvent(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for empty Event, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_UnknownEventType(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": "unknown_event_type",
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for unknown event type, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_MissingDelta(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(0),
			// delta field is missing
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for missing delta, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_UnknownDeltaType(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(0),
			"delta": map[string]any{
				"type": "unknown_delta_type",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for unknown delta type, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_MissingContentBlock(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStart,
			"index": float64(0),
			// content_block field is missing
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for missing content_block, got %T", result)
	}
}

func TestAdaptMessage_ResultMessage_NilUsage(t *testing.T) {
	resultMsg := &claudecode.ResultMessage{
		Usage:        nil,
		TotalCostUSD: nil,
	}

	result := adapter.AdaptMessage(resultMsg)

	msg, ok := result.(adapter.ResultMsg)
	if !ok {
		t.Fatalf("expected ResultMsg, got %T", result)
	}

	if msg.InputTokens != 0 {
		t.Errorf("InputTokens = %d, want 0 for nil usage", msg.InputTokens)
	}
	if msg.OutputTokens != 0 {
		t.Errorf("OutputTokens = %d, want 0 for nil usage", msg.OutputTokens)
	}
	if msg.TotalCost != 0.0 {
		t.Errorf("TotalCost = %f, want 0.0 for nil usage", msg.TotalCost)
	}
}

func TestAdaptMessage_ResultMessage_PartialUsage(t *testing.T) {
	resultMsg := &claudecode.ResultMessage{
		Usage: &map[string]any{
			"input_tokens": float64(100),
			// output_tokens is missing
		},
		TotalCostUSD: nil,
	}

	result := adapter.AdaptMessage(resultMsg)

	msg, ok := result.(adapter.ResultMsg)
	if !ok {
		t.Fatalf("expected ResultMsg, got %T", result)
	}

	if msg.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want 100", msg.InputTokens)
	}
	if msg.OutputTokens != 0 {
		t.Errorf("OutputTokens = %d, want 0 for missing field", msg.OutputTokens)
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestStreamCmd_RapidMessages(t *testing.T) {
	ctx := context.Background()
	ch := make(chan claudecode.Message, 100)

	// Send 100 messages rapidly
	for i := 0; i < 100; i++ {
		ch <- &claudecode.AssistantMessage{}
	}
	close(ch)

	// Receive all messages
	received := 0
	for {
		cmd := adapter.StreamCmd(ctx, ch)
		result := cmd()

		if _, ok := result.(adapter.StreamDoneMsg); ok {
			break
		}
		if _, ok := result.(adapter.AssistantMsg); ok {
			received++
		}
	}

	if received != 100 {
		t.Errorf("received %d messages, want 100", received)
	}
}

func TestStreamCmd_ContextCancelDuringReceive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan claudecode.Message)

	// Start receiving in goroutine
	var result tea.Msg
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := adapter.StreamCmd(ctx, ch)
		result = cmd()
	}()

	// Give goroutine time to block on channel
	time.Sleep(10 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for result
	wg.Wait()

	errMsg, ok := result.(adapter.StreamErrorMsg)
	if !ok {
		t.Fatalf("expected StreamErrorMsg, got %T", result)
	}

	if errMsg.Err != context.Canceled {
		t.Errorf("Err = %v, want %v", errMsg.Err, context.Canceled)
	}
}

func TestStreamCmd_ChannelCloseRace(t *testing.T) {
	// Test that concurrent close and cancel doesn't cause issues
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan claudecode.Message)

		var result tea.Msg
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := adapter.StreamCmd(ctx, ch)
			result = cmd()
		}()

		// Race between close and cancel
		go cancel()
		go func() { close(ch) }()

		wg.Wait()

		// Should get either StreamDoneMsg or StreamErrorMsg
		switch result.(type) {
		case adapter.StreamDoneMsg, adapter.StreamErrorMsg:
			// OK
		default:
			t.Fatalf("iteration %d: unexpected result type %T", i, result)
		}
	}
}

// ============================================================================
// Index Extraction Tests
// ============================================================================

func TestAdaptMessage_StreamEvent_IndexAsInt(t *testing.T) {
	// Test when index comes as int instead of float64
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStop,
			"index": 5, // int, not float64
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStopMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStopMsg, got %T", result)
	}

	if msg.Index != 5 {
		t.Errorf("Index = %d, want %d", msg.Index, 5)
	}
}

func TestAdaptMessage_StreamEvent_MissingIndex(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeContentBlockStop,
			// index is missing
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStopMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStopMsg, got %T", result)
	}

	// Should default to 0
	if msg.Index != 0 {
		t.Errorf("Index = %d, want 0 for missing index", msg.Index)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - UserMsg
// ============================================================================

func TestUserMsg_Fields(t *testing.T) {
	msg := adapter.UserMsg{
		UUID:            "uuid-123",
		Content:         "Hello, Claude!",
		ParentToolUseID: "tool-456",
	}

	if msg.UUID != "uuid-123" {
		t.Errorf("UUID = %q, want %q", msg.UUID, "uuid-123")
	}
	if msg.Content != "Hello, Claude!" {
		t.Errorf("Content = %v, want %q", msg.Content, "Hello, Claude!")
	}
	if msg.ParentToolUseID != "tool-456" {
		t.Errorf("ParentToolUseID = %q, want %q", msg.ParentToolUseID, "tool-456")
	}
}

func TestAdaptMessage_UserMessage(t *testing.T) {
	uuid := "user-msg-uuid"
	userMsg := &claudecode.UserMessage{
		Content: "Hello from user",
		UUID:    &uuid,
	}

	result := adapter.AdaptMessage(userMsg)

	msg, ok := result.(adapter.UserMsg)
	if !ok {
		t.Fatalf("expected UserMsg, got %T", result)
	}

	if msg.UUID != "user-msg-uuid" {
		t.Errorf("UUID = %q, want %q", msg.UUID, "user-msg-uuid")
	}
	if msg.Content != "Hello from user" {
		t.Errorf("Content = %v, want %q", msg.Content, "Hello from user")
	}
}

func TestAdaptMessage_UserMessage_WithParentToolUseID(t *testing.T) {
	uuid := "user-msg-uuid"
	parentID := "parent-tool-123"
	userMsg := &claudecode.UserMessage{
		Content:         "Tool result content",
		UUID:            &uuid,
		ParentToolUseID: &parentID,
	}

	result := adapter.AdaptMessage(userMsg)

	msg, ok := result.(adapter.UserMsg)
	if !ok {
		t.Fatalf("expected UserMsg, got %T", result)
	}

	if msg.ParentToolUseID != "parent-tool-123" {
		t.Errorf("ParentToolUseID = %q, want %q", msg.ParentToolUseID, "parent-tool-123")
	}
}

func TestAdaptMessage_UserMessage_Nil(t *testing.T) {
	var userMsg *claudecode.UserMessage = nil

	result := adapter.AdaptMessage(userMsg)

	if result != nil {
		t.Errorf("expected nil for nil UserMessage, got %T", result)
	}
}

func TestAdaptMessage_UserMessage_NilFields(t *testing.T) {
	userMsg := &claudecode.UserMessage{
		Content:         "Test content",
		UUID:            nil, // nil UUID
		ParentToolUseID: nil, // nil ParentToolUseID
	}

	result := adapter.AdaptMessage(userMsg)

	msg, ok := result.(adapter.UserMsg)
	if !ok {
		t.Fatalf("expected UserMsg, got %T", result)
	}

	if msg.UUID != "" {
		t.Errorf("UUID = %q, want empty string for nil UUID", msg.UUID)
	}
	if msg.ParentToolUseID != "" {
		t.Errorf("ParentToolUseID = %q, want empty string for nil ParentToolUseID", msg.ParentToolUseID)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - SystemInitMsg
// ============================================================================

func TestSystemInitMsg_Fields(t *testing.T) {
	msg := adapter.SystemInitMsg{
		SessionID:      "session-abc",
		Model:          "claude-sonnet-4-20250514",
		Cwd:            "/home/user/project",
		Tools:          []string{"Read", "Write", "Bash"},
		PermissionMode: "default",
	}

	if msg.SessionID != "session-abc" {
		t.Errorf("SessionID = %q, want %q", msg.SessionID, "session-abc")
	}
	if msg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %q, want %q", msg.Model, "claude-sonnet-4-20250514")
	}
	if msg.Cwd != "/home/user/project" {
		t.Errorf("Cwd = %q, want %q", msg.Cwd, "/home/user/project")
	}
	if len(msg.Tools) != 3 || msg.Tools[0] != "Read" {
		t.Errorf("Tools = %v, want [Read Write Bash]", msg.Tools)
	}
	if msg.PermissionMode != "default" {
		t.Errorf("PermissionMode = %q, want %q", msg.PermissionMode, "default")
	}
}

func TestAdaptMessage_SystemMessage_Init(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "init",
		Data: map[string]any{
			"session_id":      "sess-123",
			"model":           "claude-sonnet-4-20250514",
			"cwd":             "/workspace",
			"tools":           []any{"Bash", "Read", "Write"},
			"permission_mode": "default",
		},
	}

	result := adapter.AdaptMessage(sysMsg)

	msg, ok := result.(adapter.SystemInitMsg)
	if !ok {
		t.Fatalf("expected SystemInitMsg, got %T", result)
	}

	if msg.SessionID != "sess-123" {
		t.Errorf("SessionID = %q, want %q", msg.SessionID, "sess-123")
	}
	if msg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %q, want %q", msg.Model, "claude-sonnet-4-20250514")
	}
	if msg.Cwd != "/workspace" {
		t.Errorf("Cwd = %q, want %q", msg.Cwd, "/workspace")
	}
	if len(msg.Tools) != 3 {
		t.Errorf("Tools len = %d, want 3", len(msg.Tools))
	}
	if msg.PermissionMode != "default" {
		t.Errorf("PermissionMode = %q, want %q", msg.PermissionMode, "default")
	}
}

func TestAdaptMessage_SystemMessage_Init_PartialFields(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "init",
		Data: map[string]any{
			"session_id": "sess-456",
			// Missing model, cwd, tools, permission_mode
		},
	}

	result := adapter.AdaptMessage(sysMsg)

	msg, ok := result.(adapter.SystemInitMsg)
	if !ok {
		t.Fatalf("expected SystemInitMsg, got %T", result)
	}

	if msg.SessionID != "sess-456" {
		t.Errorf("SessionID = %q, want %q", msg.SessionID, "sess-456")
	}
	if msg.Model != "" {
		t.Errorf("Model = %q, want empty string", msg.Model)
	}
	if msg.Cwd != "" {
		t.Errorf("Cwd = %q, want empty string", msg.Cwd)
	}
	if msg.Tools != nil && len(msg.Tools) != 0 {
		t.Errorf("Tools = %v, want nil or empty", msg.Tools)
	}
}

func TestAdaptMessage_SystemMessage_Nil(t *testing.T) {
	var sysMsg *claudecode.SystemMessage = nil

	result := adapter.AdaptMessage(sysMsg)

	if result != nil {
		t.Errorf("expected nil for nil SystemMessage, got %T", result)
	}
}

func TestAdaptMessage_SystemMessage_EmptySubtype(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "",
		Data:    map[string]any{},
	}

	result := adapter.AdaptMessage(sysMsg)

	if result != nil {
		t.Errorf("expected nil for empty subtype, got %T", result)
	}
}

func TestAdaptMessage_SystemMessage_UnknownSubtype(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "unknown_subtype",
		Data:    map[string]any{},
	}

	result := adapter.AdaptMessage(sysMsg)

	if result != nil {
		t.Errorf("expected nil for unknown subtype, got %T", result)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - SystemHookResponseMsg
// ============================================================================

func TestSystemHookResponseMsg_Fields(t *testing.T) {
	msg := adapter.SystemHookResponseMsg{
		SessionID: "session-xyz",
		HookName:  "pre-commit",
		HookEvent: "PreToolUse",
		Stdout:    "Hook executed successfully",
		Stderr:    "",
		ExitCode:  0,
	}

	if msg.SessionID != "session-xyz" {
		t.Errorf("SessionID = %q, want %q", msg.SessionID, "session-xyz")
	}
	if msg.HookName != "pre-commit" {
		t.Errorf("HookName = %q, want %q", msg.HookName, "pre-commit")
	}
	if msg.HookEvent != "PreToolUse" {
		t.Errorf("HookEvent = %q, want %q", msg.HookEvent, "PreToolUse")
	}
	if msg.Stdout != "Hook executed successfully" {
		t.Errorf("Stdout = %q, want %q", msg.Stdout, "Hook executed successfully")
	}
	if msg.Stderr != "" {
		t.Errorf("Stderr = %q, want empty string", msg.Stderr)
	}
	if msg.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", msg.ExitCode)
	}
}

func TestAdaptMessage_SystemMessage_HookResponse(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "hook_response",
		Data: map[string]any{
			"session_id": "sess-hook-123",
			"hook_name":  "lint-check",
			"hook_event": "PostToolUse",
			"stdout":     "Linting passed",
			"stderr":     "Warning: unused var",
			"exit_code":  float64(0), // JSON numbers decode as float64
		},
	}

	result := adapter.AdaptMessage(sysMsg)

	msg, ok := result.(adapter.SystemHookResponseMsg)
	if !ok {
		t.Fatalf("expected SystemHookResponseMsg, got %T", result)
	}

	if msg.SessionID != "sess-hook-123" {
		t.Errorf("SessionID = %q, want %q", msg.SessionID, "sess-hook-123")
	}
	if msg.HookName != "lint-check" {
		t.Errorf("HookName = %q, want %q", msg.HookName, "lint-check")
	}
	if msg.HookEvent != "PostToolUse" {
		t.Errorf("HookEvent = %q, want %q", msg.HookEvent, "PostToolUse")
	}
	if msg.Stdout != "Linting passed" {
		t.Errorf("Stdout = %q, want %q", msg.Stdout, "Linting passed")
	}
	if msg.Stderr != "Warning: unused var" {
		t.Errorf("Stderr = %q, want %q", msg.Stderr, "Warning: unused var")
	}
	if msg.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", msg.ExitCode)
	}
}

func TestAdaptMessage_SystemMessage_HookResponse_PartialFields(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "hook_response",
		Data: map[string]any{
			"hook_name": "basic-hook",
			// Missing other fields
		},
	}

	result := adapter.AdaptMessage(sysMsg)

	msg, ok := result.(adapter.SystemHookResponseMsg)
	if !ok {
		t.Fatalf("expected SystemHookResponseMsg, got %T", result)
	}

	if msg.HookName != "basic-hook" {
		t.Errorf("HookName = %q, want %q", msg.HookName, "basic-hook")
	}
	if msg.SessionID != "" {
		t.Errorf("SessionID = %q, want empty string", msg.SessionID)
	}
	if msg.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0 for missing field", msg.ExitCode)
	}
}

func TestAdaptMessage_SystemMessage_HookResponse_NonZeroExitCode(t *testing.T) {
	sysMsg := &claudecode.SystemMessage{
		Subtype: "hook_response",
		Data: map[string]any{
			"hook_name": "failing-hook",
			"stdout":    "",
			"stderr":    "Error: hook failed",
			"exit_code": float64(1),
		},
	}

	result := adapter.AdaptMessage(sysMsg)

	msg, ok := result.(adapter.SystemHookResponseMsg)
	if !ok {
		t.Fatalf("expected SystemHookResponseMsg, got %T", result)
	}

	if msg.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", msg.ExitCode)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - MessageStartMsg
// ============================================================================

func TestMessageStartMsg_Fields(t *testing.T) {
	msg := adapter.MessageStartMsg{
		MessageID:    "msg-start-123",
		Model:        "claude-sonnet-4-20250514",
		InputTokens:  100,
		OutputTokens: 0,
	}

	if msg.MessageID != "msg-start-123" {
		t.Errorf("MessageID = %q, want %q", msg.MessageID, "msg-start-123")
	}
	if msg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %q, want %q", msg.Model, "claude-sonnet-4-20250514")
	}
	if msg.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want 100", msg.InputTokens)
	}
	if msg.OutputTokens != 0 {
		t.Errorf("OutputTokens = %d, want 0", msg.OutputTokens)
	}
}

func TestAdaptMessage_StreamEvent_MessageStart(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageStart,
			"message": map[string]any{
				"id":    "msg-abc123",
				"model": "claude-sonnet-4-20250514",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.MessageStartMsg)
	if !ok {
		t.Fatalf("expected MessageStartMsg, got %T", result)
	}

	if msg.MessageID != "msg-abc123" {
		t.Errorf("MessageID = %q, want %q", msg.MessageID, "msg-abc123")
	}
	if msg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %q, want %q", msg.Model, "claude-sonnet-4-20250514")
	}
}

func TestAdaptMessage_StreamEvent_MessageStart_WithUsage(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageStart,
			"message": map[string]any{
				"id":    "msg-usage-123",
				"model": "claude-sonnet-4-20250514",
				"usage": map[string]any{
					"input_tokens":  float64(150),
					"output_tokens": float64(0),
				},
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.MessageStartMsg)
	if !ok {
		t.Fatalf("expected MessageStartMsg, got %T", result)
	}

	if msg.InputTokens != 150 {
		t.Errorf("InputTokens = %d, want 150", msg.InputTokens)
	}
	if msg.OutputTokens != 0 {
		t.Errorf("OutputTokens = %d, want 0", msg.OutputTokens)
	}
}

func TestAdaptMessage_StreamEvent_MessageStart_MissingMessage(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageStart,
			// message field is missing
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for missing message field, got %T", result)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - MessageDeltaMsg
// ============================================================================

func TestMessageDeltaMsg_Fields(t *testing.T) {
	msg := adapter.MessageDeltaMsg{
		StopReason:   "end_turn",
		InputTokens:  100,
		OutputTokens: 250,
	}

	if msg.StopReason != "end_turn" {
		t.Errorf("StopReason = %q, want %q", msg.StopReason, "end_turn")
	}
	if msg.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want 100", msg.InputTokens)
	}
	if msg.OutputTokens != 250 {
		t.Errorf("OutputTokens = %d, want 250", msg.OutputTokens)
	}
}

func TestAdaptMessage_StreamEvent_MessageDelta(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageDelta,
			"delta": map[string]any{
				"stop_reason": "end_turn",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.MessageDeltaMsg)
	if !ok {
		t.Fatalf("expected MessageDeltaMsg, got %T", result)
	}

	if msg.StopReason != "end_turn" {
		t.Errorf("StopReason = %q, want %q", msg.StopReason, "end_turn")
	}
}

func TestAdaptMessage_StreamEvent_MessageDelta_WithUsage(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageDelta,
			"delta": map[string]any{
				"stop_reason": "tool_use",
			},
			"usage": map[string]any{
				"input_tokens":  float64(100),
				"output_tokens": float64(75),
			},
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.MessageDeltaMsg)
	if !ok {
		t.Fatalf("expected MessageDeltaMsg, got %T", result)
	}

	if msg.StopReason != "tool_use" {
		t.Errorf("StopReason = %q, want %q", msg.StopReason, "tool_use")
	}
	if msg.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want 100", msg.InputTokens)
	}
	if msg.OutputTokens != 75 {
		t.Errorf("OutputTokens = %d, want 75", msg.OutputTokens)
	}
}

func TestAdaptMessage_StreamEvent_MessageDelta_MissingDelta(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type": claudecode.StreamEventTypeMessageDelta,
			// delta field is missing
		},
	}

	result := adapter.AdaptMessage(event)

	// Should still return a MessageDeltaMsg but with empty/zero values
	msg, ok := result.(adapter.MessageDeltaMsg)
	if !ok {
		t.Fatalf("expected MessageDeltaMsg, got %T", result)
	}

	if msg.StopReason != "" {
		t.Errorf("StopReason = %q, want empty string", msg.StopReason)
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - ControlRequestMsg
// ============================================================================

func TestControlRequestMsg_Fields(t *testing.T) {
	msg := adapter.ControlRequestMsg{
		RequestID: "req-123",
		Subtype:   "can_use_tool",
	}

	if msg.RequestID != "req-123" {
		t.Errorf("RequestID = %q, want %q", msg.RequestID, "req-123")
	}
	if msg.Subtype != "can_use_tool" {
		t.Errorf("Subtype = %q, want %q", msg.Subtype, "can_use_tool")
	}
}

func TestAdaptMessage_ControlRequest(t *testing.T) {
	ctrlMsg := &claudecode.RawControlMessage{
		MessageType: "control_request",
		Data: map[string]any{
			"request_id": "ctrl-req-456",
			"subtype":    "can_use_tool",
		},
	}

	result := adapter.AdaptMessage(ctrlMsg)

	msg, ok := result.(adapter.ControlRequestMsg)
	if !ok {
		t.Fatalf("expected ControlRequestMsg, got %T", result)
	}

	if msg.RequestID != "ctrl-req-456" {
		t.Errorf("RequestID = %q, want %q", msg.RequestID, "ctrl-req-456")
	}
	if msg.Subtype != "can_use_tool" {
		t.Errorf("Subtype = %q, want %q", msg.Subtype, "can_use_tool")
	}
}

// ============================================================================
// NEW MESSAGE TYPE TESTS - ControlResponseMsg
// ============================================================================

func TestControlResponseMsg_Fields(t *testing.T) {
	msg := adapter.ControlResponseMsg{
		RequestID: "req-789",
		IsSuccess: true,
		Error:     "",
	}

	if msg.RequestID != "req-789" {
		t.Errorf("RequestID = %q, want %q", msg.RequestID, "req-789")
	}
	if !msg.IsSuccess {
		t.Errorf("IsSuccess = %v, want true", msg.IsSuccess)
	}
	if msg.Error != "" {
		t.Errorf("Error = %q, want empty string", msg.Error)
	}
}

func TestAdaptMessage_ControlResponse_Success(t *testing.T) {
	ctrlMsg := &claudecode.RawControlMessage{
		MessageType: "control_response",
		Data: map[string]any{
			"request_id": "ctrl-resp-123",
			"response": map[string]any{
				"subtype": "success",
			},
		},
	}

	result := adapter.AdaptMessage(ctrlMsg)

	msg, ok := result.(adapter.ControlResponseMsg)
	if !ok {
		t.Fatalf("expected ControlResponseMsg, got %T", result)
	}

	if msg.RequestID != "ctrl-resp-123" {
		t.Errorf("RequestID = %q, want %q", msg.RequestID, "ctrl-resp-123")
	}
	if !msg.IsSuccess {
		t.Errorf("IsSuccess = %v, want true", msg.IsSuccess)
	}
}

func TestAdaptMessage_ControlResponse_Error(t *testing.T) {
	ctrlMsg := &claudecode.RawControlMessage{
		MessageType: "control_response",
		Data: map[string]any{
			"request_id": "ctrl-resp-err",
			"response": map[string]any{
				"subtype": "error",
				"error":   "Permission denied",
			},
		},
	}

	result := adapter.AdaptMessage(ctrlMsg)

	msg, ok := result.(adapter.ControlResponseMsg)
	if !ok {
		t.Fatalf("expected ControlResponseMsg, got %T", result)
	}

	if msg.RequestID != "ctrl-resp-err" {
		t.Errorf("RequestID = %q, want %q", msg.RequestID, "ctrl-resp-err")
	}
	if msg.IsSuccess {
		t.Errorf("IsSuccess = %v, want false", msg.IsSuccess)
	}
	if msg.Error != "Permission denied" {
		t.Errorf("Error = %q, want %q", msg.Error, "Permission denied")
	}
}

func TestAdaptMessage_ControlMessage_MissingData(t *testing.T) {
	ctrlMsg := &claudecode.RawControlMessage{
		MessageType: "control_request",
		Data:        map[string]any{},
	}

	result := adapter.AdaptMessage(ctrlMsg)

	msg, ok := result.(adapter.ControlRequestMsg)
	if !ok {
		t.Fatalf("expected ControlRequestMsg, got %T", result)
	}

	// Should have empty values for missing fields
	if msg.RequestID != "" {
		t.Errorf("RequestID = %q, want empty string", msg.RequestID)
	}
	if msg.Subtype != "" {
		t.Errorf("Subtype = %q, want empty string", msg.Subtype)
	}
}

func TestAdaptMessage_ControlMessage_Nil(t *testing.T) {
	var ctrlMsg *claudecode.RawControlMessage = nil

	result := adapter.AdaptMessage(ctrlMsg)

	if result != nil {
		t.Errorf("expected nil for nil RawControlMessage, got %T", result)
	}
}

func TestAdaptMessage_ControlMessage_UnknownType(t *testing.T) {
	ctrlMsg := &claudecode.RawControlMessage{
		MessageType: "unknown_control_type",
		Data:        map[string]any{},
	}

	result := adapter.AdaptMessage(ctrlMsg)

	if result != nil {
		t.Errorf("expected nil for unknown control message type, got %T", result)
	}
}

// ============================================================================
// Additional Edge Case Tests for 100% Coverage
// ============================================================================

func TestAdaptMessage_StreamEvent_NilEvent(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: nil,
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for nil Event field, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_NilStreamEvent(t *testing.T) {
	var event *claudecode.StreamEvent = nil

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for nil StreamEvent, got %T", result)
	}
}

func TestAdaptMessage_StreamEvent_ContentBlockDelta_MissingDeltaType(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockDelta,
			"index": float64(0),
			"delta": map[string]any{
				// type field is missing
				"text": "some text",
			},
		},
	}

	result := adapter.AdaptMessage(event)

	if result != nil {
		t.Errorf("expected nil for missing delta type, got %T", result)
	}
}

func TestAdaptMessage_IndexAsInt64(t *testing.T) {
	event := &claudecode.StreamEvent{
		Event: map[string]any{
			"type":  claudecode.StreamEventTypeContentBlockStop,
			"index": int64(7), // int64 type
		},
	}

	result := adapter.AdaptMessage(event)

	msg, ok := result.(adapter.StreamBlockStopMsg)
	if !ok {
		t.Fatalf("expected StreamBlockStopMsg, got %T", result)
	}

	if msg.Index != 7 {
		t.Errorf("Index = %d, want 7", msg.Index)
	}
}

// mockMessage implements claudecode.Message for testing the default case
type mockMessage struct{}

func (m *mockMessage) Type() string { return "mock" }

func TestAdaptMessage_DefaultCase_UnknownType(t *testing.T) {
	// Use a mock message type that isn't handled by any case
	msg := &mockMessage{}

	result := adapter.AdaptMessage(msg)

	if result != nil {
		t.Errorf("expected nil for unknown message type, got %T", result)
	}
}
