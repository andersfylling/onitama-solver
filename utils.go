package onitamago

import (
	"math/bits"
	"strings"
)

func ExplainCardSlice(cards []Card) string {
	b := strings.Builder{}
	for i := range cards {
		b.WriteString(cards[i].Name() + ", ")
	}

	return b.String()[:b.Len()-2]
}

func CreateRunes(char rune, length int) []rune {
	runes := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		runes = append(runes, char)
	}

	return runes
}

func Merge(boards []Bitboard) (b Bitboard) {
	for i := range boards {
		b |= boards[i]
	}
	return
}

// MakeCompactBoard takes a bitboard and compresses it into 25 bits
func MakeCompactBoard(board Bitboard) (compact Bitboard) {
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
func MakeCompactBoardFast(board Bitboard) Bitboard {
	// TODO: SIMD
	return MakeCompactBoard(board)
}

// MakeSemiCompactBoard takes a bitboard and compresses it into 34 bits
func MakeSemiCompactBoard(board Bitboard) Bitboard {
	const down uint64 = 8
	const right uint64 = 1

	board = BoardMask & board // discard redundant info
	board = board >> (down + right)
	board = board | (board & 0x1f00000000)

	return (board & 0x1f1f1f1f) | ((board & 0x1f00000000) >> 3)
}

func CompactBoardToBitBoard(compact Bitboard) (board Bitboard) {
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

// LSB Least Significant Bit
func LSB(x Bitboard) BitboardPos {
	return BitboardPos(bits.TrailingZeros64(x))
}

// NLSB Next Least Significant Bit
func NLSB(x *Bitboard, i BitboardPos) BitboardPos {
	*x ^= 1 << i
	return LSB(*x)
}

// RemoveMSB removes most significant bit
func RemoveMSB(x Bitboard) Bitboard {
	return x & (x - 1)
}

func boardIndexToBoard(i BitboardPos) Bitboard {
	return 1 << i
}

func BoardToIndex(x Bitboard) BitboardPos {
	return LSB(x)
}

//go:inline
func CurrentMasterBitboard(st *State) Bitboard {
	return st.board[st.activePlayer*NrOfPieceTypes+MasterIndex]
}

//go:inline
func OtherMasterBitboard(st *State) Bitboard {
	return st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex]
}

//go:inline
func OtherPiecesBitboard(st *State) Bitboard {
	return st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] | st.board[st.otherPlayer*NrOfPieceTypes+StudentsIndex]
}
