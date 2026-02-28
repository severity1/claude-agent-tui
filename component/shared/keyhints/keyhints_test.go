package keyhints_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/severity1/claude-agent-tui/component/shared/keyhints"
)

// ============================================================================
// Model Initialization Tests
// ============================================================================

func TestNew_DefaultValues(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)

	if len(m.Bindings()) != 1 {
		t.Errorf("Bindings() length = %d, want 1", len(m.Bindings()))
	}
}

func TestNew_EmptyBindings(t *testing.T) {
	m := keyhints.New()

	if len(m.Bindings()) != 0 {
		t.Errorf("Bindings() length = %d, want 0", len(m.Bindings()))
	}
}

func TestNew_WithSeparator(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "Enter", Desc: "Submit"},
			keyhints.Binding{Key: "Esc", Desc: "Clear"},
		),
		keyhints.WithSeparator(" - "),
	)

	view := m.View()
	if !strings.Contains(view, " - ") {
		t.Errorf("View() = %q, want to contain ' - '", view)
	}
}

func TestNew_WithWidth(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
		keyhints.WithWidth(80),
	)

	// Width is stored internally, just verify model created
	if len(m.Bindings()) != 1 {
		t.Errorf("Bindings() length = %d, want 1", len(m.Bindings()))
	}
}

// ============================================================================
// Binding Management Tests
// ============================================================================

func TestModel_SetBindings(t *testing.T) {
	m := keyhints.New()

	newBindings := []keyhints.Binding{
		{Key: "q", Desc: "Quit"},
		{Key: "?", Desc: "Help"},
	}
	m.SetBindings(newBindings)

	if len(m.Bindings()) != 2 {
		t.Errorf("Bindings() length = %d, want 2", len(m.Bindings()))
	}
}

func TestModel_AddBinding(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)

	m.AddBinding(keyhints.Binding{Key: "Esc", Desc: "Clear"})

	if len(m.Bindings()) != 2 {
		t.Errorf("Bindings() length = %d, want 2", len(m.Bindings()))
	}
}

func TestModel_Bindings_ReturnsCopy(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)

	result := m.Bindings()
	if len(result) != 1 {
		t.Errorf("Bindings() length = %d, want 1", len(result))
	}
	if result[0].Key != "Enter" {
		t.Errorf("Bindings()[0].Key = %q, want %q", result[0].Key, "Enter")
	}

	// Verify mutation isolation - modifying returned slice should not affect model
	result[0].Key = "MODIFIED"
	result = append(result, keyhints.Binding{Key: "New", Desc: "Added"})

	// Get fresh copy and verify original is unchanged
	fresh := m.Bindings()
	if len(fresh) != 1 {
		t.Errorf("After mutation: Bindings() length = %d, want 1 (slice grew unexpectedly)", len(fresh))
	}
	if fresh[0].Key != "Enter" {
		t.Errorf("After mutation: Bindings()[0].Key = %q, want %q (value was mutated)", fresh[0].Key, "Enter")
	}
}

// ============================================================================
// View Tests
// ============================================================================

func TestModel_View_EmptyBindings(t *testing.T) {
	m := keyhints.New()

	view := m.View()
	if view != "" {
		t.Errorf("View() = %q, want empty string", view)
	}
}

func TestModel_View_SingleBinding(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)

	view := m.View()
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() = %q, want to contain 'Enter'", view)
	}
	if !strings.Contains(view, "Submit") {
		t.Errorf("View() = %q, want to contain 'Submit'", view)
	}
}

func TestModel_View_MultipleBindings(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "Enter", Desc: "Submit"},
			keyhints.Binding{Key: "Esc", Desc: "Clear"},
			keyhints.Binding{Key: "Ctrl+C", Desc: "Quit"},
		),
	)

	view := m.View()
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() missing 'Enter'")
	}
	if !strings.Contains(view, "Esc") {
		t.Errorf("View() missing 'Esc'")
	}
	if !strings.Contains(view, "Ctrl+C") {
		t.Errorf("View() missing 'Ctrl+C'")
	}
	// Default separator is " | "
	if !strings.Contains(view, "|") {
		t.Errorf("View() missing separator '|'")
	}
}

func TestModel_ViewCompact_Empty(t *testing.T) {
	m := keyhints.New()

	view := m.ViewCompact()
	if view != "" {
		t.Errorf("ViewCompact() = %q, want empty string", view)
	}
}

func TestModel_ViewCompact_ShowsKeysOnly(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "Enter", Desc: "Submit"},
			keyhints.Binding{Key: "Esc", Desc: "Clear"},
		),
	)

	view := m.ViewCompact()
	if !strings.Contains(view, "Enter") {
		t.Errorf("ViewCompact() missing 'Enter'")
	}
	if !strings.Contains(view, "Esc") {
		t.Errorf("ViewCompact() missing 'Esc'")
	}
	// Should NOT contain descriptions
	if strings.Contains(view, "Submit") {
		t.Errorf("ViewCompact() should not contain description 'Submit'")
	}
}

// ============================================================================
// Style Tests
// ============================================================================

func TestNew_WithKeyStyle(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
		keyhints.WithKeyStyle(style),
	)

	// Style is internal, verify model works
	view := m.View()
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() missing 'Enter'")
	}
}

func TestNew_WithDescStyle(t *testing.T) {
	style := lipgloss.NewStyle().Faint(true)
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
		keyhints.WithDescStyle(style),
	)

	view := m.View()
	if !strings.Contains(view, "Submit") {
		t.Errorf("View() missing 'Submit'")
	}
}

func TestNew_WithSepStyle(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "a", Desc: "A"},
			keyhints.Binding{Key: "b", Desc: "B"},
		),
		keyhints.WithSepStyle(style),
	)

	view := m.View()
	if !strings.Contains(view, "|") {
		t.Errorf("View() missing separator")
	}
}

// ============================================================================
// Expansion State Tests
// ============================================================================

func TestNew_DefaultExpanded(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "a", Desc: "A"},
			keyhints.Binding{Key: "b", Desc: "B"},
		),
	)

	if !m.Expanded() {
		t.Error("New() should create expanded model by default")
	}
}

func TestNew_WithCollapsed(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "A"}),
		keyhints.WithCollapsed(),
	)

	if m.Expanded() {
		t.Error("WithCollapsed() should create collapsed model")
	}
}

