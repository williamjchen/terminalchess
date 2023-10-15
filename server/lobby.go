package server

import (
	"log/slog"
	"time"

	"github.com/williamjchen/terminalchess/models"
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
	specs []*player
	game *Game.Game
	timer *time.Timer
	done chan string
	bot *player
	gameModel *models.Game
}

func NewLobby(done chan string, id string) *lobby {
	timer := time.NewTimer(2 * time.Hour)

	g := Game.NewGame()
	l := lobby{
		id: id,
		p1: nil,
		p2: nil,
		game: g,
		timer: timer,
		done: done,
		bot: nil,
		gameModel: &models.Game{Code:id},
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

func (l *lobby) Status() gameState{
	var s gameState
	switch l.game.Turn() {
	case Game.WhiteTurn:
		s = inProgres
	case Game.BlackTurn:
		s = inProgres
	case Game.WhiteMate:
		s = whiteWin
	case Game.BlackMate:
		s = blackWin
	case Game.Stalemate:
		s = stalemate
	}
	return s
}

func (l *lobby) SetStatus(stat gameState) {
	switch stat {
	case whiteWin:
		l.game.SetStatus(Game.WhiteMate)
	case blackWin:
		l.game.SetStatus(Game.BlackMate)
	}
}

func (l *lobby) AddBot(b bot) {
	p := &player{}
	p.name = b.Name()
	l.bot = p
	p.bot = b
	p.lob = l
	
	if l.p1 == nil {
		l.p1 = p
		p.playerType = white
		l.p1Pres = true
		l.gameModel.Player1Name = p.name
		p.common = l.p2.common
	} else if l.p2 == nil {
		l.p2 = p
		p.playerType = black
		l.p2Pres = true
		l.gameModel.Player2Name = p.name
		p.common = l.p1.common
	} 

	slog.Info("Bot added", "lobby id:", l.id, "type:", p.playerType, "name:", p.name)
}

func (l *lobby) AddPlayer(s ssh.Session, p *player) { // return 0 if white, 1 if black, 2 if spectator
	if s.User() != "" {
		p.name = s.User()
	}
	
	if l.p1 == nil {
		l.p1 = p
		p.playerType = white
		l.p1Pres = true
		l.gameModel.Player1Name = p.name
	} else if l.p2 == nil {
		l.p2 = p
		p.playerType = black
		l.p2Pres = true
		l.gameModel.Player2Name = p.name
	} else {
		l.specs = append(l.specs, p)
		p.playerType = spec
	}
	go p.common.srv.db.Games.Update(l.gameModel)

	slog.Info("Player added", "lobby id:", l.id, "type:", p.playerType, "name:", p.name)
}

func (l *lobby) RemovePlayer(p *player) {
	if l.p1 == p {
		slog.Info("Remove player 1", "id", l.id)
		l.p1Pres = false
		l.SetStatus(blackWin)
		l.SendMsg(p, finishMsg(1))
		l.p1 = nil
		l.SendMsgToSpectators(finishMsg(1))
	} else if l.p2 == p {
		slog.Info("Remove player 2", "id", l.id)
		l.p2Pres = false
		l.SetStatus(whiteWin)
		l.SendMsg(p, finishMsg(0))
		l.p2 = nil
		l.SendMsgToSpectators(finishMsg(0))
	} else {
		length := len(l.specs)
		for i := 0; i < length; i++ {
			if l.specs[i] == p {
				l.specs[i] = l.specs[length-1]
				l.specs[len(l.specs)-1] = nil
				l.specs = l.specs[:length-1]
				slog.Info("Remove spec", "id", l.id)
				return
			}
		}
		slog.Info("Player not found", "id", l.id)
	}
}

func (l *lobby) sendMove(move string, p *player) bool {
	if p != l.p1 && p != l.p2 {
		return false
	}
	if l.Status() != inProgres {
		return false
	}

	status := p.Move(move)
	if status {
		l.gameModel.Moves = l.game.GetMoveHistory()
		go p.common.srv.db.Games.Update(l.gameModel)
	}
	l.SendMsg(p, updateMsg{})
	return status
}

func (l *lobby) SendMsg(p *player, msg interface{}) { // sends message to other player that's not the argument
	if l.p2 != nil && l.p1 == p && l.p2 != l.bot {
		l.p2.common.program.Send(msg)
	} else if l.p1 != nil && l.p2 == p && l.p1 != l.bot {
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
 