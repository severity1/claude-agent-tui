package adapter

import (
	"context"
	"fmt"

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
// Note: ToolName/ToolID come from StreamBlockStartMsg, not delta events.
// The TUI layer tracks state to associate tool info with deltas by index.
type ToolUseDeltaMsg struct {
	PartialJSON string
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

// UserMsg contains a user message (typically tool results).
type UserMsg struct {
	UUID            string
	Content         any // string or content blocks
	ParentToolUseID string
}

// SystemInitMsg contains session initialization data.
type SystemInitMsg struct {
	SessionID      string
	Model          string
	Cwd            string
	Tools          []string
	PermissionMode string
}

// SystemHookResponseMsg contains hook execution results.
type SystemHookResponseMsg struct {
	SessionID string
	HookName  string
	HookEvent string
	Stdout    string
	Stderr    string
	ExitCode  int
}

// MessageStartMsg contains initial message metadata with model info.
type MessageStartMsg struct {
	MessageID    string
	Model        string
	InputTokens  int
	OutputTokens int
}

// MessageDeltaMsg contains final usage stats and stop reason.
type MessageDeltaMsg struct {
	StopReason   string
	InputTokens  int
	OutputTokens int
}

// ControlRequestMsg contains control protocol request (informational).
type ControlRequestMsg struct {
	RequestID string
	Subtype   string
}

// ControlResponseMsg contains control protocol response (informational).
type ControlResponseMsg struct {
	RequestID string
	IsSuccess bool
	Error     string
}

// StreamDoneMsg signals stream completion.
type StreamDoneMsg struct{}

// StreamErrorMsg signals a stream error.
type StreamErrorMsg struct {
	Err error
}

// UnknownMessageMsg indicates an unrecognized SDK message type.
// This is returned instead of nil for unknown types to aid debugging.
type UnknownMessageMsg struct {
	TypeName string // Go type name of the unrecognized message
}

// Error implements the error interface for StreamErrorMsg.
func (e StreamErrorMsg) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown stream error"
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
	case *claudecode.UserMessage:
		return adaptUserMessage(m)
	case *claudecode.SystemMessage:
		return adaptSystemMessage(m)
	case *claudecode.RawControlMessage:
		return adaptControlMessage(m)
	default:
		return UnknownMessageMsg{TypeName: fmt.Sprintf("%T", msg)}
	}
}

// adaptStreamEvent converts a StreamEvent to the appropriate tea.Msg
// based on the event type.
func adaptStreamEvent(event *claudecode.StreamEvent) tea.Msg {
	if event == nil || event.Event == nil {
		return UnknownMessageMsg{TypeName: "StreamEvent/nil"}
	}

	eventType, ok := event.Event["type"].(string)
	if !ok {
		return UnknownMessageMsg{TypeName: fmt.Sprintf("StreamEvent/non-string-type:%T", event.Event["type"])}
	}

	switch eventType {
	case claudecode.StreamEventTypeContentBlockStart:
		return adaptBlockStart(event)
	case claudecode.StreamEventTypeContentBlockDelta:
		return adaptBlockDelta(event)
	case claudecode.StreamEventTypeContentBlockStop:
		return adaptBlockStop(event)
	case claudecode.StreamEventTypeMessageStart:
		return adaptMessageStart(event)
	case claudecode.StreamEventTypeMessageDelta:
		return adaptMessageDelta(event)
	case claudecode.StreamEventTypeMessageStop:
		return StreamDoneMsg{}
	default:
		return UnknownMessageMsg{TypeName: fmt.Sprintf("StreamEvent/%s", eventType)}
	}
}

