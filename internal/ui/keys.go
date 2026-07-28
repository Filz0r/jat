package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

// Keybindings for the whole app, defined once here. Each binding is used both
// for behaviour (key.Matches in Update) and for the hint bar below the border
// (help.ShortHelpView), so keys and their hints can never drift apart.
var (
	// quitHard is the only quit binding active in typing contexts; it is
	// deliberately left out of the hint bar.
	quitHard = key.NewBinding(key.WithKeys("ctrl+c"))
	quitKey  = key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))

	tabNext = key.NewBinding(
		key.WithKeys("tab", "right", "l"),
		key.WithHelp("l/→/tab", "next tab"),
	)
	tabPrev = key.NewBinding(
		key.WithKeys("shift+tab", "left", "h"),
		key.WithHelp("h/←", "prev tab"),
	)

	navUp   = key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "up"))
	navDown = key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "down"))

	enterEdit  = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit"))
	confirmKey = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
	escKey     = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))

	// lockedClose is display-only: the root model already closes locked and
	// not-implemented modals on both keys, this binding just says so.
	lockedClose = key.NewBinding(
		key.WithKeys("esc", "enter"),
		key.WithHelp("esc/enter", "close"),
	)
	fieldNext = key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field"))

	newKey     = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new"))
	refreshKey = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh"))
)

// globalHelp is appended to every tab's hint bar by the root model.
func globalHelp() []key.Binding {
	return []key.Binding{tabPrev, tabNext, quitKey}
}

// modalHint renders a hint row for the edit modal and setup overlay from the
// shared bindings, using the same help styling as the global hint bar so
// hints look identical everywhere.
func modalHint(bindings ...key.Binding) string {
	h := help.New()
	h.SetWidth(60)
	return h.ShortHelpView(bindings)
}
