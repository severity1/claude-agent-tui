# Implementation Roadmap

## Overview

This roadmap outlines the phased implementation of `claude-agent-tui`. Each milestone builds on the previous, with testable deliverables at each stage.

---

## Milestone 1: Core Foundation

**Goal:** Establish project structure, core systems, and basic streaming display.

### Deliverables

1. **Project Setup**
   - Go module with dependencies (Bubble Tea, Lip Gloss, Bubbles, Harmonica)
   - Directory structure (component/, system/, adapter/, layout/, style/)
   - CI/CD configuration (GitHub Actions)
   - Linting and formatting setup

2. **Core Systems Setup**
   - `system/theme/` - Theme interface and System theme
   - `system/keymap/` - Keymap struct with defaults
   - `system/focus/` - Focus manager basics
   - Variant type definitions

3. **StreamText Component**
   - Append text via `DeltaMsg`
   - Optional blinking cursor
   - View variant support
   - Theme-aware styling

4. **Stream Adapter**
   - `StreamCmd()` for channel-to-message conversion
   - `StreamDeltaMsg`, `StreamDoneMsg`, `StreamErrorMsg`
   - Handle `StreamEvent` parsing (text, thinking, tool_use blocks)

5. **Basic Demo**
   - Minimal Bubble Tea program
   - Connect to SDK, stream response
   - Display with StreamText

### Testing
- Unit tests for StreamText component
- Integration test with mock SDK transport
- Manual test with Claude CLI

---

## Milestone 2: Theme System + Accessibility Foundation

**Goal:** Complete theming engine with built-in themes AND accessibility from day one.

> **Note:** Accessibility moved here from M11.5. It must be foundational, not an afterthought.

### Deliverables

1. **Theme Interface**
   - `Theme` interface definition (simplified: Name, IsDark, Palette)
   - `Palette` struct with all color categories
   - Theme registry (Get, Register, List)
   - `StylesFromPalette()` factory pattern per component

2. **Built-in Themes**
   - Catppuccin (Mocha, Latte, Frappe, Macchiato)
   - Dracula
   - Nord
   - Gruvbox (dark, light)
   - Tokyo Night
   - One Dark
   - System (adaptive)
   - **High Contrast** (accessibility)

3. **Accessibility Foundation**
   - `REDUCE_MOTION` environment variable support
   - `HIGH_CONTRAST` environment variable
   - Visible focus indicators (always enabled by default)
   - Color-blindness safe design (symbols + colors, never color alone)
   - Accessibility config loading from environment

4. **Component Style Integration**
   - `Styles` struct per component using Palette
   - Syntax highlighting theme adapter (chroma)
   - Focus ring styling in all interactive components

5. **Theme Switching**
   - Runtime theme change support
   - Persistence options

### Testing
- Visual tests for each theme
- Color contrast validation (WCAG AAA for high contrast)
- Light/dark mode testing
- Reduced motion behavior tests
- Focus indicator visibility tests

---

## Milestone 3: Message Display

**Goal:** Rich message rendering with multiple content types.

### Deliverables

1. **Message Component**
   - User/Assistant role styling
   - Content block iteration
   - Factory from `AssistantMessage`
   - View, Inline, Card variants

2. **Thinking Component**
   - Collapsible/expandable toggle
   - Scrollable viewport
   - Streaming support
   - View, Collapsible, Modal variants

3. **CodeBlock Component**
   - Language detection
   - Syntax highlighting (chroma)
   - Line numbers option
   - Copy button zone
   - View, Inline, Modal, Window variants

4. **ToolUse Component**
   - Tool name header with icon
   - Input JSON display (collapsible)
   - Output/error display
   - Loading spinner during execution
   - View, Card, Inline, Modal variants

5. **Status Component**
   - Token count display
   - Cost formatting
   - Model name
   - Duration
   - Bar, Inline, Compact variants

6. **Message History**
   - Viewport for scrolling
   - Message list management
   - Auto-scroll with user override

### Testing
- Visual snapshot tests
- Component unit tests
- Integration with full message flow

---

## Milestone 4: Chat Input + Client Lifecycle

**Goal:** Primary user input with mode controls, client lifecycle, and error recovery.

> **Note:** Error recovery and background tasks moved here from M11.5. Client lifecycle needs these.

### Deliverables

1. **ChatInput Component**
   - Multi-line textarea
   - Submit on Enter (configurable)
   - Mode toggle pills (auto-accept, plan mode)
   - Keybindings (Ctrl+A, Ctrl+P)
   - Placeholder and hint text
   - Character count (optional)

2. **Shared Components**
   - Button component (with variants)
   - Chip/Pill component (toggleable)
   - Card container

