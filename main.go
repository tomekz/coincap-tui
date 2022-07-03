package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	data "github.com/tomekz/tui/src"
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
	data      *data.Data
	error     error
}

type gotData struct {
	Data *data.Data
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func onEnterCmd(value string) tea.Cmd {
	return func() tea.Msg {
		data, err := data.FetchData(value)
		if err != nil {
			return err
		}
		return gotData{Data: data}
	}
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
			return m, tea.Batch(
				onEnterCmd(m.textInput.Value()),
				spinner.Tick,
			)
		}
	case data.DataFetchError:
		m.error = msg
	case gotData:
		m.data = msg.Data
		m.loading = false
		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func initialModel() Model {

	textInput := textinput.New()
	textInput.Placeholder = "Type your question here"
	textInput.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))

	return Model{
		textInput: textInput,
		spinner:   s,
	}
}

func baseView(content any) string {
	return fmt.Sprintf(
		"Question? \n\n %s",
		content,
	) + "\n\n" + help("◀ Enter: submit • q: exit ▶\n")
}

func (m Model) View() string {

	if m.error != nil {
		return baseView(fmt.Sprintf("We had some trouble: %v", m.error))
	}

	if m.loading {
		return baseView(fmt.Sprintf("%s loading...", m.spinner.View()))
	}

	if m.data != nil {
		c := keyword(fmt.Sprintf("%+v", m.data))
		return baseView(c)
	}

	return baseView(m.textInput.View())
}

func main() {
	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
