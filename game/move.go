package game

// uses 16 bits to represent moves
// AAABBBBBBCCCCCC
// A - 4 bits for special moves (only promotions for us)
// B - 6 bits for origin
// C - 6 bits for dest

// 0 - pawn  / default
// 1 - knight
// 2 - bishop
// 3 - rook
// 4 - queen


type move uint16

func (m *move) dest() int{
	return int(*m & 0x3F)
}

func (m *move) origin() int{
	return int((*m & 0xFC0) >> 6)
}

func (m *move) promotion() int{
	return int((*m & 0xF000) >> 12)
}

func (m *move) create(from, to int, prom rune) {
	promotion := 0
	if prom == 'n' {
		promotion = 1
	} else if prom == 'b' {
		promotion = 2
	} else if prom == 'r' {
		promotion = 3
	} else if prom == 'q' {
		promotion = 4
	}

	*m = *m & ^(move(0xF000)) | (move(promotion) << 12)
	*m = *m & ^(move(0xFC0)) | (move(from) << 6)
	*m = *m & ^(move(0x3F)) | move(to)
}
