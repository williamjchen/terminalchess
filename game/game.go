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

func (g *Game) WhiteMove(cmd string) bool {
	return true
}

func (g *Game) BlackMove(cmd string) bool {
	return true
}

func (g *Game) SetFlipped(flipped bool) {
	g.board.SetFlipped(flipped)
}

func (g *Game) Flip() {
	g.board.Flip()
}