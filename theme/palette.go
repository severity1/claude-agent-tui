// Package theme provides a centralized theming system for claude-agent-tui components.
//
// The package defines a Palette struct with semantic color definitions that work
// across all TUI components. All colors use lipgloss.AdaptiveColor to automatically
// adjust for light/dark terminals.
//
// # Usage
//
// Use the default theme:
//
//	p := theme.Default()
//	input := chatinput.New(chatinput.WithPalette(p))
//
// Use a specific theme:
//
//	p := theme.Get("dracula")
//	input := chatinput.New(chatinput.WithPalette(p))
//
// Register a custom theme:
//
//	theme.Register("custom", theme.Palette{
//	    Primary: lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"},
//	    // ... other fields
//	})
package theme

import "github.com/charmbracelet/lipgloss"

// Palette defines the complete color scheme for a theme.
// All fields use AdaptiveColor to support both light and dark terminals.
//
// Fields are organized into categories:
//   - Core UI (15): Primary colors for UI elements
//   - Diff (12): Colors for diff/patch display
//   - Markdown (14): Colors for markdown rendering
//   - Syntax (9): Colors for syntax highlighting
type Palette struct {
	// Core UI colors (15)

	// Primary is the main accent color for focused elements and highlights.
	Primary lipgloss.AdaptiveColor
	// Secondary is used for subtle accents and interactive elements.
	Secondary lipgloss.AdaptiveColor
	// Accent is used for tertiary highlights and special elements.
	Accent lipgloss.AdaptiveColor
	// Error indicates error states and destructive actions.
	Error lipgloss.AdaptiveColor
	// Warning indicates warning states and cautionary elements.
	Warning lipgloss.AdaptiveColor
	// Success indicates successful operations and positive states.
	Success lipgloss.AdaptiveColor
	// Info indicates informational messages and neutral highlights.
	Info lipgloss.AdaptiveColor
	// Text is the primary text color.
	Text lipgloss.AdaptiveColor
	// TextMuted is used for secondary text, placeholders, and hints.
	TextMuted lipgloss.AdaptiveColor
	// Background is the main background color.
	Background lipgloss.AdaptiveColor
	// BackgroundPanel is for panel/container backgrounds.
	BackgroundPanel lipgloss.AdaptiveColor
	// BackgroundElement is for nested element backgrounds.
	BackgroundElement lipgloss.AdaptiveColor
	// Border is the default border color.
	Border lipgloss.AdaptiveColor
	// BorderActive is the border color for focused/active elements.
	BorderActive lipgloss.AdaptiveColor
	// BorderSubtle is for subtle/de-emphasized borders.
	BorderSubtle lipgloss.AdaptiveColor

	// Diff colors (12)

	// DiffAdded is the text color for added lines.
	DiffAdded lipgloss.AdaptiveColor
	// DiffRemoved is the text color for removed lines.
	DiffRemoved lipgloss.AdaptiveColor
	// DiffContext is the text color for context lines.
	DiffContext lipgloss.AdaptiveColor
	// DiffHunkHeader is the color for hunk headers (@@ ... @@).
	DiffHunkHeader lipgloss.AdaptiveColor
	// DiffHighlightAdded is for highlighting specific added text within a line.
	DiffHighlightAdded lipgloss.AdaptiveColor
	// DiffHighlightRemoved is for highlighting specific removed text within a line.
	DiffHighlightRemoved lipgloss.AdaptiveColor
	// DiffAddedBg is the background color for added lines.
	DiffAddedBg lipgloss.AdaptiveColor
	// DiffRemovedBg is the background color for removed lines.
	DiffRemovedBg lipgloss.AdaptiveColor
	// DiffContextBg is the background color for context lines.
	DiffContextBg lipgloss.AdaptiveColor
	// DiffLineNumber is the color for line numbers in diffs.
	DiffLineNumber lipgloss.AdaptiveColor
	// DiffAddedLineNumberBg is the background for added line numbers.
	DiffAddedLineNumberBg lipgloss.AdaptiveColor
	// DiffRemovedLineNumberBg is the background for removed line numbers.
	DiffRemovedLineNumberBg lipgloss.AdaptiveColor

	// Markdown colors (14)

	// MarkdownText is the default text color in markdown.
	MarkdownText lipgloss.AdaptiveColor
	// MarkdownHeading is the color for headings (# ## ###).
	MarkdownHeading lipgloss.AdaptiveColor
	// MarkdownLink is the color for link URLs.
	MarkdownLink lipgloss.AdaptiveColor
	// MarkdownLinkText is the color for link display text.
	MarkdownLinkText lipgloss.AdaptiveColor
	// MarkdownCode is the color for inline code.
	MarkdownCode lipgloss.AdaptiveColor
	// MarkdownBlockQuote is the color for block quotes.
	MarkdownBlockQuote lipgloss.AdaptiveColor
	// MarkdownEmph is the color for emphasized (italic) text.
	MarkdownEmph lipgloss.AdaptiveColor
	// MarkdownStrong is the color for strong (bold) text.
	MarkdownStrong lipgloss.AdaptiveColor
	// MarkdownHorizontalRule is the color for horizontal rules.
	MarkdownHorizontalRule lipgloss.AdaptiveColor
	// MarkdownListItem is the color for list item markers.
	MarkdownListItem lipgloss.AdaptiveColor
	// MarkdownListEnumeration is the color for numbered list markers.
	MarkdownListEnumeration lipgloss.AdaptiveColor
	// MarkdownImage is the color for image syntax.
	MarkdownImage lipgloss.AdaptiveColor
	// MarkdownImageText is the color for image alt text.
	MarkdownImageText lipgloss.AdaptiveColor
	// MarkdownCodeBlock is the color for code block text.
	MarkdownCodeBlock lipgloss.AdaptiveColor

	// Syntax highlighting colors (9)

	// SyntaxComment is the color for code comments.
	SyntaxComment lipgloss.AdaptiveColor
	// SyntaxKeyword is the color for language keywords.
	SyntaxKeyword lipgloss.AdaptiveColor
	// SyntaxFunction is the color for function names.
	SyntaxFunction lipgloss.AdaptiveColor
	// SyntaxVariable is the color for variable names.
	SyntaxVariable lipgloss.AdaptiveColor
	// SyntaxString is the color for string literals.
	SyntaxString lipgloss.AdaptiveColor
	// SyntaxNumber is the color for numeric literals.
	SyntaxNumber lipgloss.AdaptiveColor
	// SyntaxType is the color for type names.
	SyntaxType lipgloss.AdaptiveColor
	// SyntaxOperator is the color for operators.
	SyntaxOperator lipgloss.AdaptiveColor
	// SyntaxPunctuation is the color for punctuation.
	SyntaxPunctuation lipgloss.AdaptiveColor
}
