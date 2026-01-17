package adapter

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	claudecode "github.com/severity1/claude-agent-sdk-go"
)

// Block type constants for content blocks.
const (
	BlockTypeText     = "text"
	BlockTypeThinking = "thinking"
	BlockTypeToolUse  = "tool_use"
)

// Delta type constants for content block deltas.
const (
	DeltaTypeText      = "text_delta"
	DeltaTypeThinking  = "thinking_delta"
	DeltaTypeInputJSON = "input_json_delta"
)

// StreamDeltaMsg contains incremental text from streaming.
type StreamDeltaMsg struct {
	Text      string
	BlockType string // "text", "thinking", "tool_use"
	Index     int
}

// StreamBlockStartMsg signals a new content block.
type StreamBlockStartMsg struct {
	BlockType string // "text", "thinking", "tool_use"
	Index     int
	ToolName  string // Only for tool_use blocks
	ToolID    string // Only for tool_use blocks
}

// StreamBlockStopMsg signals content block completion.
type StreamBlockStopMsg struct {
	Index int
}

// ThinkingDeltaMsg contains incremental thinking text.
type ThinkingDeltaMsg struct {
	Text  string
	Index int
}

// ToolUseDeltaMsg contains incremental tool input JSON.
type ToolUseDeltaMsg struct {
	PartialJSON string
	ToolName    string
	ToolID      string
	Index       int
}

// AssistantMsg contains a complete assistant message.
type AssistantMsg struct {
	Message *claudecode.AssistantMessage
}

// ResultMsg contains the final result with usage stats.
type ResultMsg struct {
	Message      *claudecode.ResultMessage
	InputTokens  int
	OutputTokens int
	TotalCost    float64
}

// StreamDoneMsg signals stream completion.
type StreamDoneMsg struct{}

// StreamErrorMsg signals a stream error.
type StreamErrorMsg struct {
	Err error
}

// StreamCmd returns a tea.Cmd that reads from SDK message channel.
// Call this in a loop to continuously receive messages.
func StreamCmd(ctx context.Context, ch <-chan claudecode.Message) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return StreamDoneMsg{}
		}

		select {
		case msg, ok := <-ch:
			if !ok {
				return StreamDoneMsg{}
			}
			return AdaptMessage(msg)
		case <-ctx.Done():
			return StreamErrorMsg{Err: ctx.Err()}
		}
	}
}

// AdaptMessage converts SDK message to appropriate tea.Msg.
// This function is exported for testing purposes.
func AdaptMessage(msg claudecode.Message) tea.Msg {
	if msg == nil {
		return nil
	}

	switch m := msg.(type) {
	case *claudecode.StreamEvent:
		return adaptStreamEvent(m)
	case *claudecode.AssistantMessage:
		return AssistantMsg{Message: m}
	case *claudecode.ResultMessage:
		return adaptResultMessage(m)
	default:
		return nil
	}
}

// adaptStreamEvent converts a StreamEvent to the appropriate tea.Msg
// based on the event type.
func adaptStreamEvent(event *claudecode.StreamEvent) tea.Msg {
	if event == nil || event.Event == nil {
		return nil
	}

	eventType, ok := event.Event["type"].(string)
	if !ok {
		return nil
	}

	switch eventType {
	case claudecode.StreamEventTypeContentBlockStart:
		return adaptBlockStart(event)
	case claudecode.StreamEventTypeContentBlockDelta:
		return adaptBlockDelta(event)
	case claudecode.StreamEventTypeContentBlockStop:
		return adaptBlockStop(event)
	case claudecode.StreamEventTypeMessageStop:
		return StreamDoneMsg{}
	default:
		return nil
	}
}

// adaptBlockStart handles content_block_start events.
func adaptBlockStart(event *claudecode.StreamEvent) tea.Msg {
	contentBlock, ok := event.Event["content_block"].(map[string]any)
	if !ok {
		return nil
	}

	blockType, _ := contentBlock["type"].(string)
	index := extractIndex(event.Event)

	msg := StreamBlockStartMsg{
		BlockType: blockType,
		Index:     index,
	}

	// Extract tool info for tool_use blocks
	if blockType == BlockTypeToolUse {
		if name, ok := contentBlock["name"].(string); ok {
			msg.ToolName = name
		}
		if id, ok := contentBlock["id"].(string); ok {
			msg.ToolID = id
		}
	}

	return msg
}

// adaptBlockDelta handles content_block_delta events.
func adaptBlockDelta(event *claudecode.StreamEvent) tea.Msg {
	delta, ok := event.Event["delta"].(map[string]any)
	if !ok {
		return nil
	}

	deltaType, ok := delta["type"].(string)
	if !ok {
		return nil
	}

	index := extractIndex(event.Event)

	switch deltaType {
	case DeltaTypeText:
		text, _ := delta["text"].(string)
		return StreamDeltaMsg{
			Text:      text,
			BlockType: BlockTypeText,
			Index:     index,
		}

	case DeltaTypeThinking:
		text, _ := delta["thinking"].(string)
		return ThinkingDeltaMsg{
			Text:  text,
			Index: index,
		}

	case DeltaTypeInputJSON:
		partialJSON, _ := delta["partial_json"].(string)
		return ToolUseDeltaMsg{
			PartialJSON: partialJSON,
			Index:       index,
		}

	default:
		return nil
	}
}

// adaptBlockStop handles content_block_stop events.
func adaptBlockStop(event *claudecode.StreamEvent) tea.Msg {
	return StreamBlockStopMsg{
		Index: extractIndex(event.Event),
	}
}

// adaptResultMessage converts a ResultMessage to ResultMsg with usage stats.
func adaptResultMessage(result *claudecode.ResultMessage) tea.Msg {
	msg := ResultMsg{
		Message: result,
	}

	// Extract total cost if available
	if result.TotalCostUSD != nil {
		msg.TotalCost = *result.TotalCostUSD
	}

	// Extract usage tokens if available
	if result.Usage != nil {
		usage := *result.Usage
		if inputTokens, ok := usage["input_tokens"]; ok {
			msg.InputTokens = toInt(inputTokens)
		}
		if outputTokens, ok := usage["output_tokens"]; ok {
			msg.OutputTokens = toInt(outputTokens)
		}
	}

	return msg
}

// extractIndex extracts the index field from an event map.
// Returns 0 if the field is missing or not a number.
func extractIndex(event map[string]any) int {
	return toInt(event["index"])
}

// toInt converts a value to int, handling both int and float64 types
// (JSON numbers decode as float64).
func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	case int64:
		return int(n)
	default:
		return 0
	}
}