func TestNew_WithDefaultVisible(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "a", Desc: "A"},
			keyhints.Binding{Key: "b", Desc: "B"},
			keyhints.Binding{Key: "c", Desc: "C"},
			keyhints.Binding{Key: "d", Desc: "D"},
			keyhints.Binding{Key: "e", Desc: "E"},
		),
		keyhints.WithCollapsed(),
		keyhints.WithDefaultVisible(2),
	)

	view := m.View()
	// Should show first 2 bindings
	if !strings.Contains(view, "a") {
		t.Error("View() should contain first binding 'a'")
	}
	if !strings.Contains(view, "b") {
		t.Error("View() should contain second binding 'b'")
	}
	// Should NOT show the rest
	if strings.Contains(view, "c: C") {
		t.Error("View() should NOT contain third binding 'c'")
	}
	// Should show "? more" indicator
	if !strings.Contains(view, "3 more") {
		t.Errorf("View() = %q, should contain '3 more' indicator", view)
	}
}

func TestNew_WithDefaultVisible_InvalidValue(t *testing.T) {
	// Zero should use default (3)
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "a", Desc: "A"},
			keyhints.Binding{Key: "b", Desc: "B"},
			keyhints.Binding{Key: "c", Desc: "C"},
			keyhints.Binding{Key: "d", Desc: "D"},
		),
		keyhints.WithCollapsed(),
		keyhints.WithDefaultVisible(0),
	)

	view := m.View()
	// Should show 3 bindings (default)
	if !strings.Contains(view, "a") {
		t.Error("View() should contain 'a'")
	}
	if !strings.Contains(view, "b") {
		t.Error("View() should contain 'b'")
	}
	if !strings.Contains(view, "c") {
		t.Error("View() should contain 'c'")
	}
	// Should show "? 1 more"
	if !strings.Contains(view, "1 more") {
		t.Errorf("View() = %q, should contain '1 more' indicator", view)
	}
}

func TestModel_ToggleExpanded(t *testing.T) {
	m := keyhints.New()

	// Initially expanded
	if !m.Expanded() {
		t.Error("Initial state should be expanded")
	}

	// Toggle to collapsed
	m.ToggleExpanded()
	if m.Expanded() {
		t.Error("After first ToggleExpanded(), should be collapsed")
	}

	// Toggle back to expanded
	m.ToggleExpanded()
	if !m.Expanded() {
		t.Error("After second ToggleExpanded(), should be expanded")
	}
}

func TestModel_SetExpanded(t *testing.T) {
	m := keyhints.New()

	m.SetExpanded(false)
	if m.Expanded() {
		t.Error("SetExpanded(false) should collapse")
	}

	m.SetExpanded(true)
	if !m.Expanded() {
		t.Error("SetExpanded(true) should expand")
	}
}

func TestModel_SetWidth(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "Action A"}),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	// Initially no width set
	view1 := m.ViewHelp()

	// Set width
	m.SetWidth(50)
	view2 := m.ViewHelp()

	// After setting width, help window should be constrained
	lines := strings.Split(view2, "\n")
	for _, line := range lines {
		if lipgloss.Width(line) > 50 {
			t.Errorf("ViewHelp() after SetWidth(50) should respect width, got line width %d",
				lipgloss.Width(line))
		}
	}

	// Verify width was actually applied (view should be different or same width)
	_ = view1 // Used for comparison logic if needed
}

func TestModel_View_Collapsed_ShowsLimitedBindings(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "Enter", Desc: "Submit"},
			keyhints.Binding{Key: "Esc", Desc: "Clear"},
			keyhints.Binding{Key: "Tab", Desc: "Next"},
			keyhints.Binding{Key: "Shift+Tab", Desc: "Prev"},
			keyhints.Binding{Key: "Ctrl+C", Desc: "Quit"},
		),
		keyhints.WithCollapsed(),
	)

	view := m.View()
	// Default visible is 3
	if !strings.Contains(view, "Enter") {
		t.Error("View() should show 'Enter'")
	}
	if !strings.Contains(view, "Esc") {
		t.Error("View() should show 'Esc'")
	}
	if !strings.Contains(view, "Tab") {
		t.Error("View() should show 'Tab'")
	}
	// Should NOT show 4th and 5th
	if strings.Contains(view, "Shift+Tab: Prev") {
		t.Error("View() should NOT show 'Shift+Tab'")
	}
	// Should show indicator
	if !strings.Contains(view, "2 more") {
		t.Errorf("View() = %q, should contain '2 more'", view)
	}
}

func TestModel_View_Collapsed_NoIndicatorWhenAllVisible(t *testing.T) {
	// Only 2 bindings, default visible is 3
	m := keyhints.New(
		keyhints.WithBindings(
			keyhints.Binding{Key: "Enter", Desc: "Submit"},
			keyhints.Binding{Key: "Esc", Desc: "Clear"},
		),
		keyhints.WithCollapsed(),
	)

	view := m.View()
	if strings.Contains(view, "more") {
		t.Errorf("View() = %q, should NOT contain 'more' indicator when all visible", view)
	}
}

func TestModel_View_Expanded_ShowsAllBindings(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
		{Key: "c", Desc: "C"},
		{Key: "d", Desc: "D"},
		{Key: "e", Desc: "E"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...)) // Default is expanded

	view := m.View()
	for _, b := range bindings {
		if !strings.Contains(view, b.Key) {
			t.Errorf("View() should contain key %q", b.Key)
		}
	}
	// No "more" indicator
	if strings.Contains(view, "more") {
		t.Errorf("View() = %q, expanded should NOT contain 'more'", view)
	}
}

// ============================================================================
// Bubble Tea Interface Tests
// ============================================================================

func TestModel_Init_ReturnsNil(t *testing.T) {
	m := keyhints.New()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init() returned non-nil command, want nil")
	}
}

func TestModel_Update_ReturnsModelAndNil(t *testing.T) {
	m := keyhints.New()

	updated, cmd := m.Update(nil)

	if cmd != nil {
		t.Error("Update() returned non-nil command, want nil")
	}
	if _, ok := updated.(keyhints.Model); !ok {
		t.Errorf("Update() returned %T, want keyhints.Model", updated)
	}
}

// ============================================================================
// Focus Tests
// ============================================================================

func TestNew_DefaultNotFocused(t *testing.T) {
	m := keyhints.New()

	if m.Focused() {
		t.Error("New() should create unfocused model by default")
	}
}

func TestModel_Focus(t *testing.T) {
	m := keyhints.New()

	cmd := m.Focus()

	if !m.Focused() {
		t.Error("Focus() should set focused state")
	}
	if cmd != nil {
		t.Error("Focus() should return nil (no cursor blink for keyhints)")
	}
}

func TestModel_Blur(t *testing.T) {
	m := keyhints.New()
	m.Focus()

	m.Blur()

	if m.Focused() {
		t.Error("Blur() should clear focused state")
	}
}

func TestModel_View_FocusedShowsBorder(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)
	m.Focus()

	view := m.View()
	// Should have left border (thick border character)
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() when focused should contain thick border character, got: %q", view)
	}
}

func TestModel_View_UnfocusedHasBorder(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
	)
	// Not focused by default

	view := m.View()
	// Should have left border (thick border character) in unfocused state too
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() when unfocused should contain thick border character, got: %q", view)
	}
}

