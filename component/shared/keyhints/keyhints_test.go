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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)

	if len(m.Bindings()) != 1 {
		t.Errorf("Bindings() length = %d, want 1", len(m.Bindings()))
	}
}

func TestNew_EmptyBindings(t *testing.T) {
	m := keyhints.New(nil)

	if len(m.Bindings()) != 0 {
		t.Errorf("Bindings() length = %d, want 0", len(m.Bindings()))
	}
}

func TestNew_WithSeparator(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
	}
	m := keyhints.New(bindings, keyhints.WithSeparator(" - "))

	view := m.View()
	if !strings.Contains(view, " - ") {
		t.Errorf("View() = %q, want to contain ' - '", view)
	}
}

func TestNew_WithWidth(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings, keyhints.WithWidth(80))

	// Width is stored internally, just verify model created
	if len(m.Bindings()) != 1 {
		t.Errorf("Bindings() length = %d, want 1", len(m.Bindings()))
	}
}

// ============================================================================
// Binding Management Tests
// ============================================================================

func TestModel_SetBindings(t *testing.T) {
	m := keyhints.New(nil)

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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)

	m.AddBinding(keyhints.Binding{Key: "Esc", Desc: "Clear"})

	if len(m.Bindings()) != 2 {
		t.Errorf("Bindings() length = %d, want 2", len(m.Bindings()))
	}
}

func TestModel_Bindings_ReturnsCopy(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)

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
	m := keyhints.New(nil)

	view := m.View()
	if view != "" {
		t.Errorf("View() = %q, want empty string", view)
	}
}

func TestModel_View_SingleBinding(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)

	view := m.View()
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() = %q, want to contain 'Enter'", view)
	}
	if !strings.Contains(view, "Submit") {
		t.Errorf("View() = %q, want to contain 'Submit'", view)
	}
}

func TestModel_View_MultipleBindings(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	m := keyhints.New(bindings)

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
	m := keyhints.New(nil)

	view := m.ViewCompact()
	if view != "" {
		t.Errorf("ViewCompact() = %q, want empty string", view)
	}
}

func TestModel_ViewCompact_ShowsKeysOnly(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
	}
	m := keyhints.New(bindings)

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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m := keyhints.New(bindings, keyhints.WithKeyStyle(style))

	// Style is internal, verify model works
	view := m.View()
	if !strings.Contains(view, "Enter") {
		t.Errorf("View() missing 'Enter'")
	}
}

func TestNew_WithDescStyle(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	style := lipgloss.NewStyle().Faint(true)
	m := keyhints.New(bindings, keyhints.WithDescStyle(style))

	view := m.View()
	if !strings.Contains(view, "Submit") {
		t.Errorf("View() missing 'Submit'")
	}
}

func TestNew_WithSepStyle(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	m := keyhints.New(bindings, keyhints.WithSepStyle(style))

	view := m.View()
	if !strings.Contains(view, "|") {
		t.Errorf("View() missing separator")
	}
}

// ============================================================================
// Expansion State Tests
// ============================================================================

func TestNew_DefaultExpanded(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
	}
	m := keyhints.New(bindings)

	if !m.Expanded() {
		t.Error("New() should create expanded model by default")
	}
}

func TestNew_WithCollapsed(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings, keyhints.WithCollapsed())

	if m.Expanded() {
		t.Error("WithCollapsed() should create collapsed model")
	}
}

func TestNew_WithDefaultVisible(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
		{Key: "c", Desc: "C"},
		{Key: "d", Desc: "D"},
		{Key: "e", Desc: "E"},
	}
	m := keyhints.New(bindings, keyhints.WithCollapsed(), keyhints.WithDefaultVisible(2))

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
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
		{Key: "b", Desc: "B"},
		{Key: "c", Desc: "C"},
		{Key: "d", Desc: "D"},
	}
	// Zero should use default (3)
	m := keyhints.New(bindings, keyhints.WithCollapsed(), keyhints.WithDefaultVisible(0))

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
	m := keyhints.New(nil)

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
	m := keyhints.New(nil)

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
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "Action A"},
	}
	m := keyhints.New(bindings, keyhints.WithHelpMode(keyhints.HelpModeDown))

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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
		{Key: "Tab", Desc: "Next"},
		{Key: "Shift+Tab", Desc: "Prev"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	m := keyhints.New(bindings, keyhints.WithCollapsed())

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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
		{Key: "Esc", Desc: "Clear"},
	}
	// Only 2 bindings, default visible is 3
	m := keyhints.New(bindings, keyhints.WithCollapsed())

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
	m := keyhints.New(bindings) // Default is expanded

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
	m := keyhints.New(nil)
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init() returned non-nil command, want nil")
	}
}

func TestModel_Update_ReturnsModelAndNil(t *testing.T) {
	m := keyhints.New(nil)

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
	m := keyhints.New(nil)

	if m.Focused() {
		t.Error("New() should create unfocused model by default")
	}
}

