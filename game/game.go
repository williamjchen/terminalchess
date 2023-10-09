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

func (g *Game) PrintBoard() string{
	return g.board.PrintBoard()
}