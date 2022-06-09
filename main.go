package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/tomekz/tui/tui"
)

var (
	color   = termenv.EnvColorProfile().Color
	keyword = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	help    = termenv.Style{}.Foreground(color("241")).Styled
)

type Model struct {
	textInput textinput.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, func() tea.Msg {
				return tui.GetData(m.textInput.Value())
			}
		}

	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func initialModel() Model {

	textInput := textinput.New()
	textInput.Placeholder = "Type your question here"
	textInput.Focus()

	return Model{
		textInput: textInput,
	}
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Question? \n\n %s",
		m.textInput.View(),
	) + "\n\n" + help("space: switch modes • q: exit\n")
}

func main() {
	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
