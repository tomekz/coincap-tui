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
	tui "github.com/tomekz/tui/tui/commands"
)

var (
	color       = termenv.EnvColorProfile().Color
	help        = termenv.Style{}.Foreground(color("241")).Styled
	highlighted = termenv.Style{}.Foreground(color("5")).Styled
)
var labels = map[string]string{
	"BTC":   "Bitcoin",
	"ETH":   "Ethereum",
	"USDT":  "Tether",
	"USDC":  "USD Coin",
	"BNB":   "BNB",
	"BUSD":  "Binance USD",
	"XRP":   "XRP",
	"ADA":   "Cardano",
	"SOL":   "Solana",
	"DOT":   "Polkadot",
	"Other": "Search other assets...",
}

type Model struct {
	choices  []string
	cursor   int
	selected string
	spinner  spinner.Model
	table    table.Model
	loading  bool
	data     data.Asset
	error    error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.loading = true
			m.selected = m.choices[m.cursor]
			return m, tea.Batch(
				tui.SearchCmd(m.selected),
				spinner.Tick,
			)
		}
	case data.DataFetchError:
		m.error = msg
	case tui.GotAsset:
		m.data = msg.Asset
		m.loading = false

		rows := []table.Row{}
		rows = append(rows, table.SimpleRow{m.data.AssetID, m.data.Name, m.data.PriceUsd})
		m.table.SetRows(rows)
		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func initialModel() Model {

	textInput := textinput.New()
	textInput.Placeholder = "Type your question here"
	textInput.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))

	tbl := table.New([]string{"ASSETID", "NAME", "PRICE IN USD"}, 40, 2)

	return Model{
		choices: []string{"BTC", "ETH", "USDT", "USDC", "BNB", "BUSD", "XRP", "ADA", "SOL", "DOT", "Other"},
		spinner: s,
		table:   tbl,
	}
}

func baseView(content string) string {
	return "Select asset" +
		"\n\n" +
		content +
		"\n\n" + help("◀ ↑/k: up • ↓/j: down • enter: submit • q: exit ▶\n")
}

func (m Model) View() string {
	if m.error != nil {
		return baseView(fmt.Sprintf("We had some trouble: %v", m.error))
	}

	var s string

	for i, choice := range m.choices {

		cursor := " "
		if i == m.cursor {
			cursor = ">"
			s += fmt.Sprintf(highlighted("%s %s\n"), cursor, labels[choice])
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, labels[choice])
		}

	}

	if m.loading {
		return baseView(fmt.Sprintf("%s loading...", m.spinner.View()))
	}

	if m.data.AssetID != "" {
		return baseView(m.table.View())
	}

	return baseView(s)
}

func Start() {
	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
