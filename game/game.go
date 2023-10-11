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

func (g *Game) Turn() turn {
	return g.board.pos.turn
}

func (g *Game) PrintBoard(flipped bool) string{ // flipped = false is white at bottom
	return g.board.PrintBoard(flipped)
}

func (g *Game) WhiteMove(move string) bool {
	g.board.move()
	return true
}

func (g *Game) BlackMove(move string) bool {
	g.board.move()
	return true
}
