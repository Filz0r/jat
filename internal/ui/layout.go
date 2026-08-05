package ui

//
//import (
//	"strings"
//
//	"charm.land/lipgloss/v2"
//)
//
//// Layout chrome measurements. layoutMarginX/layoutMarginY are the breathing
//// room outside the border box; the hint bar is always one reserved row
//// directly below the box (blank while an overlay is open, so the frame never
//// shifts); the border itself takes one row/column on each side.
//const (
//	layoutMarginX = 4
//	layoutMarginY = 2
//
//	hintBarHeight = 1
//	borderWidth   = 2
//
//	// tabBarHeight is a single row (labels only); dividerHeight is the dim
//	// rule separating the header (title + tab bar) from the tab content.
//	// gapAfterTitle/gapAfterDivider are blank padding rows giving the header
//	// and the content some room to breathe.
//	tabBarHeight    = 1
//	dividerHeight   = 1
//	gapAfterTitle   = 1
//	gapAfterDivider = 1
//)
//
//// bodyDims returns the size of the body region inside the bordered box: the
//// terminal minus the side margins and every fixed row of chrome (vertical
//// margins, hint bar, border, title, tab bar, divider, gaps). This is the
//// single place the layout's size math lives — renderLayout and the resize
//// path both use it.
//func bodyDims(fullW, fullH int) (int, int) {
//	w := fullW - layoutMarginX*2
//	h := fullH - layoutMarginY*2 - hintBarHeight - borderWidth -
//		titleHeight - tabBarHeight - dividerHeight - gapAfterTitle - gapAfterDivider
//	return max(w, 0), max(h, 0)
//}
//
//// renderLayout stacks the title, tab bar, divider, active tab body and the
//// hint row, wraps everything except the hint row in the bordered box, and
//// centers the whole unit in the terminal with the hint row sitting directly
//// below the bottom border. The regions are pre-rendered by their owners;
//// this function only sizes and arranges them.
//func renderLayout(m Model, title, tabBar, body, hints string) string {
//	if m.width == 0 || m.height == 0 {
//		return ""
//	}
//
//	bodyW, bodyH := bodyDims(m.width, m.height)
//	body = lipgloss.NewStyle().Width(bodyW).Height(bodyH).Render(body)
//	divider := dividerStyle.Width(bodyW).Render(strings.Repeat("─", bodyW))
//	gap := lipgloss.NewStyle().Width(bodyW).Render("")
//
//	// Width() in lipgloss v2 INCLUDES the border cells, so the box must be
//	// asked for bodyW plus the border width — otherwise the border eats two
//	// columns of interior space and every bodyW-wide region wraps by 2.
//	box := boxStyle.
//		Width(bodyW + borderWidth).
//		Render(lipgloss.JoinVertical(lipgloss.Top, title, gap, tabBar, divider, gap, body))
//
//	hintBar := lipgloss.NewStyle().
//		Width(lipgloss.Width(box)).
//		Align(lipgloss.Center).
//		Render(hints)
//
//	return lipgloss.Place(
//		m.width, m.height,
//		lipgloss.Center, lipgloss.Center,
//		lipgloss.JoinVertical(lipgloss.Left, box, hintBar),
//	)
//}
