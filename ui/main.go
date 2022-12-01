package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

	l := list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "₿"
	l.SetSpinner(spinner.Pulse)
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keymap.Exit,
		}
	}
	return mainModel{
		keymap: keymap,
		list:   l,
	}
}

type mainModel struct {
	list   list.Model
	keymap *keymap
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			log.Println("tea.KeyMsg -> ctrl+c")
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m mainModel) View() string {
	return m.list.View()
}
