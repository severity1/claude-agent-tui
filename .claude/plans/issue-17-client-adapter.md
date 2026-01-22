# TDD Implementation Plan: Issue #17 - ClientAdapter for SDK Lifecycle Management

## Issue Summary

**Title:** [MVP-2.5] Implement client adapter for SDK lifecycle management
**Labels:** mvp, phase-1, adapter

Create `adapter/client.go` to manage SDK client lifecycle and provide Bubble Tea commands for connection, queries, and disconnection.

## Requirements from Issue

### Message Types
- `ClientStateMsg` - signals client state changes
- `ClientState` enum - Disconnected, Connecting, Connected, Streaming, Error

### ClientAdapter Struct
- Wraps SDK client for Bubble Tea integration
- Thread-safe state management with sync.Mutex
- Stores tea.Program reference for message sending

### Required Methods
- `NewClientAdapter()` - constructor
- `SetProgram(p *tea.Program)` - set Bubble Tea program
- `ConnectCmd() tea.Cmd` - establishes SDK connection
- `QueryCmd(text string) tea.Cmd` - sends query and returns streaming channel
- `InterruptCmd() tea.Cmd` - interrupts current stream
- `DisconnectCmd() tea.Cmd` - cleans up resources
- `MessageChannel() <-chan claudecode.Message` - returns streaming channel

## Architecture Decisions

### State Machine
```
Disconnected -> Connecting -> Connected -> Streaming -> Connected
                    |             |            |            |
                    v             v            v            v
                 Error         Error        Error        Error
```

### Thread Safety
- Use `sync.RWMutex` to protect shared state
- Commands execute outside lock (only state reads/writes locked)
- Context cancellation for graceful shutdown

### Integration with Existing Adapter
- ClientAdapter orchestrates client lifecycle
- Returns message channel for use with existing `StreamCmd()`
- Does NOT duplicate stream message conversion (uses existing `AdaptMessage`)

---

## TDD Plan

### RED Phase - Write Failing Tests First

#### Test File: `adapter/client_test.go`

```go
package adapter_test
```

#### 1. Message Type Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientStateMsg_Fields` | Verify ClientStateMsg struct fields accessible |
| `TestClientState_String` | Verify String() method for each state |
| `TestClientState_Constants` | Verify all constants defined correctly |

#### 2. Constructor Tests

| Test Function | Description |
|---------------|-------------|
| `TestNewClientAdapter_ReturnsNonNil` | Constructor returns valid adapter |
| `TestNewClientAdapter_InitialState` | Initial state is Disconnected |
| `TestNewClientAdapter_WithOptions` | SDK options are stored |

#### 3. SetProgram Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_SetProgram_StoresProgram` | Program reference stored |
| `TestClientAdapter_SetProgram_NilSafe` | Handles nil program gracefully |

#### 4. ConnectCmd Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_ConnectCmd_ReturnsTeaCmd` | Returns tea.Cmd function |
| `TestClientAdapter_ConnectCmd_Success` | Emits ClientStateMsg{Connected} |
| `TestClientAdapter_ConnectCmd_Failure` | Emits ClientStateMsg{Error, Err} |
| `TestClientAdapter_ConnectCmd_AlreadyConnected` | Handles double connect |
| `TestClientAdapter_ConnectCmd_StateTransition` | Disconnected -> Connecting -> Connected |

#### 5. QueryCmd Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_QueryCmd_ReturnsTeaCmd` | Returns tea.Cmd function |
| `TestClientAdapter_QueryCmd_Success` | Emits ClientStateMsg{Streaming} |
| `TestClientAdapter_QueryCmd_NotConnected` | Emits ClientStateMsg{Error} |
| `TestClientAdapter_QueryCmd_EmptyPrompt` | Handles empty prompt |
| `TestClientAdapter_QueryCmd_StoresMessageChannel` | Channel accessible via MessageChannel() |

#### 6. InterruptCmd Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_InterruptCmd_ReturnsTeaCmd` | Returns tea.Cmd function |
| `TestClientAdapter_InterruptCmd_Success` | Returns to Connected state |
| `TestClientAdapter_InterruptCmd_NotStreaming` | No-op when not streaming |

#### 7. DisconnectCmd Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_DisconnectCmd_ReturnsTeaCmd` | Returns tea.Cmd function |
| `TestClientAdapter_DisconnectCmd_Success` | Emits ClientStateMsg{Disconnected} |
| `TestClientAdapter_DisconnectCmd_AlreadyDisconnected` | Handles double disconnect |
| `TestClientAdapter_DisconnectCmd_CleansUpResources` | Client and channels nil |

#### 8. MessageChannel Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_MessageChannel_ReturnsChannel` | Returns stored channel |
| `TestClientAdapter_MessageChannel_NilWhenNotStreaming` | Nil before query |

#### 9. State Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_State_ReturnsCurrentState` | State() returns current state |
| `TestClientAdapter_State_ThreadSafe` | Concurrent State() calls safe |

