package server

import (
	"log/slog"
	"crypto/rand"
	"fmt"
	"math/big"
	"math"

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

func NewLobby(done chan string) (*lobby, string) {
	id := randId(6)
	g := Game.NewGame()
	l := lobby{
		id: id,
		p1: nil,
		p2: nil,
		game: g,
		done: done,
	}
	return &l, id
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

func randId(length int) string {
	// https://stackoverflow.com/a/75518426/7361588
	bi, err := rand.Int(
        rand.Reader,
        big.NewInt(int64(math.Pow(10, float64(length)))),
    )
    if err != nil {
        panic(err)
    }
    return fmt.Sprintf("%0*d", length, bi)
}
