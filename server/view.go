package server

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
)

type (
	errMsg error
)

type state int

const (
	showMenu state = iota
	showGame
)

type commonModel struct {
	choices  []string           // items on the to-do list
    choice string
	chosen bool
	begin bool
	srv *Server
	player *player
	sess ssh.Session
	program *tea.Program
}
type parentModel struct {
	state state
	common *commonModel
	menu *menuModel
	game *gameModel
}

func GetModelOption(s ssh.Session, options []string, server *Server, sess ssh.Session) {
	model := Model(options, server, sess)
    p := tea.NewProgram(
        model,
        tea.WithInput(s),
        tea.WithOutput(s),
    )
	model.common.program = p
    _, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return
    }
}

func Model(options []string, server *Server, sess ssh.Session) *parentModel {
	common := commonModel {
		choices: options,
		choice: "",
		chosen: false,
		begin: false,
		srv: server,
		sess: sess,
	}

	common.player = NewPlayer(&common)

	p := parentModel{
		common: &common,
		menu: NewMenu(&common),
		game: NewGame(&common),
	}

	return &p
}

func (m *parentModel) Reset() {
	m.state = showMenu
	m.common.choice = ""
	m.common.chosen = false
	m.menu.cursor = 0
	m.common.player.lob = nil
}

func (m *parentModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *parentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		var cmd tea.Cmd
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			cmd = tea.Quit
			m.common.player.lob.RemovePlayer(m.common.player)
			return m, cmd
		} else if k == "ctrl+n" {
			m.common.player.lob.RemovePlayer(m.common.player)
			m.Reset()
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		slog.Info("change size", "width", msg.Width, "height", msg.Height)
	}

	switch m.state{
	case showMenu:
		men, cmd := m.menu.Update(msg)
		m.menu = men.(*menuModel)
		if m.common.chosen {
			m.state = showGame
		}
		return m, cmd
	case showGame:
		g, cmd := m.game.Update(msg)
		m.game = g.(*gameModel)
		return m, cmd
	}
	return m, nil
}

func (m *parentModel) View() string {
	switch m.state{
	case showMenu:
		return m.menu.View()
	case showGame:
		return m.game.View()
	}
	return ""
}
