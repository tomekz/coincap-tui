/*
Tui let's you check crypto prices in your terminal.
Usage:

  tui [command]

Available Commands:
	-h
		help for Tui
	-v
		version for Tui
	-refresh
		enable auto refresh (every 10 seconds)
*/
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/tui/ui"
)

func main() {

	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	p := tea.NewProgram(ui.Init(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
