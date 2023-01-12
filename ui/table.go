package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/coincap"
)

type tableKeymap struct {
	Exit    key.Binding
	Help    key.Binding
	GoBack  key.Binding
	Refresh key.Binding
	Select  key.Binding
}

type tableModel struct {
	help      help.Model
	keymap    tableKeymap
	table     table.Model
	spinner   spinner.Model
	isLoading bool
	error     error
}

func (m tableModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getAssetsCmd())
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keymap.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
		if key.Matches(msg, m.keymap.Refresh) {
			m.isLoading = true
			return m, getAssetsCmd()
		}
		if key.Matches(msg, m.keymap.Select) {
			assetId := m.table.SelectedRow()[1]
			return m, SelectAssetCmd(assetId)
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

func (m tableModel) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %s", m.error)
	}

	if m.isLoading {
		return m.spinner.View()
	}

	return m.table.View()
}

func newTable() tableModel {
	keymap := tableKeymap{
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
		GoBack: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "go back"),
		),
	}

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

	spinner := spinner.New()

	return tableModel{
		keymap:    keymap,
		table:     t,
		spinner:   spinner,
		isLoading: true,
	}
}

// ======== //
// commands //
// ======== //

type SelectAssetMsg struct {
	value string
}

func SelectAssetCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return SelectAssetMsg{value}
	}
}

type getAssetsMsg struct {
	assets []coincap.Asset
}

func getAssetsCmd() tea.Cmd {
	return func() tea.Msg {
		assets, err := coincap.GetAssets()
		if err != nil {
			return errMsg{err}
		}
		return getAssetsMsg{assets: assets}
	}
}
