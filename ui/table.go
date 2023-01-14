package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/coincap"
)

type tableKeymap struct {
	Refresh key.Binding
	Select  key.Binding
}

type tableModel struct {
	keymap tableKeymap
	table  table.Model
	error  error
}

func (m tableModel) Init() tea.Cmd {
	return getAssetsCmd()
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Refresh) {
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

	case errMsg:
		m.error = msg.error
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

	return tableModel{
		keymap: keymap,
		table:  t,
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
		assets, err := coincap.GetAssets(200)
		if err != nil {
			return errMsg{err}
		}
		return getAssetsMsg{assets: assets}
	}
}

func (k tableKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Refresh,
		k.Select,
	}
}

func (k tableKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Select},
	}
}
