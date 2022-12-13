package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Rank", Width: 10},
		{Title: "Symbol", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Supply", Width: 10},
		{Title: "MaxSupply", Width: 10},
		{Title: "MarketCapUsd", Width: 10},
		{Title: "VolumeUsd24Hr", Width: 10},
		{Title: "PriceUsd", Width: 10},
		{Title: "ChangePercent24Hr", Width: 10},
		{Title: "Vwap24Hr", Width: 10},
	}

	table := table.New(table.WithColumns(columns), table.WithFocused(true))

	return mainModel{
		keymap: keymap,
		table:  table,
	}
}

type mainModel struct {
	keymap *keymap
	error  error
	table  table.Model
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
		rows := make([]table.Row, len(msg.assets))
		for i, asset := range msg.assets {
			rows[i] = []string{
				asset.ID,
				asset.Rank,
				asset.Symbol,
				asset.Name,
				asset.Supply,
				asset.MaxSupply,
				asset.MarketCapUsd,
				asset.VolumeUsd24Hr,
				asset.PriceUsd,
				asset.ChangePercent24Hr,
				asset.Vwap24Hr,
			}
		}
		m.table.SetRows(rows)

	case errMsg:
		log.Println(msg.error)
		m.error = msg.error
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	return m.table.View()
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

// ======= //
// msgs    //
// ======= //
type getAssetsMsg struct {
	assets []coincap.Asset
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// ======= //
// models  //
// ======= //
