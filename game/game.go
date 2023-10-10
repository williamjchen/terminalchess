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

func (g *Game) PrintBoard() string{ // flipped = false is white at bottom
	return g.board.PrintBoard()
}

func (g *Game) Move(cmd string) {

}

func (g *Game) SetFlipped(flipped bool) {
	g.board.SetFlipped(flipped)
}