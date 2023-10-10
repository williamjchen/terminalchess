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
	spinner spinnerModel
	textinput textinput.Model
}

func NewGame(com *commonModel) gameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter move in algebraic notation..."
	ti.CharLimit = 5
	ti.Width = 20

	g := gameModel{
		common: com,
		stockfish: NewStockfishModel(com),
		join: NewJoinModel(com),
		create: NewCreateModel(com),
		spinner: NewSpinner(),
		textinput: ti,
	}
	return g
}

func (m gameModel) Init() tea.Cmd {
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.common.player.lob != nil {
		return gameUpdate(msg, m)
	}

	switch m.common.choice {
	case m.common.choices[1]: // join
		j, cmd := m.join.Update(msg)
		m.join = j.(joinModel)
		return m, cmd
	default:
		return m, nil
	}
}

func (m gameModel) View() string {
	if m.common.player.lob != nil {
		return gameView(m)
	}

	switch m.common.choice {
	case m.common.choices[1]: // join
		return m.join.View()
	default:
		return ""
	}
}

func gameUpdate(msg tea.Msg, m gameModel) (tea.Model, tea.Cmd) {
	m.textinput.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.textinput.Reset()
			return m, cmd
		case "ctrl+f":
			m.common.player.lob.game.Flip()
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

	s.WriteString(m.common.player.lob.game.PrintBoard())
	s.WriteString("\n\n")
	
	if m.common.player.lob.game.WhiteTurn() {
		s.WriteString("White to move\n")
	} else {
		s.WriteString("Black to Move\n")
	}

	s.WriteString(m.textinput.View())
	return s.String()
}
