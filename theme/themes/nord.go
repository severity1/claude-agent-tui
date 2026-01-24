package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("nord", Nord())
}

// Nord returns the Nord theme palette.
// An arctic, north-bluish color palette.
func Nord() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#5E81AC", Dark: "#88C0D0"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#81A1C1"},
		Accent:            lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		Error:             lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#BF616A"},
		Warning:           lipgloss.AdaptiveColor{Light: "#D08770", Dark: "#D08770"},
		Success:           lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#A3BE8C"},
		Info:              lipgloss.AdaptiveColor{Light: "#5E81AC", Dark: "#88C0D0"},
		Text:              lipgloss.AdaptiveColor{Light: "#2E3440", Dark: "#ECEFF4"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#3B4252", Dark: "#8B95A7"},
		Background:        lipgloss.AdaptiveColor{Light: "#ECEFF4", Dark: "#2E3440"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#D8DEE9", Dark: "#434C5E"},
		Border:            lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#434C5E"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#434C5E", Dark: "#4C566A"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#434C5E"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#A3BE8C"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#BF616A"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#8B95A7"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#8B95A7"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#A3BE8C"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#BF616A"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#D8DEE9", Dark: "#434C5E"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#3B4252"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#2E3440", Dark: "#D8DEE9"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#5E81AC", Dark: "#88C0D0"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#81A1C1"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#A3BE8C"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#8B95A7"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#D08770", Dark: "#D08770"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#EBCB8B", Dark: "#EBCB8B"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#8B95A7"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#5E81AC", Dark: "#88C0D0"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#81A1C1"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#2E3440", Dark: "#D8DEE9"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#4C566A", Dark: "#8B95A7"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#81A1C1"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#88C0D0", Dark: "#88C0D0"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#A3BE8C"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#B48EAD", Dark: "#B48EAD"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#8FBCBB", Dark: "#8FBCBB"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#81A1C1"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#2E3440", Dark: "#D8DEE9"},
	}
}
