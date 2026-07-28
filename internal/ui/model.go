package ui

import (
	"errors"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/services"
	"github.com/filz0r/jat/internal/utils"

	"charm.land/lipgloss/v2"
	"github.com/filz0r/jat/internal/config"
)

const (
	ListView uint = iota
	ApplicationView
	CompanyView
	LoginView
)

type setupAccountResultMsg struct {
	err error
}

type Model struct {
	ActiveView    uint
	Config        *config.ConfigFile
	width         int
	height        int
	tabs          []Tab
	active        int
	setup         Setup
	editModal     EditModal
	editModalOpen bool
	helpBar       help.Model
}

// NewModel builds the top-level model with all tabs in bar order and the
// Job Applications tab selected by default.
func NewModel(cfg *config.ConfigFile) Model {
	return Model{
		ActiveView: LoginView,
		Config:     cfg,
		tabs: []Tab{
			NewJobApplicationsTab(),
			NewCompaniesTab(),
			NewApplicationStatusTab(cfg),
			NewSettingsTab(cfg),
		},
		active:  TabJobApplications,
		setup:   NewSetup(),
		helpBar: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	if s, ok := m.tabs[TabApplicationStatus].(ApplicationStatusTab); ok {
		s.loading = true
		m.tabs[TabApplicationStatus] = s
		return tea.Batch(s.fetch(), s.spinner.Tick)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if s, ok := m.tabs[TabApplicationStatus].(ApplicationStatusTab); ok {
			var cmd tea.Cmd
			m.tabs[TabApplicationStatus], cmd = s.Update(msg)
			return m, cmd
		}
	case tea.PasteMsg:
		if !m.Config.IsInitialized() {
			var cmd tea.Cmd
			m.setup, cmd = m.setup.Update(msg)
			return m, cmd
		}
		if m.editModalOpen {
			var cmd tea.Cmd
			m.editModal, cmd = m.editModal.Update(msg)
			return m, cmd
		}
		if len(m.tabs) > 0 {
			var cmd tea.Cmd
			m.tabs[m.active], cmd = m.tabs[m.active].Update(msg)
			return m, cmd
		}
	case tea.KeyPressMsg:
		// While the config is uninitialized the setup overlay owns input.
		if !m.Config.IsInitialized() {
			return m.updateSetup(msg)
		}
		if m.editModalOpen {
			return m.updateEditModal(msg)
		}
		switch {
		case key.Matches(msg, quitHard, quitKey):
			return m, tea.Quit
		case key.Matches(msg, tabNext):
			return switchTab(m, 1), nil
		case key.Matches(msg, tabPrev):
			return switchTab(m, -1), nil
		}
		// route everything else to the active tab
		if len(m.tabs) > 0 {
			var cmd tea.Cmd
			m.tabs[m.active], cmd = m.tabs[m.active].Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		bodyW, bodyH := bodyDims(msg.Width, msg.Height)
		m.helpBar.SetWidth(bodyW)
		for i := range m.tabs {
			m.tabs[i] = m.tabs[i].Resize(bodyW, bodyH)
		}
	case editRequestMsg:
		m.editModal = newEditModal(msg.req)
		m.editModalOpen = true
		return m, nil
	case editResultMsg:
		if msg.err != nil {
			m.editModal.err = msg.err.Error()
			return m, nil
		}
		m.editModalOpen = false
		var cmd tea.Cmd
		switch msg.refreshTab {
		case TabSettings:
			if s, ok := m.tabs[m.active].(SettingsTab); ok {
				m.tabs[TabSettings] = s.refresh()
			}
		case TabApplicationStatus:
			if s, ok := m.tabs[m.active].(ApplicationStatusTab); ok {
				s, cmd = s.refresh()
				m.tabs[TabApplicationStatus] = s
			}
		}
		return m, cmd
	case applicationsDataMsg:
		if s, ok := m.tabs[TabApplicationStatus].(ApplicationStatusTab); ok {
			var cmd tea.Cmd
			m.tabs[TabApplicationStatus], cmd = s.Update(msg)
			return m, cmd
		}
	case setupAccountResultMsg:
		if msg.err != nil {
			m.setup.err = msg.err.Error()
			m.setup.step = stepAccount
			m.setup = m.setup.focusAccountField()
			return m, nil
		}
		m.Config.SetInitialized(true)
		if err := m.Config.Update(); err != nil {
			m.setup.err = err.Error()
			m.setup.step = stepAccount
			return m, nil
		}
		if s, ok := m.tabs[TabSettings].(SettingsTab); ok {
			m.tabs[TabSettings] = s.refresh()
		}
		m.setup.step = stepDone
		return m, nil
	case database.ConnectResult:
		if msg.Err != nil {
			m.setup.err = msg.Err.Error()
			if m.Config.IsClient() {
				m.setup.step = stepServerUrl
			} else {
				m.setup.step = stepDbUrl
			}
			return m, nil
		}
		if m.Config.IsStandalone() {
			m.Config.SetDB(msg.DB)
			m.Config.Services = services.NewServiceManager(msg.DB)
			m.setup.step = stepAccount
			m.setup.err = ""
			m.setup = m.setup.focusAccountField()
		}
		return m, nil
	}

	return m, nil
}

// updateSetup handles keys while the first-run overlay is up. Everything that
// is not an explicit control key is typed into the focused input.
func (m Model) updateSetup(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, quitHard):
		return m, tea.Quit
	case key.Matches(msg, confirmKey):
		return m.setupEnter()
	}
	var cmd tea.Cmd
	m.setup, cmd = m.setup.Update(msg)
	return m, cmd
}

