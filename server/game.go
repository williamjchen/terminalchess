package server

import (
	"time"
	"strings"
	"log/slog"

	Game "github.com/williamjchen/terminalchess/game"	

	"github.com/charmbracelet/bubbles/cursor"
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
	ti.Placeholder = "Enter move in long algebraic notation ex. e2e4 | e1g1 | e7e8q ..."
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
	if !m.join.idInput.Focused() {
		cmd := tea.Batch(m.join.idInput.Cursor.SetMode(cursor.CursorBlink),
				m.join.idInput.Focus(),
			)
		return cmd
	}
	return nil
}

func (m *gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.common.player.lob != nil {
		var cmd tea.Cmd
		if !m.textinput.Focused() {
			cmd = tea.Batch(m.textinput.Cursor.SetMode(cursor.CursorBlink),
					m.textinput.Focus(),
				)
		}
		m2, cmd2 := gameUpdate(msg, m)
		return m2, tea.Batch(cmd, cmd2)
	}

	var cmd tea.Cmd
	switch m.common.choice {
	case m.common.choices[1]: // join
		j, cmd := m.join.Update(msg)
		m.join = j.(*joinModel)
		return m, tea.Batch(cmd)
	default:
		return m, cmd
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := m.textinput.Value()
			m.textinput.Reset()
			return m, sendMove(val, m)
		case "ctrl+f":
			m.common.player.Flip()
			return m, cmd
		}

	case errMsg:
		return m, nil

	case chessMsg:
		m.validMove = bool(msg)
		if m.validMove && m.common.player.lob.bot != nil {
			go func() {
				time.Sleep(1 * time.Second)
				botMove :=  m.common.player.lob.bot.bot.GetMove()
				slog.Info("Bot move", "move", botMove, "bot", m.common.player.lob.bot)
				m.common.player.lob.sendMove(botMove, m.common.player.lob.bot)
			} ()
			
		}

	case updateMsg:
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func gameView(m *gameModel) string {
	s := strings.Builder{}

	s.WriteString("\n")
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
