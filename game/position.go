package game

import (
	"strings"
	"fmt"
	"strconv"
	"unicode"
)

type position struct {
	// bitboard layout
    // 63 62 61	60 59 58 57	56
    // 55 54 53	52 51 50 49	48
    // 47 46 45	44 43 42 41	40
    // 39 38 37	36 35 34 33	32
    // 31 30 29	28 27 26 25	24
    // 23 22 21	20 19 18 17	16
    // 15 14 13	12 11 10 09	08
    // 07 06 05	04 03 02 01 00
	typeBB [6]uint64
	colourBB [2]uint64

	// white pawn - 1, white knight - 2, white bishop - 3, white rook - 4, white queen - 5, white king - 6
    // black pawn - 9, black knight  - 10, black bishop - 11, black rook - 12, black queen - 13, black king - 14
    pieceToChar string

	whiteTurn bool// 0 = white, 1 = black
    castleRights int // 0 = no rights, 1 = white king, 2 = white queen, 4 = black king, 8 = black queen, 15 = all
    halfMoveClock int
    fullMoveNumber int
    enPassant uint64 // en passant square
}

func NewPosition(fen string) *position {
	p := position{}

	p.pieceToChar = " PNBRQK  pnbrqk"
	p.whiteTurn = true
	p.castleRights = 0
	p.enPassant = 0
	p.loadPosition(fen)

	return &p
}

func (p *position)getWhitePawns() uint64 {return p.typeBB[0] & p.colourBB[0]}
func (p *position)getWhiteKnights() uint64 {return p.typeBB[1] & p.colourBB[0]}
func (p *position)getWhiteBishops() uint64 {return p.typeBB[2] & p.colourBB[0]}
func (p *position)getWhiteRooks() uint64 {return p.typeBB[3] & p.colourBB[0]}
func (p *position)getWhiteQueens() uint64 {return p.typeBB[4] & p.colourBB[0]}
func (p *position)getWhiteKing() uint64 {return p.typeBB[5] & p.colourBB[0]}

func (p *position)getBlackPawns() uint64 {return p.typeBB[0] & p.colourBB[1]}
func (p *position)getBlackKnights() uint64 {return p.typeBB[1] & p.colourBB[1]}
func (p *position)getBlackBishops() uint64 {return p.typeBB[2] & p.colourBB[1]}
func (p *position)getBlackRooks() uint64 {return p.typeBB[3] & p.colourBB[1]}
func (p *position)getBlackQueens() uint64 {return p.typeBB[4] & p.colourBB[1]}
func (p *position)getBlackKing() uint64 {return p.typeBB[5] & p.colourBB[1]}

func (p *position)getPawns() uint64 {return p.typeBB[0]}
func (p *position)getKnights() uint64 {return p.typeBB[1]}
func (p *position)getBishops() uint64 {return p.typeBB[2]}
func (p *position)getRooks() uint64 {return p.typeBB[3]}
func (p *position)getQueens() uint64 {return p.typeBB[4]}
func (p *position)getKing() uint64 {return p.typeBB[5]}

func (p *position)getWhitePieces() uint64 {return p.colourBB[0]}
func (p *position)getBlackPieces() uint64 {return p.colourBB[1]}
func (p *position)getAllPieces() uint64 {return p.colourBB[0] | p.colourBB[1]}


func (p *position)validateMove() {

}

func (p *position)move() {

}

func (p *position)loadPosition(fen string) error {
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	parts := strings.Fields(fen)
	if len(parts) != 6 {
		return fmt.Errorf("Invalid FEN string %s", fen)
	}
	// 1. load piece positions
	row, col := 0, 0
	for _, a := range parts[0] {
		if a == '/' {
			col = 0
			row++
		} else if unicode.IsDigit(a) {
			col += int(a - '0')
		} else {
			piece := strings.IndexRune(p.pieceToChar, a)
			colour := piece >> 3
			pieceType := piece & 0b0111

			var index uint64 = p.fileRankToIndex(8 - row, col + 1)

			p.colourBB[colour] |= index
			p.typeBB[pieceType - 1] |= index

			col++
		}
	}

	// 2. load active colour
	if rune(parts[1][0]) == 'w' {
		p.whiteTurn = true
	} else {
		p.whiteTurn = false
	}

	// 3. load castling availability
	for _, a := range parts[2] {
		switch string(a) {
		case "K":
			p.castleRights |= 1 << 0
		case "Q":
			p.castleRights |= 1 << 1
		case "k":
			p.castleRights |= 1 << 2
		case "q":
			p.castleRights |= 1 << 3
		}
	}

	// 4. load en passant target square
	if parts[3] != "-" {
		rank := 8 - (parts[3][1] - '0')
		file := parts[3][0] - 'a' + 1
		p.enPassant = p.fileRankToIndex(int(rank), int(file))
	}

	// 5. load halfmove clock
	t, err := strconv.Atoi(parts[4])
	if err != nil {
		return err
	}
	p.halfMoveClock = t

	// 6. load fullmove number
	t, err = strconv.Atoi(parts[5])
	if err != nil {
		return err
	}
	p.fullMoveNumber = t


	return nil
}

func (p *position) pieceAtPosition(rank, file int) string {
	index := p.fileRankToIndex(rank, file)

	var piece, colour uint64 = 0, 0
	for i := uint64(0); i < 6; i++ {
		if p.typeBB[i] & index != 0 {
			piece = i + 1
		}
	}

	if index & p.colourBB[0] != 0 {
		colour = 0
	} else {
		colour = 1
	}

	if piece == 0 {
		return " "
	}
	i := piece | (colour << 3)
	return p.pieceToChar[i: i + 1]
}

func (p *position)fileRankToIndex(rank, file int) uint64 {
	return 1 << (64 - ((8 - rank) * 8 + (file - 1)) - 1);
}
