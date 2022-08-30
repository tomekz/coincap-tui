package searchui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/tui/commands"
	"github.com/tomekz/tui/tui/constants"
)

type model struct {
	label string
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func New() tea.Model {
	return model{
		label: "Bar",
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
func (m model) View() string {
	// The header
	s := "Bar"

	s += "\n " + m.label

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