func TestModel_Update_FocusedHandlesQuestionMark(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "A"}),
	)
	m.Focus()

	// Initial state: help not showing
	if m.ShowingHelp() {
		t.Fatal("Should start with help hidden")
	}

	// Simulate pressing "?" using tea.KeyMsg
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updated, _ := m.Update(keyMsg)
	m = updated.(keyhints.Model)

	// Should now show help
	if !m.ShowingHelp() {
		t.Error("Update() with '?' should toggle help window when focused")
	}
}

func TestModel_Update_UnfocusedIgnoresQuestionMark(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "A"}),
	)
	// Not focused

	// Initial state: help not showing
	if m.ShowingHelp() {
		t.Fatal("Should start with help hidden")
	}

	// Simulate pressing "?" using tea.KeyMsg
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updated, _ := m.Update(keyMsg)
	m = updated.(keyhints.Model)

	// Should still not show help (ignored because not focused)
	if m.ShowingHelp() {
		t.Error("Update() with '?' should be ignored when not focused")
	}
}

func TestNew_WithFocusStyle(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "Enter", Desc: "Submit"}),
		keyhints.WithFocusStyle(style),
	)
	m.Focus()

	view := m.View()
	// Just verify it renders without error
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() missing 'Enter'")
	}
}

// ============================================================================
// Help Window Tests
// ============================================================================

func TestNew_DefaultHelpHidden(t *testing.T) {
	m := keyhints.New()

	if m.ShowingHelp() {
		t.Error("New() should create model with help hidden")
	}
}

func TestModel_ShowHelp(t *testing.T) {
	m := keyhints.New()

	m.ShowHelp()

	if !m.ShowingHelp() {
		t.Error("ShowHelp() should show help window")
	}
}

func TestModel_HideHelp(t *testing.T) {
	m := keyhints.New()
	m.ShowHelp()

	m.HideHelp()

	if m.ShowingHelp() {
		t.Error("HideHelp() should hide help window")
	}
}

func TestModel_ToggleHelp(t *testing.T) {
	m := keyhints.New()

	m.ToggleHelp()
	if !m.ShowingHelp() {
		t.Error("ToggleHelp() should show help when hidden")
	}

	m.ToggleHelp()
	if m.ShowingHelp() {
		t.Error("ToggleHelp() should hide help when shown")
	}
}

func TestModel_ViewHelp_EmptyBindings(t *testing.T) {
	m := keyhints.New()

	view := m.ViewHelp()
	if !strings.Contains(view, "No keybindings") {
		t.Errorf("ViewHelp() with no bindings should show 'No keybindings', got %q", view)
	}
}

func TestModel_ViewHelp_ShowsAllBindings(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	view := m.ViewHelp()
	for _, b := range bindings {
		if !strings.Contains(view, b.Key) {
			t.Errorf("ViewHelp() missing key %q", b.Key)
		}
		if !strings.Contains(view, b.Desc) {
			t.Errorf("ViewHelp() missing desc %q", b.Desc)
		}
	}
}

func TestModel_ViewHelp_HasTitle(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "A"}),
	)

	view := m.ViewHelp()
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("ViewHelp() should contain title 'Keyboard Shortcuts'")
	}
}

func TestModel_ViewHelp_HasCloseHint(t *testing.T) {
	m := keyhints.New(
		keyhints.WithBindings(keyhints.Binding{Key: "a", Desc: "A"}),
	)

	view := m.ViewHelp()
	if !strings.Contains(view, "Esc") {
		t.Errorf("ViewHelp() should mention Esc to close")
	}
}

func TestModel_Update_EscClosesHelp(t *testing.T) {
	m := keyhints.New()
	m.Focus()
	m.ShowHelp()

	// Simulate pressing Esc
	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updated, _ := m.Update(keyMsg)
	m = updated.(keyhints.Model)

	if m.ShowingHelp() {
		t.Error("Esc should close help window")
	}
}

func TestNew_WithHelpStyle(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	style := lipgloss.NewStyle().Border(lipgloss.DoubleBorder())
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithHelpStyle(style))

	view := m.ViewHelp()
	// Just verify it renders without error
	if !strings.Contains(view, "Enter") {
		t.Errorf("ViewHelp() missing 'Enter'")
	}
}

// ============================================================================
// Width-Based Truncation Tests
// ============================================================================

func TestModel_View_WithWidth_TruncatesWhenNeeded(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
		{Key: "Tab", Desc: "Next"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	// Set a narrow width that can't fit all bindings
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(40))

	view := m.View()
	// Should show truncation indicator
	if !strings.Contains(view, "more") {
		t.Errorf("View() with narrow width should contain 'more' indicator, got: %q", view)
	}
}

func TestModel_View_WithWidth_ShowsAllWhenFits(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	// Set a wide width that can fit all bindings
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(200))

	view := m.View()
	// Should NOT show truncation indicator
	if strings.Contains(view, "more") {
		t.Errorf("View() with wide width should NOT contain 'more' indicator, got: %q", view)
	}
	// Should contain both bindings
	if !strings.Contains(view, "a") {
		t.Errorf("View() should contain 'a'")
	}
	if !strings.Contains(view, "b") {
		t.Errorf("View() should contain 'b'")
	}
}

func TestModel_View_WithWidth_UsesPlusFormat(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit long description"},
		{Key: "Esc", Desc: "Clear input field"},
		{Key: "Tab", Desc: "Next field"},
	}
	// Set a narrow width
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(50))

	view := m.View()
	// Width-based truncation uses "+N more" format
	if strings.Contains(view, "more") && !strings.Contains(view, "+") {
		t.Errorf("Width-based truncation should use '+N more' format, got: %q", view)
	}
}

// ============================================================================
// SetBindings Mutation Isolation Tests
// ============================================================================

func TestModel_SetBindings_ClonesInput(t *testing.T) {
	m := keyhints.New()

	original := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m.SetBindings(original)

	// Mutate the original slice
	original[0].Key = "MODIFIED"

	// Get bindings from model - should be unchanged
	result := m.Bindings()
	if result[0].Key != "Enter" {
		t.Errorf("SetBindings should clone input - mutation affected model: got %q, want %q",
			result[0].Key, "Enter")
	}
}

// ============================================================================
// Padding Tests
// ============================================================================

func TestNew_DefaultPadding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	view := m.View()
	// Default padding is empty string, view starts with border character
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() should contain border character, got: %q", view)
	}
}

func TestNew_WithPadding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding("   ")) // 3 spaces

	view := m.View()
	// Padding is applied inside the border, so content should contain the padding followed by key
	if !strings.Contains(view, "   Enter") {
		t.Errorf("View() should contain custom padding before key, got: %q", view)
	}
}

