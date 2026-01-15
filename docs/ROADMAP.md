# Implementation Roadmap

## Overview

This roadmap outlines the phased implementation of `claude-agent-tui`. Each milestone builds on the previous, with testable deliverables at each stage.

---

## Milestone 1: Core Foundation

**Goal:** Establish project structure and basic streaming display.

### Deliverables

1. **Project Setup**
   - Go module with dependencies (Bubble Tea, Lip Gloss, Bubbles)
   - Basic project structure (component/, adapter/, layout/, style/)
   - CI/CD configuration (GitHub Actions)

2. **StreamText Component**
   - Append text via `DeltaMsg`
   - Optional blinking cursor
   - Basic Lip Gloss styling

3. **Stream Adapter**
   - `StreamCmd()` for channel-to-message conversion
   - `StreamDeltaMsg`, `StreamDoneMsg`, `StreamErrorMsg`
   - Handle `StreamEvent` parsing

4. **Basic Demo**
   - Minimal Bubble Tea program
   - Connect to SDK, stream response
   - Display with StreamText

### Testing
- Unit tests for StreamText component
- Integration test with mock SDK transport
- Manual test with Claude CLI

---

## Milestone 2: Message Display

**Goal:** Rich message rendering with multiple content types.

### Deliverables

1. **MessageBubble Component**
   - User/Assistant role styling
   - Content block iteration
   - Factory from `AssistantMessage`

2. **TextBlock Rendering**
   - Markdown detection
   - Basic formatting (bold, italic, code)

3. **CodeBlock Component**
   - Language detection
   - Syntax highlighting (chroma integration)
   - Line numbers option

4. **StatusBar Component**
   - Token count display
   - Cost formatting
   - Model name
   - Duration

5. **Message History**
   - Viewport for scrolling
   - Message list management

### Testing
- Visual snapshot tests
- Component unit tests
- Integration with full message flow

---

## Milestone 3: Chat Input

**Goal:** Primary user input with mode controls.

### Deliverables

1. **ChatInput Component**
   - Multi-line textarea
   - Submit on Enter/Ctrl+Enter
   - Mode toggle pills (auto-accept, plan mode)
   - Keybindings (Ctrl+A, Ctrl+P)

2. **Client Adapter**
   - `ClientAdapter` for lifecycle management
   - `ConnectCmd`, `QueryCmd`, `DisconnectCmd`
   - Permission mode switching

3. **Chat Layout (Basic)**
   - Message history + current stream
   - Input area at bottom
   - Status bar

4. **Full Chat Flow**
   - Connect -> Query -> Stream -> Display
   - Multiple turns

### Testing
- Input component unit tests
- Full conversation integration test
- Manual multi-turn testing

---

## Milestone 4: Permission Handling

**Goal:** Tool approval workflow.

### Deliverables

1. **PermissionPrompt Component**
   - Tool name and input display
   - Allow/Deny/Always options
   - Keyboard navigation
   - Danger highlighting

2. **Hook Adapter**
   - `PreToolUseHook` integration
   - Response channel pattern
   - Block/allow decision flow

3. **Input State Machine**
   - State enum (Chat, Permission, etc.)
   - Transition handling
   - Focus management

4. **Layout Integration**
   - Prompt overlay on permission request
   - Return to chat after decision

### Testing
- Permission flow integration tests
- State machine unit tests
- Manual tool approval testing

---

## Milestone 5: Question Handling

**Goal:** AskUserQuestion tool support.

### Deliverables

1. **QuestionPrompt Component**
   - Question display
   - Option list with descriptions
   - Single/multi-select modes
   - "Other" custom input
   - Keyboard navigation

2. **Tool Call Adapter**
   - `AskUserQuestion` detection
   - Input parsing
   - Response handling

3. **Layout Integration**
   - Prompt overlay on question
   - Return to chat with response

### Testing
- Question component tests
- Multi-select behavior tests
- Integration with SDK

---

## Milestone 6: Plan Mode

**Goal:** Full plan mode workflow.

### Deliverables

1. **PlanEnterConfirm Component**
   - Confirmation message
   - Confirm/Skip buttons
   - Quick keys (y/n)

2. **PlanExitPrompt Component**
   - Scrollable plan viewport
   - Permission checkboxes
   - Approve/Reject/Changes actions
   - Feedback input for changes

3. **Planning State**
   - Read-only streaming during planning
   - Visual indicator
   - State transitions

4. **Layout Integration**
   - Enter -> Planning -> Exit flow
   - Permission pre-approval

### Testing
- Plan mode flow tests
- Component unit tests
- Full workflow manual testing

---

## Milestone 7: Advanced Features

**Goal:** Polish and advanced components.

### Deliverables

1. **ThinkingBlock Component**
   - Collapse/expand toggle
   - Scrollable content
   - Signature display

2. **ToolUseCard Component**
   - Tool name header
   - Input JSON display
   - Output/error display
   - Expandable details

3. **InterruptConfirm Component**
   - Modal overlay
   - Ctrl+C handling
   - Graceful interruption

4. **Copy Support**
   - Code block copy
   - Message copy
   - Keyboard shortcuts

### Testing
- Component tests
- Interrupt flow tests
- Accessibility review

---

## Milestone 8: Layout Refinement

**Goal:** Production-ready layouts.

### Deliverables

1. **Chat Layout (Complete)**
   - All components integrated
   - Smooth transitions
   - Error handling
   - Resize handling

2. **Agent Dashboard Layout**
   - Split view (chat + tools)
   - Tool history panel
   - Session info

3. **Default Styles**
   - Light/dark aware
   - Sensible defaults
   - Documentation

4. **Examples**
   - Basic chat example
   - Custom styling example
   - Agent dashboard example

### Testing
- Layout integration tests
- Terminal size tests
- Style documentation review

---

## Milestone 9: Documentation & Release

**Goal:** Production release.

### Deliverables

1. **API Documentation**
   - GoDoc comments
   - Usage examples
   - Migration guide

2. **README**
   - Installation
   - Quick start
   - Component gallery

3. **Examples Repository**
   - Standalone examples
   - Real-world use cases

4. **Release**
   - Version tagging
   - Changelog
   - Announcement

### Testing
- Documentation review
- Example verification
- Community feedback

---

## Timeline Flexibility

This roadmap focuses on deliverables rather than dates. Each milestone should be:
- Independently testable
- Incrementally useful
- Building toward the complete vision

Milestones can be adjusted based on:
- User feedback
- SDK changes
- Priority shifts
