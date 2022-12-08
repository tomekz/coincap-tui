package ui

import (
	"fmt"
	"io"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomekz/tui/coincap"
)

type keymap struct {
	Exit key.Binding
}

func Init() tea.Model {
	keymap := &keymap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
	}

	l := list.NewModel([]list.Item{}, itemDelegate{}, 20, 20)
	l.Title = "₿"
	l.SetSpinner(spinner.Pulse)
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keymap.Exit,
		}
	}
	return mainModel{
		keymap: keymap,
		list:   l,
	}
}

type mainModel struct {
	list   list.Model
	keymap *keymap
	error  error
}

func (m mainModel) Init() tea.Cmd {
	return getAssetsCmd()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
	case getAssetsMsg:
		items := make([]list.Item, len(msg.assets))
		for _, asset := range msg.assets {
			items = append(items, item{
				title: asset.Name,
			})
		}
		cmds = append(cmds, m.list.SetItems(items))
	case errMsg:
		log.Println(msg.error)
		m.error = msg.error
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	return m.list.View()
}

// cmds
func getAssetsCmd() tea.Cmd {
	return func() tea.Msg {
		assets, err := coincap.GetAssets()
		if err != nil {
			return errMsg{err}
		}
		return getAssetsMsg{assets: assets}
	}
}

// msgs
type getAssetsMsg struct {
	assets []coincap.Asset
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// models

type item struct {
	title string
}

func (i item) Title() string {
	// if i.end.IsZero() {
	// 	return boldStyle.Render(i.title)
	// }
	return i.title
}

func (i item) Description() string {
	return "LOL"
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}
