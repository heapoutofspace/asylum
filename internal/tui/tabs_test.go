package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTabsModelNavigation(t *testing.T) {
	tabs := []Tab{
		{Title: "Kits", Kind: StepMultiSelect, Options: []Option{{Label: "A"}, {Label: "B"}}, DefaultSel: []int{0}},
		{Title: "Creds", Kind: StepMultiSelect, Options: []Option{{Label: "X"}}, DefaultSel: []int{}},
		{Title: "Isolation", Kind: StepSelect, Options: []Option{{Label: "Shared"}, {Label: "Isolated"}}, DefaultIdx: 1},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	if m.active != 0 {
		t.Fatalf("expected active=0, got %d", m.active)
	}

	// Navigate right
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)
	if m.active != 1 {
		t.Fatalf("expected active=1 after right, got %d", m.active)
	}

	// Navigate right again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)
	if m.active != 2 {
		t.Fatalf("expected active=2 after right, got %d", m.active)
	}

	// Navigate right at edge — should not wrap
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)
	if m.active != 2 {
		t.Fatalf("expected active=2 (no wrap), got %d", m.active)
	}

	// Navigate left
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)
	if m.active != 1 {
		t.Fatalf("expected active=1 after left, got %d", m.active)
	}

	// Navigate left to first
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)
	if m.active != 0 {
		t.Fatalf("expected active=0 after left, got %d", m.active)
	}

	// Navigate left at edge — should not wrap
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)
	if m.active != 0 {
		t.Fatalf("expected active=0 (no wrap), got %d", m.active)
	}
}

func TestTabsModelSavesStateOnSwitch(t *testing.T) {
	tabs := []Tab{
		{Title: "Kits", Kind: StepMultiSelect, Options: []Option{{Label: "A"}, {Label: "B"}}, DefaultSel: []int{0}},
		{Title: "Isolation", Kind: StepSelect, Options: []Option{{Label: "S"}, {Label: "I"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	// Toggle option B (index 1)
	m.multiModel.cursor = 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(tabsModel)

	// Switch to tab 1 — should save tab 0 state
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)

	// Tab 0 should have both A and B selected
	if len(m.results[0].MultiIdx) != 2 {
		t.Errorf("expected 2 selections saved for tab 0, got %d: %v", len(m.results[0].MultiIdx), m.results[0].MultiIdx)
	}
}

func TestTabsModelRestoresStateOnSwitch(t *testing.T) {
	tabs := []Tab{
		{Title: "Tab1", Kind: StepMultiSelect, Options: []Option{{Label: "A"}, {Label: "B"}}, DefaultSel: []int{0}},
		{Title: "Tab2", Kind: StepSelect, Options: []Option{{Label: "X"}, {Label: "Y"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	// Toggle B on
	m.multiModel.cursor = 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(tabsModel)

	// Switch to tab 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)

	// Move cursor to Y
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(tabsModel)

	// Switch back to tab 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)

	// Results should have both selections saved
	if len(m.results[0].MultiIdx) != 2 {
		t.Errorf("expected saved state with 2 selections, got %v", m.results[0].MultiIdx)
	}

	// Visual model should also reflect the saved selections
	if !m.multiModel.selected[0] {
		t.Error("expected option A (index 0) to be selected in visual model")
	}
	if !m.multiModel.selected[1] {
		t.Error("expected option B (index 1) to be selected in visual model")
	}
}

func TestTabsModelEnterConfirmsAll(t *testing.T) {
	tabs := []Tab{
		{Title: "Tab1", Kind: StepMultiSelect, Options: []Option{{Label: "A"}}, DefaultSel: []int{0}},
		{Title: "Tab2", Kind: StepSelect, Options: []Option{{Label: "X"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tabsModel)

	if cmd == nil {
		t.Fatal("expected quit command on enter")
	}

	for i, r := range m.results {
		if !r.Completed {
			t.Errorf("tab %d should be marked completed", i)
		}
	}
}

func TestTabsModelEscCancels(t *testing.T) {
	tabs := []Tab{
		{Title: "Tab1", Kind: StepSelect, Options: []Option{{Label: "A"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(tabsModel)

	if cmd == nil {
		t.Fatal("expected quit command on escape")
	}
	if !m.cancelled {
		t.Error("expected cancelled=true on escape")
	}
}

func TestTabsModelEmptySelectionPreserved(t *testing.T) {
	tabs := []Tab{
		{Title: "Kits", Kind: StepMultiSelect, Options: []Option{{Label: "A"}, {Label: "B"}}, DefaultSel: []int{0, 1}},
		{Title: "Other", Kind: StepSelect, Options: []Option{{Label: "X"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	// Deselect A (index 0)
	m.multiModel.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(tabsModel)

	// Deselect B (index 1)
	m.multiModel.cursor = 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(tabsModel)

	// Switch to tab 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)

	// Switch back to tab 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)

	// Should remain empty, not revert to defaults
	if m.multiModel.selected[0] {
		t.Error("option A should remain deselected after round-trip")
	}
	if m.multiModel.selected[1] {
		t.Error("option B should remain deselected after round-trip")
	}
	if len(m.results[0].MultiIdx) != 0 {
		t.Errorf("expected empty MultiIdx, got %v", m.results[0].MultiIdx)
	}
}

func TestTabsModelSelectCursorPreserved(t *testing.T) {
	tabs := []Tab{
		{Title: "Multi", Kind: StepMultiSelect, Options: []Option{{Label: "A"}}, DefaultSel: []int{0}},
		{Title: "Select", Kind: StepSelect, Options: []Option{{Label: "X"}, {Label: "Y"}, {Label: "Z"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	// Switch to select tab
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)

	// Move cursor to Z (index 2)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(tabsModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(tabsModel)

	// Switch to tab 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(tabsModel)

	// Switch back to select tab
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(tabsModel)

	// Cursor should be at index 2, not default 0
	if m.selModel.cursor != 2 {
		t.Errorf("expected select cursor at 2, got %d", m.selModel.cursor)
	}
}

func TestTabsView(t *testing.T) {
	tabs := []Tab{
		{Title: "Kits", Kind: StepMultiSelect, Options: []Option{{Label: "A"}}, DefaultSel: []int{0}},
		{Title: "Creds", Kind: StepSelect, Options: []Option{{Label: "X"}}, DefaultIdx: 0},
	}

	results := make([]TabResult, len(tabs))
	m := tabsModel{tabs: tabs, results: results}
	m.initTab(0)

	view := m.View()
	if !containsPlainText(view, "Kits") {
		t.Error("view should contain tab title 'Kits'")
	}
	if !containsPlainText(view, "Creds") {
		t.Error("view should contain tab title 'Creds'")
	}
	if !containsPlainText(view, "switch tab") {
		t.Error("view should contain tab hint")
	}
}
