package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/tomekz/tui/tui/commands"
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
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Println("main.Update", msg, m.currentView)
	switch msg := msg.(type) {

	case tui.ChangeUiMsg:
		log.Println("ChangeUiMsg", m.currentView)
		m.currentView = searchView
	// case searchUi.ChangeUiMsg:
	// 	log.Println("ChangeUiMsg", m.state)
	// 	m.state = barUi

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case startView:
		s, c := m.start.Update(msg)
		m.start = s
		cmd = c
	case searchView:
		m.search = searchui.New()
		sc, nCmd := m.search.Update(msg)
		m.search = sc
		cmd = nCmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}
func (m Model) View() string {
	switch m.currentView {
	case searchView:
		return m.search.View()
	default:
		return m.start.View()
	}
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
