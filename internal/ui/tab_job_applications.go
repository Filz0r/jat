package ui

//
//import (
//	"fmt"
//	"io"
//	"strings"
//
//	"charm.land/bubbles/v2/key"
//	"charm.land/bubbles/v2/list"
//	"charm.land/bubbles/v2/spinner"
//	tea "charm.land/bubbletea/v2"
//	"github.com/charmbracelet/x/ansi"
//	"github.com/filz0r/jat/internal/config"
//	"github.com/filz0r/jat/internal/database"
//)
//
//type jobApplicationsMsg struct {
//	data []database.JobApplication
//	err  error
//}
//
//type applicationItem struct {
//	row database.JobApplication
//}
//
//func (i applicationItem) FilterValue() string { return i.row.Title }
//
//type applicationDelegate struct {
//	listDelegate
//	styles listStyles
//}
//
//func (d applicationDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
//	i, ok := listItem.(applicationItem)
//	if !ok {
//		return
//	}
//	str := fmt.Sprintf(
//		"%d. %s — %s (%s)",
//		index+1,
//		i.row.Title,
//		orDash(companyNames(i.row.Companies)),
//		orDash(i.row.Status.Status),
//	)
//	str = ansi.Truncate(str, max(m.Width()-4, 1), "…")
//
//	fn := d.styles.item.Render
//	if index == m.Index() {
//		fn = func(s ...string) string {
//			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
//		}
//	}
//
//	fmt.Fprint(w, fn(str))
//}
//
//// JobApplicationsTab is the Job Applications region — the default tab. Rows
//// come from the user's own applications; enter opens the view modal, e opens
//// it straight in edit mode, n opens create mode and r reloads.
//type JobApplicationsTab struct {
//	err     error
//	cfg     *config.ConfigFile
//	loading bool
//	list    list.Model
//	spinner spinner.Model
//	styles  listStyles
//	width   int
//	height  int
//}
//
//func NewJobApplicationsTab(cfg *config.ConfigFile) JobApplicationsTab {
//	t := JobApplicationsTab{cfg: cfg}
//	t.styles = newListStyles(true) // dark
//	t.spinner = newLoadingSpinner()
//	t.list = newStyledList("Job Applications", applicationDelegate{styles: t.styles}, t.styles)
//	return t
//}
//
//func (t JobApplicationsTab) ShortHelp() []key.Binding {
//	return []key.Binding{navUp, navDown, openKey, editKey, newKey, refreshKey}
//}
//
//func (t JobApplicationsTab) Resize(w, h int) Tab {
//	t.width = w
//	t.height = h
//	t.list.SetSize(w, h)
//	return t
//}
//
//func (t JobApplicationsTab) refresh() (JobApplicationsTab, tea.Cmd) {
//	t.loading = true
//	t.err = nil
//	return t, tea.Batch(t.fetch(), t.spinner.Tick)
//}
//
//func (t JobApplicationsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
//	switch msg := msg.(type) {
//	case jobApplicationsMsg:
//		t.loading = false
//		if msg.err != nil {
//			t.err = msg.err
//			return t, nil
//		}
//		t.err = nil
//		items := make([]list.Item, 0, len(msg.data))
//		for _, a := range msg.data {
//			items = append(items, applicationItem{row: a})
//		}
//		return t, t.list.SetItems(items)
//	case spinner.TickMsg:
//		if !t.loading {
//			return t, nil
//		}
//		var cmd tea.Cmd
//		t.spinner, cmd = t.spinner.Update(msg)
//		return t, cmd
//	case tea.KeyPressMsg:
//		switch {
//		case key.Matches(msg, openKey):
//			it, ok := t.list.SelectedItem().(applicationItem)
//			if !ok {
//				return t, nil
//			}
//			return t, t.openModalCmd(it.row, appModeView)
//		case key.Matches(msg, editKey):
//			it, ok := t.list.SelectedItem().(applicationItem)
//			if !ok {
//				return t, nil
//			}
//			return t, t.openModalCmd(it.row, appModeForm)
//		case key.Matches(msg, newKey):
//			return t, t.createModalCmd()
//		case key.Matches(msg, refreshKey):
//			t.loading = true
//			t.err = nil
//			return t, tea.Batch(t.fetch(), t.spinner.Tick)
//		}
//	}
//	var cmd tea.Cmd
//	t.list, cmd = t.list.Update(msg) // j/k, pgup/pgdn → list
//	return t, cmd
//}
//
//func (t JobApplicationsTab) View() string {
//	if t.loading || t.err != nil || len(t.list.Items()) == 0 {
//		return renderListState(
//			t.styles, t.width, t.height,
//			t.loading, t.spinner.View(), t.err,
//			"no applications yet — press n to create one",
//		)
//	}
//	return placeListView(t.width, t.height, t.list.View())
//}
//
//func (t JobApplicationsTab) fetch() tea.Cmd {
//	cfg := t.cfg
//	return func() tea.Msg {
//		userID, err := parseConfigUserID(cfg)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("get user id: %w", err)}
//		}
//		data, err := cfg.Services.GetJobApplications(userID)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load job applications: %w", err)}
//		}
//		return jobApplicationsMsg{data: data}
//	}
//}
//
//// openModalCmd preloads everything the application modal needs for one row —
//// notes, status history, the user's statuses and the global companies — then
//// opens it in the requested mode.
//func (t JobApplicationsTab) openModalCmd(app database.JobApplication, mode appMode) tea.Cmd {
//	cfg := t.cfg
//	return func() tea.Msg {
//		userID, err := parseConfigUserID(cfg)
//		if err != nil {
//			return jobApplicationsMsg{err: err}
//		}
//		notes, err := cfg.Services.GetNotes(userID, app.ID)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load notes: %w", err)}
//		}
//		history, err := cfg.Services.GetStatusHistory(userID, app.ID)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load history: %w", err)}
//		}
//		statuses, err := cfg.Services.GetSuggestionForApplicationStatus(userID)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load statuses: %w", err)}
//		}
//		companies, err := cfg.Services.GetAllCompaniesSuggestions()
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load companies: %w", err)}
//		}
//		return editRequestMsg{
//			req: newApplicationModalRequest(cfg, app, notes, history, statuses, companies, mode),
//		}
//	}
//}
//
//// createModalCmd loads the pickers' data and opens the application modal in
//// create mode (an empty application).
//func (t JobApplicationsTab) createModalCmd() tea.Cmd {
//	cfg := t.cfg
//	return func() tea.Msg {
//		userID, err := parseConfigUserID(cfg)
//		if err != nil {
//			return jobApplicationsMsg{err: err}
//		}
//		statuses, err := cfg.Services.GetSuggestionForApplicationStatus(userID)
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load statuses: %w", err)}
//		}
//		companies, err := cfg.Services.GetAllCompaniesSuggestions()
//		if err != nil {
//			return jobApplicationsMsg{err: fmt.Errorf("load companies: %w", err)}
//		}
//		return editRequestMsg{
//			req: newApplicationModalRequest(
//				cfg, database.JobApplication{}, nil, nil, statuses, companies, appModeForm,
//			),
//		}
//	}
//}
