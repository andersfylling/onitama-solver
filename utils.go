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
