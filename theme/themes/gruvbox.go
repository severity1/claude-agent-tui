package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("gruvbox", Gruvbox())
}

// Gruvbox returns the Gruvbox theme palette.
// A retro groove color scheme with warm earthy tones.
func Gruvbox() theme.Palette {
	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "#076678", Dark: "#83a598"},
		Secondary:         lipgloss.AdaptiveColor{Light: "#8f3f71", Dark: "#d3869b"},
		Accent:            lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#8ec07c"},
		Error:             lipgloss.AdaptiveColor{Light: "#9d0006", Dark: "#fb4934"},
		Warning:           lipgloss.AdaptiveColor{Light: "#af3a03", Dark: "#fe8019"},
		Success:           lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"},
		Info:              lipgloss.AdaptiveColor{Light: "#b57614", Dark: "#fabd2f"},
		Text:              lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#ebdbb2"},
		TextMuted:         lipgloss.AdaptiveColor{Light: "#7c6f64", Dark: "#928374"},
		Background:        lipgloss.AdaptiveColor{Light: "#fbf1c7", Dark: "#282828"},
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "#ebdbb2", Dark: "#3c3836"},
		BackgroundElement: lipgloss.AdaptiveColor{Light: "#d5c4a1", Dark: "#504945"},
		Border:            lipgloss.AdaptiveColor{Light: "#bdae93", Dark: "#665c54"},
		BorderActive:      lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#ebdbb2"},
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "#d5c4a1", Dark: "#504945"},

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#98971a"},
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "#9d0006", Dark: "#cc241d"},
		DiffContext:             lipgloss.AdaptiveColor{Light: "#7c6f64", Dark: "#928374"},
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#689d6a"},
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"},
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "#9d0006", Dark: "#fb4934"},
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "#dcd8a4", Dark: "#32302f"},
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "#e2c7c3", Dark: "#322929"},
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "#ebdbb2", Dark: "#3c3836"},
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "#bdae93", Dark: "#665c54"},
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "#cec99e", Dark: "#2a2827"},
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "#d3bdb9", Dark: "#2a2222"},

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#ebdbb2"},
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "#076678", Dark: "#83a598"},
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#8ec07c"},
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"},
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "#b57614", Dark: "#fabd2f"},
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "#7c6f64", Dark: "#928374"},
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "#8f3f71", Dark: "#d3869b"},
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "#af3a03", Dark: "#fe8019"},
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "#7c6f64", Dark: "#928374"},
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "#076678", Dark: "#83a598"},
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#8ec07c"},
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#8ec07c"},
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"},
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#ebdbb2"},

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "#7c6f64", Dark: "#928374"},
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "#9d0006", Dark: "#fb4934"},
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"},
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "#076678", Dark: "#83a598"},
		SyntaxString:      lipgloss.AdaptiveColor{Light: "#b57614", Dark: "#fabd2f"},
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "#8f3f71", Dark: "#d3869b"},
		SyntaxType:        lipgloss.AdaptiveColor{Light: "#427b58", Dark: "#8ec07c"},
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "#af3a03", Dark: "#fe8019"},
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#ebdbb2"},
	}
}
