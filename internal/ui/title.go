package ui

import (
	"charm.land/lipgloss/v2"
)

// renderTitle returns the static title bar: white text on a blue background,
// full content width, a single line, centered.
func renderTitle(width int) string {
	return lipgloss.NewStyle().
		//Background(lipgloss.Blue).
		//Foreground(lipgloss.White).
		//BorderBackground(lipgloss.Blue).
		//Margin(1, 0).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderBottom(true).
		Width(width - 10).
		Align(lipgloss.Center).
		Render("Job Application Tracker")
}
