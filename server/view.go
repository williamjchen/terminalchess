package server

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
)

type (
	errMsg error
)

type state int

const (
	showMenu state = iota
	showGame
)

type commonModel struct {
	choices  []string           // items on the to-do list
    choice string
	chosen bool
	begin bool
	srv *Server
	gameState *gameState
}
type parentModel struct {
	state state
	common *commonModel
	menu menuModel
	game gameModel
}

func GetModelOption(s ssh.Session, options []string, server *Server) {
    p := tea.NewProgram(
        Model(options, server),
        tea.WithInput(s),
        tea.WithOutput(s),
    )
    _, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return
    }
}

func Model(options []string, server *Server) parentModel {
	gs := gameState{
		lobby: nil,
	}

	common := commonModel {
		choices: options,
		choice: "",
		chosen: false,
		begin: false,
		srv: server,
		gameState: &gs,
	}

	p := parentModel{
		common: &common,
		menu: NewMenu(&common),
		game: NewGame(&common),
	}

	return p
}

func (m parentModel) Init() tea.Cmd {
    return tea.Batch(tea.EnterAltScreen, m.game.spinner.spinner.Tick, textinput.Blink)
}

func (m parentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.state{
	case showMenu:
		men, cmd := m.menu.Update(msg)
		m.menu = men.(menuModel)
		if m.common.chosen {
			m.state = showGame
		}
		return m, cmd
	case showGame:
		g, cmd := m.game.Update(msg)
		m.game = g.(gameModel)
		return m, cmd
	}
	return m, nil
}

func (m parentModel) View() string {
	switch m.state{
	case showMenu:
		return m.menu.View()
	case showGame:
		return m.game.View()
	}
	return ""
}
