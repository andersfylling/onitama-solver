package onitamago

import "math/bits"

type Piece = Move // to simplify move extraction

const (
	NrOfPieceTypes   Amount = 3 // student, master, temple. temples cannot move
	NrOfPlayerPieces Amount = 5

	StudentsIndex Index = 0
	MasterIndex   Index = 1

	Master Piece = iota
	Student
)

// LSB Least Significant Bit
func LSB(x Board) BoardIndex {
	return BoardIndex(bits.TrailingZeros64(x))
}

// NLSB Next Least Significant Bit
func NLSB(x *Board, i BoardIndex) BoardIndex {
	*x ^= 1 << i
	return LSB(*x)
}

func boardIndexToBoard(i BoardIndex) Board {
	return 1 << i
}

func BoardToIndex(x Board) BoardIndex {
	return LSB(x)
}

func pieceAtBoardIndex(b Board, i BoardIndex) bool {
	return (b & boardIndexToBoard(i)) > 0
}
