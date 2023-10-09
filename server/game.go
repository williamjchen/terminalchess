package server

import (
	"log/slog"

	Game "github.com/williamjchen/terminalchess/game"	

	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
    game Game.Game
	turn bool

}

func NewGame() *game {
    g := game {

	}
	return &g
}

func (m game) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return tea.EnterAltScreen
}

func (m game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m game) View() string {
	s := strings.Builder{}
	
	s.WriteString(m.game.)
	return s.String()
}