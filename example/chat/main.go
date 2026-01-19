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
		fmt.Println("Example: go run main.go \"say hello\"")
		os.Exit(1)
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