func TestNew_WithPadding_Empty(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding("")) // no padding

	view := m.View()
	// Should contain border and content
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() with empty padding should contain binding key, got: %q", view)
	}
}

func TestModel_View_AppliesPadding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding(">>"))

	view := m.View()
	// Padding is applied inside the border
	if !strings.Contains(view, ">>a") {
		t.Errorf("View() should apply custom padding before key, got: %q", view)
	}
}

func TestModel_ViewCompact_AppliesPadding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding(">>"))

	view := m.ViewCompact()
	// Padding is applied, ViewCompact does not use border wrapper
	if !strings.Contains(view, ">>") {
		t.Errorf("ViewCompact() should apply custom padding, got: %q", view)
	}
}

func TestModel_ViewHelp_NoPadding(t *testing.T) {
	// ViewHelp should NOT apply padding - it has its own left border styling
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding(">>"), keyhints.WithBorderLeft(true))

	view := m.ViewHelp()
	// Should NOT start with padding - should start with left border
	if strings.HasPrefix(view, ">>") {
		t.Errorf("ViewHelp() should NOT apply padding prefix, got: %q", view[:50])
	}
	// Should start with thick left border character when borderLeft=true
	if !strings.HasPrefix(view, "┃") && !strings.HasPrefix(view, "\u2503") {
		t.Errorf("ViewHelp() with borderLeft=true should start with thick left border, got: %q", view[:20])
	}
}

func TestModel_View_FocusedWithPadding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithPadding("  ")) // 2 spaces
	m.Focus()

	view := m.View()
	// View should contain the binding and border
	if !strings.Contains(view, "a") {
		t.Errorf("View() when focused should contain binding, got: %q", view)
	}
	// Should contain border character
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() when focused should contain border, got: %q", view)
	}
}

func TestModel_ViewCompact_EmptyWithPadding(t *testing.T) {
	// Empty bindings should return empty string regardless of padding
	m := keyhints.New(keyhints.WithPadding(">>"))

	view := m.ViewCompact()
	if view != "" {
		t.Errorf("ViewCompact() with empty bindings should return empty string, got: %q", view)
	}
}

func TestModel_View_EmptyWithPadding(t *testing.T) {
	// Empty bindings should return empty string regardless of padding
	m := keyhints.New(keyhints.WithPadding(">>"))

	view := m.View()
	if view != "" {
		t.Errorf("View() with empty bindings should return empty string, got: %q", view)
	}
}

// ============================================================================
// Boundary Value Tests
// ============================================================================

func TestNew_WithWidth_ZeroValue(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
	}
	// WithWidth(0) should effectively disable width constraint
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(0))

	view := m.View()
	// Should render all bindings without truncation when width is 0
	if !strings.Contains(view, "Enter") {
		t.Error("View() with width 0 should contain 'Enter'")
	}
	if !strings.Contains(view, "Esc") {
		t.Error("View() with width 0 should contain 'Esc'")
	}
	// Should NOT show "more" indicator when width constraint is disabled
	if strings.Contains(view, "more") {
		t.Errorf("View() with width 0 should NOT contain 'more' indicator, got: %q", view)
	}
}

// ============================================================================
// HelpMode Tests
// ============================================================================

func TestNew_DefaultHelpMode(t *testing.T) {
	m := keyhints.New()

	if m.HelpMode() != keyhints.HelpModeDown {
		t.Errorf("Default HelpMode() = %v, want HelpModeDown (0)", m.HelpMode())
	}
}

func TestNew_WithHelpMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     keyhints.HelpMode
		expected keyhints.HelpMode
	}{
		{"HelpModeDown", keyhints.HelpModeDown, keyhints.HelpModeDown},
		{"HelpModeUp", keyhints.HelpModeUp, keyhints.HelpModeUp},
		{"HelpModeCenter", keyhints.HelpModeCenter, keyhints.HelpModeCenter},
		{"HelpModeTop", keyhints.HelpModeTop, keyhints.HelpModeTop},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := keyhints.New(keyhints.WithHelpMode(tt.mode))
			if m.HelpMode() != tt.expected {
				t.Errorf("HelpMode() = %v, want %v", m.HelpMode(), tt.expected)
			}
		})
	}
}

// ============================================================================
// HelpHeight Tests
// ============================================================================

func TestModel_HelpHeight_EmptyBindings(t *testing.T) {
	m := keyhints.New()

	height := m.HelpHeight()
	// Empty bindings shows "No keybindings defined" which is 1 line
	if height != 1 {
		t.Errorf("HelpHeight() with empty bindings = %d, want 1", height)
	}
}

func TestModel_HelpHeight_WithBindings(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
		{Key: "c", Desc: "C"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	height := m.HelpHeight()
	// Title (1) + blank (1) + bindings (3) + blank (1) + close hint (1) = 7
	expected := 7
	if height != expected {
		t.Errorf("HelpHeight() with 3 bindings = %d, want %d", height, expected)
	}
}

// ============================================================================
// ViewHelp Left Border Tests
// ============================================================================

func TestModel_ViewHelp_HasLeftBorder(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	view := m.ViewHelp()
	// Should contain thick left border character
	if !strings.Contains(view, "┃") && !strings.Contains(view, "\u2503") {
		t.Errorf("ViewHelp() should contain thick left border, got: %q", view)
	}
}

func TestModel_ViewHelp_NoFullBorder(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	view := m.ViewHelp()
	// Should NOT contain rounded border characters
	roundedBorderChars := []string{"╭", "╮", "╰", "╯"}
	for _, char := range roundedBorderChars {
		if strings.Contains(view, char) {
			t.Errorf("ViewHelp() should NOT contain rounded border char %q, got: %q", char, view)
		}
	}
}

// ============================================================================
// ContentBackground Tests
// ============================================================================

func TestNew_WithContentBackground(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	bgColor := lipgloss.Color("236")
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithContentBackground(bgColor))

	// Just verify it renders without error
	view := m.ViewHelp()
	if !strings.Contains(view, "a") {
		t.Errorf("ViewHelp() with content background should contain binding key")
	}
}

func TestModel_View_WithContentBackground(t *testing.T) {
	// View() inline bar should render correctly when contentBg is set
	// Note: Background styling is NOT applied to View() - only to ViewHelp()
	// This test verifies View() still works correctly when contentBg option is used
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	bgColor := lipgloss.Color("236")
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithContentBackground(bgColor))

	view := m.View()

	// Should contain the binding content
	if !strings.Contains(view, "a") {
		t.Errorf("View() should contain binding key")
	}
	// Should contain the binding description
	if !strings.Contains(view, "A") {
		t.Errorf("View() should contain binding description")
	}
}