// adaptBlockStart handles content_block_start events.
func adaptBlockStart(event *claudecode.StreamEvent) tea.Msg {
	contentBlock, ok := event.Event["content_block"].(map[string]any)
	if !ok {
		return UnknownMessageMsg{TypeName: "StreamEvent/content_block_start/missing-content_block"}
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
		return UnknownMessageMsg{TypeName: "StreamEvent/content_block_delta/missing-delta"}
	}

	deltaType, ok := delta["type"].(string)
	if !ok {
		return UnknownMessageMsg{TypeName: "StreamEvent/content_block_delta/missing-delta-type"}
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
		return UnknownMessageMsg{TypeName: fmt.Sprintf("StreamEvent/content_block_delta/%s", deltaType)}
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

// adaptUserMessage converts a UserMessage to UserMsg.
func adaptUserMessage(msg *claudecode.UserMessage) tea.Msg {
	if msg == nil {
		return UnknownMessageMsg{TypeName: "UserMessage/nil"}
	}

	return UserMsg{
		UUID:            msg.GetUUID(),
		Content:         msg.Content,
		ParentToolUseID: msg.GetParentToolUseID(),
	}
}

// adaptSystemMessage converts a SystemMessage to the appropriate tea.Msg
// based on the subtype.
func adaptSystemMessage(msg *claudecode.SystemMessage) tea.Msg {
	if msg == nil {
		return UnknownMessageMsg{TypeName: "SystemMessage/nil"}
	}

	switch msg.Subtype {
	case "init":
		return adaptSystemInit(msg)
	case "hook_response":
		return adaptHookResponse(msg)
	default:
		return UnknownMessageMsg{TypeName: fmt.Sprintf("SystemMessage/%s", msg.Subtype)}
	}
}

// adaptSystemInit extracts session initialization data from SystemMessage.
func adaptSystemInit(msg *claudecode.SystemMessage) tea.Msg {
	return SystemInitMsg{
		SessionID:      toString(msg.Data["session_id"]),
		Model:          toString(msg.Data["model"]),
		Cwd:            toString(msg.Data["cwd"]),
		Tools:          toStringSlice(msg.Data["tools"]),
		PermissionMode: toString(msg.Data["permission_mode"]),
	}
}

// adaptHookResponse extracts hook execution results from SystemMessage.
func adaptHookResponse(msg *claudecode.SystemMessage) tea.Msg {
	return SystemHookResponseMsg{
		SessionID: toString(msg.Data["session_id"]),
		HookName:  toString(msg.Data["hook_name"]),
		HookEvent: toString(msg.Data["hook_event"]),
		Stdout:    toString(msg.Data["stdout"]),
		Stderr:    toString(msg.Data["stderr"]),
		ExitCode:  toInt(msg.Data["exit_code"]),
	}
}

// adaptMessageStart handles message_start stream events.
func adaptMessageStart(event *claudecode.StreamEvent) tea.Msg {
	message, ok := event.Event["message"].(map[string]any)
	if !ok {
		return UnknownMessageMsg{TypeName: "StreamEvent/message_start/missing-message"}
	}

	msg := MessageStartMsg{
		MessageID: toString(message["id"]),
		Model:     toString(message["model"]),
	}

	// Extract usage if present
	if usage, ok := message["usage"].(map[string]any); ok {
		msg.InputTokens = toInt(usage["input_tokens"])
		msg.OutputTokens = toInt(usage["output_tokens"])
	}

	return msg
}

// adaptMessageDelta handles message_delta stream events.
func adaptMessageDelta(event *claudecode.StreamEvent) tea.Msg {
	msg := MessageDeltaMsg{}

	// Extract stop reason from delta
	if delta, ok := event.Event["delta"].(map[string]any); ok {
		msg.StopReason = toString(delta["stop_reason"])
	}

	// Extract usage if present
	if usage, ok := event.Event["usage"].(map[string]any); ok {
		msg.InputTokens = toInt(usage["input_tokens"])
		msg.OutputTokens = toInt(usage["output_tokens"])
	}

	return msg
}

// adaptControlMessage converts a RawControlMessage to the appropriate tea.Msg.
func adaptControlMessage(msg *claudecode.RawControlMessage) tea.Msg {
	if msg == nil {
		return UnknownMessageMsg{TypeName: "RawControlMessage/nil"}
	}

	switch msg.MessageType {
	case claudecode.MessageTypeControlRequest:
		return adaptControlRequest(msg)
	case claudecode.MessageTypeControlResponse:
		return adaptControlResponse(msg)
	default:
		return UnknownMessageMsg{TypeName: fmt.Sprintf("RawControlMessage/%s", msg.MessageType)}
	}
}

// adaptControlRequest extracts control request data.
func adaptControlRequest(msg *claudecode.RawControlMessage) tea.Msg {
	return ControlRequestMsg{
		RequestID: toString(msg.Data["request_id"]),
		Subtype:   toString(msg.Data["subtype"]),
	}
}

// adaptControlResponse extracts control response data.
func adaptControlResponse(msg *claudecode.RawControlMessage) tea.Msg {
	requestID := toString(msg.Data["request_id"])
	isSuccess := false
	errStr := ""

	if response, ok := msg.Data["response"].(map[string]any); ok {
		subtype := toString(response["subtype"])
		isSuccess = subtype == "success"
		errStr = toString(response["error"])
	}

	return ControlResponseMsg{
		RequestID: requestID,
		IsSuccess: isSuccess,
		Error:     errStr,
	}
}

// toString safely converts an interface{} to string.
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// toStringSlice converts a []any to []string.
func toStringSlice(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}

	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
