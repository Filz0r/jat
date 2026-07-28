package ui

import tea "charm.land/bubbletea/v2"

// CompaniesTab is the Companies region. Same list-navigation structure as the
// other tabs; renders only its name for now.
type CompaniesTab struct {
	items    []string
	selected int
}

func NewCompaniesTab() CompaniesTab { return CompaniesTab{} }

func (t CompaniesTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	kp, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return t, nil
	}
	switch kp.String() {
	case "j", "down":
		t.selected = navigateList(t.items, t.selected, 1)
	case "k", "up":
		t.selected = navigateList(t.items, t.selected, -1)
	case "enter":
		// open the selected company — not implemented yet
	}
	return t, nil
}

func (t CompaniesTab) View() string { return "Companies" }