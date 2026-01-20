package adapter

import (
	"context"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	claudecode "github.com/severity1/claude-agent-sdk-go"
)

// ClientState represents the connection state of the ClientAdapter.
type ClientState int

// Client state constants.
const (
	ClientStateDisconnected ClientState = iota
	ClientStateConnecting
	ClientStateConnected
	ClientStateStreaming
	ClientStateError
)

// String returns a string representation of the ClientState.
func (s ClientState) String() string {
	switch s {
	case ClientStateDisconnected:
		return "disconnected"
	case ClientStateConnecting:
		return "connecting"
	case ClientStateConnected:
		return "connected"
	case ClientStateStreaming:
		return "streaming"
	case ClientStateError:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// ClientStateMsg signals client state changes to the Bubble Tea program.
type ClientStateMsg struct {
	State ClientState
	Error error
}

// ClientAdapter wraps the SDK client for Bubble Tea integration.
// It provides tea.Cmd methods for connection lifecycle management.
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

// NewClientAdapter creates a new ClientAdapter with the given SDK options.
// The adapter starts in the Disconnected state.
func NewClientAdapter(opts ...claudecode.Option) *ClientAdapter {
	return &ClientAdapter{
		state:   ClientStateDisconnected,
		options: opts,
	}
}

// NewClientAdapterWithClient creates a ClientAdapter with an existing client.
// This is primarily for testing with mock clients.
func NewClientAdapterWithClient(client claudecode.Client, opts ...claudecode.Option) *ClientAdapter {
	return &ClientAdapter{
		client:  client,
		state:   ClientStateDisconnected,
		options: opts,
	}
}

// SetProgram sets the Bubble Tea program for sending messages.
// This should be called after the program is created but before connecting.
func (c *ClientAdapter) SetProgram(p *tea.Program) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.program = p
}

// State returns the current connection state.
// This method is thread-safe.
func (c *ClientAdapter) State() ClientState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// setState updates the internal state (must be called with lock held or within command).
func (c *ClientAdapter) setState(state ClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = state
}

// ConnectCmd returns a tea.Cmd that establishes a connection to the SDK.
// It emits ClientStateMsg upon completion.
func (c *ClientAdapter) ConnectCmd() tea.Cmd {
	return c.ConnectCmdWithContext(context.Background())
}

// ConnectCmdWithContext returns a tea.Cmd that establishes a connection
// with the given context for cancellation support.
func (c *ClientAdapter) ConnectCmdWithContext(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		c.mu.Lock()

		// Already connected or connecting - return actual current state
		if c.state == ClientStateConnected || c.state == ClientStateConnecting {
			currentState := c.state
			c.mu.Unlock()
			return ClientStateMsg{State: currentState}
		}

		// Set state to connecting and initialize resources atomically
		c.state = ClientStateConnecting
		c.ctx, c.cancel = context.WithCancel(ctx)
		if c.client == nil {
			c.client = claudecode.NewClient(c.options...)
		}
		client := c.client
		connCtx := c.ctx
		c.mu.Unlock()

		// Connect the client (blocking operation, outside lock)
		if err := client.Connect(connCtx); err != nil {
			c.setState(ClientStateError)
			return ClientStateMsg{State: ClientStateError, Error: fmt.Errorf("connect failed: %w", err)}
		}

		c.setState(ClientStateConnected)
		return ClientStateMsg{State: ClientStateConnected}
	}
}

// QueryCmd returns a tea.Cmd that sends a query to the SDK.
// The adapter must be connected before calling this method.
// It emits ClientStateMsg upon completion and stores the message channel.
func (c *ClientAdapter) QueryCmd(prompt string) tea.Cmd {
	return func() tea.Msg {
		c.mu.RLock()
		state := c.state
		client := c.client
		ctx := c.ctx
		c.mu.RUnlock()

		// Must be connected to query
		if state != ClientStateConnected && state != ClientStateStreaming {
			return ClientStateMsg{
				State: ClientStateError,
				Error: fmt.Errorf("cannot query: client not connected (state: %s)", state),
			}
		}

		if client == nil {
			return ClientStateMsg{
				State: ClientStateError,
				Error: fmt.Errorf("cannot query: client is nil"),
			}
		}

		// Use background context if no context set
		if ctx == nil {
			ctx = context.Background()
		}

		// Send the query
		if err := client.Query(ctx, prompt); err != nil {
			c.setState(ClientStateError)
			return ClientStateMsg{State: ClientStateError, Error: fmt.Errorf("query failed: %w", err)}
		}

		// Get and store the message channel
		msgChan := client.ReceiveMessages(ctx)

		c.mu.Lock()
		c.msgChan = msgChan
		c.state = ClientStateStreaming
		c.mu.Unlock()

		return ClientStateMsg{State: ClientStateStreaming}
	}
}

// InterruptCmd returns a tea.Cmd that interrupts the current streaming operation.
// It emits ClientStateMsg upon completion.
func (c *ClientAdapter) InterruptCmd() tea.Cmd {
	return func() tea.Msg {
		c.mu.RLock()
		state := c.state
		client := c.client
		ctx := c.ctx
		c.mu.RUnlock()

		// Not connected
		if state == ClientStateDisconnected {
			return ClientStateMsg{
				State: ClientStateError,
				Error: fmt.Errorf("cannot interrupt: client not connected"),
			}
		}

		// Not streaming - no-op, return current state
		if state != ClientStateStreaming {
			return ClientStateMsg{State: state}
		}

		if client == nil {
			return ClientStateMsg{
				State: ClientStateError,
				Error: fmt.Errorf("cannot interrupt: client is nil"),
			}
		}

		// Use background context if no context set
		if ctx == nil {
			ctx = context.Background()
		}

		// Send interrupt signal
		if err := client.Interrupt(ctx); err != nil {
			c.setState(ClientStateError)
			return ClientStateMsg{State: ClientStateError, Error: fmt.Errorf("interrupt failed: %w", err)}
		}

		c.setState(ClientStateConnected)
		return ClientStateMsg{State: ClientStateConnected}
	}
}

// DisconnectCmd returns a tea.Cmd that disconnects from the SDK.
// It emits ClientStateMsg upon completion and cleans up resources.
// Resources are always cleaned up, even if disconnect fails.
func (c *ClientAdapter) DisconnectCmd() tea.Cmd {
	return func() tea.Msg {
		c.mu.Lock()

		// Already disconnected
		if c.state == ClientStateDisconnected {
			c.mu.Unlock()
			return ClientStateMsg{State: ClientStateDisconnected}
		}

		client := c.client
		cancel := c.cancel
		c.mu.Unlock()

		// Cancel context first
		if cancel != nil {
			cancel()
		}

		// Disconnect client and capture any error
		var disconnectErr error
		if client != nil {
			disconnectErr = client.Disconnect()
		}

		// Always clean up resources, even if disconnect failed
		c.mu.Lock()
		c.client = nil
		c.msgChan = nil
		c.ctx = nil
		c.cancel = nil
		if disconnectErr != nil {
			c.state = ClientStateError
			c.mu.Unlock()
			return ClientStateMsg{
				State: ClientStateError,
				Error: fmt.Errorf("disconnect failed (resources cleaned up): %w", disconnectErr),
			}
		}
		c.state = ClientStateDisconnected
		c.mu.Unlock()

		return ClientStateMsg{State: ClientStateDisconnected}
	}
}

// MessageChannel returns the current message channel for streaming.
// Returns nil if not currently streaming or if not connected.
func (c *ClientAdapter) MessageChannel() <-chan claudecode.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.msgChan
}
