package ui

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

type applicationsDataMsg struct {
	data []database.ApplicationStatus
	err  error
}

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func newStyles(darkBG bool) styles {
	var s styles
	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(colorHighlight)
	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(darkBG).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

type statusItem struct {
	row database.ApplicationStatus
}

func (i statusItem) FilterValue() string { return i.row.Status }

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(statusItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.row.Status)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// ApplicationStatusTab is the Application Status region.
type ApplicationStatusTab struct {
	err     error
	cfg     *config.ConfigFile
	loading bool
	list    list.Model
	spinner spinner.Model
	styles  styles
	width   int
	height  int
}

func NewApplicationStatusTab(cfg *config.ConfigFile) ApplicationStatusTab {
	t := ApplicationStatusTab{cfg: cfg}
	t.styles = newStyles(true) // dark
	t.spinner = spinner.New()
	t.spinner.Spinner = spinner.Dot

	l := list.New(nil, itemDelegate{styles: &t.styles}, 20, 14)
	l.Title = "Application Status"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = t.styles.title
	l.Styles.PaginationStyle = t.styles.pagination
	l.Styles.HelpStyle = t.styles.help
	t.list = l
	return t
}

func (t ApplicationStatusTab) ShortHelp() []key.Binding {
	return []key.Binding{navUp, navDown, enterEdit, newKey, refreshKey}
}

func (t ApplicationStatusTab) Resize(w, h int) Tab {
	t.width = w
	t.height = h
	t.list.SetSize(w, h)
	return t
}

func (t ApplicationStatusTab) refresh() (ApplicationStatusTab, tea.Cmd) {
	t.loading = true
	t.err = nil
	return t, tea.Batch(t.fetch(), t.spinner.Tick)
}

func (t ApplicationStatusTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case applicationsDataMsg:
		t.loading = false
		if msg.err != nil {
			t.err = msg.err
			return t, nil
		}
		t.err = nil
		items := make([]list.Item, 0, len(msg.data))
		for _, a := range msg.data {
			items = append(items, statusItem{row: a})
		}
		return t, t.list.SetItems(items)
	case spinner.TickMsg:
		if !t.loading {
			return t, nil
		}
		var cmd tea.Cmd
		t.spinner, cmd = t.spinner.Update(msg)
		return t, cmd
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, enterEdit):
			it, ok := t.list.SelectedItem().(statusItem)
			if !ok {
				return t, nil
			}
			existing := it.row
			cfg := t.cfg
			req := editRequest{
				title: "Edit Application Status",
				kind:  modalText,
				input: newTextInput(it.row.Status, "Status"),
				commit: func(value string, _ int) tea.Cmd {
					return func() tea.Msg {
						if value == "" {
							return editResultMsg{err: errors.New("status cannot be empty")}
						}
						existing.Status = value
						existing.UpdatedAt = time.Now()
						userId, err := cfg.GetUserID()
						if err != nil {
							return editResultMsg{err: err}
						}
						parsedID, err := uuid.Parse(userId)
						if err != nil {
							return editResultMsg{err: err}
						}
						_, err = cfg.Services.UpdateApplicationStatus(parsedID, existing)
						if err != nil {
							return editResultMsg{err: err}
						}
						return editResultMsg{refreshTab: TabApplicationStatus}

					}
				},
			}
			return t, func() tea.Msg { return editRequestMsg{req: req} }
		case key.Matches(msg, refreshKey):
			t.loading = true
			t.err = nil
			return t, tea.Batch(t.fetch(), t.spinner.Tick)
		case key.Matches(msg, newKey):
			cfg := t.cfg
			req := editRequest{
				title: "New Application Status",
				kind:  modalText,
				input: newTextInput("", "Status"),
				commit: func(value string, _ int) tea.Cmd {
					return func() tea.Msg {
						if value == "" {
							return editResultMsg{err: errors.New("status cannot be empty")}
						}
						id, err := cfg.GetUserID()
						if err != nil {
							return editResultMsg{err: err}
						}
						parsedID, err := uuid.Parse(id)
						if err != nil {
							return editResultMsg{err: err}
						}
						_, err = cfg.Services.CreateApplicationStatus(
							parsedID,
							database.ApplicationStatus{Status: value},
						)
						if err != nil {
							return editResultMsg{err: err}
						}
						return editResultMsg{refreshTab: TabApplicationStatus}
					}
				},
			}
			return t, func() tea.Msg { return editRequestMsg{req: req} }
		}
	}
	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg) // j/k, filter, pgup/pgdn → list
	return t, cmd
}

func (t ApplicationStatusTab) View() string {
	switch {
	case t.loading:
		content := lipgloss.JoinVertical(lipgloss.Center, t.spinner.View(), "Loading…")
		return lipgloss.Place(t.width, t.height, lipgloss.Center, lipgloss.Center, content)

	case t.err != nil:
		msg := t.styles.quitText.Render(fmt.Sprintf("error: %v\n\npress r to retry", t.err))
		return lipgloss.Place(t.width, t.height, lipgloss.Center, lipgloss.Center, msg)

	case len(t.list.Items()) == 0:
		msg := t.styles.quitText.Render("no statuses — press r to load or n to create a new one")
		return lipgloss.Place(t.width, t.height, lipgloss.Center, lipgloss.Center, msg)

	default:
		return lipgloss.Place(t.width, t.height, lipgloss.Center, lipgloss.Top, t.list.View())
	}
}

func (t ApplicationStatusTab) fetch() tea.Cmd {
	cfg := t.cfg
	return func() tea.Msg {
		id, err := cfg.GetUserID()
		if err != nil {
			return applicationsDataMsg{err: fmt.Errorf("get user id: %w", err)}
		}
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return applicationsDataMsg{err: fmt.Errorf("parse user id: %w", err)}
		}
		data, err := cfg.Services.GetAllApplicationStatus(parsedID)
		if err != nil {
			return applicationsDataMsg{err: fmt.Errorf("load statuses data: %w", err)}
		}

		return applicationsDataMsg{data: data}
	}
}
