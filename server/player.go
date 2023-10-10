package server

type playerType int

const (
	white playerType = iota // 0
	black // 1
	spec // 2
)
type player struct {
	name string
	common *commonModel
	playerType playerType
	lob *lobby
	flipped bool
}

func NewPlayer(com *commonModel) *player {
	p := player{
		name: "Anonymous",
		common: com,
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

func (p *player) SetFlipped(flipped bool) {
	p.flipped = flipped
}

func (p *player) Flip() {
	p.flipped = !p.flipped
}
