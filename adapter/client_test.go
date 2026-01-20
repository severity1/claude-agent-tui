package adapter_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	claudecode "github.com/severity1/claude-agent-sdk-go"
	"github.com/severity1/claude-agent-tui/adapter"
)

// ============================================================================
// Mock Client for Testing
// ============================================================================

// mockClient implements claudecode.Client for testing.
type mockClient struct {
	connectErr    error
	queryErr      error
	disconnectErr error
	interruptErr  error
	msgChan       chan claudecode.Message
	connected     bool
	mu            sync.Mutex
}

func newMockClient() *mockClient {
	return &mockClient{
		msgChan: make(chan claudecode.Message, 10),
	}
}

func (m *mockClient) Connect(_ context.Context, _ ...claudecode.StreamMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.connectErr != nil {
		return m.connectErr
	}
	m.connected = true
	return nil
}

func (m *mockClient) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.disconnectErr != nil {
		return m.disconnectErr
	}
	m.connected = false
	return nil
}

func (m *mockClient) Query(_ context.Context, _ string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.queryErr
}

func (m *mockClient) QueryWithSession(_ context.Context, _ string, _ string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.queryErr
}

func (m *mockClient) QueryStream(_ context.Context, _ <-chan claudecode.StreamMessage) error {
	return nil
}

func (m *mockClient) ReceiveMessages(_ context.Context) <-chan claudecode.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.msgChan
}

func (m *mockClient) ReceiveResponse(_ context.Context) claudecode.MessageIterator {
	return nil
}

func (m *mockClient) Interrupt(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.interruptErr
}

func (m *mockClient) SetModel(_ context.Context, _ *string) error {
	return nil
}

func (m *mockClient) SetPermissionMode(_ context.Context, _ claudecode.PermissionMode) error {
	return nil
}

func (m *mockClient) RewindFiles(_ context.Context, _ string) error {
	return nil
}

func (m *mockClient) GetStreamIssues() []claudecode.StreamIssue {
	return nil
}

func (m *mockClient) GetStreamStats() claudecode.StreamStats {
	return claudecode.StreamStats{}
}

func (m *mockClient) GetServerInfo(_ context.Context) (map[string]interface{}, error) {
	return nil, nil
}

// ============================================================================
// ClientState Tests
// ============================================================================