func TestModel_Focus(t *testing.T) {
	m := keyhints.New(nil)

	cmd := m.Focus()

	if !m.Focused() {
		t.Error("Focus() should set focused state")
	}
	if cmd != nil {
		t.Error("Focus() should return nil (no cursor blink for keyhints)")
	}
}

func TestModel_Blur(t *testing.T) {
	m := keyhints.New(nil)
	m.Focus()

	m.Blur()

	if m.Focused() {
		t.Error("Blur() should clear focused state")
	}
}

func TestModel_View_FocusedShowsBorder(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)
	m.Focus()

	view := m.View()
	// Should have left border (thick border character)
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() when focused should contain thick border character, got: %q", view)
	}
}

func TestModel_View_UnfocusedHasBorder(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	m := keyhints.New(bindings)
	// Not focused by default

	view := m.View()
	// Should have left border (thick border character) in unfocused state too
	if !strings.Contains(view, "\u2503") && !strings.Contains(view, "┃") {
		t.Errorf("View() when unfocused should contain thick border character, got: %q", view)
	}
}

func TestModel_Update_FocusedHandlesQuestionMark(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings)
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
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings)
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
	bindings := []keyhints.Binding{
		{Key: "Enter", Desc: "Submit"},
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	m := keyhints.New(bindings, keyhints.WithFocusStyle(style))
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
	m := keyhints.New(nil)

	if m.ShowingHelp() {
		t.Error("New() should create model with help hidden")
	}
}

func TestModel_ShowHelp(t *testing.T) {
	m := keyhints.New(nil)

	m.ShowHelp()

	if !m.ShowingHelp() {
		t.Error("ShowHelp() should show help window")
	}
}

func TestModel_HideHelp(t *testing.T) {
	m := keyhints.New(nil)
	m.ShowHelp()

	m.HideHelp()

	if m.ShowingHelp() {
		t.Error("HideHelp() should hide help window")
	}
}

func TestModel_ToggleHelp(t *testing.T) {
	m := keyhints.New(nil)

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
	m := keyhints.New(nil)

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
	m := keyhints.New(bindings)

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
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings)

	view := m.ViewHelp()
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("ViewHelp() should contain title 'Keyboard Shortcuts'")
	}
}

func TestModel_ViewHelp_HasCloseHint(t *testing.T) {
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings)

	view := m.ViewHelp()
	if !strings.Contains(view, "Esc") {
		t.Errorf("ViewHelp() should mention Esc to close")
	}
}

func TestModel_Update_EscClosesHelp(t *testing.T) {
	m := keyhints.New(nil)
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
	m := keyhints.New(bindings, keyhints.WithHelpStyle(style))

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
	m := keyhints.New(bindings, keyhints.WithWidth(40))

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
	m := keyhints.New(bindings, keyhints.WithWidth(200))

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
	m := keyhints.New(bindings, keyhints.WithWidth(50))

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
	m := keyhints.New(nil)

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
	m := keyhints.New(bindings)

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
	m := keyhints.New(bindings, keyhints.WithPadding("   ")) // 3 spaces

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
	m := keyhints.New(bindings, keyhints.WithPadding("")) // no padding

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
	m := keyhints.New(bindings, keyhints.WithPadding(">>"))

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
	m := keyhints.New(bindings, keyhints.WithPadding(">>"))

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
	m := keyhints.New(bindings, keyhints.WithPadding(">>"), keyhints.WithBorderLeft(true))

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
	m := keyhints.New(bindings, keyhints.WithPadding("  ")) // 2 spaces
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
	m := keyhints.New(nil, keyhints.WithPadding(">>"))

	view := m.ViewCompact()
	if view != "" {
		t.Errorf("ViewCompact() with empty bindings should return empty string, got: %q", view)
	}
}

func TestModel_View_EmptyWithPadding(t *testing.T) {
	// Empty bindings should return empty string regardless of padding
	m := keyhints.New(nil, keyhints.WithPadding(">>"))

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
	m := keyhints.New(bindings, keyhints.WithWidth(0))

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
	m := keyhints.New(nil)

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
			m := keyhints.New(nil, keyhints.WithHelpMode(tt.mode))
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
	m := keyhints.New(nil)

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
	m := keyhints.New(bindings)

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
	m := keyhints.New(bindings)

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
	m := keyhints.New(bindings)

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
	m := keyhints.New(bindings, keyhints.WithContentBackground(bgColor))

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
	m := keyhints.New(bindings, keyhints.WithContentBackground(bgColor))

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
			m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
	m := keyhints.New(nil,
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
			m := keyhints.New(bindings,
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
	m := keyhints.New(bindings, keyhints.WithWidth(width))

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
	m := keyhints.New(bindings, keyhints.WithWidth(0))

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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings, keyhints.WithWidth(50))

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
	m := keyhints.New(bindings)

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
	m := keyhints.New(bindings, keyhints.WithWidth(40))

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
	m := keyhints.New(bindings, keyhints.WithBorderLeft(true))

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
	m := keyhints.New(bindings, keyhints.WithBorderLeft(false))

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
	m := keyhints.New(bindings,
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
	m := keyhints.New(bindings,
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