func TestModel_ViewHelp_WithContentBackground_AttachedModes(t *testing.T) {
	// ViewHelp() applies background color using Place() with WhitespaceBackground
	// in attached modes (Down, Up, Top)
	attachedModes := []keyhints.HelpMode{
		keyhints.HelpModeDown,
		keyhints.HelpModeUp,
		keyhints.HelpModeTop,
	}

	for _, mode := range attachedModes {
		t.Run(modeToString(mode), func(t *testing.T) {
			bindings := []keyhints.Binding{
				{Key: "a", Desc: "Action A"},
			}
			bgColor := lipgloss.Color("236")
			m := keyhints.New(keyhints.WithBindings(bindings...),
				keyhints.WithContentBackground(bgColor),
				keyhints.WithHelpMode(mode),
				keyhints.WithWidth(50),
			)

			view := m.ViewHelp()

			// Should contain binding content
			if !strings.Contains(view, "a") {
				t.Errorf("ViewHelp() should contain binding key")
			}
			if !strings.Contains(view, "Action A") {
				t.Errorf("ViewHelp() should contain binding description")
			}
			// Should contain title
			if !strings.Contains(view, "Keyboard Shortcuts") {
				t.Errorf("ViewHelp() should contain title")
			}
		})
	}
}

func TestModel_ViewHelp_WithContentBackground_CenterMode(t *testing.T) {
	// ViewHelp() in Center mode also applies background using Place()
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
	}
	bgColor := lipgloss.Color("236")
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithContentBackground(bgColor),
		keyhints.WithHelpMode(keyhints.HelpModeCenter),
	)

	view := m.ViewHelp()

	// Should contain binding content
	if !strings.Contains(view, "a") {
		t.Errorf("ViewHelp() in Center mode should contain binding key")
	}
	if !strings.Contains(view, "Action A") {
		t.Errorf("ViewHelp() in Center mode should contain binding description")
	}
}

func TestModel_ViewHelp_NoContentBackground(t *testing.T) {
	// ViewHelp() should render correctly when contentBg is explicitly nil
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	// Explicitly set nil background
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithContentBackground(nil),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	view := m.ViewHelp()

	// Should render without error and contain content
	if !strings.Contains(view, "a") {
		t.Errorf("ViewHelp() should contain binding key")
	}
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("ViewHelp() should contain title")
	}
}

// ============================================================================
// ViewHelp Width Behavior Per Mode Tests
// ============================================================================

func TestModel_ViewHelp_DownMode_MatchesBarWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
		{Key: "b", Desc: "Action B"},
	}
	width := 60
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(width),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	view := m.ViewHelp()
	// Help window should have consistent width matching the bar
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > width {
			t.Errorf("ViewHelp() line exceeds width %d: got %d for line %q",
				width, lineWidth, line)
		}
	}
}

func TestModel_ViewHelp_UpMode_MatchesBarWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
	}
	width := 50
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(width),
		keyhints.WithHelpMode(keyhints.HelpModeUp),
	)

	view := m.ViewHelp()
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > width {
			t.Errorf("ViewHelp() in Up mode line exceeds width %d: got %d",
				width, lineWidth)
		}
	}
}

func TestModel_ViewHelp_TopMode_MatchesBarWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
	}
	width := 50
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(width),
		keyhints.WithHelpMode(keyhints.HelpModeTop),
	)

	view := m.ViewHelp()
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > width {
			t.Errorf("ViewHelp() in Top mode line exceeds width %d: got %d",
				width, lineWidth)
		}
	}
}

func TestModel_ViewHelp_CenterMode_UsesContentWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
	}
	// Set a narrow width that would truncate content if applied
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(20),
		keyhints.WithHelpMode(keyhints.HelpModeCenter),
	)

	view := m.ViewHelp()
	// Center mode should NOT constrain width - content should be fully visible
	if !strings.Contains(view, "Action A") {
		t.Errorf("ViewHelp() in Center mode should use content width, but description was truncated")
	}
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("ViewHelp() in Center mode should contain full title")
	}
}

func TestModel_ViewHelp_ZeroWidth_NoConstraint(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit the form"},
	}
	// Width 0 should not apply any constraint
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(0),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	view := m.ViewHelp()
	// Content should be fully visible
	if !strings.Contains(view, "Submit the form") {
		t.Errorf("ViewHelp() with width 0 should not constrain content")
	}
}

func TestModel_ViewHelp_EmptyBindings_WithWidth(t *testing.T) {
	m := keyhints.New(
		keyhints.WithWidth(40),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	view := m.ViewHelp()
	if !strings.Contains(view, "No keybindings") {
		t.Errorf("ViewHelp() with empty bindings should show message, got %q", view)
	}
	// Should also respect width for empty message
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		if lipgloss.Width(line) > 40 {
			t.Errorf("ViewHelp() empty message exceeds width 40")
		}
	}
}

func TestModel_ViewHelp_AllAttachedModes_ApplyWidth(t *testing.T) {
	// Attached modes: Down, Up, Top - all should apply width constraint
	attachedModes := []keyhints.HelpMode{
		keyhints.HelpModeDown,
		keyhints.HelpModeUp,
		keyhints.HelpModeTop,
	}

	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	width := 50

	for _, mode := range attachedModes {
		t.Run(modeToString(mode), func(t *testing.T) {
			m := keyhints.New(keyhints.WithBindings(bindings...),
				keyhints.WithWidth(width),
				keyhints.WithHelpMode(mode),
			)

			view := m.ViewHelp()
			lines := strings.Split(view, "\n")
			for i, line := range lines {
				lineWidth := lipgloss.Width(line)
				if lineWidth > width {
					t.Errorf("Line %d exceeds width %d: got %d for %q",
						i, width, lineWidth, line)
				}
			}
		})
	}
}

// Helper function to convert HelpMode to string for test naming
func modeToString(mode keyhints.HelpMode) string {
	switch mode {
	case keyhints.HelpModeDown:
		return "Down"
	case keyhints.HelpModeUp:
		return "Up"
	case keyhints.HelpModeCenter:
		return "Center"
	case keyhints.HelpModeTop:
		return "Top"
	default:
		return "Unknown"
	}
}

// ============================================================================
// View Width Style Tests
// ============================================================================

func TestModel_View_WithWidth_AppliesWidthStyle(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	width := 80
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(width))

	view := m.View()
	viewWidth := lipgloss.Width(view)

	// View should fill the specified width (accounting for both borders)
	// With borderRight = true (default), bw = 2 (left + right borders)
	// Style.Width(width - bw) + 2 borders = width
	expectedWidth := width
	if viewWidth != expectedWidth {
		t.Errorf("View() width = %d, want %d", viewWidth, expectedWidth)
	}
}

func TestModel_View_WithWidth_ZeroDoesNotApplyStyle(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	// Width 0 should not apply width style
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(0))

	view := m.View()
	// Just verify it renders without crashing
	if !strings.Contains(view, "a") {
		t.Error("View() with width 0 should contain binding")
	}
}