func (m Model) setupEnter() (tea.Model, tea.Cmd) {
	switch m.setup.step {
	case stepMode:
		choice := strings.ToLower(m.setup.modeChoices[m.setup.modeCursor])
		if err := m.Config.SetMode(config.JatMode(choice)); err != nil {
			m.setup.err = err.Error()
			return m, nil
		}
		m.setup.err = ""
		if m.Config.IsClient() {
			m.setup = m.setup.withUrlStep(stepServerUrl)
		} else {
			m.setup = m.setup.withUrlStep(stepDbUrl)
		}
		return m, nil
	case stepDbUrl:
		uri := m.setup.urlInput.Value()
		if err := utils.ValidateDbUri(uri); err != nil {
			m.setup.err = err.Error()
			return m, nil
		}
		err := m.Config.SetDbUri(uri)
		if err != nil {
			m.setup.err = err.Error()
			return m, nil
		}
		m.setup.err = ""
		m.setup.step = stepDbConnecting
		return m, func() tea.Msg {
			db, err := database.ConnectDb(uri, false)
			return database.ConnectResult{
				DB:  db,
				Err: err,
			}
		}
	case stepServerUrl:
		uri := m.setup.urlInput.Value()
		if err := utils.ValidateUrl(uri); err != nil {
			m.setup.err = err.Error()
			return m, nil
		}
		m.Config.SetServerURL(uri)
		m.setup.err = ""
		m.setup.step = stepServerConnecting
		return m, func() tea.Msg {
			// TODO: this needs to be refactored after client mode is implemented
			return database.ConnectResult{Err: errors.New("client mode is not implemented yet")}
		}
	case stepAccount:
		if m.setup.acctField < 3 {
			m.setup.acctField++
			m.setup = m.setup.focusAccountField()
			return m, nil
		}
		if m.setup.email.Value() == "" || m.setup.password.Value() == "" || m.setup.username.Value() == "" {
			m.setup.err = "all fields are required"
			return m, nil
		}
		if m.setup.password.Value() != m.setup.confirmPassword.Value() {
			m.setup.err = "passwords don't match"
			return m, nil
		}
		m.setup.err = ""
		m.setup.step = stepAccountCreating
		return m, func() tea.Msg {
			user, err := m.Config.Services.CreateUser(database.User{
				Email:    m.setup.email.Value(),
				Password: m.setup.password.Value(),
				Username: m.setup.username.Value(),
				IsAdmin:  true,
			})
			if err != nil {
				return setupAccountResultMsg{err: err}
			}
			err = m.Config.Services.CreateInitialApplicationStatus(user.ID)
			m.Config.SetUserID(user.ID.String())
			return setupAccountResultMsg{err: err}
		}
	default:
		m.setup.err = "This program is probably about to crash"
	}
	return m, nil
}

// updateEditModal handles keys while the edit modal is open. Only ctrl+c
// quits here — plain letters (including q) belong to the text input.
func (m Model) updateEditModal(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, quitHard):
		return m, tea.Quit
	case key.Matches(msg, escKey):
		m.editModalOpen = false
		return m, nil
	case key.Matches(msg, confirmKey):
		if m.editModal.req.kind == modalLocked {
			m.editModalOpen = false
			return m, nil
		}
		return m.applyEditModal()
	}
	var cmd tea.Cmd
	m.editModal, cmd = m.editModal.Update(msg)
	return m, cmd
}

func (m Model) applyEditModal() (tea.Model, tea.Cmd) {
	em := m.editModal
	if em.req.commit == nil {
		m.editModalOpen = false
		return m, nil
	}
	return m, em.req.commit(em.req.input.Value(), em.req.selected)
}

func (m Model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("")
	}
	bodyW, _ := bodyDims(m.width, m.height)
	body := ""
	if len(m.tabs) > 0 {
		body = m.tabs[m.active].View()
	}

	// The hint row lives below the border and only in normal tab mode — the
	// edit modal and setup wizard carry their own hints inside their boxes.
	var hints string
	if m.Config != nil && m.Config.IsInitialized() && !m.editModalOpen {
		bindings := make([]key.Binding, 0, 8)
		bindings = append(bindings, m.tabs[m.active].ShortHelp()...)
		bindings = append(bindings, globalHelp()...)
		hints = m.helpBar.ShortHelpView(bindings)
	}

	bg := renderLayout(m, renderTitle(bodyW), renderTabBar(m, bodyW), body, hints)

	if m.editModalOpen {
		bg = centerOverlay(bg, m.editModal.View(), m.width, m.height)
	}

	// composite the setup overlay on top while the config is uninitialized
	if m.Config != nil && !m.Config.IsInitialized() {
		bg = centerOverlay(bg, m.setup.View(), m.width, m.height)
	}

	view := tea.NewView(bg)
	view.AltScreen = true
	return view
}

// centerOverlay renders overlay centered on top of bg.
func centerOverlay(bg, overlay string, w, h int) string {
	ow := lipgloss.Width(overlay)
	oh := lipgloss.Height(overlay)
	return lipgloss.NewCompositor(
		lipgloss.NewLayer(bg).Z(0),
		lipgloss.NewLayer(overlay).X(max(0, (w-ow)/2)).Y(max(0, (h-oh)/2)).Z(1),
	).Render()
}
