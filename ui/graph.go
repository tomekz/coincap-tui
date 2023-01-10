package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guptarohit/asciigraph"
)

type graphKeymap struct {
	GoBack key.Binding
}

type graphModel struct {
	keymap        graphKeymap
	assethHistory []float64
	help          help.Model
	isLoading     bool
	spinner       spinner.Model
	height        int
	width         int
}

func (m graphModel) Init() tea.Cmd {
	return nil
}

func (m graphModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case SelectAssetMsg:
		m.isLoading = true
		return m, tea.Batch(m.spinner.Tick, getAssetHistoryCmd(msg.value))

	case getAssetHistoryMsg:
		assetHistory := make([]float64, len(msg.assetHistory))
		for i, ah := range msg.assetHistory {
			assetHistory[i] = ah.PriceUsd
		}
		m.assethHistory = assetHistory
		m.isLoading = false

	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.GoBack) {
			return m, GoBackCmd()
		}

	}
	return m, nil
}

func (m graphModel) View() string {
	if m.isLoading {
		return m.spinner.View()
	}

	tWidth, tHeight := calculateTableDimensions(m.width, m.height)

	graph := asciigraph.Plot(
		m.assethHistory,
		asciigraph.Height(tHeight),
		asciigraph.Width(tWidth),
		asciigraph.Caption("Price History"),
		asciigraph.CaptionColor(asciigraph.Red),
	)

	return graph
}

func newGraph() graphModel {
	keymap := graphKeymap{
		GoBack: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "go back"),
		),
	}
	spinner := spinner.New()

	return graphModel{
		assethHistory: []float64{},
		keymap:        keymap,
		spinner:       spinner,
	}
}

// ======== //
// commands //
// ======== //

type GoBackMsg bool

func GoBackCmd() tea.Cmd {
	return func() tea.Msg {
		return GoBackMsg(true)
	}
}

func (k graphKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		// k.Help,
		// k.Refresh,
		// k.Exit,
		// k.Select,
		k.GoBack,
	}
}

func (k graphKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{},
	}
}
