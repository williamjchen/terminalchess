package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math"
)
type lobby struct {
	id string
	players []*player
	move bool
	done chan string
}

func NewLobby(chan string) (*lobby, string) {
	id := randId(6)
	l := lobby{
		id: id,
	}
	return &l, id
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
