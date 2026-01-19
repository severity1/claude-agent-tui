// Package adapter provides integration between claude-agent-sdk-go and
// Bubble Tea components.
//
// The adapter layer bridges the SDK's channel-based streaming and callback
// architecture to Bubble Tea's message-passing model.
//
// # Files
//
// The package consists of three main files:
//
//   - stream.go: StreamCmd converts SDK ReceiveMessages channel to tea.Msg
//   - control.go: ToolControlAdapter handles canUseTool callbacks
//   - client.go: ClientAdapter manages SDK client lifecycle
//
// # Stream Messages
//
// StreamCmd reads from the SDK message channel and emits:
//
//   - StreamDeltaMsg: Incremental text content
//   - StreamBlockStartMsg: New content block (text, thinking, tool_use)
//   - StreamBlockStopMsg: Content block completion
//   - ThinkingDeltaMsg: Incremental thinking content
//   - ToolUseDeltaMsg: Incremental tool input JSON
//   - AssistantMsg: Complete assistant message
//   - ResultMsg: Final result with usage statistics
//   - UserMsg: User message (typically tool results)
//   - SystemInitMsg: Session initialization data
//   - SystemHookResponseMsg: Hook execution results
//   - MessageStartMsg: Initial message metadata with model info
//   - MessageDeltaMsg: Final usage stats and stop reason
//   - ControlRequestMsg: Control protocol request (informational)
//   - ControlResponseMsg: Control protocol response (informational)
//   - StreamDoneMsg: Stream completion signal
//   - StreamErrorMsg: Stream error
//   - UnknownMessageMsg: Unrecognized SDK message type (for debugging)
//
// # Control Messages (Planned - control.go)
//
// ToolControlAdapter will intercept canUseTool callbacks and emit:
//
//   - ShowPermissionPromptMsg: Tool permission request
//   - ShowQuestionPromptMsg: AskUserQuestion tool
//   - ShowPlanEnterPromptMsg: EnterPlanMode tool
//   - ShowPlanExitPromptMsg: ExitPlanMode tool
//   - ShowTodoUpdateMsg: TodoWrite tool updates
//
// # Thread Safety
//
// The SDK's canUseTool callback runs in a separate goroutine. The adapter
// uses channels to synchronize responses between the SDK goroutine and the
// Bubble Tea goroutine. Always call Shutdown() before tea.Program.Quit().
//
// See docs/ADAPTERS.md for complete specification.
package adapter
