package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/coincap-tui/coincap"
)

type tableKeymap struct {
	Refresh   key.Binding
	Select    key.Binding
	Favourite key.Binding
}

var favs map[string]bool = make(map[string]bool)

/*
	Adds or removes asset from favourites
*/
func favourite(favs map[string]bool, assetId string) {
	exists := favs[assetId]
	if exists {
		delete(favs, assetId)
	} else {
		favs[assetId] = true
	}
}

type tableModel struct {
	rows      []table.Row
	keymap    tableKeymap
	table     table.Model
	error     error
	spinner   spinner.Model
	isLoading bool
}

func (m tableModel) Init() tea.Cmd {
	return tea.Batch(getAssetsCmd(), LoadFavsCmd())
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Refresh) {
			m.isLoading = true
			return m, tea.Batch(m.spinner.Tick, getAssetsCmd())
		}
		if key.Matches(msg, m.keymap.Select) {
			assetId := m.table.SelectedRow()[2]
			return m, SelectAssetCmd(assetId)
		}
		if key.Matches(msg, m.keymap.Favourite) {
			assetId := m.table.SelectedRow()[2]
			return m, FavAssetCmd(assetId)
		}

	case getAssetsMsg:
		rows := make([]table.Row, len(msg.assets))
		for i, asset := range msg.assets {
			fav := ""
			if favs[asset.ID] {
				fav = "★"
			}
			rows[i] = []string{
				strconv.FormatInt(asset.Rank, 10),
				fav,
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
		m.rows = rows // save rows cause table model don't have a getter
		m.isLoading = false

	case FavAssetMsg:
		rows := make([]table.Row, len(m.rows))
		for i, row := range m.rows {
			fav := ""
			if favs[row[2]] {
				fav = "★"
			}
			rows[i] = []string{
				row[0],
				fav,
				row[2],
				row[3],
				row[4],
				row[5],
				row[6],
				row[7],
				row[8],
				row[9],
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
		Favourite: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "add to favourites"),
		),
	}

	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Fav", Width: 3},
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
	spinner := spinner.NewModel()

	return tableModel{
		keymap:  keymap,
		table:   t,
		spinner: spinner,
	}
}

// ======== //
// commands //
// ======== //

type FavAssetMsg struct {
	assetId string
}

func FavAssetCmd(assetId string) tea.Cmd {
	return func() tea.Msg {
		favourite(favs, assetId)
		content, err := json.Marshal(favs)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile("fav_assets.json", content, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return FavAssetMsg{assetId}
	}
}

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

func LoadFavsCmd() tea.Cmd {
	return func() tea.Msg {
		content, err := ioutil.ReadFile("fav_assets.json")
		if err != nil {
			return errMsg{err}
		}
		err = json.Unmarshal(content, &favs)
		log.Println("ll", favs)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
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
		k.Favourite,
	}
}

func (k tableKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Select, k.Favourite},
	}
}
