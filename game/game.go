package game

type Game struct {
	board board
}

func (g *Game) PrintBoard() {
	g.board.PrintBoard()
}