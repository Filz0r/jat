package ui

//
//import (
//	"fmt"
//	"strconv"
//	"strings"
//
//	"charm.land/bubbles/v2/key"
//	tea "charm.land/bubbletea/v2"
//	"charm.land/lipgloss/v2"
//	"github.com/charmbracelet/x/ansi"
//	"github.com/filz0r/jat/internal/config"
//	"github.com/filz0r/jat/internal/database"
//	"github.com/filz0r/jat/internal/services"
//	"github.com/filz0r/jat/internal/utils"
//)
//
//type appMode int
//
//const (
//	appModeView appMode = iota
//	appModeForm
//	appModeNote
//)
//
//// Field indices inside the application form's fields slice.
//const (
//	appFieldTitle = iota
//	appFieldUrl
//	appFieldStatus
//	appFieldCompany
//)
//
//// appNotesMsg carries the freshly reloaded notes after an add-note submit.
//// The application modal stays open when it arrives, switching back to view
//// mode with the new note visible.
//type appNotesMsg struct {
//	notes []database.ApplicationNote
//	err   error
//}
//
//// applicationModal is the single modal behind every job-application
//// interaction: view mode shows all data plus notes and the status history,
//// form mode handles both create and edit (create is just an empty app), and
//// note mode appends a note without leaving the modal.
//type applicationModal struct {
//	cfg       *config.ConfigFile
//	app       database.JobApplication
//	create    bool
//	mode      appMode
//	notes     []database.ApplicationNote
//	history   []database.StatusHistory
//	statuses  []utils.SuggestionRecord
//	companies []utils.SuggestionRecord
//	fields    []formField
//	focused   int
//	scroll    int
//	err       string
//}
//
//// newApplicationModalRequest builds the editRequest that opens the
//// application modal. All data (notes, history, statuses, companies) is
//// preloaded by the caller's command; an empty app means create mode.
//func newApplicationModalRequest(
//	cfg *config.ConfigFile,
//	app database.JobApplication,
//	notes []database.ApplicationNote,
//	history []database.StatusHistory,
//	statuses, companies []utils.SuggestionRecord,
//	mode appMode,
//) editRequest {
//	m := applicationModal{
//		cfg:       cfg,
//		app:       app,
//		create:    app.ID == 0,
//		notes:     notes,
//		history:   history,
//		statuses:  statuses,
//		companies: companies,
//	}
//	if mode == appModeView {
//		m.mode = appModeView
//	} else {
//		m = m.withFormFields()
//	}
//	return editRequest{kind: modalApplication, appModal: m, refreshTab: TabJobApplications}
//}
//
//func (a applicationModal) title() string {
//	switch a.mode {
//	case appModeForm:
//		if a.create {
//			return "New Application"
//		}
//		return "Edit Application"
//	case appModeNote:
//		return "Add Note"
//	default:
//		if a.app.Title != "" {
//			return a.app.Title
//		}
//		return "Application"
//	}
//}
//
//// withFormFields builds the create/edit form. Edit pre-fills from the
//// application; create pre-selects the default status and leaves inputs blank.
//func (a applicationModal) withFormFields() applicationModal {
//	statusSelected := 0
//	for i, s := range a.statuses {
//		if a.create {
//			if strings.EqualFold(s.Label, services.DefaultApplicationStatus) {
//				statusSelected = i
//				break
//			}
//		} else if s.Value == strconv.Itoa(int(a.app.StatusID)) {
//			statusSelected = i
//			break
//		}
//	}
//	companyName := ""
//	if len(a.app.Companies) > 0 {
//		companyName = a.app.Companies[0].Name
//	}
//
//	a.mode = appModeForm
//	a.fields = []formField{
//		{label: "Title", kind: fieldText, input: newFieldTextInput("Senior Go Developer", a.app.Title)},
//		{label: "Url", kind: fieldText, input: newFieldTextInput("https://example.com/jobs/1", a.app.Url)},
//		{label: "Status", kind: fieldPicker, choices: a.statuses, selected: statusSelected},
//		{label: "Company", kind: fieldSuggest, input: newFieldTextInput("Acme Inc", companyName), suggestions: a.companies},
//	}
//	a.focused = 0
//	a.fields[a.focused].focus()
//	return a
//}
//
//func (a applicationModal) withNoteField() applicationModal {
//	a.mode = appModeNote
//	a.fields = []formField{
//		{label: "Note", kind: fieldText, input: newFieldTextInput("Recruiter call scheduled…", "")},
//	}
//	a.focused = 0
//	a.fields[a.focused].focus()
//	return a
//}
//
//func (a applicationModal) Update(msg tea.Msg) (applicationModal, tea.Cmd) {
//	if noteMsg, ok := msg.(appNotesMsg); ok {
//		if noteMsg.err != nil {
//			a.err = noteMsg.err.Error()
//			return a, nil
//		}
//		a.err = ""
//		a.notes = noteMsg.notes
//		a.mode = appModeView
//		a.scroll = max(len(a.viewLines())-detailMaxRows, 0)
//		return a, nil
//	}
//
//	kp, ok := msg.(tea.KeyPressMsg)
//	if !ok {
//		// paste and friends go to the focused text input when a form is up
//		if a.mode != appModeView && a.focused < len(a.fields) &&
//			a.fields[a.focused].kind != fieldPicker {
//			var cmd tea.Cmd
//			a.fields[a.focused].input, cmd = a.fields[a.focused].input.Update(msg)
//			return a, cmd
//		}
//		return a, nil
//	}
//
//	if a.mode == appModeView {
//		switch {
//		case key.Matches(kp, escKey):
//			return a, func() tea.Msg { return closeModalMsg{} }
//		case key.Matches(kp, editKey):
//			a = a.withFormFields()
//			a.err = ""
//		case key.Matches(kp, addNoteKey):
//			a = a.withNoteField()
//			a.err = ""
//		case key.Matches(kp, navDown):
//			maxScroll := max(len(a.viewLines())-detailMaxRows, 0)
//			a.scroll = min(a.scroll+1, maxScroll)
//		case key.Matches(kp, navUp):
//			a.scroll = max(a.scroll-1, 0)
//		}
//		return a, nil
//	}
//
//	fields, focused, cmd, action := handleFormKey(a.fields, a.focused, kp)
//	a.fields, a.focused = fields, focused
//	switch action {
//	case formCancel:
//		a.err = ""
//		if a.create {
//			return a, func() tea.Msg { return closeModalMsg{} }
//		}
//		a.mode = appModeView
//		return a, nil
//	case formSubmit:
//		if a.mode == appModeNote {
//			return a.submitNote()
//		}
//		return a.submitForm()
//	default:
//		return a, cmd
//	}
//}
//
//func (a applicationModal) submitNote() (applicationModal, tea.Cmd) {
//	body := strings.TrimSpace(a.fields[0].input.Value())
//	if body == "" {
//		a.err = "note cannot be empty"
//		return a, nil
//	}
//	cfg, appID := a.cfg, a.app.ID
//	a.err = ""
//	return a, func() tea.Msg {
//		userID, err := parseConfigUserID(cfg)
//		if err != nil {
//			return appNotesMsg{err: err}
//		}
//		if _, err := cfg.Services.AddNote(userID, appID, body, nil); err != nil {
//			return appNotesMsg{err: err}
//		}
//		notes, err := cfg.Services.GetNotes(userID, appID)
//		if err != nil {
//			return appNotesMsg{err: err}
//		}
//		return appNotesMsg{notes: notes}
//	}
//}
//
//func (a applicationModal) submitForm() (applicationModal, tea.Cmd) {
//	values := formValues(a.fields)
//	title := values[appFieldTitle].text
//	url := values[appFieldUrl].text
//	company := values[appFieldCompany].text
//	switch {
//	case title == "":
//		a.err = "title is required"
//		return a, nil
//	case company == "":
//		a.err = "company is required"
//		return a, nil
//	}
//	if url != "" {
//		if err := utils.ValidateUrl(url); err != nil {
//			a.err = err.Error()
//			return a, nil
//		}
//	}
//	var statusID *uint
//	if parsed, err := strconv.Atoi(values[appFieldStatus].id); err == nil {
//		id := uint(parsed)
//		statusID = &id
//	}
//
//	cfg, app, create := a.cfg, a.app, a.create
//	a.err = ""
//	return a, func() tea.Msg {
//		userID, err := parseConfigUserID(cfg)
//		if err != nil {
//			return editResultMsg{err: err}
//		}
//		if create {
//			_, err = cfg.Services.CreateJobApplication(userID, title, url, company, statusID)
//		} else {
//			sid := app.StatusID
//			if statusID != nil {
//				sid = *statusID
//			}
//			_, err = cfg.Services.UpdateJobApplication(userID, app.ID, title, url, company, sid)
//		}
//		if err != nil {
//			return editResultMsg{err: err}
//		}
//		return editResultMsg{refreshTab: TabJobApplications}
//	}
//}
//
//// viewLines renders the full read-only body of view mode as one flat,
//// scrollable line buffer.
//func (a applicationModal) viewLines() []string {
//	lines := []string{
//		fmt.Sprintf("Company   %s", orDash(companyNames(a.app.Companies))),
//		fmt.Sprintf("Url       %s", orDash(a.app.Url)),
//		fmt.Sprintf("Status    %s", orDash(a.app.Status.Status)),
//		"",
//		"Notes",
//	}
//	if len(a.notes) == 0 {
//		lines = append(lines, faint("no notes yet"))
//	}
//	for _, n := range a.notes {
//		lines = append(lines, fmt.Sprintf("- %s", n.Body))
//	}
//	lines = append(lines, "", "Status history")
//	if len(a.history) == 0 {
//		lines = append(lines, faint("no history yet"))
//	}
//	for i, h := range a.history {
//		lines = append(lines, fmt.Sprintf(
//			"%d. %s — %s", i+1, h.Status.Status, h.CreatedAt.Format("2006-01-02 15:04"),
//		))
//	}
//	for i := range lines {
//		lines[i] = ansi.Truncate(lines[i], modalInnerWidth, "…")
//	}
//	return lines
//}
//
//func (a applicationModal) view() string {
//	width := lipgloss.NewStyle().Width(modalInnerWidth)
//	hintRow := func(bindings ...key.Binding) string {
//		return lipgloss.NewStyle().PaddingTop(1).Render(modalHint(bindings...))
//	}
//
//	var stack string
//	if a.mode == appModeView {
//		lines := scrolledLines(a.viewLines(), a.scroll, detailMaxRows)
//		hints := []key.Binding{editKey, addNoteKey, escKey}
//		if len(a.viewLines()) > detailMaxRows {
//			hints = []key.Binding{editKey, addNoteKey, navUp, navDown, escKey}
//		}
//		stack = lipgloss.JoinVertical(
//			lipgloss.Top,
//			lipgloss.JoinVertical(lipgloss.Top, lines...),
//			hintRow(hints...),
//		)
//	} else {
//		stack = lipgloss.JoinVertical(
//			lipgloss.Top,
//			viewFormFields(a.fields, a.focused),
//			hintRow(fieldNext, fieldPrev, confirmKey, escKey),
//		)
//	}
//	if a.err != "" {
//		stack = lipgloss.JoinVertical(
//			lipgloss.Top, stack,
//			errorStyle.PaddingTop(1).Render(a.err),
//		)
//	}
//	return width.Render(stack)
//}
//
//// companyNames joins an application's company names for one-line display.
//func companyNames(companies []database.Company) string {
//	names := make([]string, 0, len(companies))
//	for _, c := range companies {
//		names = append(names, c.Name)
//	}
//	return strings.Join(names, ", ")
//}
//
//func orDash(s string) string {
//	if s == "" {
//		return "—"
//	}
//	return s
//}
