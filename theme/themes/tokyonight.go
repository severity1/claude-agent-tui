package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("tokyonight", TokyoNight())
}

// TokyoNight returns the Tokyo Night theme palette.
// A clean, dark theme inspired by Tokyo city lights.
func TokyoNight() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#9854f1", Dark: "#c099ff"},
		Accent:            lipgloss.AdaptiveColor{Light: "#b15c00", Dark: "#ff966c"},
		Error:             lipgloss.AdaptiveColor{Light: "#f52a65", Dark: "#ff757f"},
		Warning:           lipgloss.AdaptiveColor{Light: "#b15c00", Dark: "#ff966c"},
		Success:           lipgloss.AdaptiveColor{Light: "#587539", Dark: "#c3e88d"},
		Info:              lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		Text:              lipgloss.AdaptiveColor{Light: "#3760bf", Dark: "#c8d3f5"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#8990a3", Dark: "#828bb8"},
		Background:        lipgloss.AdaptiveColor{Light: "#e1e2e7", Dark: "#1a1b26"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#d5d6db", Dark: "#1e2030"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#c8c9ce", Dark: "#222436"},
		Border:            lipgloss.AdaptiveColor{Light: "#737a8c", Dark: "#737aa2"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#5a607d", Dark: "#9099b2"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#9699a8", Dark: "#545c7e"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#1e725c", Dark: "#4fd6be"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#c53b53", Dark: "#c53b53"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#7086b5", Dark: "#828bb8"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#7086b5", Dark: "#828bb8"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#4db380", Dark: "#b8db87"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#f52a65", Dark: "#e26a75"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#d5e5d5", Dark: "#20303b"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#f7d8db", Dark: "#37222c"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#d5d6db", Dark: "#1e2030"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#c8c9ce", Dark: "#222436"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#c5d5c5", Dark: "#1b2b34"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#e7c8cb", Dark: "#2d1f26"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#3760bf", Dark: "#c8d3f5"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#9854f1", Dark: "#c099ff"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#007197", Dark: "#86e1fc"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#587539", Dark: "#c3e88d"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#8c6c3e", Dark: "#ffc777"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#8c6c3e", Dark: "#ffc777"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#b15c00", Dark: "#ff966c"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#8990a3", Dark: "#828bb8"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#007197", Dark: "#86e1fc"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#007197", Dark: "#86e1fc"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#3760bf", Dark: "#c8d3f5"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#8990a3", Dark: "#828bb8"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#9854f1", Dark: "#c099ff"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#2e7de9", Dark: "#82aaff"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#f52a65", Dark: "#ff757f"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#587539", Dark: "#c3e88d"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#b15c00", Dark: "#ff966c"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#8c6c3e", Dark: "#ffc777"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#007197", Dark: "#86e1fc"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#3760bf", Dark: "#c8d3f5"},
	}
}
