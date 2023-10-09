package game

type Game struct {
	board board
}

func (g *Game) printBoard() {
	g.board.printBoard()
}