func TestClientState_String(t *testing.T) {
	tests := []struct {
		state adapter.ClientState
		want  string
	}{
		{adapter.ClientStateDisconnected, "disconnected"},
		{adapter.ClientStateConnecting, "connecting"},
		{adapter.ClientStateConnected, "connected"},
		{adapter.ClientStateStreaming, "streaming"},
		{adapter.ClientStateError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("ClientState.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClientState_String_Unknown(t *testing.T) {
	unknown := adapter.ClientState(99)
	got := unknown.String()
	if got == "" {
		t.Error("ClientState.String() for unknown state returned empty string")
	}
}

func TestClientState_Constants(t *testing.T) {
	// Verify constants are defined and have expected ordering
	if adapter.ClientStateDisconnected != 0 {
		t.Errorf("ClientStateDisconnected = %d, want 0", adapter.ClientStateDisconnected)
	}
	if adapter.ClientStateConnecting != 1 {
		t.Errorf("ClientStateConnecting = %d, want 1", adapter.ClientStateConnecting)
	}
	if adapter.ClientStateConnected != 2 {
		t.Errorf("ClientStateConnected = %d, want 2", adapter.ClientStateConnected)
	}
	if adapter.ClientStateStreaming != 3 {
		t.Errorf("ClientStateStreaming = %d, want 3", adapter.ClientStateStreaming)
	}
	if adapter.ClientStateError != 4 {
		t.Errorf("ClientStateError = %d, want 4", adapter.ClientStateError)
	}
}

// ============================================================================
// ClientStateMsg Tests
// ============================================================================

func TestClientStateMsg_Fields(t *testing.T) {
	testErr := errors.New("test error")
	msg := adapter.ClientStateMsg{
		State: adapter.ClientStateConnected,
		Error: testErr,
	}

	if msg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", msg.State, adapter.ClientStateConnected)
	}
	if msg.Error != testErr {
		t.Errorf("Error = %v, want %v", msg.Error, testErr)
	}
}

func TestClientStateMsg_NilError(t *testing.T) {
	msg := adapter.ClientStateMsg{
		State: adapter.ClientStateConnected,
		Error: nil,
	}

	if msg.Error != nil {
		t.Errorf("Error = %v, want nil", msg.Error)
	}
}

func TestClientStateMsg_IsTeaMsg(t *testing.T) {
	// Verify ClientStateMsg implements tea.Msg interface
	var msg tea.Msg = adapter.ClientStateMsg{State: adapter.ClientStateConnected}
	if msg == nil {
		t.Error("ClientStateMsg should implement tea.Msg")
	}
}

// ============================================================================
// NewClientAdapter Tests
// ============================================================================

func TestNewClientAdapter_ReturnsNonNil(t *testing.T) {
	ca := adapter.NewClientAdapter()
	if ca == nil {
		t.Error("NewClientAdapter() returned nil")
	}
}

func TestNewClientAdapter_InitialState(t *testing.T) {
	ca := adapter.NewClientAdapter()
	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("Initial state = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}
}

func TestNewClientAdapter_WithOptions(t *testing.T) {
	// Options should be stored for later use in ConnectCmd
	model := "test-model"
	ca := adapter.NewClientAdapter(claudecode.WithModel(model))
	if ca == nil {
		t.Error("NewClientAdapter() with options returned nil")
	}
	// Options are applied internally during ConnectCmd, not directly testable here
	// The adapter should still initialize correctly
	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("Initial state = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}
}

// ============================================================================
// SetProgram Tests
// ============================================================================

func TestClientAdapter_SetProgram_StoresProgram(t *testing.T) {
	ca := adapter.NewClientAdapter()

	// We can't easily test this without a real program, but we can verify
	// the method doesn't panic
	ca.SetProgram(nil)
	// If we get here without panic, the method works
}

func TestClientAdapter_SetProgram_NilSafe(t *testing.T) {
	ca := adapter.NewClientAdapter()

	// Should not panic with nil
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetProgram(nil) panicked: %v", r)
		}
	}()

	ca.SetProgram(nil)
}

// ============================================================================
// State Tests
// ============================================================================

func TestClientAdapter_State_ReturnsCurrentState(t *testing.T) {
	ca := adapter.NewClientAdapter()

	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("State() = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}
}

func TestClientAdapter_State_ThreadSafe(t *testing.T) {
	ca := adapter.NewClientAdapter()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ca.State()
		}()
	}
	wg.Wait()
	// If we get here without race detector errors, the method is thread-safe
}

// ============================================================================
// ConnectCmd Tests
// ============================================================================

func TestClientAdapter_ConnectCmd_ReturnsTeaCmd(t *testing.T) {
	ca := adapter.NewClientAdapter()
	cmd := ca.ConnectCmd()

	if cmd == nil {
		t.Error("ConnectCmd() returned nil")
	}
}

func TestClientAdapter_ConnectCmd_WithMock_Success(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	cmd := ca.ConnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateConnected)
	}
	if stateMsg.Error != nil {
		t.Errorf("Error = %v, want nil", stateMsg.Error)
	}
}

func TestClientAdapter_ConnectCmd_WithMock_Failure(t *testing.T) {
	mock := newMockClient()
	mock.connectErr = errors.New("connection failed")
	ca := adapter.NewClientAdapterWithClient(mock)

	cmd := ca.ConnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
	if stateMsg.Error == nil {
		t.Error("Error should not be nil on connection failure")
	}
}

func TestClientAdapter_ConnectCmd_AlreadyConnected(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// First connect
	cmd := ca.ConnectCmd()
	cmd()

	// Second connect attempt
	cmd = ca.ConnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should either succeed (idempotent) or return already connected state
	if stateMsg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateConnected)
	}
}

func TestClientAdapter_ConnectCmd_StateTransition(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Initial state
	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("Initial state = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}

	// After connect
	cmd := ca.ConnectCmd()
	cmd()

	if ca.State() != adapter.ClientStateConnected {
		t.Errorf("After connect state = %v, want %v", ca.State(), adapter.ClientStateConnected)
	}
}

// ============================================================================
// QueryCmd Tests
// ============================================================================

func TestClientAdapter_QueryCmd_ReturnsTeaCmd(t *testing.T) {
	ca := adapter.NewClientAdapter()
	cmd := ca.QueryCmd("test prompt")

	if cmd == nil {
		t.Error("QueryCmd() returned nil")
	}
}

func TestClientAdapter_QueryCmd_WithMock_Success(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	connectCmd := ca.ConnectCmd()
	connectCmd()

	// Then query
	cmd := ca.QueryCmd("Hello Claude")
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateStreaming {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateStreaming)
	}
}

