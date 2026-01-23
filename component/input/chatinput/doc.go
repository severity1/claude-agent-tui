// Package chatinput provides a styled text input component for chat interfaces.
//
// The component wraps bubbles/textarea with custom keybindings and built-in styling:
//   - Background applied directly to textarea internal styles (seamless appearance)
//   - Prompt rendered separately with bold primary color (no background)
//   - Focus indicated by cursor visibility, not background color change
//   - Enter: Submit text (emits SubmitMsg)
//   - Alt+Enter / Ctrl+J: Insert newline
//   - Esc: Clear input
//
// # Styling Approach
//
// ChatInput uses the OpenCode-style seamless input design:
//   - No container wrapper - prompt and textarea are joined horizontally
//   - Background color applied to all textarea internal styles (Base, CursorLine, Placeholder, Text)
//   - Same background for both focused and blurred states
//   - Prompt uses bold text with primary accent color (no background)
//   - Focus state is indicated by cursor visibility, not color changes
//
// # API
//
// Constructor with functional options:
//   - New(opts ...Option): Create a new ChatInput model
//   - WithPlaceholder(text string): Set placeholder text
//   - WithWidth(w int): Set initial input width
//   - WithHeight(h int): Set initial input height
//   - WithPrompt(prompt string): Set prompt prefix (default "> ")
//   - WithPromptStyle(s lipgloss.Style): Customize prompt style
//
// Focus management:
//   - Focus(): Set focus and return cursor blink command
//   - Blur(): Remove focus
//   - Focused(): Check if focused
//
// Value management:
//   - Value(): Get current text
//   - SetValue(s string): Set current text
//   - Prompt(): Get current prompt prefix
//
// Runtime sizing (for responsive layouts):
//   - SetWidth(w int): Update width at runtime
//   - SetHeight(h int): Update height at runtime
//   - SetSize(width, height int): Update both dimensions
//   - Width(): Get current configured width
//
// # Messages
//
// The component emits SubmitMsg when the user presses Enter with non-empty input.
// The SubmitMsg contains the submitted text and clears the input automatically.
//
// # Example
//
//	// Create with defaults - no lipgloss import needed
//	input := chatinput.New(
//	    chatinput.WithPlaceholder("Type a message..."),
//	    chatinput.WithHeight(3),
//	)
//
//	// Or customize the prompt
//	input := chatinput.New(
//	    chatinput.WithPrompt("$ "),
//	)
//
// # Responsive Layouts
//
// For responsive layouts, handle tea.WindowSizeMsg in your parent component
// and call SetWidth() to adjust the input width:
//
//	case tea.WindowSizeMsg:
//	    inputWidth := msg.Width - 2  // account for padding
//	    if inputWidth < 20 {
//	        inputWidth = 20  // minimum width
//	    }
//	    m.input.SetWidth(inputWidth)
//	    return m, nil
//
// This follows the ResponsiveComponent pattern where the parent controls
// child dimensions rather than having children track terminal size internally.
//
// # Implementation
//
// The component implements the Bubble Tea Model interface (Init, Update, View).
// Key messages are intercepted and transformed based on the keybindings above.
// Other messages are passed through to the underlying textarea. The View() method
// uses lipgloss.JoinHorizontal for proper multi-line prompt alignment.
package chatinput
