/*
*

	Tui let's you check crypto prices in your terminal.

*
*/
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomekz/coincap-tui/ui"
)

func main() {
	if os.Getenv("DEBUG") != "true" {
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
	}

	log.SetPrefix("tui: ")
	log.SetFlags(log.Ltime | log.LUTC)

	p := tea.NewProgram(ui.Init(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
