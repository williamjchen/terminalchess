package server

type player struct {
	name string
}

func NewPlayer() player {
	p := player{
		name: "Anonymous",
	}
	
	return p
}