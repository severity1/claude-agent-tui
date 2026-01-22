// Package chatinput provides a text input component for chat interfaces.
//
// The component wraps bubbles/textarea with custom keybindings:
//   - Enter: Submit text (emits SubmitMsg)
//   - Alt+Enter / Ctrl+J: Insert newline
//   - Esc: Clear input
//
// # API
//
// Constructor with functional options:
//   - New(opts ...Option): Create a new ChatInput model
//   - WithPlaceholder(text string): Set placeholder text
//   - WithWidth(w int): Set input width
//   - WithHeight(h int): Set input height
//
// Focus management:
//   - Focus(): Set focus and return cursor blink command
//   - Blur(): Remove focus
//   - Focused(): Check if focused
//
// Value management:
//   - Value(): Get current text
//   - SetValue(s string): Set current text
//
// # Messages
//
// The component emits SubmitMsg when the user presses Enter with non-empty input.
// The SubmitMsg contains the submitted text and clears the input automatically.
//
// # Implementation
//
// The component implements the Bubble Tea Model interface (Init, Update, View).
// Key messages are intercepted and transformed based on the keybindings above.
// Other messages are passed through to the underlying textarea.
package chatinput
