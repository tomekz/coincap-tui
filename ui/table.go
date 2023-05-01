package ui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/coincap-tui/coincap"
	"github.com/tomekz/coincap-tui/favourites"
)

type tableKeymap struct {
	Refresh        key.Binding
	Select         key.Binding
	Favourite      key.Binding
	ShowFavourites key.Binding
	ShowAll        key.Binding
}

type tableModel struct {
	assets    []Asset
	keymap    tableKeymap
	table     table.Model
	error     error
	spinner   spinner.Model
	isLoading bool
}

var Favs = favourites.New()

func (m tableModel) Init() tea.Cmd {
	return tea.Batch(GetAssetsCmd(), LoadFavsCmd())
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Refresh) {
			m.error = nil
			m.isLoading = true
			return m, tea.Batch(m.spinner.Tick, GetAssetsCmd())
		}
		if key.Matches(msg, m.keymap.Select) {
			assetId := m.table.SelectedRow()[2]
			return m, SelectAssetCmd(assetId)
		}
		if key.Matches(msg, m.keymap.Favourite) {
			assetId := m.table.SelectedRow()[2]
			return m, FavAssetCmd(assetId)
		}
		if key.Matches(msg, m.keymap.ShowFavourites) {
			return m, GetFavouriteAssetsCmd(m.assets)
		}
		if key.Matches(msg, m.keymap.ShowAll) {
			return m, GetAssetsCmd()
		}

	case GetAssetsMsg:
		m.assets = msg.assets // save assets cause table model don't have a getter
		rows := getRows(msg.assets)
		m.table.SetRows(rows)
		m.table.SetCursor(0)
		m.isLoading = false

	case FavAssetMsg:
		m.assets = updateFavourite(m.assets, msg.assetId)
		rows := getRows(m.assets)
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
			key.WithHelp("f", "toggle favourite"),
		),
		ShowFavourites: key.NewBinding(
			key.WithKeys("F"),
			key.WithHelp("F", "show favourites"),
		),
		ShowAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "show all"),
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
		Favs.Favourite(assetId)
		if err := Favs.Save(); err != nil {
			return errMsg{err}
		}
		return FavAssetMsg{assetId}
	}
}

func GetFavouriteAssetsCmd(assets []Asset) tea.Cmd {
	return func() tea.Msg {
		favs := getFavouriteAssets(assets)
		if len(favs) == 0 {
			return errMsg{errors.New("no favourites, refresh to load all assets")}
		}
		return GetAssetsMsg{assets: favs}
	}
}

func LoadFavsCmd() tea.Cmd {
	return func() tea.Msg {
		if err := Favs.Load(); err != nil {
			return errMsg{err}
		}
		return nil
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

type GetAssetsMsg struct {
	assets []Asset
}

func GetAssetsCmd() tea.Cmd {
	return func() tea.Msg {
		assets, err := coincap.GetAssets(300)
		if err != nil {
			return errMsg{err}
		}
		return GetAssetsMsg{assets: mapAssets(assets)}
	}
}

func (k tableKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Refresh,
		k.Select,
		k.ShowAll,
		k.Favourite,
		k.ShowFavourites,
	}
}

func (k tableKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Select, k.Favourite, k.ShowFavourites},
	}
}
