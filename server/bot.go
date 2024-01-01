package server

import (
	"time"
	"strings"

	"github.com/williamjchen/terminalchess/stockfish"
)
type bot interface {
	Name() string
	GetMove() string
}

type basicBot struct {
	name string
	lob *lobby
}

type stockfishBot struct {
	name string
	lob *lobby
	depth int
}


// BASIC
func NewBasicBot(name string, lob *lobby) basicBot {
	b := basicBot{
		name: name,
		lob: lob,
	}

	return b
}

func (b basicBot) GetMove() string {
	time.Sleep(1 * time.Second)
	return b.lob.game.GetRandomMove()
}

func (b basicBot) Name() string {
	return b.name
}

// STOCKFISH 
func NewStockfishBot(name string, depth int, lob *lobby) stockfishBot {
	b := stockfishBot{
		name: name,
		lob: lob,
		depth: depth,
	}

	return b
}

func (b stockfishBot) GetMove() string {
	hist := b.lob.game.GetMoveHistory()
	return stockfish.Move(strings.Join(hist, " "), b.depth)
}

func (b stockfishBot) Name() string {
	return b.name
}