3. **Client Adapter**
   - `ClientAdapter` with `canUseTool` callback
   - `ConnectCmd`, `QueryCmd`, `DisconnectCmd`
   - State management (Disconnected, Connected, Streaming, etc.)

4. **Error Recovery System**
   - RecoveryManager with exponential backoff
   - Auto-reconnect logic with jitter
   - Partial message recovery (save on error, restore on reconnect)
   - Connection state indicators (healthy, unstable, lost, reconnecting, failed)
   - User notification via toasts

5. **Background Task Foundation**
   - TaskStatus indicator (shows running task count)
   - Support for `run_in_background` parameter
   - TaskStartedMsg, TaskProgressMsg, TaskCompletedMsg handling
   - Basic task tracking (full TaskManager UI in M11)

6. **Chat Layout (Basic)**
   - Message history + current stream
   - Input area at bottom
   - Status bar with connection state
   - Fullscreen, Window, View variants

7. **Full Chat Flow**
   - Connect -> Query -> Stream -> Display
   - Multiple turns
   - Error handling with recovery
   - Graceful degradation on connection loss

### Testing
- Input component unit tests
- Full conversation integration test
- Manual multi-turn testing
- Error recovery flow tests
- Connection state transition tests

---

## Milestone 5: Keybinding System

**Goal:** Customizable keyboard shortcuts with help overlay.

### Deliverables

1. **Keymap System**
   - Full `Keymap` struct with all bindings
   - Default keybindings (vim + arrows)
   - Component-specific keymaps (ChatInput, Permission, Question)

2. **Customization API**
   - `Set()`, `Get()`, `Merge()` methods
   - Conflict detection
   - JSON loading support

3. **Help Overlay**
   - Auto-generated from keybindings
   - Short and full help views
   - Context-aware (show relevant bindings)

4. **Component Integration**
   - All components use keymap
   - Consistent key handling patterns

### Testing
- Keybinding unit tests
- Help overlay rendering tests
- Custom keymap loading tests

---

## Milestone 6: Permission Handling + Diff Viewer

**Goal:** Tool approval workflow with rich display, including proper diff viewing for Edit tool.

> **Note:** Diff Viewer moved here from M11.5. Edit tool permissions need it.

### Deliverables

1. **Permission Component**
   - Tool name and description
   - Tool input display (formatted per tool type)
   - Allow/Deny/Always Allow options
   - Keyboard navigation
   - Danger highlighting for risky operations
   - Modal, Card variants

2. **Tool Control Adapter**
   - `canUseTool` callback integration
   - Auto-approval management
   - Response channel pattern

3. **DiffView Component**
   - Unified diff mode (default)
   - Split diff mode (side-by-side)
   - Inline diff mode (within text)
   - Syntax highlighting for both old and new
   - Line number display
   - View, Modal variants

4. **Tool-Specific Displays**
   - Bash: command preview with description
   - Write: file path + content preview (collapsible)
   - Edit: **DiffView** showing old vs new
   - Read: file path with line range
   - Glob/Grep: pattern + path

5. **Input State Machine**
   - State enum (Chat, Permission, Question, etc.)
   - Transition handling
   - Focus management

6. **Layout Integration**
   - Prompt overlay on permission request
   - Return to chat after decision

### Testing
- Permission flow integration tests
- State machine unit tests
- Diff rendering tests (unified, split, inline)
- Manual tool approval testing

---

## Milestone 7: Question & Plan Handling

**Goal:** AskUserQuestion and plan mode support.

### Deliverables

1. **Question Component**
   - Question display with header
   - Option list with descriptions
   - Single/multi-select modes
   - "Other" custom input option
   - Number key shortcuts (1-4)
   - Modal, Overlay, View variants

2. **PlanEnter Component**
   - Confirmation message
   - Confirm/Skip buttons
   - Quick keys (y/n)
   - Modal, Overlay variants

3. **PlanExit Component**
   - Scrollable plan viewport (markdown rendered)
   - Permission checkboxes
   - Approve/Reject/Request Changes actions
   - Feedback input for changes
   - Modal, Fullscreen, Window variants

4. **Planning State**
   - Read-only streaming during planning
   - Visual plan mode indicator
   - State transitions

5. **Todo Component**
   - Task list display
   - Status indicators (pending, in_progress, completed)
   - Active task highlighting
   - View, Sidebar, Overlay, Collapsible variants

### Testing
- Question component tests
- Plan mode flow tests
- Full workflow manual testing

---

## Milestone 8: Mouse Support

**Goal:** Full mouse interaction with click, hover, and scroll.

### Deliverables

1. **Mouse System**
   - Zone registration and management
   - Click handler routing
   - Hover state tracking
   - Scroll event handling