func TestClientAdapter_QueryCmd_NotConnected(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Query without connecting
	cmd := ca.QueryCmd("Hello Claude")
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
	if stateMsg.Error == nil {
		t.Error("Error should not be nil when not connected")
	}
}

func TestClientAdapter_QueryCmd_EmptyPrompt(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	connectCmd := ca.ConnectCmd()
	connectCmd()

	// Query with empty prompt
	cmd := ca.QueryCmd("")
	result := cmd()

	// Empty prompt should still work (SDK handles validation)
	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Either streaming or error depending on SDK behavior
	if stateMsg.State != adapter.ClientStateStreaming && stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want Streaming or Error", stateMsg.State)
	}
}

func TestClientAdapter_QueryCmd_StoresMessageChannel(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	connectCmd := ca.ConnectCmd()
	connectCmd()

	// Query
	cmd := ca.QueryCmd("Hello")
	cmd()

	// Message channel should be available
	ch := ca.MessageChannel()
	if ch == nil {
		t.Error("MessageChannel() returned nil after query")
	}
}

func TestClientAdapter_QueryCmd_Failure(t *testing.T) {
	mock := newMockClient()
	mock.queryErr = errors.New("query failed")
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	connectCmd := ca.ConnectCmd()
	connectCmd()

	// Query should fail
	cmd := ca.QueryCmd("Hello")
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
}

// ============================================================================
// InterruptCmd Tests
// ============================================================================

func TestClientAdapter_InterruptCmd_ReturnsTeaCmd(t *testing.T) {
	ca := adapter.NewClientAdapter()
	cmd := ca.InterruptCmd()

	if cmd == nil {
		t.Error("InterruptCmd() returned nil")
	}
}

func TestClientAdapter_InterruptCmd_Success(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect and query
	ca.ConnectCmd()()
	ca.QueryCmd("Hello")()

	// Interrupt
	cmd := ca.InterruptCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should return to connected state
	if stateMsg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateConnected)
	}
}

func TestClientAdapter_InterruptCmd_NotStreaming(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect but don't query
	ca.ConnectCmd()()

	// Interrupt when not streaming
	cmd := ca.InterruptCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should stay connected (no-op)
	if stateMsg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateConnected)
	}
}

func TestClientAdapter_InterruptCmd_NotConnected(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Interrupt without connecting
	cmd := ca.InterruptCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should return error state
	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
}

// ============================================================================
// DisconnectCmd Tests
// ============================================================================

func TestClientAdapter_DisconnectCmd_ReturnsTeaCmd(t *testing.T) {
	ca := adapter.NewClientAdapter()
	cmd := ca.DisconnectCmd()

	if cmd == nil {
		t.Error("DisconnectCmd() returned nil")
	}
}

func TestClientAdapter_DisconnectCmd_Success(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	ca.ConnectCmd()()

	// Disconnect
	cmd := ca.DisconnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateDisconnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateDisconnected)
	}
}

func TestClientAdapter_DisconnectCmd_AlreadyDisconnected(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Disconnect without connecting
	cmd := ca.DisconnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should stay disconnected (idempotent)
	if stateMsg.State != adapter.ClientStateDisconnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateDisconnected)
	}
}

