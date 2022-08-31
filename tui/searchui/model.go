package searchui

import (
	"fmt"
	"log"

	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomekz/tui/data"
	"github.com/tomekz/tui/tui/commands"
	"github.com/tomekz/tui/tui/constants"
)

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	table     table.Model
	loading   bool
	data      data.Asset
	error     error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func New() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Type your search here"
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))

	tbl := table.New([]string{"ASSETID", "NAME", "PRICE IN USD"}, 40, 2)

	return model{
		textInput: ti,
		spinner:   s,
		table:     tbl,
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
		case key.Matches(msg, constants.Keymap.Enter):
			m.loading = true
			m.error = nil
			return m, tea.Batch(
				commands.SearchCmd(m.textInput.Value()),
				spinner.Tick,
			)
		}
	case commands.GotAsset:
		m.data = msg.Asset
		m.loading = false

		rows := []table.Row{}
		rows = append(rows, table.SimpleRow{m.data.AssetID, m.data.Name, fmt.Sprintf("%.2f$", m.data.PriceUsd)})
		m.table.SetRows(rows)
		return m, nil
	case data.DataFetchError:
		m.error = msg
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.error != nil {
		return fmt.Sprintf("We had some trouble: %v", m.error)
	}

	if m.loading {
		return fmt.Sprintf("%s loading...", m.spinner.View())
	}

	if m.data.AssetID != "" {
		return m.table.View()
	}

	return m.textInput.View()
}
