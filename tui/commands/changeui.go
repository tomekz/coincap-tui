package commands

import tea "github.com/charmbracelet/bubbletea"

type ChangeUiMsg string

func ChangeUiCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return ChangeUiMsg(value)
	}
}
