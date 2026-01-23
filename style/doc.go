// Package style provides centralized styling for claude-agent-tui components.
//
// This package defines a semantic color palette and component-specific styles
// that provide consistent theming across all TUI components. Components use
// these styles as defaults but allow customization via functional options.
//
// # Color Palette
//
// The package exports semantic colors that work well in terminal environments:
//
//   - Primary: Accent color for focused elements (#7C3AED, purple)
//   - Secondary: Unfocused borders and subtle accents (#6B7280, gray)
//   - Text: Primary text color (#F9FAFB, light)
//   - TextMuted: Secondary text, placeholders (#9CA3AF, gray muted)
//
// # Component Styles
//
// Each component has its own style file with sensible defaults:
//
//   - chatinput.go: InputContainer, InputContainerFocused, InputPrompt
//
// # Usage
//
// Components import centralized styles and use them as defaults:
//
//	import "github.com/severity1/claude-agent-tui/style"
//
//	func New(opts ...Option) Model {
//	    m := Model{
//	        containerStyle: style.InputContainer,
//	        containerFocusedStyle: style.InputContainerFocused,
//	    }
//	    // Apply user options which can override defaults
//	    for _, opt := range opts {
//	        opt(&m)
//	    }
//	    return m
//	}
//
// Users can override styles via functional options:
//
//	input := chatinput.New(
//	    chatinput.WithContainerStyle(myCustomStyle),
//	)
//
// # Future Work
//
// Additional component styles will be added as components are developed:
//   - keyhints.go: KeyHints component styles
//   - streamtext.go: StreamText component styles
//   - Theme extraction and switching support
package style
