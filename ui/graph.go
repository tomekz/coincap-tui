package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type graphKeymap struct {
	Refresh key.Binding
}

type graphModel struct {
	keymap        graphKeymap
	assethHistory []float64
	help          help.Model
}

func (m graphModel) Init() tea.Cmd {
	return nil
}

func (m graphModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m graphModel) View() string {
	return "Hello,grap!"
}

func newGraph() graphModel {
	return graphModel{
		assethHistory: []float64{},
	}
}

func (k graphKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		// k.Help,
		k.Refresh,
		// k.Exit,
		// k.Select,
		// k.GoBack,
	}
}

func (k graphKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{},
		{},
	}
}