func TestClientAdapter_DisconnectCmd_CleansUpResources(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect and query
	ca.ConnectCmd()()
	ca.QueryCmd("Hello")()

	// Verify channel exists
	if ca.MessageChannel() == nil {
		t.Error("MessageChannel should exist before disconnect")
	}

	// Disconnect
	ca.DisconnectCmd()()

	// Channel should be nil after disconnect
	if ca.MessageChannel() != nil {
		t.Error("MessageChannel should be nil after disconnect")
	}

	// State should be disconnected
	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("State = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}
}

func TestClientAdapter_DisconnectCmd_Failure(t *testing.T) {
	mock := newMockClient()
	mock.disconnectErr = errors.New("disconnect failed")
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect first
	ca.ConnectCmd()()

	// Disconnect should report error
	cmd := ca.DisconnectCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
	if stateMsg.Error == nil {
		t.Error("Error should not be nil on disconnect failure")
	}
}

// ============================================================================
// MessageChannel Tests
// ============================================================================

func TestClientAdapter_MessageChannel_ReturnsChannel(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect and query
	ca.ConnectCmd()()
	ca.QueryCmd("Hello")()

	ch := ca.MessageChannel()
	if ch == nil {
		t.Error("MessageChannel() returned nil")
	}
}

func TestClientAdapter_MessageChannel_NilWhenNotStreaming(t *testing.T) {
	ca := adapter.NewClientAdapter()

	ch := ca.MessageChannel()
	if ch != nil {
		t.Errorf("MessageChannel() = %v, want nil when not streaming", ch)
	}
}

func TestClientAdapter_MessageChannel_NilAfterDisconnect(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect, query, then disconnect
	ca.ConnectCmd()()
	ca.QueryCmd("Hello")()
	ca.DisconnectCmd()()

	ch := ca.MessageChannel()
	if ch != nil {
		t.Error("MessageChannel() should be nil after disconnect")
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestClientAdapter_ConcurrentConnect(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := ca.ConnectCmd()
			cmd()
		}()
	}
	wg.Wait()

	// Should end up connected
	if ca.State() != adapter.ClientStateConnected {
		t.Errorf("Final state = %v, want %v", ca.State(), adapter.ClientStateConnected)
	}
}

func TestClientAdapter_ConcurrentStateReads(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect
	ca.ConnectCmd()()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ca.State()
		}()
	}
	wg.Wait()
	// If no race condition, test passes
}

func TestClientAdapter_RapidConnectDisconnect(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	for i := 0; i < 10; i++ {
		ca.ConnectCmd()()
		ca.DisconnectCmd()()
	}

	// Should end up disconnected
	if ca.State() != adapter.ClientStateDisconnected {
		t.Errorf("Final state = %v, want %v", ca.State(), adapter.ClientStateDisconnected)
	}
}

func TestClientAdapter_ConcurrentOperations(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	var wg sync.WaitGroup

	// Concurrent connects
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ca.ConnectCmd()()
		}()
	}

	// Concurrent state reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ca.State()
		}()
	}

	// Concurrent message channel reads
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ca.MessageChannel()
		}()
	}

	wg.Wait()
	// If no race condition, test passes
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestClientAdapter_QueryAfterDisconnect(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect, then disconnect
	ca.ConnectCmd()()
	ca.DisconnectCmd()()

	// Try to query
	cmd := ca.QueryCmd("Hello")
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	if stateMsg.State != adapter.ClientStateError {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateError)
	}
}

func TestClientAdapter_InterruptWithoutQuery(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect but don't query
	ca.ConnectCmd()()

	// Interrupt (should be no-op)
	cmd := ca.InterruptCmd()
	result := cmd()

	stateMsg, ok := result.(adapter.ClientStateMsg)
	if !ok {
		t.Fatalf("expected ClientStateMsg, got %T", result)
	}

	// Should stay connected
	if stateMsg.State != adapter.ClientStateConnected {
		t.Errorf("State = %v, want %v", stateMsg.State, adapter.ClientStateConnected)
	}
}

func TestClientAdapter_ContextCancellation(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	ctx, cancel := context.WithCancel(context.Background())

	// Start connect with context
	_ = ca.ConnectCmdWithContext(ctx)

	// Cancel immediately
	cancel()

	// Give time for cancellation to propagate
	time.Sleep(10 * time.Millisecond)

	// State should handle cancellation gracefully
	// (exact behavior depends on implementation)
}

func TestClientAdapter_MultipleQueries(t *testing.T) {
	mock := newMockClient()
	ca := adapter.NewClientAdapterWithClient(mock)

	// Connect
	ca.ConnectCmd()()

	// Multiple queries
	for i := 0; i < 3; i++ {
		cmd := ca.QueryCmd("Query " + string(rune('A'+i)))
		result := cmd()

		stateMsg, ok := result.(adapter.ClientStateMsg)
		if !ok {
			t.Fatalf("expected ClientStateMsg, got %T", result)
		}

		if stateMsg.State != adapter.ClientStateStreaming && stateMsg.State != adapter.ClientStateError {
			t.Errorf("Query %d: State = %v, want Streaming or Error", i, stateMsg.State)
		}
	}
}
