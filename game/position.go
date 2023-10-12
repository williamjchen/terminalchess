package game

// good post about generating moves: https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/

import (
	"math/bits"
	"strings"
	"fmt"
	"strconv"
	"unicode"

	"github.com/williamjchen/terminalchess/magic"
)

type turn int

const (
	WhiteTurn turn = iota
	BlackTurn
	whiteMate
	blackMate
	stalemate
)

type position struct {
	// bitboard layout
	// 56 57 58 59 60 61 62 63
	// 48 49 50 51 52 53 54 55
	// 40 41 42 43 44 45 46 47 
	// 32 33 34 35 36 37 38 39 
	// 24 25 26 27 28 29 30 31 
	// 16 17 18 19 20 21 22 23 
	// 08 09 10 11 12 13 14 15 
	// 00 01 02 03 04 05 06 07
	// ^ corresponds to rank and file index at 1
	// 8
	// 7
	// ... (files)
	// 3
	// 2
	// 1  2  3  4  5  6  7  8 (rank)
	// A  B  C  D  E  F  G  H
	// Compass for bitboard
	// <<7 <<8 <<9
	// >>1  0  <<1
	// >>9 >>8 >>7
	
	typeBB [6]uint64
	colourBB [2]uint64

	// white pawn - 1, white knight - 2, white bishop - 3, white rook - 4, white queen - 5, white king - 6
    // black pawn - 9, black knight  - 10, black bishop - 11, black rook - 12, black queen - 13, black king - 14
    pieceToChar string

	turn turn
	castleRights int // 0 = no rights, 1 = white king, 2 = white queen, 4 = black king, 8 = black queen, 15 = all
    halfMoveClock int
    fullMoveNumber int
    enPassant uint64 // en passant square
}

func NewPosition(fen string) *position {
	p := position{}

	p.pieceToChar = " PNBRQK  pnbrqk"
	p.turn = WhiteTurn
	p.castleRights = 0
	p.enPassant = 0
	p.loadPosition(fen)

	return &p
}

func (p *position) getWhitePawns() uint64 {return p.typeBB[0] & p.colourBB[0]}
func (p *position) getWhiteKnights() uint64 {return p.typeBB[1] & p.colourBB[0]}
func (p *position) getWhiteBishops() uint64 {return p.typeBB[2] & p.colourBB[0]}
func (p *position) getWhiteRooks() uint64 {return p.typeBB[3] & p.colourBB[0]}
func (p *position) getWhiteQueens() uint64 {return p.typeBB[4] & p.colourBB[0]}
func (p *position) getWhiteKing() uint64 {return p.typeBB[5] & p.colourBB[0]}

func (p *position) getBlackPawns() uint64 {return p.typeBB[0] & p.colourBB[1]}
func (p *position) getBlackKnights() uint64 {return p.typeBB[1] & p.colourBB[1]}
func (p *position) getBlackBishops() uint64 {return p.typeBB[2] & p.colourBB[1]}
func (p *position) getBlackRooks() uint64 {return p.typeBB[3] & p.colourBB[1]}
func (p *position) getBlackQueens() uint64 {return p.typeBB[4] & p.colourBB[1]}
func (p *position) getBlackKing() uint64 {return p.typeBB[5] & p.colourBB[1]}

func (p *position) getPawns() uint64 {return p.typeBB[0]}
func (p *position) getKnights() uint64 {return p.typeBB[1]}
func (p *position) getBishops() uint64 {return p.typeBB[2]}
func (p *position) getRooks() uint64 {return p.typeBB[3]}
func (p *position) getQueens() uint64 {return p.typeBB[4]}
func (p *position) getKing() uint64 {return p.typeBB[5]}

func (p *position) getWhitePieces() uint64 {return p.colourBB[0]}
func (p *position) getBlackPieces() uint64 {return p.colourBB[1]}
func (p *position) getAllPieces() uint64 {return p.colourBB[0] | p.colourBB[1]}


// we use these functions to create simple move.

