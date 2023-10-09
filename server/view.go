package server

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
)

type (
	errMsg error
)

const (
	showMenu int = iota
	showGame
)

type commonModel struct {
	choices  []string           // items on the to-do list
    choice string
	chosen bool
}
type parentModel struct {
	state int
	common *commonModel
	menu menuModel
	game gameModel
}

func GetModelOption(s ssh.Session, options []string) {
    p := tea.NewProgram(
        Model(options),
        tea.WithInput(s),
        tea.WithOutput(s),
    )
    _, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return
    }
}

func Model(options []string) parentModel {
	common := commonModel {
		choices: options,
		choice: "",
		chosen: false,
	}

	p := parentModel{
		common: &common,
		menu: NewMenu(&common),
		game: NewGame(&common),
	}

	return p
}

func (m parentModel) Init() tea.Cmd {
    return tea.EnterAltScreen
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
