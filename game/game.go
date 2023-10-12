package game

import (
	"log/slog"
	"strings"
)

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
	valid, origin, dest := parseMove(move)
	slog.Info("move", "valid", valid, "origin", origin, "dest", dest)
	if !valid {
		return false
	}
	g.board.move(origin, dest)
	return true
}

func (g *Game) BlackMove(move string) bool {
	valid, origin, dest := parseMove(move)
	if !valid {
		return false
	}
	g.board.move(origin, dest)
	return true
}


// TODO - promotion
func parseMove(move string) (bool, int, int) {
	move = strings.ToLower(move)
	if len(move) >= 4 {
		fromRank := rune(move[0])
		fromFile := rune(move[1])
		toRank := rune(move[2])
		toFile := rune(move[3])

		if fromRank < 'a' || fromRank > 'h' || fromFile < '1' || fromFile > '8' {
			return false, -1, -1
		}
		if toRank < 'a' || toRank > 'h' || toFile < '1' || toFile > '8' {
			return false, -1, -1
		}

		return true, int((fromFile - '1') * 8 + (fromRank - 'a')), int((toFile - '1') * 8 + (toRank - 'a'))
	}
	return false, -1, -1
}