2. **Component Mouse Integration**
   - Button hover/click states
   - List item selection
   - Scrollable areas (wheel support)
   - Collapsible toggle clicks

3. **Visual Feedback**
   - Hover styles per component
   - Click animation triggers
   - Cursor style hints

4. **Double-Click Detection**
   - Configurable threshold
   - Per-zone tracking

### Testing
- Mouse event routing tests
- Hover state tests
- Click zone accuracy tests

---

## Milestone 9: Animation System

**Goal:** Smooth animations for state changes.

### Deliverables

1. **Animation Engine**
   - `Animation` interface
   - Spring-based motion (harmonica)
   - Animation configs and presets

2. **Built-in Animations**
   - FadeIn/FadeOut
   - SlideIn/SlideOut (4 directions)
   - Pulse, Shake, Bounce
   - Breathe (idle)
   - Blink (cursor)
   - Spin (loading)
   - TypeWriter

3. **Animation Triggers**
   - User input (key, mouse)
   - SDK events (stream, tool)
   - State changes (idle, error)

4. **Component Animation Config**
   - Per-trigger animations
   - Enter/exit transitions
   - Idle delay configuration

5. **Toast Component**
   - Toast manager with queue
   - Slide-in animations
   - Auto-dismiss with duration
   - Multiple positions
   - Toast levels (info, success, warning, error)

6. **Reduced Motion Support**
   - Environment variable check
   - Disable animations option

### Testing
- Animation timing tests
- Spring physics validation
- Reduced motion behavior tests

---

## Milestone 10: Responsive System

**Goal:** Adaptive layouts for different terminal sizes.

### Deliverables

1. **Breakpoint System**
   - Width breakpoints (XS, SM, MD, LG, XL)
   - Height breakpoints
   - Current breakpoint detection

2. **Responsive Config**
   - Per-component responsive settings
   - Variant per breakpoint
   - Padding per breakpoint
   - Visibility rules

3. **Layout Constraints**
   - Min/max width/height
   - Clamp functions
   - Available space calculations

4. **Resize Handling**
   - WindowSizeMsg processing
   - Component size updates
   - Layout recalculation

5. **Graceful Degradation**
   - "Too small" fallback
   - Minimum size requirements
   - Hidden elements at small sizes

### Testing
- Breakpoint detection tests
- Resize behavior tests
- Various terminal size testing

---

## Milestone 11: Advanced Features

**Goal:** Additional components and polish.

### Deliverables

1. **FilePicker Component**
   - Directory navigation
   - File filtering
   - Selection (single/multi)
   - Modal, Overlay, Window variants

2. **Confirm Component**
   - Generic yes/no prompt
   - Custom button labels
   - Modal, Inline variants

3. **Interrupt Component**
   - Ctrl+C handling
   - Confirmation modal
   - Graceful interruption flow

4. **Overlay Component**
   - Backdrop rendering
   - Focus trapping
   - Click-outside dismiss

5. **Copy Support**
   - Code block copy
   - Message copy
   - Keyboard shortcuts
   - Clipboard integration

### Testing
- Component tests
- Interrupt flow tests
- Accessibility review

---

## Milestone 11.5: Remaining Gap Coverage

**Goal:** Address remaining coverage gaps not covered by earlier milestones.

> **Note:** Several items were moved to earlier milestones:
> - Accessibility Foundation -> M2
> - Error Recovery System -> M4
> - Background Task Foundation -> M4
> - Diff Viewer -> M6

### Deliverables

1. **Full TaskManager UI**
   - TaskManager component (Modal, Collapsible variants)
   - Background task output viewer (scrollable, virtualized)
   - KillShell confirmation dialog
   - Task history panel

2. **Session Management**
   - SessionPicker component (Modal variants)
   - Session persistence and restore
   - Resume from previous conversations
   - Export/import sessions
   - Context window usage tracking

3. **File Result Viewers**
   - GlobResults component (file tree/list display)
   - GrepResults component (match highlighting)
   - Navigate to file:line on selection

4. **Enhanced Status**
   - Context window indicator (progress bar)
   - Session info display

5. **Screen Reader Announcements**
   - Announcer component for screen reader events
   - OSC 9 notification support (with limitations noted)
   - Event-to-announcement mapping

6. **Search in History**
   - Search component (Modal variants)
   - Live search as you type
   - Result highlighting
   - Navigate to message on selection

7. **New Tool Support**
   - TaskOutput tool integration
   - KillShell tool integration
   - Skill tool integration

### Testing
- Full TaskManager flow tests
- Session save/load tests
- Search functionality tests

### Dependencies
- Milestone 11 (Advanced Features)

---

## Milestone 12: Layout Refinement

