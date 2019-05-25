package onitamago

type Move = uint16
type MoveAction = Board

// TypeUndoMove can be used to signify the state should undo the current move
// may it be to go backwards in a game tree or for other reasons.
const MoveUndo = ^Move(0)

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
	action = 7 & action // 0b0111 & 0bxxxx
	action = action << 12
	return m | Move(action)
}

func getMoveWin(m Move) Index {
	return Index(getMoveAction(m) & 1)
}

func getMoveFriendlyBoardIndex(m Move) Index {
	action := getMoveAction(m)
	return Index((action & 2) >> 1)
}

func getMoveHostileBoardIndex(m Move) Index {
	action := getMoveAction(m)
	master := action & 1
	temple := action & 4
	return Index((temple >> 1) | (((temple >> 2) | master) ^ (temple >> 2)))
}

// getMoveActionAttackedPieceBoard if a piece was attacked, this returns the
// board where the related bit is set.
// This should be used together with getMoveHostileBoardIndex to handle
// temple attacks.
func getMoveActionAttackedPieceBoard(m Move) Board {
	action := getMoveAction(m)
	if (action & 4) > 0 {
		return 0
	}

	return boardIndexToBoard(getMoveTo(m))
}

func getMoveAction(m Move) Move {
	return (m & MoveMaskAction) >> 12
}

func getPieceMoved(m Move) Piece {
	return Piece(getMoveFriendlyBoardIndex(m)) // Master == 1, Student == 0
}

func setMoveCardIndex(m Move, index BoardIndex) Move {
	return m | Move(index<<15)
}

func getMoveCardIndex(m Move) BoardIndex {
	return BoardIndex(m&MoveMaskCardIndex) >> 15
}

func encodeMove(st *State, fromIndex, toIndex, cardIndex BoardIndex) (move Move) {
	move = resetMove(move) // in case code gets change, and re-used populated bits later on

	///////////////////
	// From
	///////////////////
	move = setMoveFrom(move, fromIndex)

	///////////////////
	// To
	///////////////////
	move = setMoveTo(move, toIndex)

	///////////////////
	// Action
	///////////////////
	// mark the win bit if the temple or a master is taken
	win := st.temples[st.otherPlayer] >> toIndex
	win |= st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex] >> toIndex
	win &= 1 // filter out unwanted bits

	// set which piece type was moved
	master := st.board[st.activePlayer*NrOfPieceTypes+MasterIndex] >> fromIndex
	master &= 1 // filter out unwanted bits

	// weakAttack regards attacks that kills no student, nor master.
	attack := st.board[st.otherPlayer*NrOfPieceTypes+StudentsIndex] | st.board[st.otherPlayer*NrOfPieceTypes+MasterIndex]
	attack = attack >> toIndex
	attack &= 1                   // filter out unwanted bits
	noAttack := 4 ^ (attack << 2) // 0b100 ^ (0b001 << 2) == 0 or 0b100 ^ (0b000 << 2) == 0b100

	// merge actions
	action := noAttack | (master << 1) | win
	move = setMoveAction(move, action)

	///////////////////
	// Card Selection
	// -> relative to current state. Card is either 0 or 1
	///////////////////
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

	card := cardsBeforeMove[playerIndex*NrOfPlayerCards+getMoveCardIndex(m)]

	var piece string
	if getPieceMoved(m) == Master {
		piece = "master"
	} else {
		piece = "student"
	}

	var winner string
	if getMoveWin(m) == 1 {
		winner = ", WINNER"
	}

	return piece + "{" + BoardPos(col1, row1) + " => " + BoardPos(col2, row2) + ", " + CardName(card) + "}" + winner
}
