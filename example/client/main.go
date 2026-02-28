// Package main demonstrates the ClientAdapter for SDK lifecycle management.
//
// This example shows how to use the ClientAdapter to manage the connection
// lifecycle with the Claude Code CLI, send queries, and receive streaming
// messages via Bubble Tea commands.
//
// Usage:
//
//	go run example/client/main.go              # Simulated demo (no SDK required)
//	go run example/client/main.go --live       # Live mode (requires claude CLI)
//
// In simulated mode, it demonstrates the ClientAdapter state machine without
// requiring a real Claude connection. In live mode, it connects to the actual
// Claude Code CLI.
//
// Controls:
//   - c: Connect (when disconnected)
//   - q: Send query (when connected)
//   - i: Interrupt (when streaming)
//   - d: Disconnect (when connected)
//   - Ctrl+C: Quit
package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	claudecode "github.com/severity1/claude-agent-sdk-go"

	"github.com/severity1/claude-agent-tui/adapter"
	"github.com/severity1/claude-agent-tui/component/output/streamtext"
)

// simulatedClient provides a mock implementation for demonstration.
type simulatedClient struct {
	connected bool
	msgChan   chan claudecode.Message
	mu        sync.Mutex
}

func newSimulatedClient() *simulatedClient {
	return &simulatedClient{
		msgChan: make(chan claudecode.Message, 10),
	}
}

func (s *simulatedClient) Connect(_ context.Context, _ ...claudecode.StreamMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	time.Sleep(100 * time.Millisecond) // Simulate connection delay
	s.connected = true
	return nil
}

func (s *simulatedClient) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connected = false
	return nil
}

func (s *simulatedClient) Query(_ context.Context, prompt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Simulate sending some messages
	go func() {
		// Simulate streaming response
		chunks := []string{
			"Hello! ",
			"I received your query: ",
			fmt.Sprintf("%q. ", prompt),
			"Here is my response.\n\n",
			"This is a simulated demonstration ",
			"of the ClientAdapter lifecycle. ",
			"The adapter manages:\n",
			"- Connection state\n",
			"- Query execution\n",
			"- Message streaming\n",
			"- Resource cleanup\n",
		}
		for _, chunk := range chunks {
			time.Sleep(50 * time.Millisecond)
			s.msgChan <- &claudecode.StreamEvent{
				Event: map[string]any{
					"type":  "content_block_delta",
					"index": float64(0),
					"delta": map[string]any{
						"type": "text_delta",
						"text": chunk,
					},
				},
			}
		}
		// Signal completion
		time.Sleep(50 * time.Millisecond)
		s.msgChan <- &claudecode.StreamEvent{
			Event: map[string]any{
				"type": "message_stop",
			},
		}
	}()
	return nil
}

func (s *simulatedClient) QueryWithSession(_ context.Context, _ string, _ string) error {
	return nil
}

func (s *simulatedClient) QueryStream(_ context.Context, _ <-chan claudecode.StreamMessage) error {
	return nil
}

func (s *simulatedClient) ReceiveMessages(_ context.Context) <-chan claudecode.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.msgChan
}

func (s *simulatedClient) ReceiveResponse(_ context.Context) claudecode.MessageIterator {
	return nil
}

func (s *simulatedClient) Interrupt(_ context.Context) error {
	return nil
}

func (s *simulatedClient) SetModel(_ context.Context, _ *string) error {
	return nil
}

func (s *simulatedClient) SetPermissionMode(_ context.Context, _ claudecode.PermissionMode) error {
	return nil
}

func (s *simulatedClient) RewindFiles(_ context.Context, _ string) error {
	return nil
}

func (s *simulatedClient) GetStreamIssues() []claudecode.StreamIssue {
	return nil
}

func (s *simulatedClient) GetStreamStats() claudecode.StreamStats {
	return claudecode.StreamStats{}
}

