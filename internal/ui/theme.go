package ui

import "charm.land/lipgloss/v2"

// Palette — named lipgloss color constants only, never hex/ANSI numbers.
var (
	colorAccent    = lipgloss.Yellow        // active tab, selected values, accent text
	colorHighlight = lipgloss.BrightMagenta // selected rows in lists
	colorInfo      = lipgloss.Blue          // headers
	colorSurface   = lipgloss.BrightBlack   // selected-row background, dim chrome
	colorDanger    = lipgloss.Red           // errors
)

// Shared styles used across views so every region looks the same.
var (
	hintStyle  = lipgloss.NewStyle().Faint(true)
	errorStyle = lipgloss.NewStyle().Foreground(colorDanger)

	activeTabStyle = lipgloss.NewStyle().
			PaddingLeft(1).PaddingRight(1).
			Background(colorAccent).Foreground(lipgloss.Black).
			Bold(true).Underline(true)
	inactiveTabStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface)

	dividerStyle = lipgloss.NewStyle().Foreground(colorSurface)
)

// faint renders s with the terminal faint attribute. Hand-rolled escapes on
// purpose: bubbles' table measures row width from the raw string, and wrapping
// rows in a lipgloss style breaks that measurement.
func faint(s string) string { return "\x1b[2m" + s + "\x1b[22m" }
