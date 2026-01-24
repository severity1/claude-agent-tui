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
	// ViewHelp should NOT apply padding - the box has its own styling
	// and prepending padding would break the border alignment
	bindings := []keyhints.Binding{
		{Key: "a", Desc: "A"},
	}
	m := keyhints.New(bindings, keyhints.WithPadding(">>"))

	view := m.ViewHelp()
	// Should NOT start with padding - should start with box border
	if strings.HasPrefix(view, ">>") {
		t.Errorf("ViewHelp() should NOT apply padding prefix, got: %q", view[:50])
	}
	// Should start with box border character
	if !strings.HasPrefix(view, "╭") {
		t.Errorf("ViewHelp() should start with box border, got: %q", view[:20])
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
