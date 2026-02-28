package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("dracula", Dracula())
}

// Dracula returns the Dracula theme palette.
// A classic dark theme with vibrant colors.
func Dracula() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
		Accent:            lipgloss.AdaptiveColor{Light: "#8be9fd", Dark: "#8be9fd"},
		Error:             lipgloss.AdaptiveColor{Light: "#ff5555", Dark: "#ff5555"},
		Warning:           lipgloss.AdaptiveColor{Light: "#f1fa8c", Dark: "#f1fa8c"},
		Success:           lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
		Info:              lipgloss.AdaptiveColor{Light: "#ffb86c", Dark: "#ffb86c"},
		Text:              lipgloss.AdaptiveColor{Light: "#282a36", Dark: "#f8f8f2"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		Background:        lipgloss.AdaptiveColor{Light: "#f8f8f2", Dark: "#282a36"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#e8e8e2", Dark: "#21222c"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#d8d8d2", Dark: "#44475a"},
		Border:            lipgloss.AdaptiveColor{Light: "#c8c8c2", Dark: "#44475a"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#e0e0e0", Dark: "#191a21"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#ff5555", Dark: "#ff5555"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#ff5555", Dark: "#ff5555"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#e0ffe0", Dark: "#1a3a1a"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#ffe0e0", Dark: "#3a1a1a"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#e8e8e2", Dark: "#21222c"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#c8c8c2", Dark: "#44475a"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#e0ffe0", Dark: "#1a3a1a"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#ffe0e0", Dark: "#3a1a1a"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#282a36", Dark: "#f8f8f2"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#8be9fd", Dark: "#8be9fd"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#f1fa8c", Dark: "#f1fa8c"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#ffb86c", Dark: "#ffb86c"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#8be9fd", Dark: "#8be9fd"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#8be9fd", Dark: "#8be9fd"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#282a36", Dark: "#f8f8f2"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#282a36", Dark: "#f8f8f2"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#f1fa8c", Dark: "#f1fa8c"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#8be9fd", Dark: "#8be9fd"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#282a36", Dark: "#f8f8f2"},
	}
}
