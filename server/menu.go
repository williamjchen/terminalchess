package server

import (
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menuModel struct {
	common *commonModel
	cursor int
}

func NewMenu(com *commonModel) menuModel {
	m := menuModel {
		common: com,
	}
	return m
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.common.choice = m.common.choices[m.cursor]
			m.common.chosen = true

			switch m.common.choice {
			case m.common.choices[0]: // stockfish
				l := m.common.srv.mng.CreateLobby()
				m.common.player.lob = l
				return m, nil
			case m.common.choices[2]: // create
				l := m.common.srv.mng.CreateLobby()
				slog.Info("Create lobby", "id:", l.id)

				l.AddPlayer(m.common.sess, m.common.player)
				m.common.player.lob = l
				m.common.player.lob.game.SetFlipped(m.common.player.playerType != black)
				return m, nil
			default:
				return m, nil
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.common.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.common.choices) - 1
			}
		}
	}
	return m, nil
}



func (m menuModel) View() string {
	s := strings.Builder{}
	s.WriteString("What chess mode would you like to play?\n\n")

	for i := 0; i < len(m.common.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.common.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press esc to quit)\n")

	return s.String()
}