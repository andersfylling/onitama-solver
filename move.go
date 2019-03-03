package onitamago

type Move = uint16
type MoveAction = Board

const MovePositionMask Move = 0x3f
const MoveMaskTo Move = MovePositionMask << 0
const MoveMaskFrom Move = MovePositionMask << 6
const MoveMaskAction Move = 0x7 << 12
const MoveMaskCardIndex Move = 0x1 << 15

// resetMove always call this before starting to write content to a move
func resetMove(m Move) Move {
	return 0
}

func setMoveTo(m Move, pos BoardIndex) Move {
	return m | Move(pos)
}

func getMoveTo(m Move) BoardIndex {
	return BoardIndex(m & MoveMaskTo)
}

func setMoveFrom(m Move, pos BoardIndex) Move {
	return m | Move(pos<<6)
}

func getMoveFrom(m Move) BoardIndex {
	return BoardIndex((m & MoveMaskFrom) >> 6)
}

func setMoveAction(m Move, action MoveAction) Move {
	action = 0x7 & action
	action = action << 12
	return m | Move(action)
}

func getMoveWin(m Move) Index {
	return Index(getMoveAction(m) & 0x1)
}

func getMoveFriendlyBoardIndex(m Move) Index {
	action := getMoveAction(m)
	return Index((action & 0x2) >> 0x1)
}

func getMoveHostileBoardIndex(m Move) Index {
	action := getMoveAction(m)
	master := action & 0x1
	temple := action & 0x4
	return Index((temple >> 1) | (((temple >> 2) | master) ^ (temple >> 2)))
}

func getMoveAction(m Move) Move {
	return (m & MoveMaskAction) >> 12
}

func getPieceMoved(m Move) Piece {
	action := getMoveAction(m)
	if action == 0 || action == 1 || action == 4 || action == 5 {
		return Student
	}

	return Master
}

func setMoveCardIndex(m Move, index BoardIndex) Move {
	return m | Move(index<<15)
}

func getMoveCardIndex(m Move) BoardIndex {
	return BoardIndex(m&MoveMaskCardIndex) >> 15
}

func encodeMove(st *State, fromIndex, toIndex, cardIndex BoardIndex) (move Move) {
	move = resetMove(move) // in case code gets change, and re-used populated bits later on
	move = setMoveFrom(move, fromIndex)
	move = setMoveTo(move, toIndex)

	// mark the win bit if the temple or a master is taken
	win := st.temples[st.otherPlayer] >> toIndex
	win |= st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] >> toIndex
	win &= 0x1

	pieceType := st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] >> fromIndex
	pieceType &= 0x1 // redundant

	action := (0x4 & (st.temples[st.otherPlayer] >> (toIndex + 2))) | (pieceType << 1) | win
	move = setMoveAction(move, action)

	move = setMoveCardIndex(move, cardIndex)

	return move
}

func explainMove(m Move, playerIndex BoardIndex, cardsBeforeMove []Card) string {
	from := bitboardIndexToOnitamaIndex(getMoveFrom(m))
	to := bitboardIndexToOnitamaIndex(getMoveTo(m))

	row1 := from / 5
	col1 := 4 - (from % 5) // reversed, due to string
	row2 := to / 5
	col2 := 4 - (to % 5) // reversed, due to string

	card := cardsBeforeMove[playerIndex * NrOfPlayerCards + getMoveCardIndex(m)]

	var piece string
	if getPieceMoved(m) == Master {
		piece = "master"
	} else {
		piece = "student"
	}

	return piece + "{" + BoardPos(col1, row1) + " => " + BoardPos(col2, row2) + ", " + CardName(card) + "}"
}