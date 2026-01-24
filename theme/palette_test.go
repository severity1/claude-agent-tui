package theme_test

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
	_ "github.com/severity1/claude-agent-tui/theme/themes"
)

// TestPaletteFieldCount verifies the Palette struct has exactly 50 fields.
// Categories: Core UI (15) + Diff (12) + Markdown (14) + Syntax (9) = 50
func TestPaletteFieldCount(t *testing.T) {
	p := theme.Palette{}
	v := reflect.ValueOf(p)
	if v.NumField() != 50 {
		t.Errorf("Palette should have 50 fields, got %d", v.NumField())
	}
}

// TestAllThemesHaveNonEmptyColors verifies all registered themes have non-empty color values.
func TestAllThemesHaveNonEmptyColors(t *testing.T) {
	themes := theme.List()
	if len(themes) == 0 {
		t.Fatal("No themes registered")
	}

	for _, name := range themes {
		t.Run(name, func(t *testing.T) {
			p := theme.Get(name)
			if p == nil {
				t.Fatalf("theme.Get(%q) returned nil", name)
			}

			v := reflect.ValueOf(*p)
			typ := v.Type()

			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				fieldName := typ.Field(i).Name

				// Each field should be AdaptiveColor
				ac, ok := field.Interface().(lipgloss.AdaptiveColor)
				if !ok {
					t.Errorf("field %s is not AdaptiveColor", fieldName)
					continue
				}

				// Both Light and Dark should be non-empty
				if ac.Light == "" {
					t.Errorf("field %s has empty Light value", fieldName)
				}
				if ac.Dark == "" {
					t.Errorf("field %s has empty Dark value", fieldName)
				}
			}
		})
	}
}

// TestPaletteFieldCategories verifies the expected fields exist in each category.
func TestPaletteFieldCategories(t *testing.T) {
	coreUI := []string{
		"Primary", "Secondary", "Accent", "Error", "Warning", "Success", "Info",
		"Text", "TextMuted", "Background", "BackgroundPanel", "BackgroundElement",
		"Border", "BorderActive", "BorderSubtle",
	}
	diff := []string{
		"DiffAdded", "DiffRemoved", "DiffContext", "DiffHunkHeader",
		"DiffHighlightAdded", "DiffHighlightRemoved",
		"DiffAddedBg", "DiffRemovedBg", "DiffContextBg",
		"DiffLineNumber", "DiffAddedLineNumberBg", "DiffRemovedLineNumberBg",
	}
	markdown := []string{
		"MarkdownText", "MarkdownHeading", "MarkdownLink", "MarkdownLinkText",
		"MarkdownCode", "MarkdownBlockQuote", "MarkdownEmph", "MarkdownStrong",
		"MarkdownHorizontalRule", "MarkdownListItem", "MarkdownListEnumeration",
		"MarkdownImage", "MarkdownImageText", "MarkdownCodeBlock",
	}
	syntax := []string{
		"SyntaxComment", "SyntaxKeyword", "SyntaxFunction", "SyntaxVariable",
		"SyntaxString", "SyntaxNumber", "SyntaxType", "SyntaxOperator", "SyntaxPunctuation",
	}

	allFields := make(map[string]bool)
	v := reflect.TypeOf(theme.Palette{})
	for i := 0; i < v.NumField(); i++ {
		allFields[v.Field(i).Name] = true
	}

	checkFields := func(category string, fields []string) {
		for _, f := range fields {
			if !allFields[f] {
				t.Errorf("category %s: missing field %s", category, f)
			}
		}
	}

	checkFields("Core UI", coreUI)
	checkFields("Diff", diff)
	checkFields("Markdown", markdown)
	checkFields("Syntax", syntax)

	// Verify counts: 15 + 12 + 14 + 9 = 50 total
	if len(coreUI) != 15 {
		t.Errorf("Core UI should have 15 fields, got %d", len(coreUI))
	}
	if len(diff) != 12 {
		t.Errorf("Diff should have 12 fields, got %d", len(diff))
	}
	if len(markdown) != 14 {
		t.Errorf("Markdown should have 14 fields, got %d", len(markdown))
	}
	if len(syntax) != 9 {
		t.Errorf("Syntax should have 9 fields, got %d", len(syntax))
	}

	total := len(coreUI) + len(diff) + len(markdown) + len(syntax)
	if total != 50 {
		t.Errorf("Total fields should be 50, got %d", total)
	}
}
