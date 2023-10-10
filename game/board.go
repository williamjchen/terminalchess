package game

import (
	"strings"
	"strconv"
)

type board struct {
	flipped bool
	pieces_unicode map[string] string
	upT, cross, botT, leftT, rightT, leftBotEnd, rightBotEnd, rightUpEnd, leftUpEnd, horiz, vert string
	padding int
	pos *position
}

func NewBoard() *board {
	b := board{}
	b.pieces_unicode = map[string]string {
		"R": "♜", 
        "N": "♞", 
        "B": "♝", 
        "Q": "♛", 
        "K": "♚", 
        "P": "♟", 
        "r": "♖", 
        "n": "♘", 
        "b": "♗", 
        "q": "♕", 
        "k": "♔",
        "p": "♙",
        " ": " ",
	}
	b.upT = "┬";
    b.cross = "┼";
    b.botT = "┴";
    b.leftT = "├";
    b.rightT = "┤";
    b.leftBotEnd = "└";
    b.rightBotEnd = "┘";
    b.rightUpEnd = "┐";
    b.leftUpEnd = "┌";
    b.horiz = "─";
    b.vert = "│";

	b.padding = 1

	b.pos = NewPosition("")

	return &b
}

func (b *board) move() {

}

func (b *board) buildBorderRow(row int) string {
	var leftEdge, middle, rightEdge string
	switch row{
	case 0:
		leftEdge = b.leftUpEnd
		middle = b.upT
		rightEdge = b.rightUpEnd
	case 9:
		leftEdge = b.leftBotEnd
		middle = b.botT
		rightEdge = b.rightBotEnd
	default:
		leftEdge = b.leftT
		middle = b.cross
		rightEdge = b.rightT
	}

	s := strings.Builder{}
	s.WriteString("  ")
	s.WriteString(leftEdge)
	for i := 0; i < 7; i++ {
		for j := 0; j < b.padding * 2 + 1; j++ {
			s.WriteString(b.horiz)
		}
		s.WriteString(middle)
	}
	for j := 0; j < b.padding * 2 + 1; j++ {
		s.WriteString(b.horiz)
	}
	s.WriteString(rightEdge)
	return s.String()
}

func (b *board) buildChessRow(row int, flipped bool) string {
	s := strings.Builder{}
	if flipped {
		s.WriteString(strconv.Itoa(row + 1))
	} else {
		s.WriteString(strconv.Itoa(8 - row))
	}
	s.WriteString(" ")
	s.WriteString(b.vert)
	for i := 0; i < 8; i++ {
		s.WriteString(strings.Repeat(" ", b.padding))
		if flipped {
			s.WriteString(b.pieces_unicode[b.pos.pieceAtPosition(row + 1, 8 - i)])
		} else {
			s.WriteString(b.pieces_unicode[b.pos.pieceAtPosition(8 - row, i + 1)])
		}
		s.WriteString(strings.Repeat(" ", b.padding))
		s.WriteString(b.vert)
	}
	return s.String()
}

func (b *board) buildPaddingRow(row int) string {
	return ""
	s := strings.Builder{}
	s.WriteString("  ")
	s.WriteString(b.vert)
	for i := 0; i < 8; i++ {
		for j := 0; j < 5; j++ {
			s.WriteString(" ")
		}
		s.WriteString(b.vert)
	}
	s.WriteString("\n")
	return s.String()
}

func (b *board) PrintBoard(flipped bool) string { // flipped = false is white at bottom
	s := strings.Builder{}
	for i := 0; i < 8; i++ {
		s.WriteString(b.buildBorderRow(i))
		s.WriteString("\n")
		s.WriteString(b.buildPaddingRow(i))
		s.WriteString(b.buildChessRow(i, flipped))
		s.WriteString("\n")
		s.WriteString(b.buildPaddingRow(i))
	}
	s.WriteString(b.buildBorderRow(9))
	s.WriteString("\n")

	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	labels := strings.Builder{}
	labels.WriteString("  ")
	
	for i := 0; i < 8; i++ {
		labels.WriteString(strings.Repeat(" ", b.padding + 1))
		if flipped {
			labels.WriteString(letters[7 - i])
		} else {
			labels.WriteString(letters[i])
		}
		labels.WriteString(strings.Repeat(" ", b.padding))
	}

	s.WriteString(labels.String())
	return s.String()
}

