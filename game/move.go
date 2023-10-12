package game

// uses 16 bits to represent moves
// AAABBBBBBCCCCCC
// A - 3 bits for promotin
// B - 6 bits for origin
// C - 6 bits for dest

type move uint16

func (m *move) dest() int{
	return int(*m & 0x3F)
}

func (m *move) origin() int{
	return int((*m & 0xFC0) >> 6)
}

func (m *move) promotion() int{
	return int((*m & 0x7000) >> 12)
}

func (m *move) create(from, to int) {
	*m = *m & ^(move(0xFC0)) | (move(from) << 6)
	*m = *m & ^(move(0x3F)) | move(to)
}
