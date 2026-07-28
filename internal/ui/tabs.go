package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Tab is a single switchable region rendered below the tab bar. Each tab owns
// its own state and behaviour; the layout only renders the active one.
type Tab interface {
	Update(tea.Msg) (Tab, tea.Cmd)
	View() string
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

// renderTabBar renders the row of tab labels; the active tab is bold and
// underlined, the rest are plain. The bar is sized to width so it spans the
// full content width.
func renderTabBar(m Model, width int) string {
	labels := make([]string, len(tabLabels))
	for i, name := range tabLabels {
		style := lipgloss.NewStyle().PaddingRight(1).PaddingLeft(1)
		if i == m.active {
			style = style.Background(lipgloss.Yellow).Bold(true).Underline(true)
		}
		labels[i] = style.Render(name)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, labels...)
	return lipgloss.NewStyle().Width(width).PaddingRight(10).PaddingTop(1).PaddingBottom(1).Align(lipgloss.Center).Render(row)
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

// navigateList moves a selection index by delta within items, clamped to the
// valid range. Shared by every tab so they all honour the same j/k/arrow
// contract.
func navigateList(items []string, selected, delta int) int {
	if len(items) == 0 {
		return 0
	}
	selected += delta
	if selected < 0 {
		selected = 0
	}
	if selected >= len(items) {
		selected = len(items) - 1
	}
	return selected
}
