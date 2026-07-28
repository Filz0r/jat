package ui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type modalKind int

const (
	modalText modalKind = iota
	modalList
	modalNotImplemented
	modalLocked
)

type editRequest struct {
	title      string
	kind       modalKind
	input      textinput.Model
	choices    []string
	selected   int
	message    string
	commit     func(value string, selected int) tea.Cmd
	refreshTab int
}

type editRequestMsg struct {
	req editRequest
}

type editResultMsg struct {
	err        error
	refreshTab int
}

type EditModal struct {
	req editRequest
	err string
}

func newEditModal(req editRequest) EditModal {
	return EditModal{req: req}
}

func newTextInput(initial, prompt string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt + " > "
	ti.SetValue(initial)
	ti.Focus()
	ti.SetWidth(40) // TODO: make this expandable?
	return ti
}

func (em EditModal) Update(msg tea.Msg) (EditModal, tea.Cmd) {
	switch em.req.kind {
	case modalText:
		var cmd tea.Cmd
		em.req.input, cmd = em.req.input.Update(msg)
		return em, cmd
	case modalList:
		if kp, ok := msg.(tea.KeyPressMsg); ok {
			switch {
			case key.Matches(kp, navDown):
				em.req.selected = navigateList(em.req.choices, em.req.selected, 1)
			case key.Matches(kp, navUp):
				em.req.selected = navigateList(em.req.choices, em.req.selected, -1)
			}
		}
	}
	return em, nil
}

func (em EditModal) View() string {
	title := lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render(em.req.title)
	// Every modal kind ends with a hint row rendered from the shared
	// keybindings, so modal hints match the global hint bar's format.
	hintRow := func(bindings ...key.Binding) string {
		return lipgloss.NewStyle().PaddingTop(1).Render(modalHint(bindings...))
	}
	var body string
	switch em.req.kind {
	case modalText:
		body = lipgloss.JoinVertical(
			lipgloss.Top,
			em.req.input.View(),
			hintRow(confirmKey, escKey),
		)
	case modalList:
		rows := make([]string, len(em.req.choices))
		for i, choice := range em.req.choices {
			style := lipgloss.NewStyle()
			if i == em.req.selected {
				style = style.Bold(true).Underline(true)
			}
			rows[i] = style.Render(choice)
		}
		body = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Top, rows...),
			hintRow(navUp, navDown, confirmKey, escKey),
		)
	case modalLocked:
		body = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().Foreground(colorAccent).Render(em.req.message),
			hintRow(lockedClose),
		)
	case modalNotImplemented:
		body = lipgloss.JoinVertical(
			lipgloss.Top,
			hintStyle.Render("Not implemented yet"),
			hintRow(lockedClose),
		)
	}
	stack := lipgloss.JoinVertical(lipgloss.Top, title, body)
	if em.err != "" {
		errLine := errorStyle.PaddingTop(1).Render(em.err)
		stack = lipgloss.JoinVertical(lipgloss.Top, stack, errLine)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(stack)
}
