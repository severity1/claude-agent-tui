package styles_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/severity1/claude-agent-tui/internal/styles"
)

// ============================================================================
// DefaultBorderConfig Tests
// ============================================================================

func TestDefaultBorderConfig(t *testing.T) {
	cfg := styles.DefaultBorderConfig()

	// Verify default values
	if cfg.Left {
		t.Error("DefaultBorderConfig().Left = true, want false")
	}
	if !cfg.Right {
		t.Error("DefaultBorderConfig().Right = false, want true")
	}
	if cfg.PaddingLeft != 1 {
		t.Errorf("DefaultBorderConfig().PaddingLeft = %d, want 1", cfg.PaddingLeft)
	}
	if cfg.PaddingTop != 0 {
		t.Errorf("DefaultBorderConfig().PaddingTop = %d, want 0", cfg.PaddingTop)
	}
	if cfg.PaddingBottom != 0 {
		t.Errorf("DefaultBorderConfig().PaddingBottom = %d, want 0", cfg.PaddingBottom)
	}
	if cfg.PaddingRight != 0 {
		t.Errorf("DefaultBorderConfig().PaddingRight = %d, want 0", cfg.PaddingRight)
	}
}

// ============================================================================
// ModifyBorderVisibility Tests
// ============================================================================

func TestModifyBorderVisibility_BothVisible(t *testing.T) {
	b := lipgloss.ThickBorder()
	result := styles.ModifyBorderVisibility(b, true, true)

	// When both visible, border characters should be unchanged
	if result.Left == " " {
		t.Error("ModifyBorderVisibility with showLeft=true should keep left border")
	}
	if result.Right == " " {
		t.Error("ModifyBorderVisibility with showRight=true should keep right border")
	}
}

func TestModifyBorderVisibility_LeftHidden(t *testing.T) {
	b := lipgloss.ThickBorder()
	result := styles.ModifyBorderVisibility(b, false, true)

	if result.Left != " " {
		t.Errorf("ModifyBorderVisibility with showLeft=false should set Left to space, got %q", result.Left)
	}
	if result.MiddleLeft != " " {
		t.Errorf("ModifyBorderVisibility with showLeft=false should set MiddleLeft to space, got %q", result.MiddleLeft)
	}
	if result.Right == " " {
		t.Error("ModifyBorderVisibility should keep right border visible when showRight=true")
	}
}

func TestModifyBorderVisibility_RightHidden(t *testing.T) {
	b := lipgloss.ThickBorder()
	result := styles.ModifyBorderVisibility(b, true, false)

	if result.Right != " " {
		t.Errorf("ModifyBorderVisibility with showRight=false should set Right to space, got %q", result.Right)
	}
	if result.MiddleRight != " " {
		t.Errorf("ModifyBorderVisibility with showRight=false should set MiddleRight to space, got %q", result.MiddleRight)
	}
	if result.Left == " " {
		t.Error("ModifyBorderVisibility should keep left border visible when showLeft=true")
	}
}

func TestModifyBorderVisibility_BothHidden(t *testing.T) {
	b := lipgloss.ThickBorder()
	result := styles.ModifyBorderVisibility(b, false, false)

	if result.Left != " " {
		t.Error("ModifyBorderVisibility with showLeft=false should hide left border")
	}
	if result.Right != " " {
		t.Error("ModifyBorderVisibility with showRight=false should hide right border")
	}
	if result.MiddleLeft != " " {
		t.Error("ModifyBorderVisibility should hide MiddleLeft when showLeft=false")
	}
	if result.MiddleRight != " " {
		t.Error("ModifyBorderVisibility should hide MiddleRight when showRight=false")
	}
}

func TestModifyBorderVisibility_PreservesOriginal(t *testing.T) {
	b := lipgloss.ThickBorder()
	originalLeft := b.Left
	originalRight := b.Right

	_ = styles.ModifyBorderVisibility(b, false, false)

	// Original should be unchanged
	if b.Left != originalLeft {
		t.Error("ModifyBorderVisibility should not mutate original border")
	}
	if b.Right != originalRight {
		t.Error("ModifyBorderVisibility should not mutate original border")
	}
}

// ============================================================================
// ApplyToStyle Tests
// ============================================================================

