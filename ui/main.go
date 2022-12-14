package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/coincap"
)

type keymap struct {
	Exit key.Binding
	Help key.Binding
}

func Init() tea.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 6},
		{Title: "Price USD", Width: 10},
		{Title: "Change (24hr)", Width: 10},
		{Title: "Supply", Width: 10},
		{Title: "Max Supply", Width: 10},
		{Title: "Market Cap", Width: 10},
		{Title: "Volume (24hr)", Width: 10},
	}

	table := table.New(table.WithColumns(columns), table.WithFocused(true))
	keymap := keymap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
	}
	return mainModel{
		keymap: keymap,
		table:  table,
		help:   help.New(),
	}
}

type mainModel struct {
	keymap keymap
	error  error
	table  table.Model
	help   help.Model
}

func (m mainModel) Init() tea.Cmd {
	return getAssetsCmd()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keymap.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}

	case getAssetsMsg:
		rows := make([]table.Row, len(msg.assets))
		for i, asset := range msg.assets {
			rows[i] = []string{
				asset.Rank,
				asset.Name,
				asset.Symbol,
				asset.PriceUsd,
				asset.ChangePercent24Hr,
				asset.Supply,
				asset.MaxSupply,
				asset.MarketCapUsd,
				asset.VolumeUsd24Hr,
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

	return baseStyle.Render(m.table.View() + "\n\n" + m.help.View(m.keymap))
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
// Help functions. Used in creating the help menu
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Exit,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{table.DefaultKeyMap().GotoTop},
		{table.DefaultKeyMap().GotoBottom},
		{table.DefaultKeyMap().PageUp},
		{table.DefaultKeyMap().PageDown},
		{table.DefaultKeyMap().HalfPageUp},
		{table.DefaultKeyMap().HalfPageDown},
		{k.Exit},
	}
}
