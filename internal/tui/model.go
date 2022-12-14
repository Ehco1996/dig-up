package tui

import (
	"time"

	"github.com/Ehco1996/dig-up/pkg/bc"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	c   *bc.Client
	err error

	upUID, favID           int
	currentPage, totalPage int

	autoCheck bool

	table     table.Model
	tableRows *[]table.Row
}

func InitialModel(curlString string, upUID, favID int) (model, error) {
	c, err := bc.NewClient(curlString)
	if err != nil {
		return model{}, err
	}
	return model{c: c, upUID: upUID, favID: favID}, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			m.err = m.FetchVideoPage(m.currentPage+1, PageSize)
		case "enter":
			m.err = m.CheckAndSave()
		case "p":
			m.CheckAndMoveCursor()
		case "o":
			m.autoCheck = !m.autoCheck
			if m.autoCheck {
				return m, tick()
			}
		}
	case tickMsg:
		if m.autoCheck {
			m.CheckAndMoveCursor()
			return m, tick()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
