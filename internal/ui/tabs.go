package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Tab is a single switchable region rendered below the tab bar. Every tab
// owns its state, reports the size it renders into, and declares the
// keybindings shown in the hint bar below the border while it is active.
type Tab interface {
	Update(tea.Msg) (Tab, tea.Cmd)
	View() string
	// ShortHelp returns the tab-specific bindings for the hint bar. The root
	// model appends the global bindings after these.
	ShortHelp() []key.Binding
	// Resize gives the tab the exact body dimensions computed by the layout.
	Resize(w, h int) Tab
}

// Tab identifiers, in the order they appear in the tab bar.
const (
	TabJobApplications = iota
	TabCompanies
	TabApplicationStatus
	TabSettings
)

// tabLabels is the order the tabs appear left-to-right in the bar.
var tabLabels = []string{
	"Job Applications",
	"Companies",
	"Application Status",
	"Settings",
}

// renderTabBar renders the row of tab labels centered on a single line; the
// active tab is highlighted, the rest are plain.
func renderTabBar(m Model, width int) string {
	labels := make([]string, len(tabLabels))
	for i, name := range tabLabels {
		style := inactiveTabStyle
		if i == m.active {
			style = activeTabStyle
		}
		labels[i] = style.Render(name)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, labels...)
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(row)
}

// switchTab advances the active tab by delta, wrapping around the ends.
func switchTab(m Model, delta int) Model {
	n := len(m.tabs)
	if n == 0 {
		return m
	}
	m.active = (m.active + delta) % n
	if m.active < 0 {
		m.active += n
	}
	return m
}
