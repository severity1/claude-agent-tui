package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
)

func init() {
	theme.Register("system", System())
}

// System returns the System theme palette.
// Uses ANSI 16-color codes for maximum terminal compatibility.
// This theme works on any terminal that supports basic ANSI colors.
func System() theme.Palette {
	// ANSI color codes (0-15)
	// Basic colors: 0=black, 1=red, 2=green, 3=yellow, 4=blue, 5=magenta, 6=cyan, 7=white
	// Bright colors: 8=bright black, 9=bright red, 10=bright green, 11=bright yellow,
	//                12=bright blue, 13=bright magenta, 14=bright cyan, 15=bright white

	return theme.Palette{
		// Core UI
		Primary:           lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		Secondary:         lipgloss.AdaptiveColor{Light: "5", Dark: "13"}, // magenta
		Accent:            lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		Error:             lipgloss.AdaptiveColor{Light: "1", Dark: "9"},  // red
		Warning:           lipgloss.AdaptiveColor{Light: "3", Dark: "11"}, // yellow
		Success:           lipgloss.AdaptiveColor{Light: "2", Dark: "10"}, // green
		Info:              lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		Text:              lipgloss.AdaptiveColor{Light: "0", Dark: "15"}, // black/bright white
		TextMuted:         lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		Background:        lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // bright white/black
		BackgroundPanel:   lipgloss.AdaptiveColor{Light: "7", Dark: "8"},  // white/bright black
		BackgroundElement: lipgloss.AdaptiveColor{Light: "7", Dark: "8"},  // white/bright black
		Border:            lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		BorderActive:      lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		BorderSubtle:      lipgloss.AdaptiveColor{Light: "7", Dark: "8"},  // white/bright black

		// Diff
		DiffAdded:               lipgloss.AdaptiveColor{Light: "2", Dark: "10"}, // green
		DiffRemoved:             lipgloss.AdaptiveColor{Light: "1", Dark: "9"},  // red
		DiffContext:             lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		DiffHunkHeader:          lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		DiffHighlightAdded:      lipgloss.AdaptiveColor{Light: "2", Dark: "10"}, // green
		DiffHighlightRemoved:    lipgloss.AdaptiveColor{Light: "1", Dark: "9"},  // red
		DiffAddedBg:             lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // background
		DiffRemovedBg:           lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // background
		DiffContextBg:           lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // background
		DiffLineNumber:          lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		DiffAddedLineNumberBg:   lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // background
		DiffRemovedLineNumberBg: lipgloss.AdaptiveColor{Light: "15", Dark: "0"}, // background

		// Markdown
		MarkdownText:            lipgloss.AdaptiveColor{Light: "0", Dark: "15"}, // text
		MarkdownHeading:         lipgloss.AdaptiveColor{Light: "5", Dark: "13"}, // magenta
		MarkdownLink:            lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		MarkdownLinkText:        lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		MarkdownCode:            lipgloss.AdaptiveColor{Light: "2", Dark: "10"}, // green
		MarkdownBlockQuote:      lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		MarkdownEmph:            lipgloss.AdaptiveColor{Light: "3", Dark: "11"}, // yellow
		MarkdownStrong:          lipgloss.AdaptiveColor{Light: "3", Dark: "11"}, // yellow
		MarkdownHorizontalRule:  lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		MarkdownListItem:        lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		MarkdownListEnumeration: lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		MarkdownImage:           lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		MarkdownImageText:       lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		MarkdownCodeBlock:       lipgloss.AdaptiveColor{Light: "0", Dark: "15"}, // text

		// Syntax
		SyntaxComment:     lipgloss.AdaptiveColor{Light: "8", Dark: "7"},  // bright black/white
		SyntaxKeyword:     lipgloss.AdaptiveColor{Light: "5", Dark: "13"}, // magenta
		SyntaxFunction:    lipgloss.AdaptiveColor{Light: "4", Dark: "12"}, // blue
		SyntaxVariable:    lipgloss.AdaptiveColor{Light: "1", Dark: "9"},  // red
		SyntaxString:      lipgloss.AdaptiveColor{Light: "2", Dark: "10"}, // green
		SyntaxNumber:      lipgloss.AdaptiveColor{Light: "5", Dark: "13"}, // magenta
		SyntaxType:        lipgloss.AdaptiveColor{Light: "6", Dark: "14"}, // cyan
		SyntaxOperator:    lipgloss.AdaptiveColor{Light: "3", Dark: "11"}, // yellow
		SyntaxPunctuation: lipgloss.AdaptiveColor{Light: "0", Dark: "15"}, // text
	}
}
