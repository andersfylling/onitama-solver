package onitamago

func CreateRunes(char rune, length int) []rune {
	runes := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		runes = append(runes, char)
	}

	return runes
}

func Merge(boards []Board) (b Board) {
	for i := range boards {
		b |= boards[i]
	}
	return
}

// MakeCompactBoard takes a bitboard and compresses it into 25 bits
func MakeCompactBoard(board Board) (compact Board) {
	board = BoardMask & board // discard redundant info
	for i, mask := range rows {
		row := board & mask
		row = row >> 8 // one down
		row = row >> 1 // one right

		row = row >> (uint64(i) * 8)

		compact |= row << (uint64(i) * 5)
	}

	return compact
}

// MakeCompactBoard takes a bitboard and compresses it into 25 bits using ASM
func MakeCompactBoardFast(board Board) Board {
	// TODO: SIMD
	return MakeCompactBoard(board)
}

// MakeSemiCompactBoard takes a bitboard and compresses it into 34 bits
func MakeSemiCompactBoard(board Board) Board {
	const down uint64 = 8
	const right uint64 = 1

	board = BoardMask & board // discard redundant info
	board = board >> (down + right)
	board = board | (board & 0x1f00000000)

	return (board & 0x1f1f1f1f) | ((board & 0x1f00000000) >> 3)
}

func CompactBoardToBitBoard(compact Board) (board Board) {
	compact = MaskKeyBoards & compact

	for i := range rows {
		row := compact & (0x1f << (uint64(i) * 5))
		row = row >> (uint64(i) * 5)

		row = row << 1
		row = row << 8
		row = row << (8 * uint64(i))
		board |= row
	}

	return board
}
