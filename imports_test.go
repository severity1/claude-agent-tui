package tui_test

import (
	"testing"

	// Charmbracelet dependencies
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// SDK dependency
	claudecode "github.com/severity1/claude-agent-sdk-go"
)

func TestDependencyImports(t *testing.T) {
	t.Run("bubbletea imports", func(t *testing.T) {
		var _ tea.Model
		var _ tea.Cmd
	})

	t.Run("lipgloss imports", func(t *testing.T) {
		_ = lipgloss.NewStyle()
	})

	t.Run("bubbles textarea imports", func(t *testing.T) {
		_ = textarea.New()
	})

	t.Run("bubbles viewport imports", func(t *testing.T) {
		_ = viewport.New(80, 24)
	})

	t.Run("claude-agent-sdk-go imports", func(t *testing.T) {
		var _ *claudecode.Client
	})
}
