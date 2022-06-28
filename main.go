package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
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
	spinner   spinner.Model
	loading   bool
	data      *tui.Data
	error     error
}

type gotData struct {
	Data *tui.Data
}

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string {
	return e.err.Error()
}

const url = "https://jsonplaceholder.typicode.com"

func fetchData(id string) tea.Msg {
	data := &tui.Data{}
	err := tui.GetJson(fmt.Sprintf("%s/todos/%s", url, id), data)
	if err != nil {
		return errorMsg{err}
	}
	return gotData{Data: data}
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
			m.loading = true
			m.spinner, cmd = m.spinner.Update(msg)
			return m, func() tea.Msg {
				return fetchData(m.textInput.Value())
			}
		}
	case errorMsg:
		m.error = msg
	case gotData:
		m.data = msg.Data
		m.loading = false
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func initialModel() Model {

	textInput := textinput.New()
	textInput.Placeholder = "Type your question here"
	textInput.Focus()

	spinner := spinner.New()

	return Model{
		textInput: textInput,
		spinner:   spinner,
	}
}

func (m Model) View() string {
	if m.loading {
		// return m.spinner.View()
		return fmt.Sprint("Loading...")
	}

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
