package constants

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/muesli/termenv"
)

var (
	color = termenv.EnvColorProfile().Color
	// HelpStyle styling for help context menu
	HelpStyle = termenv.Style{}.Foreground(color("241")).Styled
	// HighligtedStyle styling for help context menu
	HighlightedStyle = termenv.Style{}.Foreground(color("5")).Styled
)

type keymap struct {
	Change  key.Binding
	Enter   key.Binding
	Restart key.Binding
	Quit    key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Change: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "change view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Restart: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "restart"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
