package tui

import tea "github.com/charmbracelet/bubbletea"

type ChangeUiMsg bool

func ChangeUiCmd(value bool) tea.Cmd {
	return func() tea.Msg {
		return ChangeUiMsg(value)
	}
}
