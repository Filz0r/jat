package ui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/filz0r/jat/internal/utils"
)

type modalKind int

const (
	modalText modalKind = iota
	modalList
	modalNotImplemented
	modalLocked
	modalForm
	modalDetail
	modalApplication
)

// closeModalMsg closes the modal without refreshing any tab data. Composite
// modals (form/detail/application) emit it on their esc paths since they own
// their key handling.
type closeModalMsg struct{}

type editRequest struct {
	title      string
	kind       modalKind
	input      textinput.Model
	choices    []string
	selected   int
	message    string
	commit     func(value string, selected int) tea.Cmd
	refreshTab int

	// modalForm: multi-field form with a single submit.
	fields     []formField
	focused    int
	commitForm func(values []formValue) tea.Cmd

	// modalDetail: read-only scrollable body. onKey receives every
	// non-scroll, non-esc key so builders can wire shortcuts (e.g. e → edit).
	detail       []string
	detailHints  []key.Binding
	detailScroll int
	onKey        func(kp tea.KeyPressMsg) tea.Cmd

	// modalApplication: the job-application view/edit/create modal, which
	// owns all of its state and key handling.
	appModal applicationModal
}

type editRequestMsg struct {
	req editRequest
}

type editResultMsg struct {
	err        error
	refreshTab int
}

// --- form fields ------------------------------------------------------------

type fieldKind int

const (
	fieldText fieldKind = iota
	fieldPicker
	fieldSuggest
)

// formField is one row of a modal form: a text input, a picker over a fixed
// choice set, or a text input with live case-insensitive suggestions.
type formField struct {
	label string
	kind  fieldKind

	input textinput.Model

	// fieldPicker
	choices  []utils.SuggestionRecord
	selected int

	// fieldSuggest
	suggestions []utils.SuggestionRecord
	suggestSel  int
}

// formValue is a submitted field: the typed/selected text plus the resolved
// suggestion id when the text matches a known record (empty otherwise).
type formValue struct {
	label string
	text  string
	id    string
}

type formAction int

const (
	formNoop formAction = iota
	formAdopted
	formSubmit
	formCancel
)

const (
	// formLabelWidth aligns every form's label column, mirroring setup.go.
	formLabelWidth = 12
	// suggestMaxMatches caps the live suggestion list under a suggest field.
	suggestMaxMatches = 5
	// detailMaxRows caps scrollable read-only content inside a modal so the
	// box can never outgrow the screen.
	detailMaxRows = 12
)

func newFieldTextInput(placeholder, initial string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(initial)
	ti.SetWidth(modalInnerWidth - formLabelWidth - 4)
	return ti
}

func (f *formField) focus() {
	if f.kind != fieldPicker {
		f.input.Focus()
	}
}

func (f *formField) blur() { f.input.Blur() }

// matches filters the field's suggestions by the typed text,
// case-insensitively, capped at suggestMaxMatches.
func (f formField) matches() []utils.SuggestionRecord {
	typed := strings.ToLower(strings.TrimSpace(f.input.Value()))
	out := make([]utils.SuggestionRecord, 0, suggestMaxMatches)
	for _, s := range f.suggestions {
		if typed == "" || strings.Contains(strings.ToLower(s.Label), typed) {
			out = append(out, s)
			if len(out) == suggestMaxMatches {
				break
			}
		}
	}
	return out
}

// adoptableMatch reports the highlighted suggestion when it differs from the
// typed text, so enter first completes the input before it submits the form.
func (f formField) adoptableMatch() (utils.SuggestionRecord, bool) {
	matches := f.matches()
	if len(matches) == 0 {
		return utils.SuggestionRecord{}, false
	}
	sel := f.suggestSel
	if sel >= len(matches) {
		sel = len(matches) - 1
	}
	match := matches[sel]
	if strings.EqualFold(match.Label, strings.TrimSpace(f.input.Value())) {
		return utils.SuggestionRecord{}, false
	}
	return match, true
}

