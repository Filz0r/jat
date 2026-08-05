package ui

//import (
//	"errors"
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
//	"github.com/filz0r/jat/internal/utils"
//)
//
//type companiesMsg struct {
//	data []database.Company
//	err  error
//}
//
//type companyItem struct {
//	row database.Company
//}
//
//func (i companyItem) FilterValue() string { return i.row.Name }
//
//type companyDelegate struct {
//	listDelegate
//	styles listStyles
//}
//
//func (d companyDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
//	i, ok := listItem.(companyItem)
//	if !ok {
//		return
//	}
//	str := fmt.Sprintf("%d. %s — %s", index+1, i.row.Name, orDash(companyWebsite(i.row)))
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
//// companyWebsite dereferences a company's optional website for display.
//func companyWebsite(c database.Company) string {
//	if c.Website == nil {
//		return ""
//	}
//	return *c.Website
//}
//
//// CompaniesTab is the Companies region. Companies are global records; enter
//// opens a detail modal with the current user's applications at that company,
//// e jumps straight to the edit form and n creates a company.
//type CompaniesTab struct {
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
//func NewCompaniesTab(cfg *config.ConfigFile) CompaniesTab {
//	t := CompaniesTab{cfg: cfg}
//	t.styles = newListStyles(true) // dark
//	t.spinner = newLoadingSpinner()
//	t.list = newStyledList("Companies", companyDelegate{styles: t.styles}, t.styles)
//	return t
//}
//
//func (t CompaniesTab) ShortHelp() []key.Binding {
//	return []key.Binding{navUp, navDown, openKey, editKey, newKey, refreshKey}
//}
//
//func (t CompaniesTab) Resize(w, h int) Tab {
//	t.width = w
//	t.height = h
//	t.list.SetSize(w, h)
//	return t
//}
//
//func (t CompaniesTab) refresh() (CompaniesTab, tea.Cmd) {
//	t.loading = true
//	t.err = nil
//	return t, tea.Batch(t.fetch(), t.spinner.Tick)
//}
//
//func (t CompaniesTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
//	switch msg := msg.(type) {
//	case companiesMsg:
//		t.loading = false
//		if msg.err != nil {
//			t.err = msg.err
//			return t, nil
//		}
//		t.err = nil
//		items := make([]list.Item, 0, len(msg.data))
//		for _, c := range msg.data {
//			items = append(items, companyItem{row: c})
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
//			it, ok := t.list.SelectedItem().(companyItem)
//			if !ok {
//				return t, nil
//			}
//			return t, t.openDetailCmd(it.row)
//		case key.Matches(msg, editKey):
//			it, ok := t.list.SelectedItem().(companyItem)
//			if !ok {
//				return t, nil
//			}
//			req := newCompanyFormRequest(t.cfg, it.row)
//			return t, func() tea.Msg { return editRequestMsg{req: req} }
//		case key.Matches(msg, newKey):
//			req := newCompanyFormRequest(t.cfg, database.Company{})
//			return t, func() tea.Msg { return editRequestMsg{req: req} }
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
//func (t CompaniesTab) View() string {
//	if t.loading || t.err != nil || len(t.list.Items()) == 0 {
//		return renderListState(
//			t.styles, t.width, t.height,
//			t.loading, t.spinner.View(), t.err,
//			"no companies yet — press n to create one",
//		)
//	}
//	return placeListView(t.width, t.height, t.list.View())
//}
//
//func (t CompaniesTab) fetch() tea.Cmd {
//	cfg := t.cfg
//	return func() tea.Msg {
//		data, err := cfg.Services.GetAllCompanies()
//		if err != nil {
//			return companiesMsg{err: fmt.Errorf("load companies: %w", err)}
//		}
//		return companiesMsg{data: data}
//	}
//}
//
//// openDetailCmd loads the current user's applications at the selected company
//// and opens the detail modal.
//func (t CompaniesTab) openDetailCmd(company database.Company) tea.Cmd {
//	//cfg := t.cfg
//	return func() tea.Msg {
//		//userID, err := parseConfigUserID(cfg)
//		//if err != nil {
//		return companiesMsg{err: errors.New("open companies")}
//		//}
//		//apps, err := cfg.Services.GetApplicationsForCompany(userID, company.ID)
//		//if err != nil {
//		//	return companiesMsg{err: fmt.Errorf("load applications: %w", err)}
//		//}
//
//		//return editRequestMsg{req: newCompanyDetailRequest(cfg, company, apps)}
//	}
//}
//
//// newCompanyDetailRequest builds the read-only company modal: its website and
//// the current user's applications there. The e shortcut swaps the modal for
//// the edit form via the standard editRequestMsg flow.
//func newCompanyDetailRequest(
//	cfg *config.ConfigFile,
//	company database.Company,
//	apps []database.JobApplication,
//) editRequest {
//	lines := []string{
//		fmt.Sprintf("Website   %s", orDash(companyWebsite(company))),
//		"",
//		"Your applications here",
//	}
//	if len(apps) == 0 {
//		lines = append(lines, faint("none yet"))
//	}
//	for i, app := range apps {
//		lines = append(lines, ansi.Truncate(
//			fmt.Sprintf("%d. %s (%s)", i+1, app.Title, orDash(app.Status.Status)),
//			modalInnerWidth, "…",
//		))
//	}
//	hints := []key.Binding{editKey, escKey}
//	if len(lines) > detailMaxRows {
//		hints = []key.Binding{editKey, navUp, navDown, escKey}
//	}
//	return editRequest{
//		title:       company.Name,
//		kind:        modalDetail,
//		detail:      lines,
//		detailHints: hints,
//		refreshTab:  TabCompanies,
//		onKey: func(kp tea.KeyPressMsg) tea.Cmd {
//			if key.Matches(kp, editKey) {
//				req := newCompanyFormRequest(cfg, company)
//				return func() tea.Msg { return editRequestMsg{req: req} }
//			}
//			return nil
//		},
//	}
//}
//
//// newCompanyFormRequest builds the company create/edit form. An empty company
//// value means create mode.
//func newCompanyFormRequest(cfg *config.ConfigFile, company database.Company) editRequest {
//	create := company.ID == 0
//	title := "Edit Company"
//	if create {
//		title = "New Company"
//	}
//	return editRequest{
//		title: title,
//		kind:  modalForm,
//		fields: []formField{
//			{label: "Name", kind: fieldText, input: newFieldTextInput("Acme Inc", company.Name)},
//			{label: "Website", kind: fieldText, input: newFieldTextInput("https://acme.com", companyWebsite(company))},
//		},
//		commitForm: func(values []formValue) tea.Cmd {
//			name := values[0].text
//			site := values[1].text
//			return func() tea.Msg {
//				userID, err := parseConfigUserID(cfg)
//				if err != nil {
//					return editResultMsg{err: err}
//				}
//				if name == "" {
//					return editResultMsg{err: errors.New("name is required")}
//				}
//				var sitePtr *string
//				if site != "" {
//					if err := utils.ValidateUrl(site); err != nil {
//						return editResultMsg{err: err}
//					}
//					sitePtr = &site
//				}
//				if create {
//					_, err = cfg.Services.CreateCompany(userID, database.Company{Name: name, Website: sitePtr})
//				} else {
//					record := company
//					record.Name = name
//					record.Website = sitePtr
//					_, err = cfg.Services.UpdateCompany(userID, record)
//				}
//				if err != nil {
//					return editResultMsg{err: err}
//				}
//				return editResultMsg{refreshTab: TabCompanies}
//			}
//		},
//		refreshTab: TabCompanies,
//	}
//}
