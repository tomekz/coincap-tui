package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomekz/tui/coincap"
)

type keymap struct {
	Exit   key.Binding
	Help   key.Binding
	Select key.Binding
	GoBack key.Binding
}

type mainModel struct {
	currView  currentView
	graphView graphModel
	tableView tableModel
	keymap    keymap
	error     error
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

	case SelectAssetMsg:
		m.currView = graphView

	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
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

func baseView(contentView string, helpView string, width int, height int) string {

	return lipgloss.JoinVertical(
		lipgloss.Center,
		tableStyles.
			// Align(lipgloss.Center).
			Width(width).
			Height(height).
			Render(contentView),
		helpStyles.Align(lipgloss.Center).Width(width).Render(helpView),
	)
}

func (m mainModel) View() string {

	if m.error != nil {
		return baseView(fmt.Sprintf("Error: %s", m.error), m.help.View(m.keymap), m.width, m.height)
	}
	if m.error != nil {
		return baseView(m.error.Error(), m.help.View(m.keymap), m.width, m.height)
	}

	tWidth, tHeight := calculateTableDimensions(m.width, m.height)

	switch m.currView {
	// helpStyles.Align(lipgloss.Center).Width(tWidth).Render(m.help.View(m.keymap)),
	case graphView:
		return baseView(
			m.graphView.View(),
			m.graphView.help.View(m.graphView.keymap),
			tWidth,
			tHeight,
		)
	default:
		return baseView(
			m.tableView.View(),
			m.tableView.help.View(m.tableView.keymap),
			tWidth,
			tHeight,
		)
	}
}

// ======= //
// cmds    //
// ======= //
func getAssetsCmd() tea.Cmd {
	return func() tea.Msg {
		assets, err := coincap.GetAssets()
		if err != nil {
			return errMsg{err}
		}
		return getAssetsMsg{assets: assets}
	}
}

func getAssetHistoryCmd(assetId string) tea.Cmd {
	return func() tea.Msg {
		assetHistory, err := coincap.GetAssetHistory(assetId)
		if err != nil {
			return errMsg{err}
		}
		return getAssetHistoryMsg{assetHistory: assetHistory}
	}
}

// ======= //
// msgs    //
// ======= //
type getAssetsMsg struct {
	assets []coincap.Asset
}

type getAssetHistoryMsg struct {
	assetHistory []coincap.AssetHistory
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// ======= //
// models  //
// ======= //
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		// k.Refresh,
		k.Exit,
		k.Select,
		k.GoBack,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{table.DefaultKeyMap().GotoTop, table.DefaultKeyMap().GotoBottom, table.DefaultKeyMap().PageUp, table.DefaultKeyMap().PageDown},
		{table.DefaultKeyMap().HalfPageUp, table.DefaultKeyMap().HalfPageDown,
			// k.Refresh,
			k.Exit},
	}
}