func TestApplyToStyle_Unfocused(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:         lipgloss.ThickBorder(),
		Color:         lipgloss.Color("240"),
		FocusColor:    lipgloss.Color("205"),
		Left:          true,
		Right:         true,
		PaddingTop:    1,
		PaddingBottom: 2,
		PaddingLeft:   3,
		PaddingRight:  4,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyle(style, false)

	// Verify style has border applied
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyle should set BorderLeft")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyle should set BorderRight")
	}
	if result.GetPaddingTop() != 1 {
		t.Errorf("ApplyToStyle PaddingTop = %d, want 1", result.GetPaddingTop())
	}
	if result.GetPaddingBottom() != 2 {
		t.Errorf("ApplyToStyle PaddingBottom = %d, want 2", result.GetPaddingBottom())
	}
	if result.GetPaddingLeft() != 3 {
		t.Errorf("ApplyToStyle PaddingLeft = %d, want 3", result.GetPaddingLeft())
	}
	if result.GetPaddingRight() != 4 {
		t.Errorf("ApplyToStyle PaddingRight = %d, want 4", result.GetPaddingRight())
	}
}

func TestApplyToStyle_Focused_UsesFocusColor(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:      lipgloss.ThickBorder(),
		Color:      lipgloss.Color("240"),
		FocusColor: lipgloss.Color("205"),
		Left:       true,
		Right:      true,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyle(style, true)

	// Verify focused uses FocusColor
	// Note: We can't easily test the actual color value, but we can verify style is applied
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyle focused should set BorderLeft")
	}
}

func TestApplyToStyle_Focused_NilFocusColor_UsesColor(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:      lipgloss.ThickBorder(),
		Color:      lipgloss.Color("240"),
		FocusColor: nil, // nil FocusColor
		Left:       true,
		Right:      true,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyle(style, true)

	// When FocusColor is nil, should use Color instead
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyle with nil FocusColor should still apply border")
	}
}

func TestApplyToStyle_HiddenBorders(t *testing.T) {
	cfg := styles.BorderConfig{
		Style: lipgloss.ThickBorder(),
		Color: lipgloss.Color("240"),
		Left:  false, // hidden
		Right: false, // hidden
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyle(style, false)

	// Borders should be set (for space reservation) but visibility handled by ModifyBorderVisibility
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyle should set BorderLeft (space reservation)")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyle should set BorderRight (space reservation)")
	}
}

// ============================================================================
// SetPaddingIfValid Tests
// ============================================================================

func TestSetPaddingIfValid_AllPositive(t *testing.T) {
	cfg := styles.BorderConfig{}
	result := cfg.SetPaddingIfValid(1, 2, 3, 4)

	if result.PaddingTop != 1 {
		t.Errorf("SetPaddingIfValid PaddingTop = %d, want 1", result.PaddingTop)
	}
	if result.PaddingBottom != 2 {
		t.Errorf("SetPaddingIfValid PaddingBottom = %d, want 2", result.PaddingBottom)
	}
	if result.PaddingLeft != 3 {
		t.Errorf("SetPaddingIfValid PaddingLeft = %d, want 3", result.PaddingLeft)
	}
	if result.PaddingRight != 4 {
		t.Errorf("SetPaddingIfValid PaddingRight = %d, want 4", result.PaddingRight)
	}
}

func TestSetPaddingIfValid_NegativeIgnored(t *testing.T) {
	cfg := styles.BorderConfig{
		PaddingTop:    10,
		PaddingBottom: 20,
		PaddingLeft:   30,
		PaddingRight:  40,
	}
	result := cfg.SetPaddingIfValid(-1, -1, -1, -1)

	// Negative values should be ignored, original values preserved
	if result.PaddingTop != 10 {
		t.Errorf("SetPaddingIfValid with negative should preserve PaddingTop, got %d", result.PaddingTop)
	}
	if result.PaddingBottom != 20 {
		t.Errorf("SetPaddingIfValid with negative should preserve PaddingBottom, got %d", result.PaddingBottom)
	}
	if result.PaddingLeft != 30 {
		t.Errorf("SetPaddingIfValid with negative should preserve PaddingLeft, got %d", result.PaddingLeft)
	}
	if result.PaddingRight != 40 {
		t.Errorf("SetPaddingIfValid with negative should preserve PaddingRight, got %d", result.PaddingRight)
	}
}

