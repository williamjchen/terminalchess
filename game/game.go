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

func (g *Game) Move(move string) bool {
	valid, moveObj := parseMove(move)
	slog.Info("parse command", "valid", valid, "origin", moveObj.origin(), "dest", moveObj.dest(), "promotion", moveObj.promotion())
	if !valid {
		return false
	}
	return g.board.move(moveObj)
}

func (g *Game) SetStatus(stat turn) {
	g.board.SetStatus(stat)
}

func (g *Game) GetRandomMove() string {
	return g.board.getRandomMove()
}

func (g *Game) GetMoveHistory() []string {
	return g.board.getMoveHistory()
}
// TODO - promotion
func parseMove(moveToParse string) (bool, move) {
	moveToParse = strings.ToLower(moveToParse)
	promotion := 'p'
	var newMove move

	if len(moveToParse) >= 4 {
		fromRank := rune(moveToParse[0])
		fromFile := rune(moveToParse[1])
		toRank := rune(moveToParse[2])
		toFile := rune(moveToParse[3])

		if fromRank < 'a' || fromRank > 'h' || fromFile < '1' || fromFile > '8' {
			return false, newMove
		}
		if toRank < 'a' || toRank > 'h' || toFile < '1' || toFile > '8' {
			return false, newMove
		}

		if len(moveToParse) >= 5 {
			promotion = rune(moveToParse[4])
		}

		newMove.create(int((fromFile - '1') * 8 + (fromRank - 'a')), int((toFile - '1') * 8 + (toRank - 'a')), promotion)
		return true, newMove
	}
	return false, newMove
}
