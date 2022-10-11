package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Ehco1996/dig-up/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func runTUI(debug bool, curl string) error {
	if debug {
		f, err := tea.LogToFile("debug.log", "")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	m, err := tui.InitialModel(curl, upUID, favID)
	if err != nil {
		return err
	}

	if err = m.FetchVideoPage(startPage, tui.PageSize); err != nil {
		return err
	}
	p := tea.NewProgram(m)
	return p.Start()
}