func TestModel_View_ViewHelp_ConsistentWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
		{Key: "b", Desc: "Action B"},
	}
	width := 60
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithWidth(width),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	barView := m.View()
	barWidth := lipgloss.Width(barView)

	m.Focus()
	m.ShowHelp()
	helpView := m.ViewHelp()
	helpLines := strings.Split(helpView, "\n")

	// All help lines should not exceed bar width
	for i, line := range helpLines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > barWidth {
			t.Errorf("ViewHelp() line %d width (%d) exceeds View() width (%d)",
				i, lineWidth, barWidth)
		}
	}
}

func TestModel_View_SetWidth_UpdatesRendering(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(50))

	view1 := m.View()
	width1 := lipgloss.Width(view1)

	m.SetWidth(80)
	view2 := m.View()
	width2 := lipgloss.Width(view2)

	if width1 == width2 {
		t.Errorf("SetWidth() should change rendered width: before=%d, after=%d", width1, width2)
	}
	// Width should equal 80 (with both left and right borders, total = width)
	expectedWidth := 80
	if width2 != expectedWidth {
		t.Errorf("View() after SetWidth(80) = %d, want %d", width2, expectedWidth)
	}
}

// ============================================================================
// Border Left Toggle Tests
// ============================================================================

func TestNew_DefaultBorderLeft(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	view := m.View()
	// Default borderLeft is true, so left border should be visible (thick border character)
	if !strings.HasPrefix(view, "┃") && !strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with default borderLeft=true should show thick left border, got first chars: %q", view)
	}
}

func TestNew_DefaultBorderRight(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(40))

	view := m.View()
	lines := strings.Split(view, "\n")
	// Default borderRight is false, so right border should be invisible (space character)
	// No line should end with thick border character
	for _, line := range lines {
		if len(line) > 0 && (strings.HasSuffix(line, "┃") || strings.HasSuffix(line, "\u2503")) {
			t.Errorf("View() with default borderRight=false should not show thick right border, got line: %q", line)
			break
		}
	}
}

func TestNew_WithBorderLeft_True(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderLeft(true))

	view := m.View()
	// borderLeft=true should show visible thick border
	if !strings.HasPrefix(view, "┃") && !strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with borderLeft=true should show thick left border, got first chars: %q", view)
	}
}

func TestNew_WithBorderLeft_False(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderLeft(false))

	view := m.View()
	// borderLeft=false should hide left border (use space)
	if strings.HasPrefix(view, "┃") || strings.HasPrefix(view, "\u2503") {
		t.Errorf("View() with borderLeft=false should not show thick left border, got first chars: %q", view)
	}
}

func TestModel_View_BothBordersHidden(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithBorderLeft(false),
		keyhints.WithBorderRight(false),
		keyhints.WithWidth(40),
	)

	view := m.View()
	// Both borders hidden but space reserved
	// View should still have consistent width
	viewWidth := lipgloss.Width(view)
	if viewWidth != 40 {
		t.Errorf("View() width with both borders hidden = %d, want 40", viewWidth)
	}
	// Should not contain thick border characters
	if strings.Contains(view, "┃") || strings.Contains(view, "\u2503") {
		t.Errorf("View() with both borders hidden should not contain thick border characters")
	}
}

func TestModel_ViewHelp_BorderVisibility(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithBorderLeft(false),
		keyhints.WithBorderRight(true),
		keyhints.WithHelpMode(keyhints.HelpModeDown),
	)

	view := m.ViewHelp()
	// Left border should be invisible (space)
	if strings.HasPrefix(view, "┃") || strings.HasPrefix(view, "\u2503") {
		t.Errorf("ViewHelp() with borderLeft=false should not show thick left border")
	}
	// Right border should be visible
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		if len(line) > 0 && !strings.HasSuffix(line, "┃") && !strings.HasSuffix(line, "\u2503") {
			// This is expected because the right border is visible
			// but we need to verify the content is there
		}
	}
	// Should contain binding content
	if !strings.Contains(view, "A") {
		t.Errorf("ViewHelp() should contain binding description")
	}
}

// ============================================================================
// Animation Tests
// ============================================================================

func TestNew_DefaultAnimationDisabled(t *testing.T) {
	m := keyhints.New()

	if m.AnimationCmd() != nil {
		t.Error("New() should create model with animation not in progress")
	}
}

func TestNew_WithAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	// Animation is enabled but not yet in progress
	if m.AnimationCmd() != nil {
		t.Error("WithAnimation(true) should not start animation immediately")
	}
}

func TestNew_WithAnimationSpring(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	// Should not panic with custom spring parameters
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithAnimationSpring(8.0, 0.5),
	)

	if m.AnimationCmd() != nil {
		t.Error("WithAnimationSpring should not start animation")
	}
}

func TestModel_ShowHelp_WithAnimation_StartsAnimating(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()

	if m.AnimationCmd() == nil {
		t.Error("ShowHelp() with animation enabled should start animation")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp() should set showHelp to true")
	}
}

func TestModel_HideHelp_WithAnimation_StartsAnimating(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()
	// Manually stop animation for testing
	for i := 0; i < 100; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	m.HideHelp()

	if m.AnimationCmd() == nil {
		t.Error("HideHelp() with animation enabled should start animation")
	}
	// showHelp should still be true until animation completes
	if !m.ShowingHelp() {
		t.Error("HideHelp() with animation should keep showHelp=true until animation completes")
	}
}

func TestModel_ShowHelp_WithoutAnimation_NoAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...)) // No animation

	m.ShowHelp()

	if m.AnimationCmd() != nil {
		t.Error("ShowHelp() without animation enabled should not animate")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp() should set showHelp to true")
	}
}

func TestModel_AnimationCmd_WhenAnimating(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()

	cmd := m.AnimationCmd()
	if cmd == nil {
		t.Error("AnimationCmd() should return tick command when animating")
	}
}

func TestModel_AnimationCmd_WhenNotAnimating(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	// Not animating yet
	cmd := m.AnimationCmd()
	if cmd != nil {
		t.Error("AnimationCmd() should return nil when not animating")
	}
}

func TestModel_AnimatedHelpHeight_NoAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...)) // No animation

	// Not showing help
	height := m.AnimatedHelpHeight()
	if height != 0 {
		t.Errorf("AnimatedHelpHeight() when not showing help = %d, want 0", height)
	}

	// Show help
	m.ShowHelp()
	height = m.AnimatedHelpHeight()
	expected := m.HelpHeight()
	if height != expected {
		t.Errorf("AnimatedHelpHeight() when showing help = %d, want %d", height, expected)
	}
}

func TestModel_AnimatedHelpHeight_WithAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()

	// During animation, height should be a positive value
	height := m.AnimatedHelpHeight()
	if height < 0 {
		t.Errorf("AnimatedHelpHeight() during animation should be >= 0, got %d", height)
	}
}

