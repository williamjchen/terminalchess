package server

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type joinModel struct {
	common *commonModel
	idInput textinput.Model
	tried bool
}

type lobMsg *lobby

func NewJoinModel(com *commonModel) *joinModel {
	ii := textinput.New()
	ii.Placeholder = "XXXXXX"
	ii.CharLimit = 6
	ii.Width = 20

	j := joinModel {
		common: com,
		idInput: ii,
		tried: false,
	}

	return &j
}

func getLobby(id string, srv *Server) tea.Cmd {
	return func() tea.Msg {
		l := srv.mng.FindLobby(id)
		return lobMsg(l)
	}
}

func (m joinModel) Init() tea.Cmd {
    return nil
}

func (m *joinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.idInput.Focus()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String(){
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			s := m.idInput.Value()
			m.idInput.Reset()
			return m, getLobby(s, m.common.srv)
		default:
			var cmd tea.Cmd
			slog.Info("blink")
			m.idInput, cmd = m.idInput.Update(msg)
			return m, tea.Batch(textinput.Blink, cmd)
		} // switch KeyMsg

	case lobMsg:
		if msg == nil {
			m.tried = true
			return m, nil
		} else {
			var l *lobby = msg

			l.AddPlayer(m.common.sess, m.common.player)
			m.common.player.lob = l
			m.common.player.SetFlipped(m.common.player.playerType == black)

			l.SendMsg(m.common.player, updateMsg{})
			return m, nil
		}

	case errMsg:
		return m, nil

	} // switch msg
	return m, nil
}

func (m *joinModel) View() string {
	s := strings.Builder{}
	s.WriteString("Enter Room Code...\n")
	s.WriteString(m.idInput.View())
	if m.tried {
		s.WriteString("\nInvalid! Try again or create your own room!")
	}
	return s.String()
}
