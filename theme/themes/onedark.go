package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("onedark", OneDark())
}

// OneDark returns the One Dark theme palette.
// The iconic Atom editor theme.
func OneDark() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#a626a4", Dark: "#c678dd"},
		Accent:            lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		Error:             lipgloss.AdaptiveColor{Light: "#e45649", Dark: "#e06c75"},
		Warning:           lipgloss.AdaptiveColor{Light: "#c18401", Dark: "#e5c07b"},
		Success:           lipgloss.AdaptiveColor{Light: "#50a14f", Dark: "#98c379"},
		Info:              lipgloss.AdaptiveColor{Light: "#986801", Dark: "#d19a66"},
		Text:              lipgloss.AdaptiveColor{Light: "#383a42", Dark: "#abb2bf"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"},
		Background:        lipgloss.AdaptiveColor{Light: "#fafafa", Dark: "#282c34"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#f0f0f1", Dark: "#21252b"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#eaeaeb", Dark: "#353b45"},
		Border:            lipgloss.AdaptiveColor{Light: "#d1d1d2", Dark: "#393f4a"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#e0e0e1", Dark: "#2c313a"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#50a14f", Dark: "#98c379"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#e45649", Dark: "#e06c75"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#489447", Dark: "#aad482"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#d65145", Dark: "#e8828b"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#eafbe9", Dark: "#2c382b"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#fce9e8", Dark: "#3a2d2f"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#f0f0f1", Dark: "#21252b"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#c9c9ca", Dark: "#495162"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#e1f3df", Dark: "#283427"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#f5e2e1", Dark: "#36292b"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#383a42", Dark: "#abb2bf"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#a626a4", Dark: "#c678dd"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#50a14f", Dark: "#98c379"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#c18401", Dark: "#e5c07b"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#986801", Dark: "#d19a66"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#383a42", Dark: "#abb2bf"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#a0a1a7", Dark: "#5c6370"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#a626a4", Dark: "#c678dd"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#4078f2", Dark: "#61afef"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#e45649", Dark: "#e06c75"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#50a14f", Dark: "#98c379"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#986801", Dark: "#d19a66"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#c18401", Dark: "#e5c07b"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#0184bc", Dark: "#56b6c2"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#383a42", Dark: "#abb2bf"},
	}
}
