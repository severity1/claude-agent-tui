package theme_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/theme"
	_ "github.com/severity1/claude-agent-tui/theme/themes"
)

// TestDefault verifies Default() returns a non-nil palette.
func TestDefault(t *testing.T) {
	p := theme.Default()
	if p == nil {
		t.Fatal("Default() returned nil")
	}
}

// TestGetWithDefault verifies GetWithDefault() returns default for unknown themes.
func TestGetWithDefault(t *testing.T) {
	tests := []struct {
		name       string
		themeName  string
		wantNonNil bool
	}{
		{"known theme", "catppuccin", true},
		{"unknown theme", "nonexistent", true},
		{"empty name", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := theme.GetWithDefault(tt.themeName)
			if tt.wantNonNil && p == nil {
				t.Errorf("GetWithDefault(%q) returned nil, expected non-nil", tt.themeName)
			}
		})
	}
}

// TestGetWithDefault_ReturnsCorrectTheme verifies GetWithDefault returns the
// named theme when it exists, and the default when it doesn't.
func TestGetWithDefault_ReturnsCorrectTheme(t *testing.T) {
	// Known theme should return that theme
	dracula := theme.GetWithDefault("dracula")
	if dracula.Primary.Dark != "#bd93f9" {
		t.Errorf("GetWithDefault(dracula).Primary.Dark = %q, expected #bd93f9", dracula.Primary.Dark)
	}

	// Unknown theme should return default (catppuccin)
	unknown := theme.GetWithDefault("nonexistent")
	catppuccin := theme.Get("catppuccin")
	if unknown.Primary.Dark != catppuccin.Primary.Dark {
		t.Errorf("GetWithDefault(nonexistent) should return default theme")
	}
}

// TestGet verifies Get() returns correct values.
func TestGet(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
	}{
		{"catppuccin", true},
		{"dracula", true},
		{"gruvbox", true},
		{"nord", true},
		{"tokyonight", true},
		{"onedark", true},
		{"system", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := theme.Get(tt.name)
			if tt.exists && p == nil {
				t.Errorf("Get(%q) returned nil, expected palette", tt.name)
			}
			if !tt.exists && p != nil {
				t.Errorf("Get(%q) returned palette, expected nil", tt.name)
			}
		})
	}
}

// TestList verifies List() returns all 7 built-in themes.
func TestList(t *testing.T) {
	themes := theme.List()
	if len(themes) != 7 {
		t.Errorf("List() returned %d themes, expected 7", len(themes))
	}

	expected := map[string]bool{
		"catppuccin": true,
		"dracula":    true,
		"gruvbox":    true,
		"nord":       true,
		"tokyonight": true,
		"onedark":    true,
		"system":     true,
	}

	for _, name := range themes {
		if !expected[name] {
			t.Errorf("unexpected theme in list: %s", name)
		}
		delete(expected, name)
	}

	for name := range expected {
		t.Errorf("missing theme from list: %s", name)
	}
}

// TestListIsSorted verifies List() returns themes in sorted order.
func TestListIsSorted(t *testing.T) {
	themes := theme.List()
	for i := 1; i < len(themes); i++ {
		if themes[i-1] > themes[i] {
			t.Errorf("List() not sorted: %q > %q", themes[i-1], themes[i])
		}
	}
}

// TestRegisterPanicsOnDuplicate verifies Register() panics on duplicate names.
func TestRegisterPanicsOnDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register() did not panic on duplicate")
		}
	}()

	// Try to register a theme that already exists
	theme.Register("catppuccin", theme.Palette{})
}

// TestGetReturnsPointerToStoredPalette verifies Get() returns palette data correctly.
func TestGetReturnsPointerToStoredPalette(t *testing.T) {
	p := theme.Get("catppuccin")
	if p == nil {
		t.Fatal("Get() returned nil")
	}

	// Verify some specific catppuccin colors (Mocha dark variant)
	if p.Primary.Dark != "#89b4fa" {
		t.Errorf("Primary.Dark = %q, expected #89b4fa", p.Primary.Dark)
	}
	if p.Error.Dark != "#f38ba8" {
		t.Errorf("Error.Dark = %q, expected #f38ba8", p.Error.Dark)
	}
}

// TestDifferentThemesHaveDifferentColors verifies themes are distinct.
func TestDifferentThemesHaveDifferentColors(t *testing.T) {
	catppuccin := theme.Get("catppuccin")
	dracula := theme.Get("dracula")

	if catppuccin == nil || dracula == nil {
		t.Fatal("failed to get themes")
	}

	// These themes should have different primary colors
	if catppuccin.Primary.Dark == dracula.Primary.Dark {
		t.Error("catppuccin and dracula have same Primary.Dark color")
	}
}

// TestSystemThemeUsesANSICodes verifies the system theme uses ANSI color codes.
func TestSystemThemeUsesANSICodes(t *testing.T) {
	p := theme.Get("system")
	if p == nil {
		t.Fatal("system theme not found")
	}

	// System theme should use ANSI codes (0-15) not hex colors
	// Primary.Dark should be "12" (bright blue)
	if p.Primary.Dark != "12" {
		t.Errorf("system Primary.Dark = %q, expected ANSI code 12", p.Primary.Dark)
	}
	// Error.Dark should be "9" (bright red)
	if p.Error.Dark != "9" {
		t.Errorf("system Error.Dark = %q, expected ANSI code 9", p.Error.Dark)
	}
}