func (s *simulatedClient) GetServerInfo(_ context.Context) (map[string]interface{}, error) {
	return nil, nil
}

// model is the Bubble Tea application model.
type model struct {
	client     *adapter.ClientAdapter
	streamtext streamtext.Model
	ctx        context.Context
	cancel     context.CancelFunc
	msgChan    <-chan claudecode.Message
	logs       []string
	maxLogs    int
	liveMode   bool
}

// newModel creates a new application model.
func newModel(liveMode bool) model {
	ctx, cancel := context.WithCancel(context.Background())

	var ca *adapter.ClientAdapter
	if liveMode {
		// Live mode - use real SDK client
		ca = adapter.NewClientAdapter()
	} else {
		// Simulated mode - use mock client
		ca = adapter.NewClientAdapterWithClient(newSimulatedClient())
	}

	return model{
		client:     ca,
		streamtext: streamtext.New(),
		ctx:        ctx,
		cancel:     cancel,
		logs:       make([]string, 0),
		maxLogs:    10,
		liveMode:   liveMode,
	}
}

// addLog adds a log entry with timestamp.
func (m *model) addLog(format string, args ...any) {
	entry := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
	m.logs = append(m.logs, entry)
	if len(m.logs) > m.maxLogs {
		m.logs = m.logs[len(m.logs)-m.maxLogs:]
	}
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	m.addLog("Application started")
	return nil
}

// Update handles messages and returns the updated model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			if m.client.State() != adapter.ClientStateDisconnected {
				if msg := m.client.DisconnectCmd()(); msg != nil {
					if sm, ok := msg.(adapter.ClientStateMsg); ok && sm.Error != nil {
						m.addLog("Disconnect error: %v", sm.Error)
					}
				}
			}
			return m, tea.Quit

		case "c":
			if m.client.State() == adapter.ClientStateDisconnected {
				m.addLog("Connecting...")
				return m, m.client.ConnectCmdWithContext(m.ctx)
			}
			m.addLog("Already connected (state: %s)", m.client.State())

		case "q":
			if m.client.State() == adapter.ClientStateConnected {
				m.addLog("Sending query...")
				m.streamtext.Clear()
				m.streamtext.SetStreaming(true)
				return m, m.client.QueryCmd("Hello Claude, tell me about yourself")
			}
			m.addLog("Cannot query (state: %s)", m.client.State())

		case "i":
			if m.client.State() == adapter.ClientStateStreaming {
				m.addLog("Interrupting...")
				return m, m.client.InterruptCmd()
			}
			m.addLog("Not streaming (state: %s)", m.client.State())

		case "d":
			if m.client.State() != adapter.ClientStateDisconnected {
				m.addLog("Disconnecting...")
				m.streamtext.SetStreaming(false)
				m.msgChan = nil
				return m, m.client.DisconnectCmd()
			}
			m.addLog("Already disconnected")
		}

	case adapter.ClientStateMsg:
		m.addLog("State changed: %s", msg.State)
		if msg.Error != nil {
			m.addLog("Error: %v", msg.Error)
		}

		switch msg.State {
		case adapter.ClientStateConnected:
			m.addLog("Ready for queries")
		case adapter.ClientStateStreaming:
			m.msgChan = m.client.MessageChannel()
			if m.msgChan != nil {
				m.addLog("Streaming started, receiving messages")
				return m, adapter.StreamCmd(m.ctx, m.msgChan)
			}
		case adapter.ClientStateError:
			m.streamtext.SetStreaming(false)
		case adapter.ClientStateDisconnected:
			m.streamtext.SetStreaming(false)
			m.msgChan = nil
		}

	case adapter.StreamDeltaMsg:
		m.streamtext.Append(msg.Text)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.StreamDoneMsg:
		m.addLog("Stream complete")
		m.streamtext.SetStreaming(false)
		// Return to connected state
		return m, func() tea.Msg {
			return adapter.ClientStateMsg{State: adapter.ClientStateConnected}
		}

	case adapter.StreamErrorMsg:
		m.addLog("Stream error: %v", msg.Err)
		m.streamtext.SetStreaming(false)

	// Handle other stream message types
	case adapter.MessageStartMsg:
		m.addLog("Message started: model=%s", msg.Model)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.MessageDeltaMsg:
		m.addLog("Message delta: stop=%s tokens=%d", msg.StopReason, msg.OutputTokens)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.StreamBlockStartMsg:
		m.addLog("Block started: type=%s", msg.BlockType)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.StreamBlockStopMsg:
		m.addLog("Block stopped: index=%d", msg.Index)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.ThinkingDeltaMsg:
		m.addLog("Thinking: %s", truncate(msg.Text, 30))
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.ToolUseDeltaMsg:
		m.addLog("Tool use: %s", truncate(msg.PartialJSON, 30))
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.UnknownMessageMsg:
		m.addLog("Unknown message: %s", msg.TypeName)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	// Handle complete message types (non-streaming SDK output)
	case adapter.AssistantMsg:
		// Extract text content from assistant message
		if msg.Message != nil {
			for _, block := range msg.Message.Content {
				if textBlock, ok := block.(*claudecode.TextBlock); ok {
					m.streamtext.Append(textBlock.Text)
				}
			}
		}
		m.addLog("Assistant message received")
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.ResultMsg:
		m.addLog("Result: tokens_in=%d tokens_out=%d", msg.InputTokens, msg.OutputTokens)
		m.streamtext.SetStreaming(false)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.UserMsg:
		m.addLog("User message received")
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.SystemInitMsg:
		m.addLog("System init: session=%s, tools=%d", msg.SessionID, len(msg.Tools))
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.SystemHookResponseMsg:
		m.addLog("Hook response: %s", msg.HookName)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.ControlRequestMsg:
		m.addLog("Control request: id=%s", msg.RequestID)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)

	case adapter.ControlResponseMsg:
		m.addLog("Control response: id=%s", msg.RequestID)
		return m, adapter.StreamCmd(m.ctx, m.msgChan)
	}

	return m, nil
}

