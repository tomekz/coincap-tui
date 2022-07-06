package tui

import (
	"fmt"
	"log"

	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/tomekz/tui/data"
)

var (
	color = termenv.EnvColorProfile().Color
	help  = termenv.Style{}.Foreground(color("241")).Styled
)

type Model struct {
	textInput textinput.Model
	spinner   spinner.Model
	table     table.Model
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

		rows := []table.Row{}
		rows = append(rows, table.SimpleRow{m.data.UserId, m.data.Id, m.data.Title, m.data.Completed})
		m.table.SetRows(rows)
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

	tbl := table.New([]string{"USERID", "ID", "TITLE", "COMPLETED"}, 40, 2)

	return Model{
		textInput: textInput,
		spinner:   s,
		table:     tbl,
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
		return baseView(m.table.View())
	}

	return baseView(m.textInput.View())
}

func Start() {
	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
