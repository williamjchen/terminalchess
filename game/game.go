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

func (g *Game) PrintBoard() string{
	return g.board.PrintBoard()
}

func (g *Game) Move(cmd string) {

}
