package server

type playerType int

const (
	white playerType = iota
	black 
	spec
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