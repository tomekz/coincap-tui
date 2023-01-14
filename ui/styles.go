package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("#EBDBB2")
	secondaryColor = lipgloss.Color("#504945")

	tableStyles = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Foreground(primaryColor)

	helpStyles = lipgloss.NewStyle()
)

func calculateTableDimensions(width int, height int) (int, int) {
	tableWidth := (width / 7) * 6 // 6/7 of the width
	tableHeight := 20             // 20  rows

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
	return baseStyles
}
