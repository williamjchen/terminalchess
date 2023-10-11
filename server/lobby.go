package server

import (
	"log/slog"
	"time"

	Game "github.com/williamjchen/terminalchess/game"	

	"github.com/charmbracelet/ssh"
)

type gameState int

const (
	inProgres gameState = iota
	whiteWin
	blackWin
	stalemate
)
type lobby struct {
	id string
	p1 *player // white
	p2 *player // black
	status gameState
	specs []*player
	game *Game.Game
	timer *time.Timer
	done chan string
}

func NewLobby(done chan string, id string) *lobby {
	timer := time.NewTimer(2 * time.Hour)

	g := Game.NewGame()
	l := lobby{
		id: id,
		p1: nil,
		p2: nil,
		status: inProgres,
		game: g,
		timer: timer,
		done: done,
	}
	
	go func() {
        <-timer.C
		done<-"done"
    }()

	return &l
}

func (l *lobby) End() {
	l.timer.Stop()
	l.done<-"done"
	l.SendMsgEveryone(finishMsg{})
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

func (l *lobby) SendMsgEveryone(msg struct{}) {
	l.p1.common.program.Send(msg)
	l.p2.common.program.Send(msg)
	l.SendMsgToSpectators(msg)
}
 