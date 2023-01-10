package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// light palette: https://colorhunt.co/palette/201882
// dark palette:  https://colorhunt.co/palette/273948
var (
	// defaultStyles = list.NewDefaultItemStyles()

	// activeColor    = defaultStyles.SelectedTitle.GetForeground()
	// secondaryColor = defaultStyles.NormalTitle.GetForeground()

	errorColor = lipgloss.AdaptiveColor{
		Light: "#e94560",
		Dark:  "#f05945",
	}

	// secondaryForeground = lipgloss.NewStyle().Foreground(secondaryColor)
	// boldStyle             = lipgloss.NewStyle().Bold(true)
	// activeForegroundBold  = lipgloss.NewStyle().Bold(true).Foreground(activeColor)
	// errorFaintForeground  = lipgloss.NewStyle().Foreground(errorColor).Faint(true)
	// errorForegroundPadded = lipgloss.NewStyle().Padding(4).Foreground(errorColor)
	// separator             = secondaryForeground.Render(" • ")
	// listStyle             = lipgloss.NewStyle().Margin(6, 2, 0, 2)
	primaryColor   = lipgloss.Color("#EBDBB2")
	secondaryColor = lipgloss.Color("#504945")

	tableStyles = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Foreground(primaryColor)

	helpStyles = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Margin(0)
)

func calculateTableDimensions(width int, height int) (int, int) {
	tableWidth := width - 2 // 2 for padding
	tableHeight := 20       // 20  rows

	return tableWidth, tableHeight
}

func TableStyles(baseStyles table.Styles) table.Styles {
	baseStyles.Header = baseStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Foreground(secondaryColor).
		Background(primaryColor).
		BorderBottom(true)
	baseStyles.Selected = baseStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("124")).
		Bold(false)
	// baseStyles.Cell = baseStyles.Cell.BorderBottom(true).
	// 	BorderStyle(lipgloss.NormalBorder())
	return baseStyles
}