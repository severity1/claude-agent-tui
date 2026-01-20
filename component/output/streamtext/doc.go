// Package streamtext provides a streaming text display component with cursor.
//
// SDK Mapping: StreamDeltaMsg, AssistantMsg content rendering
//
// # Features
//
// The component supports:
//   - Streaming text accumulation via Append()
//   - Optional block cursor during active streaming
//   - Content clearing via Clear()
//   - Width configuration for text wrapping via WithWidth()
//
// # API
//
// The component is controlled via method calls rather than messages:
//   - Append(text): Add text to the content buffer
//   - Clear(): Reset the content buffer
//   - SetStreaming(bool): Toggle cursor visibility
//   - Content(): Return accumulated text
//
// # Implementation
//
// The component implements the Bubble Tea Model interface (Init, Update, View).
// The Update method passes through all messages unchanged. Parent components
// control the text content via the Append, Clear, and SetStreaming methods.
//
// See docs/COMPONENTS.md for complete specification.
package streamtext