#### 10. Concurrency Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_ConcurrentConnect` | Multiple connect calls safe |
| `TestClientAdapter_ConcurrentStateReads` | Multiple state reads safe |
| `TestClientAdapter_RapidConnectDisconnect` | Rapid lifecycle changes safe |

#### 11. Edge Case Tests

| Test Function | Description |
|---------------|-------------|
| `TestClientAdapter_QueryAfterDisconnect` | Error when querying after disconnect |
| `TestClientAdapter_InterruptWithoutQuery` | No-op when not streaming |
| `TestClientAdapter_ContextCancellation` | Context cancel handled gracefully |

---

### GREEN Phase - Implementation

#### File: `adapter/client.go`

```go
package adapter

// 1. Constants
const (
    ClientStateDisconnected ClientState = iota
    ClientStateConnecting
    ClientStateConnected
    ClientStateStreaming
    ClientStateError
)

// 2. Types
type ClientState int
type ClientStateMsg struct {
    State ClientState
    Error error
}

// 3. ClientAdapter struct
type ClientAdapter struct {
    client  claudecode.Client
    program *tea.Program
    msgChan <-chan claudecode.Message
    state   ClientState
    ctx     context.Context
    cancel  context.CancelFunc
    mu      sync.RWMutex
    options []claudecode.Option
}

// 4. Constructor
func NewClientAdapter(opts ...claudecode.Option) *ClientAdapter

// 5. Methods
func (c *ClientAdapter) SetProgram(p *tea.Program)
func (c *ClientAdapter) State() ClientState
func (c *ClientAdapter) ConnectCmd() tea.Cmd
func (c *ClientAdapter) QueryCmd(prompt string) tea.Cmd
func (c *ClientAdapter) InterruptCmd() tea.Cmd
func (c *ClientAdapter) DisconnectCmd() tea.Cmd
func (c *ClientAdapter) MessageChannel() <-chan claudecode.Message
func (c ClientState) String() string
```

#### Implementation Order
1. Define constants and types
2. Implement `ClientState.String()`
3. Implement `NewClientAdapter()`
4. Implement `SetProgram()`
5. Implement `State()` (thread-safe read)
6. Implement `ConnectCmd()` with state transitions
7. Implement `QueryCmd()` with channel storage
8. Implement `InterruptCmd()`
9. Implement `DisconnectCmd()` with cleanup
10. Implement `MessageChannel()`

---

### BLUE Phase - Refactoring

After tests pass:
1. Extract common error handling patterns
2. Ensure consistent error messages match stream.go style
3. Add doc comments matching existing adapter pattern
4. Verify thread safety with race detector
5. Update doc.go with ClientAdapter documentation

---

## Acceptance Criteria Mapping

| Requirement | Test(s) |
|-------------|---------|
| ClientStateMsg and ClientState types defined | TestClientStateMsg_Fields, TestClientState_Constants |
| ClientAdapter struct with SDK client wrapper | TestNewClientAdapter_ReturnsNonNil |
| NewClientAdapter() constructor | TestNewClientAdapter_* tests |
| SetProgram() for Bubble Tea integration | TestClientAdapter_SetProgram_* |
| ConnectCmd() establishes SDK connection | TestClientAdapter_ConnectCmd_* |
| QueryCmd() sends queries and returns streaming channel | TestClientAdapter_QueryCmd_* |
| InterruptCmd() cancels ongoing stream | TestClientAdapter_InterruptCmd_* |
| DisconnectCmd() cleans up resources | TestClientAdapter_DisconnectCmd_* |
| Thread-safe access to client state | TestClientAdapter_State_ThreadSafe, Concurrency tests |

---

## Files to Create/Modify

### Create
- `adapter/client.go` - Main implementation
- `adapter/client_test.go` - Comprehensive tests

### Modify
- `adapter/doc.go` - Add ClientAdapter to package documentation

---

## Testing Strategy

### Mock Approach
Since the SDK's `claudecode.Client` is an interface, create a mock for testing:

```go
type mockClient struct {
    connectErr    error
    queryErr      error
    disconnectErr error
    interruptErr  error
    msgChan       chan claudecode.Message
    connected     bool
}
```

This allows testing:
- Success paths (mock returns no error)
- Failure paths (mock returns specific errors)
- Channel behavior (mock provides test channel)

### Race Detection
Run all tests with `-race` flag to catch concurrency issues.

### Coverage Target
Aim for 90%+ coverage on adapter package (matching existing standard).

---

## Dependencies

- `github.com/charmbracelet/bubbletea` (tea.Cmd, tea.Msg, tea.Program)
- `github.com/severity1/claude-agent-sdk-go` (claudecode.Client, claudecode.Message, options)
- `context` (cancellation)
- `sync` (RWMutex)

---

## Risk Assessment

### Low Risk
- Message type definitions (straightforward structs)
- State enum (simple iota constants)
- Constructor (standard Go pattern)

### Medium Risk
- State transitions (need careful synchronization)
- Command pattern integration (must match tea.Cmd semantics)

### Mitigation
- Extensive table-driven tests for state transitions
- Concurrent tests to verify thread safety
- Mock client for isolated unit testing
