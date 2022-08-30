package searchui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/tui/commands"
	"github.com/tomekz/tui/tui/constants"
)

type model struct {
	textInput textinput.Model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func New() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Type your search here"
	ti.Focus()

	return model{
		textInput: ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Println("searchui.Update", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Change):
			return m, commands.ChangeUiCmd("start")
			// default:
			// 	m.viewport, cmd = m.viewport.Update(msg)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}
func (m model) View() string {

	return m.textInput.View()
}
