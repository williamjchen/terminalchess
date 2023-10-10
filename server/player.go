package server

type playerType int

const (
	white playerType = iota // 0
	black // 1
	spec // 2
)
type player struct {
	name string
	playerType playerType
	lob *lobby
}

func NewPlayer() *player {
	p := player{
		name: "Anonymous",
		lob: nil,
	}

	return &p
}

func (p *player) Move(cmd string) bool {
	if (p.playerType == white) {
		return p.lob.game.WhiteMove(cmd)
	} else if (p.playerType == black) {
		return p.lob.game.BlackMove(cmd)
	}
	return false
}
