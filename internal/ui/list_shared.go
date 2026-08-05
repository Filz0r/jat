package ui

//
//import (
//	"fmt"
//
//	"charm.land/bubbles/v2/list"
//	"charm.land/bubbles/v2/spinner"
//	tea "charm.land/bubbletea/v2"
//	"charm.land/lipgloss/v2"
//)
//
//// listDelegate provides the ItemDelegate plumbing every list in the app
//// shares. Each tab embeds it and only implements Render for its own row
//// format.
//type listDelegate struct{}
//
//func (d listDelegate) Height() int                             { return 1 }
//func (d listDelegate) Spacing() int                            { return 0 }
//func (d listDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
//
//// newStyledList builds a list.Model with the option set every tab uses: no
//// status bar, no filtering (data comes from the db), no built-in help or quit
//// keys (the global hint bar owns those), and the shared theme styles.
//func newStyledList(title string, delegate list.ItemDelegate, styles listStyles) list.Model {
//	l := list.New(nil, delegate, 20, 14)
//	l.Title = title
//	l.SetShowStatusBar(false)
//	l.SetFilteringEnabled(false)
//	l.SetShowHelp(false)
//	l.DisableQuitKeybindings()
//	l.Styles.Title = styles.title
//	l.Styles.PaginationStyle = styles.pagination
//	l.Styles.HelpStyle = styles.help
//	return l
//}
//
//func newLoadingSpinner() spinner.Model {
//	s := spinner.New()
//	s.Spinner = spinner.Dot
//	return s
//}
//
//// renderListState renders a data tab's non-list body: the spinner while
//// loading, an error with the retry hint, or the empty-state message. Every
//// data tab uses it so these states look identical across the app.
//func renderListState(
//	styles listStyles,
//	w, h int,
//	loading bool,
//	spinnerView string,
//	err error,
//	emptyMsg string,
//) string {
//	var content string
//	switch {
//	case loading:
//		content = lipgloss.JoinVertical(lipgloss.Center, spinnerView, "Loading…")
//	case err != nil:
//		content = styles.quitText.Render(fmt.Sprintf("error: %v\n\npress r to retry", err))
//	default:
//		content = styles.quitText.Render(emptyMsg)
//	}
//	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
//}
//
//// placeListView centers a populated list horizontally at the top of the body
//// region.
//func placeListView(w, h int, view string) string {
//	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Top, view)
//}
//
//// windowed returns the slice of lines visible in a maxRows window starting at
//// scroll (clamped), plus whether earlier/later rows exist so callers can
//// render scroll indicators. Used by every scrollable modal body.
//func windowed(lines []string, scroll, maxRows int) ([]string, bool, bool) {
//	if len(lines) <= maxRows {
//		return lines, false, false
//	}
//	if scroll > len(lines)-maxRows {
//		scroll = len(lines) - maxRows
//	}
//	if scroll < 0 {
//		scroll = 0
//	}
//	end := scroll + maxRows
//	return lines[scroll:end], scroll > 0, end < len(lines)
//}
//
//// scrolledLines wraps windowed with the ↑/↓ indicator rows every scrollable
//// modal body shows.
//func scrolledLines(lines []string, scroll, maxRows int) []string {
//	visible, moreUp, moreDown := windowed(lines, scroll, maxRows)
//	out := make([]string, 0, len(visible)+2)
//	if moreUp {
//		out = append(out, faint("↑ …"))
//	}
//	out = append(out, visible...)
//	if moreDown {
//		out = append(out, faint("↓ …"))
//	}
//	return out
//}
