package streamtext

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Cursor is the character displayed at the end of streaming text.
const Cursor = "▌"

// Model represents the streaming text component state.
type Model struct {
	content   strings.Builder
	streaming bool
	width     int
}

// Option configures the Model.
type Option func(*Model)

// New creates a new StreamText model with the given options.
func New(opts ...Option) Model {
	m := Model{}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// WithWidth sets the display width for text wrapping.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// Append adds text to the content buffer.
func (m *Model) Append(text string) {
	m.content.WriteString(text)
}

// Clear resets the content buffer.
func (m *Model) Clear() {
	m.content.Reset()
}

// SetStreaming controls cursor visibility.
func (m *Model) SetStreaming(streaming bool) {
	m.streaming = streaming
}

// Content returns the accumulated text.
func (m Model) Content() string {
	return m.content.String()
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model - renders content with optional cursor.
func (m Model) View() string {
	text := m.content.String()
	if m.streaming {
		text += Cursor
	}
	return text
}
