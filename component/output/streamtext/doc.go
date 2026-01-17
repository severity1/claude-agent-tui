// Package streamtext provides a live streaming text component for displaying
// Claude's responses as they arrive character by character.
//
// SDK Mapping: StreamEvent (content_block_delta with text_delta type)
//
// # Features
//
// The component supports:
//   - Streaming text with typing cursor animation
//   - Markdown rendering via glamour
//   - Theming and style customization via WithStyles(), WithTheme()
//   - Multiple display variants via WithVariant(): View, Inline
//   - Animation respect for REDUCE_MOTION environment variable
//
// # Messages
//
// The component handles these tea.Msg types:
//   - DeltaMsg: Appends streaming text
//   - DoneMsg: Signals stream completion, hides cursor
//
// # Implementation
//
// The component implements the Bubble Tea Model interface (Init, Update, View).
// See docs/COMPONENTS.md for complete specification.
package streamtext
