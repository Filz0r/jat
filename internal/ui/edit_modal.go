package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type editSettingsMsg struct {
	Field string
}

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
	ti.Prompt = prompt + " >"
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
			switch kp.String() {
			case "j", "down":
				em.req.selected = navigateList(em.req.choices, em.req.selected, 1)
			case "k", "up":
				em.req.selected = navigateList(em.req.choices, em.req.selected, -1)
			}
		}
	}
	return em, nil
}

func (em EditModal) View() string {
	title := lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render(em.req.title)
	var body string
	switch em.req.kind {
	case modalText:
		body = em.req.input.View()
	case modalList:
		rows := make([]string, len(em.req.choices))
		for i, choice := range em.req.choices {
			style := lipgloss.NewStyle()
			if i == em.req.selected {
				style = style.Bold(true).Underline(true)
			}
			rows[i] = style.Render(choice)
		}
		hint := lipgloss.NewStyle().Faint(true).PaddingTop(1).
			Render("↑/↓ or j/k to select, Enter to confirm, Esc to cancel")
		body = lipgloss.JoinVertical(lipgloss.Top, append(rows, hint)...)
	case modalLocked:
		body = lipgloss.JoinVertical(
			lipgloss.Top,
			title,
			lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render(em.req.message),
			lipgloss.NewStyle().Faint(true).PaddingTop(1).Render("Press Esc to close"),
		)
	case modalNotImplemented:
		body = "Not implemented yet - Press Esc to close"
	}
	stack := lipgloss.JoinVertical(lipgloss.Top, title, body)
	if em.err != "" {
		errLine := lipgloss.NewStyle().Foreground(lipgloss.Red).PaddingTop(1).Render(em.err)
		stack = lipgloss.JoinVertical(lipgloss.Top, stack, errLine)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(stack)
}

//func NewEditModal(
//	cfg *config.ConfigFile,
//	field string,
//	statusTarget *statusItem,
//) EditModal {
//	em := EditModal{
//		//kind:  kind,
//		field:        field,
//		cfg:          cfg,
//		statusTarget: statusTarget,
//	}
//	switch field {
//	case "Mode":
//		em.kind = modalList
//		em.choices = []string{"Standalone", "Client"}
//		switch {
//		case cfg.IsClient():
//			em.selected = 1
//		case cfg.IsStandalone():
//			em.selected = 0
//		}
//	case "Database URL":
//		if !cfg.IsStandalone() {
//			em.kind = modalLocked
//			em.message = "Database URL is only editable in Standalone mode"
//			return em
//		}
//		em.kind = modalText
//		em.input = newTextInput(cfg.DbUri(), field)
//	case "Server URL":
//		if !cfg.IsClient() {
//			em.kind = modalLocked
//			em.message = "Server URL is only editable in Client mode"
//			return em
//		}
//		em.kind = modalText
//		em.input = newTextInput(cfg.ServerURL(), field)
//	default:
//		em.kind = modalNotImplemented
//	}
//	return em
//}
//
//func newTextInput(initial, field string) textinput.Model {
//	ti := textinput.New()
//	ti.Prompt = field + ">"
//	ti.SetValue(initial)
//	ti.Focus()
//	// change bellow
//	ti.SetWidth(40)
//	return ti
//}
//
//func (em EditModal) Update(msg tea.Msg) (EditModal, tea.Cmd) {
//	switch em.kind {
//	case modalText:
//		var cmd tea.Cmd
//		em.input, cmd = em.input.Update(msg)
//		return em, cmd
//	case modalList:
//		if kp, ok := msg.(tea.KeyPressMsg); ok {
//			switch kp.String() {
//			case "j", "down":
//				em.selected = navigateList(em.choices, em.selected, 1)
//			case "k", "up":
//				em.selected = navigateList(em.choices, em.selected, -1)
//			}
//		}
//	default:
//		return em, nil
//	}
//	return em, nil
//}
//
//func (em EditModal) View() string {
//	// Title: which setting is being edited.
//	title := lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render(em.field)
//
//	var body string
//	switch em.kind {
//	case modalText:
//		// textinput renders its own prompt + value + cursor.
//		body = em.input.View()
//
//	case modalList:
//		// One choice per line; the selected one is bold + underlined.
//		rows := make([]string, len(em.choices))
//		for i, c := range em.choices {
//			style := lipgloss.NewStyle()
//			if i == em.selected {
//				style = style.Bold(true).Underline(true)
//			}
//			rows[i] = style.Render(c)
//		}
//		hint := lipgloss.NewStyle().Faint(true).PaddingTop(1).
//			Render("↑/↓ or j/k to select, Enter to confirm, Esc to cancel")
//		body = lipgloss.JoinVertical(lipgloss.Top, append(rows, hint)...)
//	case modalLocked:
//		body = lipgloss.JoinVertical(lipgloss.Top,
//			lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render(em.message),
//			lipgloss.NewStyle().Faint(true).PaddingTop(1).Render("Press Esc to close"),
//		)
//	case modalNotImplemented:
//		body = "Not implemented yet — press Esc to close"
//	}
//
//	stack := lipgloss.JoinVertical(lipgloss.Top, title, body)
//
//	if em.err != "" {
//		errLine := lipgloss.NewStyle().Foreground(lipgloss.Red).PaddingTop(1).Render(em.err)
//		stack = lipgloss.JoinVertical(lipgloss.Top, stack, errLine)
//	}
//
//	return lipgloss.NewStyle().
//		Border(lipgloss.RoundedBorder()).
//		Padding(1, 2).
//		Render(stack)
//}
