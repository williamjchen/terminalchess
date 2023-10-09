package server

import (
	"strings"

	Game "github.com/williamjchen/terminalchess/game"	

	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
    game *Game.Game
}

func NewGame() game {
	g := game{
		game: Game.NewGame(),
	}
	return g
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch m.choice {
	case m.choices[0]: // stockfish
		return stockfishUpdate(msg, m)
	case m.choices[1]: // join
		return joinUpdate(msg, m)
	case m.choices[2]: // create
		return createUpdate(msg, m)
	default:
		return m, nil
	}
}

func chosenView(m model) string {
	switch m.choice {
	case m.choices[0]: // stockfish
		return stockfishView(m)
	case m.choices[1]: // join
		return joinView(m)
	case m.choices[2]: // create
		return createView(m)
	default:
		return ""
	}
}

func stockfishUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

func stockfishView(m model) string {
	return ""
}

func joinUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	m.idInput.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.idInput.Reset()
			return m, cmd
		}

	case errMsg:
		return m, nil
	}

	m.idInput, cmd = m.idInput.Update(msg)
	return m, cmd
}

func joinView(m model) string {
	s := strings.Builder{}
	s.WriteString("Enter Room Code...\n")
	s.WriteString(m.idInput.View())
	return s.String()
}

func createUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

func createView(m model) string {
	return ""
}

func gameUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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

func gameView(m model) string {
	s := strings.Builder{}

	s.WriteString(m.lobby.game.game.PrintBoard())
	s.WriteString("\n\n")
	
	if m.lobby.game.game.WhiteTurn() {
		s.WriteString("White to move\n")
	} else {
		s.WriteString("Black to Move\n")
	}

	s.WriteString(m.textinput.View())
	return s.String()
}
