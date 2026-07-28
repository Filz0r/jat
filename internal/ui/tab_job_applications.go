package ui

import tea "charm.land/bubbletea/v2"

// JobApplicationsTab is the Job Applications region — the default tab. Same
// list-navigation structure as the other tabs; renders only its name for now.
type JobApplicationsTab struct {
	items    []string
	selected int
}

func NewJobApplicationsTab() JobApplicationsTab { return JobApplicationsTab{} }

func (t JobApplicationsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
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
		// open the selected application — not implemented yet
	}
	return t, nil
}

func (t JobApplicationsTab) View() string { return "Job Applications" }