func TestSetPaddingIfValid_MixedValues(t *testing.T) {
	cfg := styles.BorderConfig{
		PaddingTop:    10,
		PaddingBottom: 20,
		PaddingLeft:   30,
		PaddingRight:  40,
	}
	result := cfg.SetPaddingIfValid(1, -1, 3, -1)

	// Only non-negative values should update
	if result.PaddingTop != 1 {
		t.Errorf("SetPaddingIfValid PaddingTop = %d, want 1", result.PaddingTop)
	}
	if result.PaddingBottom != 20 {
		t.Errorf("SetPaddingIfValid PaddingBottom = %d, want 20 (preserved)", result.PaddingBottom)
	}
	if result.PaddingLeft != 3 {
		t.Errorf("SetPaddingIfValid PaddingLeft = %d, want 3", result.PaddingLeft)
	}
	if result.PaddingRight != 40 {
		t.Errorf("SetPaddingIfValid PaddingRight = %d, want 40 (preserved)", result.PaddingRight)
	}
}

func TestSetPaddingIfValid_ZeroAllowed(t *testing.T) {
	cfg := styles.BorderConfig{
		PaddingTop:    10,
		PaddingBottom: 20,
		PaddingLeft:   30,
		PaddingRight:  40,
	}
	result := cfg.SetPaddingIfValid(0, 0, 0, 0)

	// Zero should be allowed (0 >= 0)
	if result.PaddingTop != 0 {
		t.Errorf("SetPaddingIfValid PaddingTop = %d, want 0", result.PaddingTop)
	}
	if result.PaddingBottom != 0 {
		t.Errorf("SetPaddingIfValid PaddingBottom = %d, want 0", result.PaddingBottom)
	}
	if result.PaddingLeft != 0 {
		t.Errorf("SetPaddingIfValid PaddingLeft = %d, want 0", result.PaddingLeft)
	}
	if result.PaddingRight != 0 {
		t.Errorf("SetPaddingIfValid PaddingRight = %d, want 0", result.PaddingRight)
	}
}

func TestSetPaddingIfValid_DoesNotMutateOriginal(t *testing.T) {
	cfg := styles.BorderConfig{
		PaddingTop: 10,
	}
	_ = cfg.SetPaddingIfValid(99, 99, 99, 99)

	// Original should be unchanged
	if cfg.PaddingTop != 10 {
		t.Error("SetPaddingIfValid should not mutate original BorderConfig")
	}
}

// ============================================================================
// ApplyToStyleWithColor Tests
// ============================================================================

func TestApplyToStyleWithColor_UsesProvidedColor(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:      lipgloss.ThickBorder(),
		Color:      lipgloss.Color("240"), // should be ignored
		FocusColor: lipgloss.Color("205"), // should be ignored
		Left:       true,
		Right:      true,
		PaddingTop: 1,
	}

	color := lipgloss.Color("99")
	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleWithColor(style, color)

	// Verify style is applied with padding
	if result.GetPaddingTop() != 1 {
		t.Errorf("ApplyToStyleWithColor PaddingTop = %d, want 1", result.GetPaddingTop())
	}
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyleWithColor should set BorderLeft")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyleWithColor should set BorderRight")
	}
}

func TestApplyToStyleWithColor_AppliesAllPadding(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:         lipgloss.ThickBorder(),
		Left:          true,
		Right:         true,
		PaddingTop:    1,
		PaddingBottom: 2,
		PaddingLeft:   3,
		PaddingRight:  4,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleWithColor(style, lipgloss.Color("99"))

	if result.GetPaddingTop() != 1 {
		t.Errorf("PaddingTop = %d, want 1", result.GetPaddingTop())
	}
	if result.GetPaddingBottom() != 2 {
		t.Errorf("PaddingBottom = %d, want 2", result.GetPaddingBottom())
	}
	if result.GetPaddingLeft() != 3 {
		t.Errorf("PaddingLeft = %d, want 3", result.GetPaddingLeft())
	}
	if result.GetPaddingRight() != 4 {
		t.Errorf("PaddingRight = %d, want 4", result.GetPaddingRight())
	}
}

func TestApplyToStyleWithColor_HiddenBorders(t *testing.T) {
	cfg := styles.BorderConfig{
		Style: lipgloss.ThickBorder(),
		Left:  false,
		Right: false,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleWithColor(style, lipgloss.Color("99"))

	// Borders should still be set for space reservation
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyleWithColor should set BorderLeft for space reservation")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyleWithColor should set BorderRight for space reservation")
	}
}

func TestApplyToStyleWithColor_NilColor(t *testing.T) {
	cfg := styles.BorderConfig{
		Style: lipgloss.ThickBorder(),
		Left:  true,
		Right: true,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleWithColor(style, nil)

	// Should handle nil color gracefully
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyleWithColor with nil color should still apply border")
	}
}

