package ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/tomekz/tui/coincap"
)

type keymap struct {
	Exit    key.Binding
	Help    key.Binding
	Refresh key.Binding
	Select  key.Binding
	GoBack  key.Binding
}

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
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
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
		keymap:        keymap,
		table:         t,
		help:          help.New(),
		spinner:       spinner,
		isLoading:     true,
		assethHistory: []float64{},
		showGraph:     false,
	}
}

type mainModel struct {
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

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getAssetsCmd())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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
		if key.Matches(msg, m.keymap.Refresh) {
			m.error = nil
			m.isLoading = true
			m.showGraph = false
			return m, getAssetsCmd()
		}
		if key.Matches(msg, m.keymap.Select) {
			m.isLoading = true
			m.showGraph = true
			m.assethHistory = []float64{}

			return m, getAssetHistoryCmd(m.table.SelectedRow()[1])
		}
		if key.Matches(msg, m.keymap.GoBack) {
			log.Println("go back")
			m.showGraph = false
			m.error = nil
			m.assethHistory = []float64{}
			m.isLoading = false
			return m, nil
		}

	case getAssetsMsg:
		rows := make([]table.Row, len(msg.assets))
		for i, asset := range msg.assets {
			rows[i] = []string{
				strconv.FormatInt(asset.Rank, 10),
				asset.ID,
				asset.Symbol,
				strconv.FormatFloat(asset.PriceUsd, 'f', 2, 64),
				formatPercent(asset.ChangePercent24Hr),
				formatFloat(asset.Supply),
				formatFloat(asset.MaxSupply),
				formatFloat(asset.MarketCapUsd),
				formatFloat(asset.VolumeUsd24Hr),
			}
		}
		m.table.SetRows(rows)
		m.isLoading = false

	case getAssetHistoryMsg:
		assetHistory := make([]float64, len(msg.assetHistory))
		for i, ah := range msg.assetHistory {
			assetHistory[i] = ah.PriceUsd
		}
		m.assethHistory = assetHistory
		m.isLoading = false

	case errMsg:
		m.error = msg.error
	default:
		var spinnerUpdateCmd tea.Cmd
		m.spinner, spinnerUpdateCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerUpdateCmd)

	}
	var tableUpdateCmd tea.Cmd
	m.table, tableUpdateCmd = m.table.Update(msg)
	cmds = append(cmds, tableUpdateCmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {

	tWidth, tHeight := calculateTableDimensions(m.width, m.height)
	var v = m.table.View()

	if m.isLoading {
		v = m.spinner.View()
	}

	if len(m.assethHistory) > 0 && m.showGraph {
		graph := asciigraph.Plot(
			m.assethHistory,
			asciigraph.Height(tHeight),
			asciigraph.Width(tWidth),
			asciigraph.Caption("Price History"),
			asciigraph.CaptionColor(asciigraph.Red),
		)
		v = graph
	}

	if m.error != nil {
		v = fmt.Sprintf("Error: %s", m.error)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		tableStyles.
			// Align(lipgloss.Center).
			Width(tWidth).
			Height(tHeight).
			Render(v),
		helpStyles.Align(lipgloss.Center).Width(tWidth).Render(m.help.View(m.keymap)),
	)
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
		k.Refresh,
		k.Exit,
		k.Select,
		k.GoBack,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{table.DefaultKeyMap().GotoTop, table.DefaultKeyMap().GotoBottom, table.DefaultKeyMap().PageUp, table.DefaultKeyMap().PageDown},
		{table.DefaultKeyMap().HalfPageUp, table.DefaultKeyMap().HalfPageDown, k.Refresh, k.Exit},
	}
}
