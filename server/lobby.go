package server

import (
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

func (l *lobby) AddPlayer(s ssh.Session) *player { // return 0 if white, 1 if black, 2 if spectator
	name := "Anonymous"
	if s.User() != "" {
		name = s.User()
	}

	p := player{
		name: name,
	}
	
	if l.p1 == nil {
		l.p1 = &p
		p.playerType = white
		return &p
	} else if l.p2 == nil {
		l.p2 = &p
		p.playerType = black
		return &p
	} else {
		l.specs = append(l.specs, &p)
		p.playerType = spec
		return &p
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
