package onitamago

type Piece = Move // to simplify move extraction

const (
	NrOfPieceTypes   uint64 = 3 // student, master, temple. temples cannot move
	NrOfPlayerPieces uint64 = 5

	StudentsIndex uint64 = 0
	MasterIndex   uint64 = 1
	TrashIndex    uint64 = 2

	Student Piece = 0
	Master  Piece = 1
)

func pieceAtBoardIndex(b Bitboard, i BitboardPos) bool {
	return (b & boardIndexToBoard(i)) > 0
}
