package server

import (
	"strings"

	Game "github.com/williamjchen/terminalchess/game"	
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chessMsg bool
type updateMsg struct{}
type finishMsg int // 0 white, 1 black, 2 stalemate

type gameModel struct {
	common *commonModel
	info *infoModel
	stockfish *stockfishModel
	join *joinModel
	create *createModel
	spinner spinnerModel
	textinput textinput.Model
	validMove bool
}

func NewGame(com *commonModel) *gameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter move in algebraic notation..."
	ti.CharLimit = 5
	ti.Width = 20

	g := gameModel{
		common: com,
		info: NewInfoModel(com),
		stockfish: NewStockfishModel(com),
		join: NewJoinModel(com),
		create: NewCreateModel(com),
		spinner: NewSpinner(),
		textinput: ti,
		validMove: true,
	}
	return &g
}

func sendMove(move string, m *gameModel) tea.Cmd {
	return func() tea.Msg {
		status := m.common.player.lob.sendMove(move, m.common.player)
		return chessMsg(status)
	}
}

func (m *gameModel) Init() tea.Cmd {
	return nil
}

func (m *gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.common.player.lob != nil {
		return gameUpdate(msg, m)
	}

	switch m.common.choice {
	case m.common.choices[1]: // join
		j, cmd := m.join.Update(msg)
		m.join = j.(*joinModel)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *gameModel) View() string {
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

func gameUpdate(msg tea.Msg, m *gameModel) (tea.Model, tea.Cmd) {
	m.textinput.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.textinput.Reset()
			return m, sendMove("", m)
		case "ctrl+f":
			m.common.player.Flip()
			return m, cmd
		}

	case errMsg:
		return m, nil

	case chessMsg:
		m.validMove = bool(msg)

	case updateMsg:
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func gameView(m *gameModel) string {
	s := strings.Builder{}

	s.WriteString(m.common.player.lob.game.PrintBoard(m.common.player.flipped))
	s.WriteString("\n\n")
	
	if m.common.player.lob.game.Turn() == Game.WhiteTurn {
		s.WriteString("White to move\n")
	} else {
		s.WriteString("Black to Move\n")
	}

	b := strings.Builder{}
	b.WriteString(m.textinput.View())

	if !m.validMove {
		b.WriteString("\nInvalid Move! Try again...")
	}

	return lipgloss.JoinHorizontal(0.38, s.String(), m.info.View(m.common.player.flipped)) + "\n" + b.String()
}
