package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
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
	currView      currentView
	graphView     graphModel
	tableView     tableModel
	keymap        keymap
	error         error
	table         table.Model
	help          help.Model
	spinner       spinner.Model
	width         int
	height        int
	isLoading     bool
	assethHistory []float64
	showGraph     bool
}

type currentView int

const (
	tableView currentView = iota
	graphView
)

func Init() tea.Model {
	columns := []table.Column{

		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 6},
		{Title: "Price USD", Width: 10},
		{Title: "Change (24hr)", Width: 15},
		{Title: "Supply", Width: 10},
		{Title: "Max Supply", Width: 10},
		{Title: "Market Cap", Width: 10},
		{Title: "Volume (24hr)", Width: 15},
	}

	t := table.New(table.WithColumns(columns), table.WithFocused(true))
	s := table.DefaultStyles()
	t.SetStyles(TableStyles(s))

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
	spinner := spinner.New()

	return mainModel{
		currView:      tableView,
		tableView:     newTable(),
		graphView:     newGraph(),
		keymap:        keymap,
		table:         t,
		help:          help.New(),
		spinner:       spinner,
		isLoading:     true,
		assethHistory: []float64{},
		showGraph:     false,
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getAssetsCmd())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.help.Width = msg.Width

	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keymap.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
		// if key.Matches(msg, m.keymap.Select) {
		// 	m.isLoading = true
		// 	m.showGraph = true
		// 	m.assethHistory = []float64{}
		//
		// 	return m, getAssetHistoryCmd(m.table.SelectedRow()[1])
		// }
		if key.Matches(msg, m.keymap.GoBack) {
			m.showGraph = false
			m.error = nil
			m.assethHistory = []float64{}
			m.isLoading = false
			return m, nil
		}

		// case getAssetHistoryMsg:
		// 	assetHistory := make([]float64, len(msg.assetHistory))
		// 	for i, ah := range msg.assetHistory {
		// 		assetHistory[i] = ah.PriceUsd
		// 	}
		// 	m.assethHistory = assetHistory
		// 	m.isLoading = false
		//
		// 	// case errMsg:
		// 	// 	m.error = msg.error
		// 	// default:
		// 	// 	var spinnerUpdateCmd tea.Cmd
		// 	// 	m.spinner, spinnerUpdateCmd = m.spinner.Update(msg)
		// 	// 	cmds = append(cmds, spinnerUpdateCmd)
		// 	//
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

	// if m.isLoading {
	// 	v = m.spinner.View()
	// }
	//
	// if len(m.assethHistory) > 0 && m.showGraph {
	// 	graph := asciigraph.Plot(
	// 		m.assethHistory,
	// 		asciigraph.Height(tHeight),
	// 		asciigraph.Width(tWidth),
	// 		asciigraph.Caption("Price History"),
	// 		asciigraph.CaptionColor(asciigraph.Red),
	// 	)
	// 	v = graph
	// }
	//
	// if m.error != nil {
	// 	v = fmt.Sprintf("Error: %s", m.error)
	// }
	//
	// return lipgloss.JoinVertical(
	// 	lipgloss.Center,
	// 	tableStyles.
	// 		// Align(lipgloss.Center).
	// 		Width(tWidth).
	// 		Height(tHeight).
	// 		Render(v),
	// 	helpStyles.Align(lipgloss.Center).Width(tWidth).Render(m.help.View(m.keymap)),
	// )
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
