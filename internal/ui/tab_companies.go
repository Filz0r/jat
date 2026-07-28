package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// CompaniesTab is the Companies region. Same list-navigation structure as the
// other tabs; renders only its name for now.
type CompaniesTab struct {
	items    []string
	selected int
	width    int
	height   int
}

func NewCompaniesTab() CompaniesTab { return CompaniesTab{} }

func (t CompaniesTab) ShortHelp() []key.Binding {
	return []key.Binding{navUp, navDown}
}

func (t CompaniesTab) Resize(w, h int) Tab {
	t.width = w
	t.height = h
	return t
}

func (t CompaniesTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	kp, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return t, nil
	}
	switch {
	case key.Matches(kp, navDown):
		t.selected = navigateList(t.items, t.selected, 1)
	case key.Matches(kp, navUp):
		t.selected = navigateList(t.items, t.selected, -1)
	case key.Matches(kp, confirmKey):
		// open the selected company — not implemented yet
	}
	return t, nil
}

func (t CompaniesTab) View() string { return "Companies" }
