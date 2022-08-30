package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/tui/commands"
	"github.com/tomekz/tui/tui/constants"
	"github.com/tomekz/tui/tui/searchui"
	"github.com/tomekz/tui/tui/startui"
)

type sessionState int

const (
	startView sessionState = iota
	searchView
)

// MainModel the main Model of the program; holds other models and bubbles
type Model struct {
	currentView sessionState
	start       tea.Model
	search      tea.Model
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
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
	log.Println("main.Update", msg, m.currentView)
	switch msg := msg.(type) {

	case commands.ChangeUiMsg:
		if msg == "restart" {
			m.currentView = startView
			m.start = startui.New()
			m.search = searchui.New()
		} else {
			m.currentView = searchView
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Restart):
			return m, commands.ChangeUiCmd("restart")
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case startView:
		s, c := m.start.Update(msg)
		m.start = s
		cmd = c
	case searchView:
		sc, nCmd := m.search.Update(msg)
		m.search = sc
		cmd = nCmd
	}

	return m, cmd
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
		"\n\n" + constants.HelpStyle("◀ ↑: up • ↓: down • enter: submit • esc: restart • ctrl+c: exit ▶\n")
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
