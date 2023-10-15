package server

import (
	"time"
)
type bot interface {
	Name() string
	GetMove() string
}

type basicBot struct {
	name string
	lob *lobby
}

func NewBasicBot(name string, lob *lobby) basicBot {
	b := basicBot{
		name: "Rudolph Bot",
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
