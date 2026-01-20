// Package main demonstrates the StreamText component for streaming text display.
//
// Usage:
//
//	go run example/streamtext/main.go
//
// Controls:
//   - r: Restart streaming
//   - q/Ctrl+C: Quit
package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/severity1/claude-agent-tui/component/output/streamtext"
)

// chunkDelay is the time between streaming each chunk.
const chunkDelay = 50 * time.Millisecond

// tickMsg triggers the next chunk of text.
type tickMsg struct{}

// sampleChunks contains the predefined text chunks to stream.
var sampleChunks = []string{
	"Hello! I'm happy to help you with your question.\n\n",
	"Let me think about this for a moment...\n\n",
	"Based on my analysis, here are the key points:\n\n",
	"1. First, consider the context\n",
	"2. Then, evaluate the options\n",
	"3. Finally, make a decision\n\n",
	"Is there anything else you'd like to know?",
}

// model holds the application state.
type model struct {
	streamtext streamtext.Model
	chunks     []string
	index      int
	done       bool
}

// newModel creates a new model with default state.
func newModel() model {
	return model{
		streamtext: streamtext.New(),
		chunks:     sampleChunks,
		index:      0,
		done:       false,
	}
}

// tick returns a command that sends a tickMsg after the delay.
func tick() tea.Cmd {
	return tea.Tick(chunkDelay, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Init starts the streaming by setting streaming state and scheduling first tick.
func (m model) Init() tea.Cmd {
	m.streamtext.SetStreaming(true)
	return tick()
}

// Update handles messages and updates the model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			// Restart streaming
			m.streamtext.Clear()
			m.streamtext.SetStreaming(true)
			m.index = 0
			m.done = false
			return m, tick()
		}

	case tickMsg:
		if m.index < len(m.chunks) {
			m.streamtext.Append(m.chunks[m.index])
			m.index++
			if m.index >= len(m.chunks) {
				// Streaming complete
				m.streamtext.SetStreaming(false)
				m.done = true
				return m, nil
			}
			return m, tick()
		}
	}

	return m, nil
}

// View renders the current state.
func (m model) View() string {
	var status string
	if m.done {
		status = "[Done]"
	} else {
		status = "[Streaming...]"
	}

	return fmt.Sprintf(
		"StreamText Component Demo\n"+
			"=========================\n\n"+
			"%s\n\n"+
			"%s\n\n"+
			"Controls: r = restart, q = quit\n",
		m.streamtext.View(),
		status,
	)
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
