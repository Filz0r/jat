package ui

//
//import (
//	"fmt"
//	"io"
//	"strings"
//
//	"charm.land/bubbles/v2/key"
//	"charm.land/bubbles/v2/list"
//	tea "charm.land/bubbletea/v2"
//	"charm.land/lipgloss/v2"
//	"github.com/charmbracelet/x/ansi"
//	"github.com/filz0r/jat/internal/config"
//	"github.com/filz0r/jat/internal/utils"
//	"golang.org/x/text/cases"
//	"golang.org/x/text/language"
//)
//
//type settingEntry struct {
//	Label    string
//	Value    func() string
//	Editable func() bool
//}
//
//// settingLabelWidth is the fixed width of the label column, matching what the
//// old table used for its Setting column.
//const settingLabelWidth = 18
//
//type settingItem struct {
//	entry settingEntry
//}
//
//func (i settingItem) FilterValue() string { return i.entry.Label }
//
//func buildSettingItems(entries []settingEntry) []list.Item {
//	items := make([]list.Item, 0, len(entries))
//	for _, entry := range entries {
//		items = append(items, settingItem{entry: entry})
//	}
//	return items
//}
//
//type settingsDelegate struct {
//	listDelegate
//	styles listStyles
//}
//
//func (d settingsDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
//	i, ok := listItem.(settingItem)
//	if !ok {
//		return
//	}
//	label, value := i.entry.Label, i.entry.Value()
//
//	// Rows get 4 columns of left padding from the item styles, so the value
//	// must fit in what remains after the label column.
//	valueMax := max(m.Width()-4-settingLabelWidth, 1)
//	value = ansi.Truncate(value, valueMax, "…")
//	if !i.entry.Editable() {
//		label = faint(label)
//		value = faint(value)
//	}
//	str := lipgloss.NewStyle().Width(settingLabelWidth).Render(label) + value
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
//type SettingsTab struct {
//	items  []settingEntry
//	cfg    *config.ConfigFile
//	list   list.Model
//	styles listStyles
//	width  int
//	height int
//}
//
//func NewSettingsTab(cfg *config.ConfigFile) SettingsTab {
//	items := []settingEntry{
//		{Label: "Mode", Value: func() string {
//			val := cfg.Mode().String()
//			uppercased := cases.Title(language.AmericanEnglish)
//			return uppercased.String(val)
//		}, Editable: func() bool { return true }},
//		{Label: "Database URL", Value: func() string {
//			val := cfg.DbUri()
//			if val == "" {
//				return "<unset>"
//			}
//			return val
//		}, Editable: func() bool { return cfg.IsStandalone() }},
//		{Label: "Server URL", Value: func() string {
//			val := cfg.ServerURL()
//			if val == "" {
//				return "<unset>"
//			}
//			return val
//		}, Editable: func() bool { return cfg.IsClient() }},
//	}
//	t := SettingsTab{cfg: cfg, items: items}
//	t.styles = newListStyles(true)
//	t.list = newStyledList("Settings", settingsDelegate{styles: t.styles}, t.styles)
//	t.list.SetItems(buildSettingItems(items))
//	return t
//}
//
//func (t SettingsTab) ShortHelp() []key.Binding {
//	return []key.Binding{navUp, navDown, enterEdit}
//}
//
//func (t SettingsTab) editRequestFor(label string) editRequest {
//	cfg := t.cfg
//	switch label {
//	case "Mode":
//		choices := []string{"Standalone", "Client"}
//		selected := 0
//		if cfg.IsClient() {
//			selected = 1
//		}
//		return editRequest{
//			title: "Mode", kind: modalList, choices: choices, selected: selected,
//			commit: func(_ string, sel int) tea.Cmd {
//				return func() tea.Msg {
//					converted := strings.ToLower(choices[sel])
//					if err := cfg.SetMode(config.JatMode(converted)); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					if err := cfg.Update(); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					return editResultMsg{refreshTab: TabSettings}
//				}
//			},
//		}
//	case "Database URL":
//		if !cfg.IsStandalone() {
//			return editRequest{
//				title:   "Database URL",
//				kind:    modalLocked,
//				message: "Database URL is only editable in Standalone mode",
//			}
//		}
//		return editRequest{
//			title: "Database URL",
//			kind:  modalText,
//			input: newTextInput(cfg.DbUri(), "Database URL"),
//			commit: func(value string, _ int) tea.Cmd {
//				return func() tea.Msg {
//					if err := utils.ValidateDbUri(value); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					err := cfg.SetDbUri(value)
//					if err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					if err := cfg.Update(); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					return editResultMsg{refreshTab: TabSettings}
//				}
//			},
//		}
//	case "Server URL":
//		if !cfg.IsClient() {
//			return editRequest{
//				title:   "Server URL",
//				kind:    modalLocked,
//				message: "Server URL is only editable in Client mode",
//			}
//		}
//		return editRequest{
//			title: "Server URL",
//			kind:  modalText,
//			input: newTextInput(cfg.ServerURL(), "Server URL"),
//			commit: func(value string, _ int) tea.Cmd {
//				return func() tea.Msg {
//					if err := utils.ValidateUrl(value); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					cfg.SetServerURL(value)
//					if err := cfg.Update(); err != nil {
//						return editResultMsg{err: err, refreshTab: TabSettings}
//					}
//					return editResultMsg{refreshTab: TabSettings}
//				}
//			},
//		}
//	default:
//		return editRequest{
//			title: label,
//			kind:  modalNotImplemented,
//		}
//	}
//}
//
//func (t SettingsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
//	if kp, ok := msg.(tea.KeyPressMsg); ok && key.Matches(kp, enterEdit) {
//		req := t.editRequestFor(t.items[t.list.Index()].Label)
//		return t, func() tea.Msg { return editRequestMsg{req: req} }
//	}
//	var cmd tea.Cmd
//	t.list, cmd = t.list.Update(msg)
//	return t, cmd
//}
//
//func (t SettingsTab) View() string {
//	return placeListView(t.width, t.height, t.list.View())
//}
//
//func (t SettingsTab) refresh() SettingsTab {
//	t.list.SetItems(buildSettingItems(t.items))
//	return t
//}
//
//// Resize sizes the list to fill the available body dimensions; the delegate
//// truncates the value column against the list width.
//func (t SettingsTab) Resize(w, h int) Tab {
//	if w < 30 {
//		w = 30
//	}
//	t.width = w
//	t.height = h
//	t.list.SetSize(w, h)
//	return t
//}
