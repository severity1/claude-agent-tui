// Package main provides an example demonstrating the stream adapter
// by parsing live Claude Code output and showing adapted message types.
//
// NOTE: This example uses `claude --output-format stream-json` which outputs
// complete messages (system, assistant, user, result). The granular streaming
// events (content_block_start, content_block_delta, message_start, etc.) are
// only available when using the SDK with IncludePartialMessages option enabled.
// Those events would produce MessageStartMsg, MessageDeltaMsg, StreamDeltaMsg, etc.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	claudecode "github.com/severity1/claude-agent-sdk-go"
	"github.com/severity1/claude-agent-tui/adapter"
)

func main() {
	fmt.Println("Stream Adapter Validation Example")
	fmt.Println("==================================")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"<prompt>\"")
		fmt.Println("       go run main.go --validate")
		fmt.Println()
		fmt.Println("Example: go run main.go \"say hello\"")
		fmt.Println("         go run main.go --validate  (tests all message types)")
		os.Exit(1)
	}

	// Validation mode - test all message types synthetically
	if os.Args[1] == "--validate" {
		runValidation()
		return
	}

	prompt := os.Args[1]
	fmt.Printf("Prompt: %s\n\n", prompt)

	// Run Claude Code with stream-json output
	cmd := exec.Command("claude", "-p", "--output-format", "stream-json", "--verbose", prompt)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting claude: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stdout)
	// Increase buffer size for large messages
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	msgCounts := make(map[string]int)

	fmt.Println("--- Streaming Messages ---")
	fmt.Println()

	for scanner.Scan() {
		line := scanner.Text()

		// Parse JSON to determine message type
		var raw map[string]any
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		msgType, _ := raw["type"].(string)

		// Create appropriate SDK message and adapt
		var sdkMsg claudecode.Message
		switch msgType {
		case "system":
			sdkMsg = parseSystemMessage(raw)
		case "user":
			sdkMsg = parseUserMessage(raw)
		case "assistant":
			sdkMsg = parseAssistantMessage(raw)
		case "result":
			sdkMsg = parseResultMessage(raw)
		case "stream_event":
			sdkMsg = parseStreamEvent(raw)
		case "control_request":
			sdkMsg = parseControlRequest(raw)
		case "control_response":
			sdkMsg = parseControlResponse(raw)
		default:
			fmt.Printf("[UNKNOWN] type=%s\n", msgType)
			continue
		}

		// Adapt and display
		adapted := adapter.AdaptMessage(sdkMsg)
		if adapted != nil {
			typeName := fmt.Sprintf("%T", adapted)
			msgCounts[typeName]++
			printAdaptedMessage(adapted)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner error: %v\n", err)
	}

	cmd.Wait()

	// Summary
	fmt.Println()
	fmt.Println("--- Message Type Summary ---")
	for t, count := range msgCounts {
		fmt.Printf("  %s: %d\n", t, count)
	}
}

