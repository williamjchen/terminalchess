package server

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
	return b.lob.game.GetRandomMove()
}

func (b basicBot) Name() string {
	return b.name
}
