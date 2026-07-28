package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// JobApplicationsTab is the Job Applications region — the default tab. Same
// list-navigation structure as the other tabs; renders only its name for now.
type JobApplicationsTab struct {
	items    []string
	selected int
	width    int
	height   int
}

func NewJobApplicationsTab() JobApplicationsTab { return JobApplicationsTab{} }

func (t JobApplicationsTab) ShortHelp() []key.Binding {
	return []key.Binding{navUp, navDown}
}

func (t JobApplicationsTab) Resize(w, h int) Tab {
	t.width = w
	t.height = h
	return t
}

func (t JobApplicationsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
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
		// open the selected application — not implemented yet
	}
	return t, nil
}

func (t JobApplicationsTab) View() string { return "Job Applications" }
