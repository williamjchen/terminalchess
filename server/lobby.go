package server

import (
	"log/slog"

	Game "github.com/williamjchen/terminalchess/game"	

	"github.com/charmbracelet/ssh"
)
type lobby struct {
	id string
	p1 *player // white
	p2 *player // black
	specs []*player
	game *Game.Game
	done chan string
}

func NewLobby(done chan string, id string) *lobby {
	g := Game.NewGame()
	l := lobby{
		id: id,
		p1: nil,
		p2: nil,
		game: g,
		done: done,
	}
	return &l
}

func (l *lobby) AddPlayer(s ssh.Session, p *player) { // return 0 if white, 1 if black, 2 if spectator
	if s.User() != "" {
		p.name = s.User()
	}
	
	if l.p1 == nil {
		l.p1 = p
		p.playerType = white
	} else if l.p2 == nil {
		l.p2 = p
		p.playerType = black
	} else {
		l.specs = append(l.specs, p)
		p.playerType = spec
	}

	slog.Info("Player added", "lobby id:", l.id, "type:", p.playerType, "name:", p.name)

}

func (l *lobby) SendMsg(p *player, msg struct{}) { // sends message to other player that's not the argument
	if l.p1 == p {
		l.p2.common.program.Send(msg)
	} else {
		l.p1.common.program.Send(msg)
	}
	l.SendMsgToSpectators(msg)
}

func (l *lobby) SendMsgToSpectators(msg struct{}) {
	for _, p := range l.specs {
		p.common.program.Send(msg)
	}
}
