package server

import (
	"strings"

	Game "github.com/williamjchen/terminalchess/game"	

	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
    game *Game.Game
	turn bool

}

func NewGame() *game {
    g := game {
		game: Game.NewGame(),
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
		}
	}

	return m, nil
}

func (m game) View() string {
	s := strings.Builder{}
	
	s.WriteString(m.game.PrintBoard())
	return s.String()
}