// handleFormKey processes one key against a form's fields and reports what
// the key meant: nothing, a suggestion was adopted (or text edited), the user
// submitted, or the user cancelled. The host decides how to handle
// formSubmit/formCancel, which is what lets the same machinery back the
// generic modalForm, the application modal's form/note modes and any future
// form.
func handleFormKey(
	fields []formField,
	focused int,
	kp tea.KeyPressMsg,
) ([]formField, int, tea.Cmd, formAction) {
	if len(fields) == 0 {
		return fields, focused, nil, formNoop
	}
	switch {
	case key.Matches(kp, escKey):
		return fields, focused, nil, formCancel
	case key.Matches(kp, fieldNext), key.Matches(kp, fieldPrev):
		delta := 1
		if key.Matches(kp, fieldPrev) {
			delta = -1
		}
		fields[focused].blur()
		focused = (focused + delta + len(fields)) % len(fields)
		fields[focused].focus()
		return fields, focused, nil, formNoop
	case key.Matches(kp, confirmKey):
		f := fields[focused]
		if f.kind == fieldSuggest {
			if match, ok := f.adoptableMatch(); ok {
				f.input.SetValue(match.Label)
				f.suggestSel = 0
				fields[focused] = f
				return fields, focused, nil, formAdopted
			}
		}
		return fields, focused, nil, formSubmit
	}

	f := fields[focused]
	var cmd tea.Cmd
	switch f.kind {
	case fieldPicker:
		switch {
		case key.Matches(kp, navDown):
			f.selected = navigateIndex(len(f.choices), f.selected, 1)
		case key.Matches(kp, navUp):
			f.selected = navigateIndex(len(f.choices), f.selected, -1)
		}
	case fieldSuggest:
		// j/k are real input here (company names contain them), so only the
		// arrow keys move the highlight.
		matches := f.matches()
		switch {
		case key.Matches(kp, arrowDown):
			f.suggestSel = navigateIndex(len(matches), f.suggestSel, 1)
		case key.Matches(kp, arrowUp):
			f.suggestSel = navigateIndex(len(matches), f.suggestSel, -1)
		default:
			f.input, cmd = f.input.Update(kp)
			f.suggestSel = 0
		}
	default:
		f.input, cmd = f.input.Update(kp)
	}
	fields[focused] = f
	return fields, focused, cmd, formNoop
}

// formValues collapses the fields into their submitted values.
func formValues(fields []formField) []formValue {
	out := make([]formValue, len(fields))
	for i, f := range fields {
		v := formValue{label: f.label}
		switch f.kind {
		case fieldPicker:
			if len(f.choices) > 0 && f.selected < len(f.choices) {
				v.text = f.choices[f.selected].Label
				v.id = f.choices[f.selected].Value
			}
		default:
			v.text = strings.TrimSpace(f.input.Value())
			if f.kind == fieldSuggest {
				for _, s := range f.suggestions {
					if strings.EqualFold(s.Label, v.text) {
						v.id = s.Value
						break
					}
				}
			}
		}
		out[i] = v
	}
	return out
}

// viewFormFields renders the label column + field rows shared by every form.
// The focused label is accented (setup.go precedent); pickers expand to a
// choice list while focused and collapse to the selected value otherwise;
// suggest fields show their live matches while focused.
func viewFormFields(fields []formField, focused int) string {
	rows := make([]string, 0, len(fields))
	for i, f := range fields {
		label := lipgloss.NewStyle().Width(formLabelWidth)
		if i == focused {
			label = label.Bold(true).Foreground(colorAccent)
		}
		switch {
		case f.kind == fieldPicker && i == focused && len(f.choices) > 0:
			rows = append(rows, label.Render(f.label))
			for c, choice := range f.choices {
				style := lipgloss.NewStyle().PaddingLeft(formLabelWidth)
				if c == f.selected {
					style = style.Bold(true).Underline(true)
				}
				rows = append(rows, style.Render(choice.Label))
			}
		case f.kind == fieldPicker:
			value := ""
			if len(f.choices) > 0 && f.selected < len(f.choices) {
				value = f.choices[f.selected].Label
			}
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, label.Render(f.label), value))
		default:
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, label.Render(f.label), f.input.View()))
			if f.kind == fieldSuggest && i == focused {
				for mi, match := range f.matches() {
					style := lipgloss.NewStyle().PaddingLeft(formLabelWidth)
					if mi == f.suggestSel {
						style = style.Bold(true).Underline(true)
					} else {
						style = style.Faint(true)
					}
					rows = append(rows, style.Render(match.Label))
				}
			}
		}
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

// --- the modal itself ---------------------------------------------------------

type EditModal struct {
	req editRequest
	err string
}