// TestRegisterCustomTheme verifies custom themes can be registered and retrieved.
func TestRegisterCustomTheme(t *testing.T) {
	customPalette := theme.Palette{
		Primary: lipgloss.AdaptiveColor{Light: "#custom1", Dark: "#custom2"},
		// Note: In real usage, all 45 fields should be populated
	}

	// Register with unique name
	theme.Register("test-custom-theme-unique", customPalette)

	p := theme.Get("test-custom-theme-unique")
	if p == nil {
		t.Fatal("custom theme not found after registration")
	}

	if p.Primary.Light != "#custom1" || p.Primary.Dark != "#custom2" {
		t.Error("custom theme colors not preserved")
	}
}

// TestDefaultName verifies DefaultName() returns the current default theme name.
func TestDefaultName(t *testing.T) {
	name := theme.DefaultName()
	if name != "catppuccin" {
		t.Errorf("DefaultName() = %q, expected %q", name, "catppuccin")
	}

	// Verify DefaultName matches what Default() uses
	defaultPalette := theme.Default()
	namedPalette := theme.Get(name)
	if defaultPalette.Primary.Dark != namedPalette.Primary.Dark {
		t.Error("DefaultName() does not match Default() palette")
	}
}

// TestGetInfo verifies GetInfo() returns correct metadata.
func TestGetInfo(t *testing.T) {
	tests := []struct {
		name     string
		wantNil  bool
		wantDesc string
	}{
		{theme.Catppuccin, false, "Soothing pastel theme (Mocha/Latte)"},
		{theme.Dracula, false, "Classic dark theme with vibrant colors"},
		{theme.Gruvbox, false, "Retro groove with warm earthy tones"},
		{theme.Nord, false, "Arctic north-bluish color palette"},
		{theme.OneDark, false, "Atom editor's iconic dark theme"},
		{theme.System, false, "ANSI 16-color for max compatibility"},
		{theme.TokyoNight, false, "Clean theme inspired by Tokyo city lights"},
		{"nonexistent", true, ""},
		{"", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := theme.GetInfo(tt.name)
			if tt.wantNil {
				if info != nil {
					t.Errorf("GetInfo(%q) returned info, expected nil", tt.name)
				}
				return
			}
			if info == nil {
				t.Fatalf("GetInfo(%q) returned nil, expected info", tt.name)
			}
			if info.Name != tt.name {
				t.Errorf("GetInfo(%q).Name = %q, expected %q", tt.name, info.Name, tt.name)
			}
			if info.Description != tt.wantDesc {
				t.Errorf("GetInfo(%q).Description = %q, expected %q", tt.name, info.Description, tt.wantDesc)
			}
		})
	}
}

// TestListInfo verifies ListInfo() returns metadata for all themes.
func TestListInfo(t *testing.T) {
	infos := theme.ListInfo()

	// Should have at least 7 built-in themes
	if len(infos) < 7 {
		t.Errorf("ListInfo() returned %d themes, expected at least 7", len(infos))
	}

	// Verify sorted order
	for i := 1; i < len(infos); i++ {
		if infos[i-1].Name > infos[i].Name {
			t.Errorf("ListInfo() not sorted: %q > %q", infos[i-1].Name, infos[i].Name)
		}
	}

	// Verify built-in themes have descriptions
	builtInThemes := map[string]bool{
		theme.Catppuccin: true,
		theme.Dracula:    true,
		theme.Gruvbox:    true,
		theme.Nord:       true,
		theme.OneDark:    true,
		theme.System:     true,
		theme.TokyoNight: true,
	}

	for _, info := range infos {
		if builtInThemes[info.Name] && info.Description == "" {
			t.Errorf("ListInfo(): built-in theme %q has empty description", info.Name)
		}
	}
}

// TestListInfoMatchesList verifies ListInfo() and List() return same themes.
func TestListInfoMatchesList(t *testing.T) {
	names := theme.List()
	infos := theme.ListInfo()

	if len(names) != len(infos) {
		t.Fatalf("List() returned %d themes, ListInfo() returned %d", len(names), len(infos))
	}

	for i, name := range names {
		if infos[i].Name != name {
			t.Errorf("Mismatch at index %d: List() has %q, ListInfo() has %q", i, name, infos[i].Name)
		}
	}
}

// TestThemeConstants verifies the exported constants match expected values.
func TestThemeConstants(t *testing.T) {
	tests := []struct {
		constant string
		value    string
	}{
		{theme.Catppuccin, "catppuccin"},
		{theme.Dracula, "dracula"},
		{theme.Gruvbox, "gruvbox"},
		{theme.Nord, "nord"},
		{theme.OneDark, "onedark"},
		{theme.System, "system"},
		{theme.TokyoNight, "tokyonight"},
	}

	for _, tt := range tests {
		if tt.constant != tt.value {
			t.Errorf("constant value = %q, expected %q", tt.constant, tt.value)
		}
		// Verify constant works with Get()
		if theme.Get(tt.constant) == nil {
			t.Errorf("Get(%s) returned nil", tt.constant)
		}
	}
}
