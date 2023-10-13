package server

import (
	"github.com/charmbracelet/lipgloss"
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
	height int
	width int
}

func GetModelOption(s ssh.Session, options []string, server *Server, sess ssh.Session) *tea.Program {
	model := Model(options, server, sess)
    p := tea.NewProgram(
        model,
        tea.WithInput(s),
        tea.WithOutput(s),
    )
	model.common.program = p
	return p
    // _, err := p.Run()
    // if err != nil {
    //     slog.Error("failed to run menu", err)
    //     return
    // }
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
		switch msg.String() {
		case "esc", "ctrl+c":
			if m.common.player.lob != nil {
				m.common.player.lob.RemovePlayer(m.common.player)
			}
			return m, tea.Quit
		case "ctrl+n":
			if m.common.player.lob != nil {
				m.common.player.lob.RemovePlayer(m.common.player)
			}
			m.Reset()
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	switch m.state{
	case showMenu:
		var cmd2 tea.Cmd
		men, cmd := m.menu.Update(msg)
		m.menu = men.(*menuModel)
		if m.common.chosen {
			m.state = showGame
			cmd2 = m.game.Init()
		}
		return m, tea.Batch(cmd, cmd2)
	case showGame:
		g, cmd := m.game.Update(msg)
		m.game = g.(*gameModel)
		return m, cmd
	}
	return m, nil
}

func (m *parentModel) View() string {
	s := lipgloss.NewStyle().Height(m.height).Width(m.width)
	switch m.state{
	case showMenu:
		return s.Render(m.menu.View())
	case showGame:
		return s.Render(m.game.View())
	}
	return ""
}
