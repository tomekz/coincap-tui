package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/data"
)

type GotAsset struct {
	Asset data.Asset
}

func SearchCmd(value string) tea.Cmd {
	return func() tea.Msg {
		data, err := data.SearchAssets(value)

		if err != nil {
			return err
		}
		return GotAsset{Asset: data[0]}
	}
}