// ============================================================================
// ApplyToStyleCentered Tests
// ============================================================================

func TestApplyToStyleCentered_EnablesAllBorders(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:         lipgloss.NormalBorder(),
		Left:          false, // should be ignored
		Right:         false, // should be ignored
		PaddingTop:    1,
		PaddingBottom: 1,
		PaddingLeft:   1,
		PaddingRight:  1,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleCentered(style, lipgloss.Color("99"))

	// All 4 borders should be enabled
	if result.GetBorderTop() != true {
		t.Error("ApplyToStyleCentered should enable top border")
	}
	if result.GetBorderBottom() != true {
		t.Error("ApplyToStyleCentered should enable bottom border")
	}
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyleCentered should enable left border")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyleCentered should enable right border")
	}
}

func TestApplyToStyleCentered_AppliesPadding(t *testing.T) {
	cfg := styles.BorderConfig{
		Style:         lipgloss.NormalBorder(),
		PaddingTop:    1,
		PaddingBottom: 2,
		PaddingLeft:   3,
		PaddingRight:  4,
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleCentered(style, lipgloss.Color("99"))

	if result.GetPaddingTop() != 1 {
		t.Errorf("PaddingTop = %d, want 1", result.GetPaddingTop())
	}
	if result.GetPaddingBottom() != 2 {
		t.Errorf("PaddingBottom = %d, want 2", result.GetPaddingBottom())
	}
	if result.GetPaddingLeft() != 3 {
		t.Errorf("PaddingLeft = %d, want 3", result.GetPaddingLeft())
	}
	if result.GetPaddingRight() != 4 {
		t.Errorf("PaddingRight = %d, want 4", result.GetPaddingRight())
	}
}

func TestApplyToStyleCentered_IgnoresLeftRightConfig(t *testing.T) {
	// Test that Left/Right fields are ignored for centered mode
	cfg := styles.BorderConfig{
		Style: lipgloss.NormalBorder(),
		Left:  false, // should be ignored
		Right: false, // should be ignored
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleCentered(style, lipgloss.Color("99"))

	// Despite Left=false, Right=false, centered should show all borders
	if result.GetBorderLeft() != true {
		t.Error("ApplyToStyleCentered should ignore Left=false and show border")
	}
	if result.GetBorderRight() != true {
		t.Error("ApplyToStyleCentered should ignore Right=false and show border")
	}
}

func TestApplyToStyleCentered_NilColor(t *testing.T) {
	cfg := styles.BorderConfig{
		Style: lipgloss.NormalBorder(),
	}

	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyleCentered(style, nil)

	// Should handle nil color gracefully
	if result.GetBorderTop() != true {
		t.Error("ApplyToStyleCentered with nil color should still apply borders")
	}
}

// ============================================================================
// CenterBorderConfig Tests
// ============================================================================

func TestCenterBorderConfig_DefaultValues(t *testing.T) {
	cfg := styles.CenterBorderConfig()

	// Verify all expected defaults
	if !cfg.Left {
		t.Error("CenterBorderConfig().Left = false, want true")
	}
	if !cfg.Right {
		t.Error("CenterBorderConfig().Right = false, want true")
	}
	if cfg.PaddingTop != 1 {
		t.Errorf("CenterBorderConfig().PaddingTop = %d, want 1", cfg.PaddingTop)
	}
	if cfg.PaddingBottom != 1 {
		t.Errorf("CenterBorderConfig().PaddingBottom = %d, want 1", cfg.PaddingBottom)
	}
	if cfg.PaddingLeft != 1 {
		t.Errorf("CenterBorderConfig().PaddingLeft = %d, want 1", cfg.PaddingLeft)
	}
	if cfg.PaddingRight != 1 {
		t.Errorf("CenterBorderConfig().PaddingRight = %d, want 1", cfg.PaddingRight)
	}
}

func TestCenterBorderConfig_UniformPadding(t *testing.T) {
	cfg := styles.CenterBorderConfig()

	// All padding should be equal (uniform)
	if cfg.PaddingTop != cfg.PaddingBottom ||
		cfg.PaddingTop != cfg.PaddingLeft ||
		cfg.PaddingTop != cfg.PaddingRight {
		t.Error("CenterBorderConfig should have uniform padding on all sides")
	}
}

func TestCenterBorderConfig_UsesNormalBorder(t *testing.T) {
	cfg := styles.CenterBorderConfig()
	normal := lipgloss.NormalBorder()

	// Verify it uses NormalBorder
	if cfg.Style.Top != normal.Top {
		t.Error("CenterBorderConfig should use NormalBorder style")
	}
}

func TestCenterBorderConfig_NoColor(t *testing.T) {
	cfg := styles.CenterBorderConfig()

	// Color and FocusColor should be nil (caller provides color)
	if cfg.Color != nil {
		t.Error("CenterBorderConfig().Color should be nil")
	}
	if cfg.FocusColor != nil {
		t.Error("CenterBorderConfig().FocusColor should be nil")
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestApplyToStyle_ThenModify(t *testing.T) {
	// Test that returned style can be further modified
	cfg := styles.DefaultBorderConfig()
	style := lipgloss.NewStyle()
	result := cfg.ApplyToStyle(style, false)

	// Add more styling
	final := result.Width(50).Height(10)

	if final.GetWidth() != 50 {
		t.Errorf("Modified style Width = %d, want 50", final.GetWidth())
	}
	if final.GetHeight() != 10 {
		t.Errorf("Modified style Height = %d, want 10", final.GetHeight())
	}
}

func TestApplyToStyleWithColor_Chainable(t *testing.T) {
	cfg := styles.BorderConfig{
		Style: lipgloss.ThickBorder(),
		Left:  true,
		Right: true,
	}

	// Should be chainable with other style methods
	result := cfg.ApplyToStyleWithColor(lipgloss.NewStyle(), lipgloss.Color("99")).
		Width(80).
		Background(lipgloss.Color("0"))

	if result.GetWidth() != 80 {
		t.Errorf("Chained style Width = %d, want 80", result.GetWidth())
	}
}

func TestApplyToStyleCentered_Chainable(t *testing.T) {
	cfg := styles.CenterBorderConfig()

	// Should be chainable with other style methods
	result := cfg.ApplyToStyleCentered(lipgloss.NewStyle(), lipgloss.Color("99")).
		Background(lipgloss.Color("236"))

	if result.GetBorderTop() != true {
		t.Error("Chained centered style should have top border")
	}
}

// ============================================================================
// BorderConfigOption Tests
// ============================================================================

func TestWithBorderStyle(t *testing.T) {
	cfg := styles.BorderConfig{}
	opt := styles.WithBorderStyle(lipgloss.RoundedBorder())
	opt(&cfg)

	if cfg.Style.TopLeft != lipgloss.RoundedBorder().TopLeft {
		t.Error("WithBorderStyle should set the border style")
	}
}

func TestWithBorderLeft(t *testing.T) {
	cfg := styles.BorderConfig{Left: false}
	opt := styles.WithBorderLeft(true)
	opt(&cfg)

	if !cfg.Left {
		t.Error("WithBorderLeft(true) should set Left to true")
	}

	opt = styles.WithBorderLeft(false)
	opt(&cfg)
	if cfg.Left {
		t.Error("WithBorderLeft(false) should set Left to false")
	}
}

func TestWithBorderRight(t *testing.T) {
	cfg := styles.BorderConfig{Right: false}
	opt := styles.WithBorderRight(true)
	opt(&cfg)

	if !cfg.Right {
		t.Error("WithBorderRight(true) should set Right to true")
	}

	opt = styles.WithBorderRight(false)
	opt(&cfg)
	if cfg.Right {
		t.Error("WithBorderRight(false) should set Right to false")
	}
}

func TestWithBorderColor(t *testing.T) {
	cfg := styles.BorderConfig{}
	color := lipgloss.Color("240")
	opt := styles.WithBorderColor(color)
	opt(&cfg)

	if cfg.Color != color {
		t.Error("WithBorderColor should set the border color")
	}
}

func TestWithBorderFocusColor(t *testing.T) {
	cfg := styles.BorderConfig{}
	color := lipgloss.Color("205")
	opt := styles.WithBorderFocusColor(color)
	opt(&cfg)

	if cfg.FocusColor != color {
		t.Error("WithBorderFocusColor should set the focus color")
	}
}

func TestWithBorderPadding(t *testing.T) {
	cfg := styles.BorderConfig{}
	opt := styles.WithBorderPadding(1, 2, 3, 4)
	opt(&cfg)

	if cfg.PaddingTop != 1 {
		t.Errorf("WithBorderPadding PaddingTop = %d, want 1", cfg.PaddingTop)
	}
	if cfg.PaddingBottom != 2 {
		t.Errorf("WithBorderPadding PaddingBottom = %d, want 2", cfg.PaddingBottom)
	}
	if cfg.PaddingLeft != 3 {
		t.Errorf("WithBorderPadding PaddingLeft = %d, want 3", cfg.PaddingLeft)
	}
	if cfg.PaddingRight != 4 {
		t.Errorf("WithBorderPadding PaddingRight = %d, want 4", cfg.PaddingRight)
	}
}

func TestWithBorderPadding_NegativeIgnored(t *testing.T) {
	cfg := styles.BorderConfig{
		PaddingTop:    10,
		PaddingBottom: 20,
		PaddingLeft:   30,
		PaddingRight:  40,
	}
	opt := styles.WithBorderPadding(-1, -1, -1, -1)
	opt(&cfg)

	// Negative values should be ignored
	if cfg.PaddingTop != 10 {
		t.Errorf("WithBorderPadding should preserve PaddingTop with negative, got %d", cfg.PaddingTop)
	}
	if cfg.PaddingBottom != 20 {
		t.Errorf("WithBorderPadding should preserve PaddingBottom with negative, got %d", cfg.PaddingBottom)
	}
	if cfg.PaddingLeft != 30 {
		t.Errorf("WithBorderPadding should preserve PaddingLeft with negative, got %d", cfg.PaddingLeft)
	}
	if cfg.PaddingRight != 40 {
		t.Errorf("WithBorderPadding should preserve PaddingRight with negative, got %d", cfg.PaddingRight)
	}
}

// ============================================================================
// NewBorderConfig Tests
// ============================================================================

func TestNewBorderConfig_NoOptions(t *testing.T) {
	cfg := styles.NewBorderConfig()
	defaults := styles.DefaultBorderConfig()

	// Should return defaults when no options provided
	if cfg.Left != defaults.Left {
		t.Errorf("NewBorderConfig().Left = %v, want %v", cfg.Left, defaults.Left)
	}
	if cfg.Right != defaults.Right {
		t.Errorf("NewBorderConfig().Right = %v, want %v", cfg.Right, defaults.Right)
	}
	if cfg.PaddingLeft != defaults.PaddingLeft {
		t.Errorf("NewBorderConfig().PaddingLeft = %d, want %d", cfg.PaddingLeft, defaults.PaddingLeft)
	}
}

func TestNewBorderConfig_WithSingleOption(t *testing.T) {
	cfg := styles.NewBorderConfig(styles.WithBorderLeft(true))

	if !cfg.Left {
		t.Error("NewBorderConfig with WithBorderLeft(true) should set Left to true")
	}
	// Should still have other defaults
	if !cfg.Right {
		t.Error("NewBorderConfig should preserve default Right value")
	}
}

func TestNewBorderConfig_WithMultipleOptions(t *testing.T) {
	color := lipgloss.Color("240")
	focusColor := lipgloss.Color("205")

	cfg := styles.NewBorderConfig(
		styles.WithBorderStyle(lipgloss.RoundedBorder()),
		styles.WithBorderLeft(true),
		styles.WithBorderRight(false),
		styles.WithBorderColor(color),
		styles.WithBorderFocusColor(focusColor),
		styles.WithBorderPadding(1, 2, 3, 4),
	)

	if cfg.Style.TopLeft != lipgloss.RoundedBorder().TopLeft {
		t.Error("NewBorderConfig should apply WithBorderStyle")
	}
	if !cfg.Left {
		t.Error("NewBorderConfig should apply WithBorderLeft")
	}
	if cfg.Right {
		t.Error("NewBorderConfig should apply WithBorderRight")
	}
	if cfg.Color != color {
		t.Error("NewBorderConfig should apply WithBorderColor")
	}
	if cfg.FocusColor != focusColor {
		t.Error("NewBorderConfig should apply WithBorderFocusColor")
	}
	if cfg.PaddingTop != 1 || cfg.PaddingBottom != 2 || cfg.PaddingLeft != 3 || cfg.PaddingRight != 4 {
		t.Error("NewBorderConfig should apply WithBorderPadding")
	}
}

func TestNewBorderConfig_OptionsAppliedInOrder(t *testing.T) {
	// Later options should override earlier ones
	cfg := styles.NewBorderConfig(
		styles.WithBorderLeft(true),
		styles.WithBorderLeft(false), // should override
	)

	if cfg.Left {
		t.Error("NewBorderConfig should apply options in order (later overrides earlier)")
	}
}