func printAdaptedMessage(msg any) {
	switch m := msg.(type) {
	case adapter.SystemInitMsg:
		fmt.Printf("[SystemInitMsg] session=%s model=%s cwd=%s tools=%d\n",
			truncate(m.SessionID, 20), m.Model, truncate(m.Cwd, 30), len(m.Tools))

	case adapter.SystemHookResponseMsg:
		fmt.Printf("[SystemHookResponseMsg] hook=%s event=%s exit=%d\n",
			m.HookName, m.HookEvent, m.ExitCode)

	case adapter.UserMsg:
		fmt.Printf("[UserMsg] uuid=%s parent=%s\n",
			truncate(m.UUID, 20), truncate(m.ParentToolUseID, 20))

	case adapter.MessageStartMsg:
		fmt.Printf("[MessageStartMsg] id=%s model=%s input=%d output=%d\n",
			truncate(m.MessageID, 20), m.Model, m.InputTokens, m.OutputTokens)

	case adapter.MessageDeltaMsg:
		fmt.Printf("[MessageDeltaMsg] stop=%s input=%d output=%d\n",
			m.StopReason, m.InputTokens, m.OutputTokens)

	case adapter.StreamBlockStartMsg:
		fmt.Printf("[StreamBlockStartMsg] type=%s index=%d tool=%s\n",
			m.BlockType, m.Index, m.ToolName)

	case adapter.StreamDeltaMsg:
		fmt.Printf("[StreamDeltaMsg] type=%s index=%d text=%q\n",
			m.BlockType, m.Index, truncate(m.Text, 50))

	case adapter.ThinkingDeltaMsg:
		fmt.Printf("[ThinkingDeltaMsg] index=%d text=%q\n",
			m.Index, truncate(m.Text, 50))

	case adapter.ToolUseDeltaMsg:
		fmt.Printf("[ToolUseDeltaMsg] index=%d json=%q\n",
			m.Index, truncate(m.PartialJSON, 50))

	case adapter.StreamBlockStopMsg:
		fmt.Printf("[StreamBlockStopMsg] index=%d\n", m.Index)

	case adapter.StreamDoneMsg:
		fmt.Printf("[StreamDoneMsg]\n")

	case adapter.AssistantMsg:
		fmt.Printf("[AssistantMsg] model=%s\n", m.Message.Model)

	case adapter.ResultMsg:
		fmt.Printf("[ResultMsg] input=%d output=%d cost=$%.4f\n",
			m.InputTokens, m.OutputTokens, m.TotalCost)

	case adapter.ControlRequestMsg:
		fmt.Printf("[ControlRequestMsg] id=%s subtype=%s\n",
			truncate(m.RequestID, 20), m.Subtype)

	case adapter.ControlResponseMsg:
		fmt.Printf("[ControlResponseMsg] id=%s success=%v\n",
			truncate(m.RequestID, 20), m.IsSuccess)

	case adapter.StreamErrorMsg:
		fmt.Printf("[StreamErrorMsg] err=%v\n", m.Err)

	default:
		fmt.Printf("[%T] %+v\n", m, m)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Parse functions convert raw JSON to SDK message types

func parseSystemMessage(raw map[string]any) *claudecode.SystemMessage {
	subtype, _ := raw["subtype"].(string)
	return &claudecode.SystemMessage{
		Subtype: subtype,
		Data:    raw,
	}
}

func parseUserMessage(raw map[string]any) *claudecode.UserMessage {
	msg := &claudecode.UserMessage{
		Content: raw["content"],
	}
	if uuid, ok := raw["uuid"].(string); ok {
		msg.UUID = &uuid
	}
	if parentID, ok := raw["parent_tool_use_id"].(string); ok {
		msg.ParentToolUseID = &parentID
	}
	return msg
}

func parseAssistantMessage(raw map[string]any) *claudecode.AssistantMessage {
	// Model is nested in message.model for stream-json output
	msg := &claudecode.AssistantMessage{}
	if message, ok := raw["message"].(map[string]any); ok {
		if model, ok := message["model"].(string); ok {
			msg.Model = model
		}
	}
	return msg
}

func parseResultMessage(raw map[string]any) *claudecode.ResultMessage {
	msg := &claudecode.ResultMessage{}
	if cost, ok := raw["total_cost_usd"].(float64); ok {
		msg.TotalCostUSD = &cost
	}
	if usage, ok := raw["usage"].(map[string]any); ok {
		msg.Usage = &usage
	}
	return msg
}

func parseStreamEvent(raw map[string]any) *claudecode.StreamEvent {
	event, _ := raw["event"].(map[string]any)
	msg := &claudecode.StreamEvent{
		Event: event,
	}
	if uuid, ok := raw["uuid"].(string); ok {
		msg.UUID = uuid
	}
	if sessionID, ok := raw["session_id"].(string); ok {
		msg.SessionID = sessionID
	}
	return msg
}

func parseControlRequest(raw map[string]any) *claudecode.RawControlMessage {
	return &claudecode.RawControlMessage{
		MessageType: "control_request",
		Data:        raw,
	}
}

func parseControlResponse(raw map[string]any) *claudecode.RawControlMessage {
	return &claudecode.RawControlMessage{
		MessageType: "control_response",
		Data:        raw,
	}
}

// runValidation tests all adapter message types synthetically
func runValidation() {
	fmt.Println("Validating all 16 adapter message types...")
	fmt.Println()

	passed := 0
	failed := 0

	// Test each message type
	tests := []struct {
		name    string
		sdkMsg  claudecode.Message
		check   func(any) bool
	}{
		{
			name: "UserMsg",
			sdkMsg: func() *claudecode.UserMessage {
				uuid := "test-uuid"
				parent := "test-parent"
				return &claudecode.UserMessage{
					Content:         "test content",
					UUID:            &uuid,
					ParentToolUseID: &parent,
				}
			}(),
			check: func(m any) bool {
				msg, ok := m.(adapter.UserMsg)
				return ok && msg.UUID == "test-uuid" && msg.ParentToolUseID == "test-parent"
			},
		},
		{
			name: "SystemInitMsg",
			sdkMsg: &claudecode.SystemMessage{
				Subtype: "init",
				Data: map[string]any{
					"session_id": "sess-123",
					"model":      "claude-test",
					"cwd":        "/test",
					"tools":      []any{"Read", "Write"},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.SystemInitMsg)
				return ok && msg.SessionID == "sess-123" && msg.Model == "claude-test" && len(msg.Tools) == 2
			},
		},
		{
			name: "SystemHookResponseMsg",
			sdkMsg: &claudecode.SystemMessage{
				Subtype: "hook_response",
				Data: map[string]any{
					"hook_name":  "test-hook",
					"hook_event": "TestEvent",
					"exit_code":  float64(0),
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.SystemHookResponseMsg)
				return ok && msg.HookName == "test-hook" && msg.HookEvent == "TestEvent"
			},
		},
		{
			name: "AssistantMsg",
			sdkMsg: &claudecode.AssistantMessage{
				Model: "claude-test",
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.AssistantMsg)
				return ok && msg.Message.Model == "claude-test"
			},
		},
		{
			name: "ResultMsg",
			sdkMsg: func() *claudecode.ResultMessage {
				cost := 0.05
				usage := map[string]any{"input_tokens": float64(100), "output_tokens": float64(50)}
				return &claudecode.ResultMessage{
					TotalCostUSD: &cost,
					Usage:        &usage,
				}
			}(),
			check: func(m any) bool {
				msg, ok := m.(adapter.ResultMsg)
				return ok && msg.TotalCost == 0.05 && msg.InputTokens == 100 && msg.OutputTokens == 50
			},
		},
		{
			name: "MessageStartMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type": "message_start",
					"message": map[string]any{
						"id":    "msg-123",
						"model": "claude-test",
						"usage": map[string]any{"input_tokens": float64(50)},
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.MessageStartMsg)
				return ok && msg.MessageID == "msg-123" && msg.Model == "claude-test" && msg.InputTokens == 50
			},
		},
		{
			name: "MessageDeltaMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type": "message_delta",
					"delta": map[string]any{
						"stop_reason": "end_turn",
					},
					"usage": map[string]any{"output_tokens": float64(75)},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.MessageDeltaMsg)
				return ok && msg.StopReason == "end_turn" && msg.OutputTokens == 75
			},
		},
		{
			name: "StreamBlockStartMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_start",
					"index": float64(0),
					"content_block": map[string]any{
						"type": "tool_use",
						"name": "Bash",
						"id":   "tool-123",
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.StreamBlockStartMsg)
				return ok && msg.BlockType == "tool_use" && msg.ToolName == "Bash" && msg.ToolID == "tool-123"
			},
		},
		{
			name: "StreamDeltaMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_delta",
					"index": float64(0),
					"delta": map[string]any{
						"type": "text_delta",
						"text": "Hello world",
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.StreamDeltaMsg)
				return ok && msg.Text == "Hello world" && msg.BlockType == "text"
			},
		},
		{
			name: "ThinkingDeltaMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_delta",
					"index": float64(0),
					"delta": map[string]any{
						"type":     "thinking_delta",
						"thinking": "Let me think...",
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.ThinkingDeltaMsg)
				return ok && msg.Text == "Let me think..."
			},
		},
		{
			name: "ToolUseDeltaMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_delta",
					"index": float64(1),
					"delta": map[string]any{
						"type":         "input_json_delta",
						"partial_json": `{"cmd":"ls"}`,
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.ToolUseDeltaMsg)
				return ok && msg.PartialJSON == `{"cmd":"ls"}` && msg.Index == 1
			},
		},
		{
			name: "StreamBlockStopMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_stop",
					"index": float64(2),
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.StreamBlockStopMsg)
				return ok && msg.Index == 2
			},
		},
		{
			name: "StreamDoneMsg",
			sdkMsg: &claudecode.StreamEvent{
				Event: map[string]any{
					"type": "message_stop",
				},
			},
			check: func(m any) bool {
				_, ok := m.(adapter.StreamDoneMsg)
				return ok
			},
		},
		{
			name: "ControlRequestMsg",
			sdkMsg: &claudecode.RawControlMessage{
				MessageType: "control_request",
				Data: map[string]any{
					"request_id": "req-123",
					"subtype":    "can_use_tool",
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.ControlRequestMsg)
				return ok && msg.RequestID == "req-123" && msg.Subtype == "can_use_tool"
			},
		},
		{
			name: "ControlResponseMsg",
			sdkMsg: &claudecode.RawControlMessage{
				MessageType: "control_response",
				Data: map[string]any{
					"request_id": "req-456",
					"response": map[string]any{
						"subtype": "success",
					},
				},
			},
			check: func(m any) bool {
				msg, ok := m.(adapter.ControlResponseMsg)
				return ok && msg.RequestID == "req-456" && msg.IsSuccess
			},
		},
	}

	// Run all tests
	for _, tt := range tests {
		adapted := adapter.AdaptMessage(tt.sdkMsg)
		if adapted == nil {
			fmt.Printf("  FAIL: %s - AdaptMessage returned nil\n", tt.name)
			failed++
			continue
		}

		if tt.check(adapted) {
			fmt.Printf("  PASS: %s\n", tt.name)
			printAdaptedMessage(adapted)
			passed++
		} else {
			fmt.Printf("  FAIL: %s - check failed, got %T\n", tt.name, adapted)
			failed++
		}
	}

	// Test StreamErrorMsg separately (it's created internally, not from SDK message)
	fmt.Println()
	fmt.Println("Testing StreamErrorMsg (internal type):")
	errMsg := adapter.StreamErrorMsg{Err: fmt.Errorf("test error")}
	if errMsg.Error() == "test error" {
		fmt.Printf("  PASS: StreamErrorMsg\n")
		printAdaptedMessage(errMsg)
		passed++
	} else {
		fmt.Printf("  FAIL: StreamErrorMsg\n")
		failed++
	}

	// Summary
	fmt.Println()
	fmt.Println("--- Validation Summary ---")
	fmt.Printf("  Passed: %d/16\n", passed)
	fmt.Printf("  Failed: %d/16\n", failed)
	if failed == 0 {
		fmt.Println("  Status: ALL TESTS PASSED")
	} else {
		fmt.Println("  Status: SOME TESTS FAILED")
		os.Exit(1)
	}
}
