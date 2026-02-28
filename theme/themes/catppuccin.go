package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("catppuccin", Catppuccin())
}

// Catppuccin returns the Catppuccin theme palette.
// Dark mode uses Mocha variant, light mode uses Latte variant.
func Catppuccin() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#cba6f7"},
		Accent:            lipgloss.AdaptiveColor{Light: "#ea76cb", Dark: "#f5c2e7"},
		Error:             lipgloss.AdaptiveColor{Light: "#d20f39", Dark: "#f38ba8"},
		Warning:           lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f9e2af"},
		Success:           lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"},
		Info:              lipgloss.AdaptiveColor{Light: "#179299", Dark: "#94e2d5"},
		Text:              lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cdd6f4"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#5c5f77", Dark: "#bac2de"},
		Background:        lipgloss.AdaptiveColor{Light: "#eff1f5", Dark: "#1e1e2e"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#e6e9ef", Dark: "#181825"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#dce0e8", Dark: "#11111b"},
		Border:            lipgloss.AdaptiveColor{Light: "#ccd0da", Dark: "#313244"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#bcc0cc", Dark: "#45475a"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#acb0be", Dark: "#585b70"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#d20f39", Dark: "#f38ba8"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#7c7f93", Dark: "#9399b2"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#fe640b", Dark: "#fab387"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#d20f39", Dark: "#f38ba8"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#d6f0d9", Dark: "#24312b"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#f6dfe2", Dark: "#3c2a32"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#e6e9ef", Dark: "#181825"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#bcc0cc", Dark: "#45475a"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#c9e3cb", Dark: "#1e2a25"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#e9d3d6", Dark: "#32232a"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cdd6f4"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#cba6f7"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#04a5e5", Dark: "#89dceb"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f9e2af"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f9e2af"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#fe640b", Dark: "#fab387"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#6c6f85", Dark: "#a6adc8"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#04a5e5", Dark: "#89dceb"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#04a5e5", Dark: "#89dceb"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cdd6f4"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#7c7f93", Dark: "#9399b2"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#cba6f7"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#d20f39", Dark: "#f38ba8"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#fe640b", Dark: "#fab387"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f9e2af"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#04a5e5", Dark: "#89dceb"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cdd6f4"},
	}
}
