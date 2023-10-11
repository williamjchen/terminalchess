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
	p1Pres bool
	p2Pres bool
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
	l.SendMsgEveryone(finishMsg(2))
}

func (l *lobby) AddPlayer(s ssh.Session, p *player) { // return 0 if white, 1 if black, 2 if spectator
	if s.User() != "" {
		p.name = s.User()
	}
	
	if l.p1 == nil {
		l.p1 = p
		p.playerType = white
		l.p1Pres = true
	} else if l.p2 == nil {
		l.p2 = p
		p.playerType = black
		l.p1Pres = true
	} else {
		l.specs = append(l.specs, p)
		p.playerType = spec
	}

	slog.Info("Player added", "lobby id:", l.id, "type:", p.playerType, "name:", p.name)

}

func (l *lobby) RemovePlayer(p *player) {
	if l.p1 == p {
		l.p1 = nil
		l.p1Pres = true
		l.SendMsg(p, finishMsg(1))
		l.SendMsgToSpectators(finishMsg(1))
	} else if l.p2 == p {
		l.p2 = nil
		l.p2Pres = true
		l.SendMsg(p, finishMsg(0))
		l.SendMsgToSpectators(finishMsg(0))
	} else {
		length := len(l.specs)
		for i := 0; i < length; i++ {
			if l.specs[i] == p {
				l.specs[i] = l.specs[length-1]
				l.specs[len(l.specs)-1] = nil
				l.specs = l.specs[:length-1]
				return
			}
		}
	}
}

func (l *lobby) SendMsg(p *player, msg interface{}) { // sends message to other player that's not the argument
	if l.p1 == p {
		l.p2.common.program.Send(msg)
	} else {
		l.p1.common.program.Send(msg)
	}
	l.SendMsgToSpectators(msg)
}

func (l *lobby) SendMsgToSpectators(msg interface{}) {
	for _, p := range l.specs {
		p.common.program.Send(msg)
	}
}

func (l *lobby) SendMsgEveryone(msg interface{}) {
	l.p1.common.program.Send(msg)
	l.p2.common.program.Send(msg)
	l.SendMsgToSpectators(msg)
}
 