func newEditModal(req editRequest) EditModal {
	if req.kind == modalForm && len(req.fields) > 0 {
		req.fields[req.focused].focus()
	}
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
	case modalForm:
		if kp, ok := msg.(tea.KeyPressMsg); ok {
			fields, focused, cmd, action := handleFormKey(em.req.fields, em.req.focused, kp)
			em.req.fields, em.req.focused = fields, focused
			switch action {
			case formCancel:
				return em, func() tea.Msg { return closeModalMsg{} }
			case formSubmit:
				if em.req.commitForm != nil {
					return em, em.req.commitForm(formValues(em.req.fields))
				}
				return em, nil
			default:
				return em, cmd
			}
		}
		// non-key messages (paste) go to the focused text input
		if em.req.focused < len(em.req.fields) && em.req.fields[em.req.focused].kind != fieldPicker {
			var cmd tea.Cmd
			em.req.fields[em.req.focused].input, cmd = em.req.fields[em.req.focused].input.Update(msg)
			return em, cmd
		}
	case modalDetail:
		kp, ok := msg.(tea.KeyPressMsg)
		if !ok {
			return em, nil
		}
		switch {
		case key.Matches(kp, escKey):
			return em, func() tea.Msg { return closeModalMsg{} }
		case key.Matches(kp, navDown):
			maxScroll := max(len(em.req.detail)-detailMaxRows, 0)
			em.req.detailScroll = min(em.req.detailScroll+1, maxScroll)
		case key.Matches(kp, navUp):
			em.req.detailScroll = max(em.req.detailScroll-1, 0)
		default:
			if em.req.onKey != nil {
				return em, em.req.onKey(kp)
			}
		}
	case modalApplication:
		m, cmd := em.req.appModal.Update(msg)
		em.req.appModal = m
		return em, cmd
	}
	return em, nil
}

// modalInnerWidth is the fixed content width inside the modal box. Locking it
// keeps the box from shrink-wrapping each kind's content, which is what makes
// horizontal centering of the contents possible (vertical centering on screen
// is handled by the compositor in model.go).
const modalInnerWidth = 56

func (em EditModal) View() string {
	title := em.req.title
	if em.req.kind == modalApplication {
		title = em.req.appModal.title()
	}
	// Long titles (job postings are wordy) must not stretch the modal box.
	title = ansi.Truncate(title, modalInnerWidth, "…")
	center := lipgloss.NewStyle().Width(modalInnerWidth).Align(lipgloss.Center)
	titleLine := center.Bold(true).PaddingBottom(1).Render(title)
	// Every modal kind ends with a hint row rendered from the shared
	// keybindings, so modal hints match the global hint bar's format.
	hintRow := func(bindings ...key.Binding) string {
		return center.PaddingTop(1).Render(modalHint(bindings...))
	}
	var body string
	switch em.req.kind {
	case modalText:
		body = lipgloss.JoinVertical(
			lipgloss.Center,
			center.Render(em.req.input.View()),
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
			lipgloss.Center,
			center.Render(lipgloss.JoinVertical(lipgloss.Center, rows...)),
			hintRow(navUp, navDown, confirmKey, escKey),
		)
	case modalLocked:
		body = lipgloss.JoinVertical(
			lipgloss.Center,
			center.Foreground(colorAccent).Render(em.req.message),
			hintRow(lockedClose),
		)
	case modalNotImplemented:
		body = lipgloss.JoinVertical(
			lipgloss.Center,
			center.Render(hintStyle.Render("Not implemented yet")),
			hintRow(lockedClose),
		)
	case modalForm:
		hints := []key.Binding{fieldNext, fieldPrev, confirmKey, escKey}
		body = lipgloss.JoinVertical(
			lipgloss.Center,
			center.Render(viewFormFields(em.req.fields, em.req.focused)),
			hintRow(hints...),
		)
	case modalDetail:
		lines := scrolledLines(em.req.detail, em.req.detailScroll, detailMaxRows)
		body = lipgloss.JoinVertical(
			lipgloss.Center,
			center.Render(lipgloss.JoinVertical(lipgloss.Top, lines...)),
			hintRow(em.req.detailHints...),
		)
	case modalApplication:
		body = center.Render(em.req.appModal.view())
	}
	stack := lipgloss.JoinVertical(lipgloss.Center, titleLine, body)
	if em.err != "" {
		errLine := center.Foreground(colorDanger).PaddingTop(1).Render(em.err)
		stack = lipgloss.JoinVertical(lipgloss.Center, stack, errLine)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(stack)
}