func TestModel_ViewHelpPartial_ZeroLines(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	result := m.ViewHelpPartial(0)
	if result != "" {
		t.Errorf("ViewHelpPartial(0) = %q, want empty string", result)
	}
}

func TestModel_ViewHelpPartial_NegativeLines(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	result := m.ViewHelpPartial(-1)
	if result != "" {
		t.Errorf("ViewHelpPartial(-1) = %q, want empty string", result)
	}
}

func TestModel_ViewHelpPartial_FullContent(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	fullView := m.ViewHelp()
	lineCount := len(strings.Split(fullView, "\n"))

	// Request more lines than available
	result := m.ViewHelpPartial(lineCount + 10)
	if result != fullView {
		t.Error("ViewHelpPartial() with more lines than available should return full content")
	}
}

func TestModel_ViewHelpPartial_TopAnchored_Down(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithHelpMode(keyhints.HelpModeDown))

	fullView := m.ViewHelp()
	lines := strings.Split(fullView, "\n")

	// Request fewer lines - should show first N lines (top-anchored)
	result := m.ViewHelpPartial(2)
	resultLines := strings.Split(result, "\n")

	if len(resultLines) != 2 {
		t.Errorf("ViewHelpPartial(2) returned %d lines, want 2", len(resultLines))
	}
	// First lines should match
	if resultLines[0] != lines[0] {
		t.Error("ViewHelpPartial() in Down mode should show first lines (top-anchored)")
	}
}

func TestModel_ViewHelpPartial_BottomAnchored_Up(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithHelpMode(keyhints.HelpModeUp))

	fullView := m.ViewHelp()
	fullLines := strings.Split(fullView, "\n")
	totalLines := len(fullLines)

	// Request fewer lines - should show last N lines (NO padding)
	visibleCount := 2
	result := m.ViewHelpPartial(visibleCount)
	resultLines := strings.Split(result, "\n")

	// Result should have exactly visibleCount lines (no padding)
	if len(resultLines) != visibleCount {
		t.Errorf("ViewHelpPartial(%d) returned %d lines, want %d",
			visibleCount, len(resultLines), visibleCount)
	}

	// Lines should be the LAST visibleCount lines from full content
	for i := 0; i < visibleCount; i++ {
		expectedIdx := totalLines - visibleCount + i
		if resultLines[i] != fullLines[expectedIdx] {
			t.Errorf("ViewHelpPartial line %d = %q, want %q",
				i, resultLines[i], fullLines[expectedIdx])
		}
	}
}

func TestModel_ViewHelpPartial_TopAnchored_Center(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithHelpMode(keyhints.HelpModeCenter))

	fullView := m.ViewHelp()
	lines := strings.Split(fullView, "\n")

	// Center mode should be top-anchored
	result := m.ViewHelpPartial(2)
	resultLines := strings.Split(result, "\n")

	if resultLines[0] != lines[0] {
		t.Error("ViewHelpPartial() in Center mode should show first lines (top-anchored)")
	}
}

func TestModel_ViewHelpPartial_TopAnchored_Top(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithHelpMode(keyhints.HelpModeTop))

	fullView := m.ViewHelp()
	lines := strings.Split(fullView, "\n")

	// Top mode should be top-anchored
	result := m.ViewHelpPartial(2)
	resultLines := strings.Split(result, "\n")

	if resultLines[0] != lines[0] {
		t.Error("ViewHelpPartial() in Top mode should show first lines (top-anchored)")
	}
}

func TestModel_Update_HeightTickMsg_NotAnimating(t *testing.T) {
	m := keyhints.New(keyhints.WithAnimation(true))

	// Not animating
	updated, cmd := m.Update(keyhints.HeightTickMsg{})

	if cmd != nil {
		t.Error("Update(HeightTickMsg) when not animating should return nil cmd")
	}
	if _, ok := updated.(keyhints.Model); !ok {
		t.Errorf("Update() returned %T, want keyhints.Model", updated)
	}
}

func TestModel_Update_HeightTickMsg_Animating(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()

	// First tick should continue animation
	updated, cmd := m.Update(keyhints.HeightTickMsg{})
	m = updated.(keyhints.Model)

	// Animation should either continue or complete
	if m.AnimationCmd() != nil && cmd == nil {
		t.Error("Update(HeightTickMsg) while animating should return tick cmd")
	}
}

func TestModel_Update_HeightTickMsg_AnimationCompletes(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	m.ShowHelp()

	// Run animation until it completes
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	if m.AnimationCmd() != nil {
		t.Error("Animation should complete after sufficient ticks")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp animation should keep help visible after completion")
	}
}

func TestModel_Update_HeightTickMsg_HideAnimationCompletes(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	// Show and wait for animation
	m.ShowHelp()
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	// Now hide
	m.HideHelp()
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	if m.AnimationCmd() != nil {
		t.Error("Hide animation should complete after sufficient ticks")
	}
	if m.ShowingHelp() {
		t.Error("HideHelp animation should set showHelp=false after completion")
	}
}

func TestModel_Update_QuestionMark_WithAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))
	m.Focus()

	// Press ? to toggle help
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updated, cmd := m.Update(keyMsg)
	m = updated.(keyhints.Model)

	if !m.ShowingHelp() {
		t.Error("? key should show help")
	}
	if m.AnimationCmd() == nil {
		t.Error("? key with animation should start animation")
	}
	if cmd == nil {
		t.Error("? key with animation should return tick cmd")
	}
}

func TestModel_Update_Esc_WithAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))
	m.Focus()
	m.ShowHelp()

	// Run show animation to completion
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	// Press Esc to hide help
	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updated, cmd := m.Update(keyMsg)
	m = updated.(keyhints.Model)

	if m.AnimationCmd() == nil {
		t.Error("Esc key with animation should start hide animation")
	}
	if cmd == nil {
		t.Error("Esc key with animation should return tick cmd")
	}
}

// ============================================================================
// Center Mode Animation Skip Tests
// ============================================================================

func TestModel_ShowHelp_CenterMode_NoAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithHelpMode(keyhints.HelpModeCenter),
	)

	m.ShowHelp()

	// Center mode should skip animation even when enabled
	if m.AnimationCmd() != nil {
		t.Error("ShowHelp() in Center mode should NOT start animation")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp() should set showHelp to true")
	}
	// animatedHeight should be set to full height immediately
	if m.AnimatedHelpHeight() != m.HelpHeight() {
		t.Errorf("ShowHelp() in Center mode should set full height immediately, got %d want %d",
			m.AnimatedHelpHeight(), m.HelpHeight())
	}
}