func (p *position) move(origin, dest int) bool {
	origin_pos := uint64(1) << origin
	dest_pos := uint64(1) << dest
	if p.validateMove(origin, dest) {
		origin_piece := p.pieceAtSquare(origin)
		switch origin_piece {
		case "P":
			p.movePiece(0, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "N":
			p.movePiece(1, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "B":
			p.movePiece(2, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "R":
			p.movePiece(3, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "Q":
			p.movePiece(4, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "K":
			p.movePiece(5, 0, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "p":
			p.movePiece(0, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "n":
			p.movePiece(1, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "b":
			p.movePiece(2, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "r":
			p.movePiece(3, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "q":
			p.movePiece(4, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		case "k":
			p.movePiece(5, 1, origin_pos, dest_pos)
			p.removeAt(dest, true)
		}

		return true
	}
	return false
}

func (p *position) validateMove(origin, dest int) bool {
	moves := p.generateLegalMoves()
	for _, move := range moves {
		from := move.origin()
		to := move.dest()
		if origin == from && dest == to {
			return true
		}
	}
	return false
}

func (p *position) movePiece(pieceIndex, colourIndex int, origin, dest uint64) {
	p.typeBB[pieceIndex] &= ^origin
	p.typeBB[pieceIndex] |= dest // add new
	p.colourBB[colourIndex] &= ^origin // remove origin
	p.colourBB[colourIndex] |= dest // add new
}

func (p *position) removeAt(square int, white bool) {
	pos := uint64(1) << square
	if white {
		p.colourBB[0] &= ^pos
	} else {
		p.colourBB[1] &= ^pos
	}
	for i := 0; i < 6; i++ {
		p.typeBB[i] &= ^pos
	}
}

// This is used to generate all moves. Required for engine
func (p *position) generateLegalMoves() []move {
	var moves []move
	var kingPos, pawnPos, knightPos, bishopPos, rookPos, queenPos, allPieces uint64
	var e_kingPos, e_pawnPos, e_knightPos, e_bishopPos, e_rookPos, e_queenPos uint64

	if p.turn == WhiteTurn {
		kingPos = p.getWhiteKing()
		pawnPos = p.getWhitePawns()
		knightPos = p.getWhiteKnights()
		bishopPos = p.getWhiteBishops()
		rookPos = p.getWhiteRooks()
		queenPos = p.getWhiteQueens()
		e_kingPos = p.getBlackKing()
		e_pawnPos = p.getBlackPawns()
		e_knightPos = p.getBlackKnights()
		e_bishopPos = p.getBlackBishops()
		e_rookPos = p.getBlackRooks()
		e_queenPos = p.getBlackQueens()
		allPieces = p.getWhitePieces()
	} else {
		kingPos = p.getBlackKing()
		pawnPos = p.getBlackPawns()
		knightPos = p.getBlackKnights()
		bishopPos = p.getBlackBishops()
		rookPos = p.getBlackRooks()
		queenPos = p.getBlackQueens()
		e_kingPos = p.getWhiteKing()
		e_pawnPos = p.getWhitePawns()
		e_knightPos = p.getWhiteKnights()
		e_bishopPos = p.getWhiteBishops()
		e_rookPos = p.getWhiteRooks()
		e_queenPos = p.getWhiteQueens()
		allPieces = p.getBlackPieces()
	}

	num_attackers, allowed_dests := p.numAttacks(p.turn, kingPos, e_kingPos, e_pawnPos, e_knightPos, e_bishopPos, e_rookPos, e_queenPos)
	if num_attackers > 1 { // multiple check - only king moves
		moves = append(moves, p.kingPushes(kingPos, allowed_dests)...) // king captures?
		return moves
	} else if num_attackers == 1 { // single check
		pinned := p.generatePinnedSquares(kingPos, e_rookPos, e_bishopPos, e_queenPos, allPieces)
		nonPinned := ^pinned

		allowed_dests := nonPinned
		moves = append(moves, p.pawnCaptures(pawnPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.pawnPushes(pawnPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.knightPushes(knightPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.rookPushes(rookPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.bishopPushes(bishopPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.queenPushes(queenPos, nonPinned, allowed_dests)...)
		moves = append(moves, p.kingPushes(kingPos, allowed_dests)...)
		return moves
	}

	// non-check moves
	pinned := p.generatePinnedSquares(kingPos, e_rookPos, e_bishopPos, e_queenPos, allPieces)
	nonPinned := ^pinned
	moves = append(moves, p.pawnCaptures(pawnPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.pawnPushes(pawnPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.knightPushes(knightPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.rookPushes(rookPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.bishopPushes(bishopPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.queenPushes(queenPos, nonPinned, magic.Everything)...)
	moves = append(moves, p.kingPushes(kingPos, magic.Everything)...)

	return moves
}

func (p *position) numAttacks(defender turn, kingPos, e_kingPos, e_pawnPos, e_knightPos, e_bishopPos, e_rookPos, e_queenPos uint64) (int, uint64) {
	square := bits.TrailingZeros64(kingPos)
	num_attackers := 0
	var attackers_mask uint64 = 0

	var attackers uint64 = 0
	// pawns
	if defender == WhiteTurn {
		attackers = kingPos << 7
		attackers |= kingPos << 9
	} else {
		attackers = kingPos >> 7
		attackers |= kingPos >> 9
	}
	attackers &= e_pawnPos
	num_attackers += bits.OnesCount64(attackers)
	attackers_mask |= attackers

	// kings

	// kights
	attackers = magic.KnightMasks[square] & e_knightPos
	num_attackers += bits.OnesCount64(attackers)
	attackers_mask |= attackers

	// bishops + queen
	diag_blocked := magic.MagicBishopBlockerMasks[square] & p.getAllPieces()
	diag_idx := magic.BishopHash(magic.Square(square), diag_blocked)
	diag_ends := magic.MagicMovesBishop[square][diag_idx]
	diag_attackers := diag_ends & (e_bishopPos | e_queenPos)
	num_attackers += bits.OnesCount64(diag_attackers)
	attackers_mask |= diag_attackers

	// rooks + queen
	straight_blocked := magic.MagicRookBlockerMasks[square] & p.getAllPieces()
	straight_idx := magic.RookHash(magic.Square(square), straight_blocked)
	straight_ends := magic.MagicMovesRook[square][straight_idx]
	straight_attackers := straight_ends & (e_rookPos | e_queenPos)
	num_attackers += bits.OnesCount64(straight_attackers)
	attackers_mask |= straight_attackers

	return num_attackers, attackers_mask
}

func (p *position) generatePinnedSquares (origin, e_rookPos, e_bishopPos, e_queenPos, myPieces uint64) uint64 {
	kingSquare := bits.TrailingZeros64(origin)
	var opponent_slide uint64 = 0
	var king_slide uint64 = 0
	king_slide = p.generateDiagonalSquares(kingSquare, p.getAllPieces()) | p.generateStraightSquares(kingSquare, p.getAllPieces())

	// rook + queen
	op := e_rookPos | e_queenPos
	for op != 0 {
		rookIdx := bits.TrailingZeros64(op)
		op &= op - 1
		targs := p.generateStraightSquares(rookIdx, p.getAllPieces())

		opponent_slide |= targs
	}

	// bishop + queen
	op = e_bishopPos | e_queenPos
	for op != 0 {
		bishopIdx := bits.TrailingZeros64(op)
		op &= op - 1
		targs := p.generateDiagonalSquares(bishopIdx, p.getAllPieces())

		opponent_slide |= targs
	}

	return king_slide & opponent_slide & myPieces
}

func (p *position) generateDiagonalSquares(origin int, pieces uint64) uint64 {
	blockers := magic.MagicBishopBlockerMasks[origin] & pieces
	index := magic.BishopHash(magic.Square(origin), blockers)
	return magic.MagicMovesBishop[origin][index]
}

func (p *position) generateStraightSquares(origin int, pieces uint64) uint64 {
	blockers := magic.MagicRookBlockerMasks[origin] & pieces
	index := magic.RookHash(magic.Square(origin), blockers)
	return magic.MagicMovesRook[origin][index]
}

func (p *position) kingPushes(king uint64, dest uint64) []move{
	return []move{}
}

func (p *position) pawnPushes(pawns, allowed, dest uint64) []move{
	return []move{}
}

func (p *position) pawnCaptures(pawns, allowed, dest uint64) []move{
	return []move{}
}

func (p *position) knightPushes(knights, allowed, dest uint64) []move{
	return []move{}
}

func (p *position) bishopPushes(bishops, allowed, dest uint64) []move {
	return []move{}
}

func (p *position) rookPushes(rook, allowed, dest uint64) []move{
	return []move{}
}

func (p *position) queenPushes(queen, allowed, dest uint64) []move{
	return []move{}
}

func (p *position) loadPosition(fen string) error {
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
		p.turn = WhiteTurn
	} else {
		p.turn = BlackTurn
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

func (p *position) pieceAtSquare(square int) string {
	rank := square / 8
	file := square % 8 + 1
	return p.pieceAtPosition(rank, file)
}

func (p *position) pieceAtPosition(rank, file int) string { //rank = 8, file = 1 is top left
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

func (p *position) fileRankToIndex(rank, file int) uint64 { // vertical is file(A-H). rank is horizontal(1-8)
	i := (rank * 8 - 1) - (8 - file)
	return 1 << i
}