**Goal:** Production-ready composite layouts.

### Deliverables

1. **Chat Layout (Complete)**
   - All components integrated
   - All variants (Fullscreen, Window, View, Overlay)
   - Smooth transitions
   - Error handling with toasts
   - Responsive behavior

2. **Agent Dashboard Layout**
   - Split view (chat + tools)
   - Tool history panel
   - Session info
   - Todo sidebar

3. **Default Styles**
   - Theme defaults for all components
   - Documentation

4. **Examples**
   - Basic chat example
   - Custom theme example
   - Custom keybindings example
   - Agent dashboard example
   - Mouse-enabled example

### Testing
- Layout integration tests
- Terminal size tests
- Example verification

---

## Milestone 13: Documentation & Release

**Goal:** Production release.

### Deliverables

1. **API Documentation**
   - GoDoc comments (all public APIs)
   - Usage examples per component
   - Option pattern documentation

2. **Guides**
   - Getting Started guide
   - Theming guide
   - Custom components guide
   - Migration guide

3. **README**
   - Installation
   - Quick start
   - Component gallery (screenshots)
   - Feature highlights

4. **Examples Repository**
   - Standalone examples
   - Real-world use cases
   - Integration patterns

5. **Release**
   - Version tagging (semver)
   - Changelog
   - Release notes
   - Community announcement

### Testing
- Documentation review
- Example verification
- Community feedback

---

## Dependency Summary

| Milestone | Dependencies | Rationale |
|-----------|-------------|-----------|
| 1 | None | Foundation - must be first |
| 2 | 1 | Themes + accessibility need base styles |
| 3 | 1, 2 | Messages need themes |
| 4 | 1, 2, 3 | Input + client lifecycle needs message display |
| 5 | 1, 4 | Keybindings need input to test |
| 6 | 1, 2, 4, 5 | Permissions + diff viewer need input + keys |
| 7 | 1, 2, 4, 5, 6 | Questions build on permission pattern |
| 8 | 1, 2, 4, 6 | Mouse needs clickable components to test |
| 9 | 1, 2, 3, 4 | Animations need components to animate |
| 10 | 1, 2, 3, 4 | Responsive needs full layout |
| 11 | 1-10 | Advanced features need core complete |
| 11.5 | 11 | Remaining gap coverage |
| 12 | 1-11.5 | Final integration |
| 13 | 1-12 | Documentation |

---

## Applied Reordering

The following features were moved from M11.5 to earlier milestones:

| Feature | Original | Moved To | Reason |
|---------|----------|----------|--------|
| Accessibility Foundation | M11.5 | **M2** | Must be foundational, not afterthought |
| Error Recovery | M11.5 | **M4** | Client lifecycle needs recovery |
| Background Task Foundation | M11.5 | **M4** | Part of client lifecycle |
| Diff Viewer | M11.5 | **M6** | Needed for Edit tool permissions |

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
- Community contributions

---

## MVP vs Full Roadmap

**For MVP (ship something first):** Focus on M1 + M4 + M6 only.
See [MVP.md](MVP.md) for the minimal implementation plan.

**Full roadmap:** Follow milestones in order after MVP works.

---

## Feature Matrix

| Feature | Milestone | Priority | Notes |
|---------|-----------|----------|-------|
| Streaming text | 1 | Critical | MVP |
| Theme system | 2 | High | Basic theme in M1, full in M2 |
| **Accessibility foundation** | **2** | **High** | **Moved from M11.5** |
| Message display | 3 | Critical | MVP (simplified) |
| Chat input | 4 | Critical | MVP |
| **Error recovery** | **4** | **High** | **Moved from M11.5** |
| **Background task foundation** | **4** | **High** | **Moved from M11.5** |
| Keybindings | 5 | High | Basic in M1, full in M5 |
| Permission prompts | 6 | Critical | MVP |
| **Diff viewer** | **6** | **High** | **Moved from M11.5** |
| Question prompts | 7 | Critical | |
| Plan mode | 7 | High | |
| Todo display | 7 | Medium | |
| Mouse support | 8 | Medium | |
| Animations | 9 | Low | Defer until requested |
| Toast notifications | 9 | High | Error feedback |
| Responsive layout | 10 | High | |
| FilePicker | 11 | Low | |
| Copy support | 11 | Medium | |
| Full TaskManager UI | 11.5 | Medium | Builds on M4 foundation |
| Session management | 11.5 | Medium | Complex, defer |
| File result viewers | 11.5 | Medium | |
| Screen reader announcements | 11.5 | Low | Limited terminal support |
| Search in history | 11.5 | Low | |
| TaskOutput/KillShell/Skill | 11.5 | High | |
| Dashboard layout | 12 | Medium | |
