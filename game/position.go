package game

// good post about generating moves: https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/

import (
	"math/bits"
	"strings"
	"fmt"
	"strconv"
	"unicode"
	"log/slog"

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


	// STORES THE CURRENT TURN'S POSITIONS
	pawnPos uint64
	kingPos uint64
	knightPos uint64
	bishopPos uint64
	rookPos uint64
	queenPos uint64
	e_kingPos uint64
	e_pawnPos uint64
	e_knightPos uint64
	e_bishopPos uint64
	e_rookPos uint64
	e_queenPos uint64
	allPieces uint64
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
	slog.Info("origin", "square", origin, "pos", origin_pos)
	slog.Info("dest", "square", dest, "pos", dest_pos)
	if p.validateMove(origin, dest) {
		origin_piece := p.pieceAtSquare(origin)
		slog.Info("original", "piece", origin_piece)
		switch origin_piece {
		case "P":
			p.removeAt(dest, false)
			p.movePiece(0, 0, origin_pos, dest_pos)
		case "N":
			p.removeAt(dest, false)
			p.movePiece(1, 0, origin_pos, dest_pos)
		case "B":
			p.removeAt(dest, false)
			p.movePiece(2, 0, origin_pos, dest_pos)
		case "R":
			p.removeAt(dest, false)
			p.movePiece(3, 0, origin_pos, dest_pos)
		case "Q":
			p.removeAt(dest, false)
			p.movePiece(4, 0, origin_pos, dest_pos)
		case "K":
			p.removeAt(dest, false)
			p.movePiece(5, 0, origin_pos, dest_pos)
		case "p":
			p.removeAt(dest, true)
			p.movePiece(0, 1, origin_pos, dest_pos)
		case "n":
			p.removeAt(dest, true)
			p.movePiece(1, 1, origin_pos, dest_pos)
		case "b":
			p.removeAt(dest, true)
			p.movePiece(2, 1, origin_pos, dest_pos)
		case "r":
			p.removeAt(dest, true)
			p.movePiece(3, 1, origin_pos, dest_pos)
		case "q":
			p.removeAt(dest, true)
			p.movePiece(4, 1, origin_pos, dest_pos)
		case "k":
			p.removeAt(dest, true)
			p.movePiece(5, 1, origin_pos, dest_pos)
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
	
	if p.turn == WhiteTurn {
		p.kingPos = p.getWhiteKing()
		p.pawnPos = p.getWhitePawns()
		p.knightPos = p.getWhiteKnights()
		p.bishopPos = p.getWhiteBishops()
		p.rookPos = p.getWhiteRooks()
		p.queenPos = p.getWhiteQueens()
		p.e_kingPos = p.getBlackKing()
		p.e_pawnPos = p.getBlackPawns()
		p.e_knightPos = p.getBlackKnights()
		p.e_bishopPos = p.getBlackBishops()
		p.e_rookPos = p.getBlackRooks()
		p.e_queenPos = p.getBlackQueens()
		p.allPieces = p.getWhitePieces()
	} else {
		p.kingPos = p.getBlackKing()
		p.pawnPos = p.getBlackPawns()
		p.knightPos = p.getBlackKnights()
		p.bishopPos = p.getBlackBishops()
		p.rookPos = p.getBlackRooks()
		p.queenPos = p.getBlackQueens()
		p.e_kingPos = p.getWhiteKing()
		p.e_pawnPos = p.getWhitePawns()
		p.e_knightPos = p.getWhiteKnights()
		p.e_bishopPos = p.getWhiteBishops()
		p.e_rookPos = p.getWhiteRooks()
		p.e_queenPos = p.getWhiteQueens()
		p.allPieces = p.getBlackPieces()
	}

	num_attackers, allowed_dests := p.numAttacks(p.turn, p.kingPos)
	if num_attackers > 1 { // multiple check - only king moves
		moves = append(moves, p.kingPushes(allowed_dests)...) // king captures?
		return moves
	} else if num_attackers == 1 { // single check
		pinned := p.generatePinnedSquares()
		nonPinned := ^pinned

		allowed_dests := nonPinned
		moves = append(moves, p.pawnCaptures(nonPinned, allowed_dests)...)
		moves = append(moves, p.pawnPushes(nonPinned, allowed_dests)...)
		moves = append(moves, p.knightMoves(nonPinned, allowed_dests)...)
		moves = append(moves, p.rookMoves(nonPinned, allowed_dests)...)
		moves = append(moves, p.bishopMoves(nonPinned, allowed_dests)...)
		moves = append(moves, p.queenMoves(nonPinned, allowed_dests)...)
		moves = append(moves, p.kingPushes(allowed_dests)...)
		return moves
	}

	// non-check moves
	pinned := p.generatePinnedSquares()
	nonPinned := ^pinned
	moves = append(moves, p.pawnCaptures(nonPinned, magic.Everything)...)
	moves = append(moves, p.pawnPushes(nonPinned, magic.Everything)...)
	moves = append(moves, p.knightMoves(nonPinned, magic.Everything)...)
	moves = append(moves, p.rookMoves(nonPinned, magic.Everything)...)
	moves = append(moves, p.bishopMoves(nonPinned, magic.Everything)...)
	moves = append(moves, p.queenMoves(nonPinned, magic.Everything)...)
	moves = append(moves, p.kingMoves(magic.Everything)...)

	return moves
}

func (p *position) numAttacks(defender turn, kingPos uint64) (int, uint64) {
	square := bits.TrailingZeros64(kingPos)
	num_attackers := 0
	var attackers_mask uint64 = 0

	var attackers uint64 = 0
	// pawns
	if defender == WhiteTurn {
		attackers = p.kingPos << 7
		attackers |= p.kingPos << 9
	} else {
		attackers = p.kingPos >> 7
		attackers |= p.kingPos >> 9
	}
	attackers &= p.e_pawnPos
	num_attackers += bits.OnesCount64(attackers)
	attackers_mask |= attackers

	// kings

	// kights
	attackers = magic.KnightMasks[square] & p.e_knightPos
	num_attackers += bits.OnesCount64(attackers)
	attackers_mask |= attackers

	// bishops + queen
	diag_blocked := magic.MagicBishopBlockerMasks[square] & p.getAllPieces()
	diag_idx := magic.BishopHash(magic.Square(square), diag_blocked)
	diag_ends := magic.MagicMovesBishop[square][diag_idx]
	diag_attackers := diag_ends & (p.e_bishopPos | p.e_queenPos)
	num_attackers += bits.OnesCount64(diag_attackers)
	attackers_mask |= diag_attackers

	// rooks + queen
	straight_blocked := magic.MagicRookBlockerMasks[square] & p.getAllPieces()
	straight_idx := magic.RookHash(magic.Square(square), straight_blocked)
	straight_ends := magic.MagicMovesRook[square][straight_idx]
	straight_attackers := straight_ends & (p.e_rookPos | p.e_queenPos)
	num_attackers += bits.OnesCount64(straight_attackers)
	attackers_mask |= straight_attackers

	return num_attackers, attackers_mask
}

func (p *position) generatePinnedSquares () uint64 {
	kingSquare := bits.TrailingZeros64(p.kingPos)
	var opponent_slide uint64 = 0
	var king_slide uint64 = 0
	king_slide = p.generateDiagonalSquares(kingSquare, p.getAllPieces()) | p.generateStraightSquares(kingSquare, p.getAllPieces())

	// rook + queen
	op := p.e_rookPos | p.e_queenPos
	for op != 0 {
		rookIdx := bits.TrailingZeros64(op)
		op &= op - 1
		targs := p.generateStraightSquares(rookIdx, p.getAllPieces())

		opponent_slide |= targs
	}

	// bishop + queen
	op = p.e_bishopPos | p.e_queenPos
	for op != 0 {
		bishopIdx := bits.TrailingZeros64(op)
		op &= op - 1
		targs := p.generateDiagonalSquares(bishopIdx, p.getAllPieces())

		opponent_slide |= targs
	}

	return king_slide & opponent_slide & p.allPieces
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

// FOR FOLLOWING FUNCTIONS- allowed is non-pinned squares and dest is allowed locations to move to

// sans castling (used when under check)
func (p *position) kingPushes(dest uint64) []move{
	moves := []move{}
	kingSquare := bits.TrailingZeros64(p.kingPos)
	temp := p.kingPos
	p.kingPos = 0
	if p.turn == WhiteTurn {
		p.colourBB[0] &= ^temp
	} else {
		p.colourBB[1] &= ^temp
	}

	targets := magic.KingMasks[kingSquare]
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1

		if n, _:= p.numAttacks(p.turn, uint64(1) << target); n > 0 {
			continue
		}

		var move move
		move.create(kingSquare, target)
		moves = append(moves, move)
	}

	if p.turn == WhiteTurn {
		p.colourBB[0] |= temp
	} else {
		p.colourBB[1] |= temp
	}
	p.kingPos = temp
	return moves
}

// with castling
func (p *position) kingMoves(dest uint64) []move{
	return []move{}
}

func (p *position) pawnPushes(allowed, dest uint64) []move{
	return []move{}
}

func (p *position) pawnCaptures(allowed, dest uint64) []move{
	return []move{}
}

func (p *position) knightMoves(allowed, dest uint64) []move{
	moves := []move{}
	knights := p.knightPos & allowed
	for knights != 0 {
		knight := bits.TrailingZeros64(knights)
		knights &= knights - 1
		targets := magic.KnightMasks[knight] & dest
		moves = append(moves, generateMoves(knight, targets)...)
	}
	return moves
}

func (p *position) bishopMoves(allowed, dest uint64) []move {
	moves := []move{}
	bishops := p.bishopPos & allowed
	for bishops != 0 {
		bishop := bits.TrailingZeros64(bishops)
		bishops &= bishops - 1
		// magic bitboards
		blockers := magic.MagicBishopBlockerMasks[bishop] & p.getAllPieces()
		idx := magic.BishopHash(magic.Square(bishop), blockers)
		targets := magic.MagicMovesBishop[bishop][idx]
		moves = append(moves, generateMoves(bishop, targets)...)
	}
	return moves
}

func (p *position) rookMoves(allowed, dest uint64) []move{
	moves := []move{}
	rooks := p.rookPos & allowed
	for rooks != 0 {
		rook := bits.TrailingZeros64(rooks)
		rooks &= rooks - 1
		// magic bitboards
		blockers := magic.MagicRookBlockerMasks[rook] & p.getAllPieces()
		idx := magic.RookHash(magic.Square(rook), blockers)
		targets := magic.MagicMovesRook[rook][idx]
		moves = append(moves, generateMoves(rook, targets)...)
	}
	return moves
}

func (p *position) queenMoves(allowed, dest uint64) []move{
	return []move{}
}

func generateMoves(square int, targets uint64) []move{
	moves := []move{}
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1
		var move move
		move.create(square, target)
		moves = append(moves, move)
	}
	return moves
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
	rank := (square / 8) + 1
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
