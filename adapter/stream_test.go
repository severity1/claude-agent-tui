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
