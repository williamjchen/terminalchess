package server

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type gameModel struct {
	common *commonModel
	stockfish stockfishModel
	join joinModel
	create createModel
	textinput textinput.Model
    lobby *lobby
}

func NewGame(com *commonModel) gameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter move in algebraic notation..."
	ti.CharLimit = 5
	ti.Width = 20

	t := make(chan string)
	l, _ := NewLobby(t)

	g := gameModel{
		common: com,
		stockfish: NewStockfishModel(com),
		join: NewJoinModel(com),
		create: NewCreateModel(com),
		textinput: ti,
		lobby: l,
	}
	return g
}

func (m gameModel) Init() tea.Cmd {
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.common.begin {
		return gameUpdate(msg, m)
	}

	switch m.common.choice {
	case m.common.choices[0]: // stockfish
		s, cmd := m.stockfish.Update(msg)
		m.stockfish = s.(stockfishModel)
		return m, cmd
	case m.common.choices[1]: // join
		j, cmd := m.join.Update(msg)
		m.join = j.(joinModel)
		return m, cmd
	case m.common.choices[2]: // create
		c, cmd := m.create.Update(msg)
		m.create = c.(createModel)
		return m, cmd
	default:
		return m, nil
	}
}

func (m gameModel) View() string {
	if m.common.begin {
		return gameView(m)
	}

	switch m.common.choice {
	case m.common.choices[0]: // stockfish
		return m.stockfish.View()
	case m.common.choices[1]: // join
		return m.join.View()
	case m.common.choices[2]: // create
		return m.create.View()
	default:
		return ""
	}
}

func gameUpdate(msg tea.Msg, m gameModel) (tea.Model, tea.Cmd) {
	m.textinput.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.textinput.Reset()
			return m, cmd
		}

	case errMsg:
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func gameView(m gameModel) string {
	s := strings.Builder{}

	s.WriteString(m.lobby.game.PrintBoard())
	s.WriteString("\n\n")
	
	if m.lobby.game.WhiteTurn() {
		s.WriteString("White to move\n")
	} else {
		s.WriteString("Black to Move\n")
	}

	s.WriteString(m.textinput.View())
	return s.String()
}
