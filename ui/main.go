package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	Exit    key.Binding
	Help    key.Binding
	Select  key.Binding
	GoBack  key.Binding
	Refresh key.Binding
}

type mainModel struct {
	currView  currentView
	graphView graphModel
	tableView tableModel
	keymap    keymap
	help      help.Model
	width     int
	height    int
}

type currentView int

const (
	tableView currentView = iota
	graphView
)

func Init() tea.Model {
	keymap := keymap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		GoBack: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "go back"),
		),
	}

	return mainModel{
		currView:  tableView,
		tableView: newTable(),
		graphView: newGraph(),
		keymap:    keymap,
		help:      help.New(),
	}
}

func (m mainModel) Init() tea.Cmd {
	return newTable().Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.help.Width = msg.Width

	case GoBackMsg:
		m.currView = tableView
		log.Debug("going back to table view")

	case SelectAssetMsg:
		m.currView = graphView

	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keymap.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	switch m.currView {
	case tableView:
		newTableModel, newTableCmd := m.tableView.Update(msg)
		m.tableView = newTableModel.(tableModel)
		cmd = newTableCmd
	case graphView:
		newGraphModel, newGraphCmd := m.graphView.Update(msg)
		m.graphView = newGraphModel.(graphModel)
		cmd = newGraphCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	tWidth, tHeight := calculateTableDimensions(m.width, m.height)

	switch m.currView {
	case graphView:
		return lipgloss.JoinVertical(
			lipgloss.Center,
			tableStyles.Render(m.graphView.View()), helpStyles.Width(tWidth).Render(m.help.View(m.graphView.keymap)),
		)

	default:
		return lipgloss.JoinVertical(
			lipgloss.Center,
			tableStyles.Width(tWidth).Height(tHeight).Render(m.tableView.View()),
			helpStyles.Width(tWidth).Render(m.help.View(m.tableView.keymap)),
			helpStyles.Width(tWidth).Render(m.help.View(m.keymap)),
		)
	}
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// ======= //
// models  //
// ======= //
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Exit,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Exit, table.DefaultKeyMap().GotoTop, table.DefaultKeyMap().GotoBottom, table.DefaultKeyMap().PageUp, table.DefaultKeyMap().PageDown},
		{table.DefaultKeyMap().HalfPageUp, table.DefaultKeyMap().HalfPageDown},
	}
}
