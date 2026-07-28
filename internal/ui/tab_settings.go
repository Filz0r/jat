package ui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type settingEntry struct {
	Label    string
	Value    func() string
	Editable func() bool
}

type SettingsTab struct {
	items    []settingEntry
	selected int
	cfg      *config.ConfigFile
	table    table.Model
	width    int
	height   int
}

func buildSettingsRows(items []settingEntry) []table.Row {
	res := make([]table.Row, 0, len(items))
	for _, entry := range items {
		label, value := entry.Label, entry.Value()
		if !entry.Editable() {
			label = faint(label)
			value = faint(value)
		}
		res = append(res, table.Row{label, value})
	}
	return res
}

func NewSettingsTab(cfg *config.ConfigFile) SettingsTab {
	items := []settingEntry{
		{Label: "Mode", Value: func() string {
			val := cfg.Mode().String()
			uppercased := cases.Title(language.AmericanEnglish)
			return uppercased.String(val)
		}, Editable: func() bool { return true }},
		{Label: "Database URL", Value: func() string {
			val := cfg.DbUri()
			if val == "" {
				return "<unset>"
			}
			return val
		}, Editable: func() bool { return cfg.IsStandalone() }},
		{Label: "Server URL", Value: func() string {
			val := cfg.ServerURL()
			if val == "" {
				return "<unset>"
			}
			return val
		}, Editable: func() bool { return cfg.IsClient() }},
	}
	columns := []table.Column{
		{Title: "Setting", Width: 18},
		{Title: "Value", Width: 32},
	}
	rows := buildSettingsRows(items)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		Foreground(colorInfo).
		PaddingTop(1).
		MarginBottom(1).
		Bold(false)
	styles.Selected = styles.Selected.
		Background(colorSurface).
		Foreground(colorAccent).
		PaddingLeft(1).
		Bold(true).PaddingChar('>')
	styles.Cell = styles.Cell.PaddingLeft(1)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(18+32),
		table.WithHeight(len(rows)+1),
		table.WithStyles(styles),
	)
	return SettingsTab{cfg: cfg, items: items, table: t}
}

func (t SettingsTab) ShortHelp() []key.Binding {
	return []key.Binding{navUp, navDown, enterEdit}
}

func (t SettingsTab) editRequestFor(label string) editRequest {
	cfg := t.cfg
	switch label {
	case "Mode":
		choices := []string{"Standalone", "Client"}
		selected := 0
		if cfg.IsClient() {
			selected = 1
		}
		return editRequest{
			title: "Mode", kind: modalList, choices: choices, selected: selected,
			commit: func(_ string, sel int) tea.Cmd {
				return func() tea.Msg {
					converted := strings.ToLower(choices[sel])
					if err := cfg.SetMode(config.JatMode(converted)); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					if err := cfg.Update(); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					return editResultMsg{refreshTab: TabSettings}
				}
			},
		}
	case "Database URL":
		if !cfg.IsStandalone() {
			return editRequest{
				title:   "Database URL",
				kind:    modalLocked,
				message: "Database URL is only editable in Standalone mode",
			}
		}
		return editRequest{
			title: "Database URL",
			kind:  modalText,
			input: newTextInput(cfg.DbUri(), "Database URL"),
			commit: func(value string, _ int) tea.Cmd {
				return func() tea.Msg {
					if err := utils.ValidateDbUri(value); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					err := cfg.SetDbUri(value)
					if err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					if err := cfg.Update(); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					return editResultMsg{refreshTab: TabSettings}
				}
			},
		}
	case "Server URL":
		if !cfg.IsClient() {
			return editRequest{
				title:   "Server URL",
				kind:    modalLocked,
				message: "Server URL is only editable in Client mode",
			}
		}
		return editRequest{
			title: "Server URL",
			kind:  modalText,
			input: newTextInput(cfg.ServerURL(), "Server URL"),
			commit: func(value string, _ int) tea.Cmd {
				return func() tea.Msg {
					if err := utils.ValidateUrl(value); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					cfg.SetServerURL(value)
					if err := cfg.Update(); err != nil {
						return editResultMsg{err: err, refreshTab: TabSettings}
					}
					return editResultMsg{refreshTab: TabSettings}
				}
			},
		}
	default:
		return editRequest{
			title: label,
			kind:  modalNotImplemented,
		}
	}
}

func (t SettingsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	if kp, ok := msg.(tea.KeyPressMsg); ok && key.Matches(kp, enterEdit) {
		idx := t.table.Cursor()
		req := t.editRequestFor(t.items[idx].Label)
		return t, func() tea.Msg { return editRequestMsg{req: req} }
	}
	var cmd tea.Cmd
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t SettingsTab) View() string {
	tableStyle := lipgloss.NewStyle().MarginLeft(2)
	return tableStyle.Render(t.table.View())
}

func (t SettingsTab) refresh() SettingsTab {
	t.table.SetRows(buildSettingsRows(t.items))
	return t
}

// Resize sizes the table to fill the available body width. The Value column
// absorbs the remaining width after the 18-char Setting column and the cell
// padding.
func (t SettingsTab) Resize(w, h int) Tab {
	if w < 30 {
		w = 30
	}
	t.width = w
	t.height = h
	t.table.SetHeight(h - 10)
	t.table.SetWidth(w)
	t.table.SetColumns([]table.Column{
		{Title: "Setting", Width: 18},
		{Title: "Value", Width: w - 30},
	})
	return t
}
