package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/tui/commands"
	"github.com/tomekz/tui/tui/constants"
	"github.com/tomekz/tui/tui/searchui"
	"github.com/tomekz/tui/tui/startui"
)

type currentView int

const (
	startView currentView = iota
	searchView
)

// the main Model of the program; holds other models and bubbles
type Model struct {
	currentView currentView
	start       tea.Model
	search      tea.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func initialModel() Model {
	return Model{
		currentView: startView,
		start:       startui.New(),
		search:      searchui.New(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	log.Println("main.Update", msg, m.currentView)
	switch msg := msg.(type) {

	case commands.ChangeUiMsg:
		if msg == "restart" {
			m.currentView = startView
			m.start = startui.New()
			m.search = searchui.New()
		} else {
			m.currentView = searchView
			cmds = append(cmds, textinput.Blink)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Restart):
			return m, commands.ChangeUiCmd("restart")
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case startView:
		newStartModel, newStartCmd := m.start.Update(msg)
		m.start = newStartModel
		cmd = newStartCmd
	case searchView:
		newSearchModel, newSearchCmd := m.search.Update(msg)
		m.search = newSearchModel
		cmd = newSearchCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.currentView {
	case searchView:
		return baseView(m.search.View())
	default:
		return baseView(m.start.View())
	}
}

func baseView(content string) string {
	return "Select asset" +
		"\n\n" +
		content +
		"\n\n" + constants.HelpStyle(
		fmt.Sprintf("◀ ↑: up • ↓: down • enter: %v • esc: %v • ctrl+c: %v ▶ ",
			constants.Keymap.Enter.Help().Desc,
			constants.Keymap.Restart.Help().Desc,
			constants.Keymap.Quit.Help().Desc,
		))
}

func Start() {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
