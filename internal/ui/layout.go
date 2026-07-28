package ui

import "charm.land/lipgloss/v2"

const layoutMargin = 8

// bodyWidth returns the width of the body region that tabs render into.
func bodyWidth(full int) int  { return full - layoutMargin }
func bodyHeight(full int) int { return full - layoutMargin }

func bodyDims(fullW, fullH int) (int, int) {
	return bodyWidth(fullW), bodyHeight(fullH)
}

// renderLayout stacks the title, tab bar, and active tab body vertically and
// wraps the result in the bordered box, sized to fill the terminal. The three
// regions are pre-rendered by their owners; this function only sizes and
// arranges them.
func renderLayout(m Model, title, tabBar, body string) string {
	if m.width == 0 {
		return ""
	}

	contentW := bodyWidth(m.width)
	contentH := m.height - layoutMargin

	//bodyH := max(contentH-2, 0)
	body = lipgloss.NewStyle().Width(contentW).Height(contentH).Render(body)

	stacked := lipgloss.JoinVertical(lipgloss.Top, title, tabBar, body)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Margin(1, 4, 1).
		//MarginRight(2).
		//MarginTop(3).
		Width(contentW).
		Height(contentH + 6).
		Render(stacked)
}