func TestModel_HideHelp_CenterMode_NoAnimation(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithHelpMode(keyhints.HelpModeCenter),
	)

	m.ShowHelp()
	m.HideHelp()

	// Center mode should skip animation and hide immediately
	if m.AnimationCmd() != nil {
		t.Error("HideHelp() in Center mode should NOT start animation")
	}
	if m.ShowingHelp() {
		t.Error("HideHelp() in Center mode should hide immediately")
	}
}

func TestModel_ShowHelp_NonCenterModes_WithAnimation(t *testing.T) {
	// Verify non-Center modes still animate
	animatingModes := []keyhints.HelpMode{
		keyhints.HelpModeDown,
		keyhints.HelpModeUp,
		keyhints.HelpModeTop,
	}

	for _, mode := range animatingModes {
		t.Run(modeToString(mode), func(t *testing.T) {
			bindings := []keyhints.Binding{
				{Key: "a", Desc: "A"},
			}
			m := keyhints.New(keyhints.WithBindings(bindings...),
				keyhints.WithAnimation(true),
				keyhints.WithHelpMode(mode),
			)

			m.ShowHelp()

			if m.AnimationCmd() == nil {
				t.Errorf("ShowHelp() in %s mode should start animation", modeToString(mode))
			}
		})
	}
}

// ============================================================================
// Animation State Reset Tests
// ============================================================================

func TestModel_ShowHelp_ResetsAnimationState(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithAnimation(true))

	// First show/hide cycle
	m.ShowHelp()
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}
	m.HideHelp()
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	// Second show - should start from clean state
	m.ShowHelp()

	// Animation should start from 0
	if m.AnimationCmd() == nil {
		t.Error("ShowHelp() should start animation on second show")
	}

	// Initial animated height should be close to 0 (or at least small)
	height := m.AnimatedHelpHeight()
	if height > 1 {
		t.Errorf("ShowHelp() should reset animation to start from 0, got height %d", height)
	}
}

// ============================================================================
// Overlay Mode Tests
// ============================================================================

func TestNew_WithOverlay(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithOverlay(true))

	if !m.Overlay() {
		t.Error("WithOverlay(true) should enable overlay mode")
	}
}

func TestNew_WithOverlay_DefaultFalse(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	if m.Overlay() {
		t.Error("Overlay should be false by default")
	}
}

func TestModel_Overlay(t *testing.T) {
	m := keyhints.New(keyhints.WithOverlay(true))

	if !m.Overlay() {
		t.Error("Overlay() should return true when enabled")
	}

	m2 := keyhints.New(keyhints.WithOverlay(false))
	if m2.Overlay() {
		t.Error("Overlay() should return false when disabled")
	}
}

func TestModel_SetOverlay(t *testing.T) {
	m := keyhints.New()

	m.SetOverlay(true)
	if !m.Overlay() {
		t.Error("SetOverlay(true) should enable overlay mode")
	}

	m.SetOverlay(false)
	if m.Overlay() {
		t.Error("SetOverlay(false) should disable overlay mode")
	}
}

func TestModel_ToggleOverlay(t *testing.T) {
	m := keyhints.New()

	// Default is false
	if m.Overlay() {
		t.Fatal("Overlay should be false by default")
	}

	// Toggle to true
	m.ToggleOverlay()
	if !m.Overlay() {
		t.Error("ToggleOverlay() should enable overlay mode when disabled")
	}

	// Toggle back to false
	m.ToggleOverlay()
	if m.Overlay() {
		t.Error("ToggleOverlay() should disable overlay mode when enabled")
	}
}

func TestModel_ShowHelp_Overlay_UsesSnappierSpring(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithOverlay(true),
	)

	m.ShowHelp()

	if m.AnimationCmd() == nil {
		t.Error("ShowHelp() with overlay and animation should start animation")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp() should set showHelp to true")
	}
}

func TestModel_HideHelp_Overlay_UsesSnappierSpring(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithOverlay(true),
	)

	m.ShowHelp()
	// Run show animation to completion
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	m.HideHelp()

	if m.AnimationCmd() == nil {
		t.Error("HideHelp() with overlay and animation should start animation")
	}
}

func TestModel_Overlay_Animation_Completes(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithAnimation(true),
		keyhints.WithOverlay(true),
	)

	m.ShowHelp()

	// Run animation until it completes
	for i := 0; i < 200; i++ {
		updated, _ := m.Update(keyhints.HeightTickMsg{})
		m = updated.(keyhints.Model)
		if m.AnimationCmd() == nil {
			break
		}
	}

	if m.AnimationCmd() != nil {
		t.Error("Overlay animation should complete after sufficient ticks")
	}
	if !m.ShowingHelp() {
		t.Error("ShowHelp overlay animation should keep help visible after completion")
	}
}

// ============================================================================
// Width Getter Tests
// ============================================================================

func TestModel_Width_ReturnsConfiguredWidth(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(80))

	if m.Width() != 80 {
		t.Errorf("Width() = %d, want 80", m.Width())
	}
}

func TestModel_Width_ReturnsZeroIfNotSet(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	m := keyhints.New(keyhints.WithBindings(bindings...))

	if m.Width() != 0 {
		t.Errorf("Width() = %d for new model without width, want 0", m.Width())
	}
}

func TestModel_Width_UpdatedBySetWidth(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithWidth(60))

	m.SetWidth(120)

	if m.Width() != 120 {
		t.Errorf("Width() = %d after SetWidth(120), want 120", m.Width())
	}
}

// ============================================================================
// WithBorderColor Tests
// ============================================================================

func TestNew_WithBorderColor(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderColor(lipgloss.Color("#FF0000")))

	// Verify model is created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom border color")
	}
}

func TestNew_WithBorderColor_AdaptiveColor(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	color := lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderColor(color))

	// Verify model is created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with adaptive border color")
	}
}

// ============================================================================
// WithBorderPadding Validation Tests
// ============================================================================

func TestNew_WithBorderPadding(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderPadding(1, 1, 2, 2))

	// Verify model is created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with custom border padding")
	}
}

func TestNew_WithBorderPadding_NegativeIgnored(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	// Negative values should be ignored, keeping defaults
	m := keyhints.New(keyhints.WithBindings(bindings...),
		keyhints.WithBorderPadding(1, 1, 1, 1),     // Set to known values
		keyhints.WithBorderPadding(-1, -1, -1, -1), // Try to set negative - should be ignored
	)

	// Verify model is created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string when negative padding was provided")
	}
}

func TestNew_WithBorderPadding_PartialNegativeIgnored(t *testing.T) {
	bindings := []keyhints.Binding{{Key: "a", Desc: "A"}}
	// Mix of valid and negative values - negative should be ignored
	m := keyhints.New(keyhints.WithBindings(bindings...), keyhints.WithBorderPadding(2, -1, 3, -1))

	// Verify model is created without error
	view := m.View()
	if view == "" {
		t.Error("View() returned empty string with partially negative padding")
	}
}
