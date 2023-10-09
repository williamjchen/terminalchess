package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math"

	"github.com/charmbracelet/ssh"
)
type lobby struct {
	id string
	p1 *player
	p2 *player
	specs []*player
	game *game
	done chan string
}

func NewLobby(done chan string) (*lobby, string) {
	id := randId(6)
	l := lobby{
		id: id,
		p1: nil,
		p2: nil,
		game: nil,
		done: done,
	}
	return &l, id
}

func (l *lobby) AddPlayer(s ssh.Session) {
	name := "Anonymous"
	if s.User() != "" {
		name = s.User()
	}

	p := player{
		name: name,
	}
	
	if l.p1 == nil {
		l.p1 = &p
	} else if l.p2 == nil {
		l.p2 = &p
	} else {
		l.specs = append(l.specs, &p)
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
