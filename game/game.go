package game

type Game struct {
	board *board
}

func NewGame() *Game{
	g := Game{
		board: NewBoard(),
	}
	return &g
}

func (g *Game) WhiteTurn() bool {
	return g.board.pos.whiteTurn
}

func (g *Game) PrintBoard(flipped bool) string{ // flipped = false is white at bottom
	return g.board.PrintBoard(flipped)
}

func (g *Game) WhiteMove(cmd string) bool {
	return true
}

func (g *Game) BlackMove(cmd string) bool {
	return true
}

func (g *Game) whiteTurn() bool {
	return g.board.pos.whiteTurn
}