// View renders the current state.
func (m model) View() string {
	mode := "SIMULATED"
	if m.liveMode {
		mode = "LIVE"
	}

	state := m.client.State()
	stateColor := ""
	switch state {
	case adapter.ClientStateConnected:
		stateColor = "\033[32m" // green
	case adapter.ClientStateStreaming:
		stateColor = "\033[33m" // yellow
	case adapter.ClientStateError:
		stateColor = "\033[31m" // red
	default:
		stateColor = "\033[37m" // white
	}

	controls := "c=connect"
	switch state {
	case adapter.ClientStateConnected:
		controls = "q=query, d=disconnect"
	case adapter.ClientStateStreaming:
		controls = "i=interrupt, d=disconnect"
	case adapter.ClientStateError:
		controls = "c=reconnect, d=disconnect"
	}

	// Build log view
	logView := ""
	for _, log := range m.logs {
		logView += "  " + log + "\n"
	}

	return fmt.Sprintf(
		"ClientAdapter Demo (%s)\n"+
			"================================\n\n"+
			"State: %s%s\033[0m\n\n"+
			"--- Output ---\n%s\n\n"+
			"--- Event Log ---\n%s\n"+
			"Controls: %s, Ctrl+C=quit\n",
		mode,
		stateColor, state,
		m.streamtext.View(),
		logView,
		controls,
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func main() {
	liveMode := len(os.Args) > 1 && os.Args[1] == "--live"

	fmt.Println("ClientAdapter Example")
	fmt.Println("====================")
	if liveMode {
		fmt.Println("Mode: LIVE (connecting to Claude CLI)")
	} else {
		fmt.Println("Mode: SIMULATED (no Claude CLI required)")
	}
	fmt.Println()

	p := tea.NewProgram(newModel(liveMode))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
