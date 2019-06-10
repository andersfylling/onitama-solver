package onitamago

type Move = uint16
type MoveAction = Bitboard

// TypeUndoMove can be used to signify the state should undo the current move
// may it be to go backwards in a game tree or for other reasons.
const MoveUndo = ^Move(0)

// MovePassBase is the base to build a pass move.
// A Pass move is simply a normal move where the master piece,
// the piece that will always exist when a player can move, has the same
// from as to value; meaning they do not relocate.
// The card index must be set. Since you must create a pass situation for each card
// there must exist at least two pass moves when _no_ normal move can be executed.
const MovePassBase Move = 2 << 12 // action: master move

const MovePositionMask Move = 0x3f
const MoveMaskTo Move = MovePositionMask << 0
const MoveMaskFrom Move = MovePositionMask << 6
const MoveMaskAction Move = 0x7 << 12
const MoveMaskCardIndex Move = 0x1 << 15

// resetMove always call this before starting to write content to a move
func resetMove(m Move) Move {
	return 0
}

func setMoveTo(m Move, pos BitboardPos) Move {
	return m | Move(pos)
}

func getMoveTo(m Move) BitboardPos {
	return BitboardPos(m & MoveMaskTo)
}

func setMoveFrom(m Move, pos BitboardPos) Move {
	return m | Move(pos<<6)
}

func getMoveFrom(m Move) BitboardPos {
	return BitboardPos((m & MoveMaskFrom) >> 6)
}

func setMoveAction(m Move, action MoveAction) Move {
	action = 7 & action // 0b111 & 0bxxx
	action = action << 12
	return m | Move(action)
}

func getMoveWin(m Move) BitboardPos {
	return BitboardPos(getMoveAction(m) & 1)
}

func getMoveFriendlyBoardIndex(m Move) BitboardPos {
	action := getMoveAction(m)
	return BitboardPos((action & 2) >> 1)
}

func getMoveHostileBoardIndex(m Move) BitboardPos {
	action := getMoveAction(m)
	master := action & 1
	temple := action & 4
	return BitboardPos((temple >> 1) | (((temple >> 2) | master) ^ (temple >> 2)))
}

// getMoveActionAttackedPieceBoard if a piece was attacked, this returns the
// board where the related bit is set.
// This should be used together with getMoveHostileBoardIndex to handle
// temple attacks.
func getMoveActionAttackedPieceBoard(m Move) Bitboard {
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

func setMoveCardIndex(m Move, index BitboardPos) Move {
	return m | Move(index<<15)
}

func getMoveCardIndex(m Move) BitboardPos {
	return BitboardPos(m&MoveMaskCardIndex) >> 15
}

func IsPassMove(m Move) bool {
	action := getMoveAction(m)
	from := getMoveFrom(m)
	to := getMoveTo(m)

	return action == 2 && from == to
}

func encodeMove(st *State, fromIndex, toIndex, cardIndex BitboardPos) (move Move) {
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
	// Action - Moved piece type
	//  If it's not master, then student
	///////////////////
	master := CurrentMasterBitboard(st) >> fromIndex
	master &= 1 // filter out unwanted bits

	///////////////////
	// Action - Pacifist move
	///////////////////
	attack := OtherPiecesBitboard(st) >> toIndex
	attack &= 1            // filter out unwanted bits. Redundant!
	pacifist := 1 ^ attack // 0b01 ^ 0b001 == 0 or 0b001 ^ 0b000 == 0b001

	///////////////////
	// Action - Master took temple
	///////////////////
	win := (st.temples[st.otherPlayer] >> toIndex) & master

	///////////////////
	// Action - Defeated Master
	///////////////////
	win |= OtherMasterBitboard(st) >> toIndex
	win &= 1 // filter out unwanted bits

	// Action merge
	action := (pacifist << 2) | (master << 1) | win
	move = setMoveAction(move, action)

	///////////////////
	// Card Selection
	// -> relative to current state. Card is either 0 or 1
	///////////////////
	move = setMoveCardIndex(move, cardIndex)

	return move
}

func explainMove(m Move, playerIndex BitboardPos, cardsBeforeMove []Card) string {
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

	var pass string
	if IsPassMove(m) {
		pass = " (pass)"
	}

	return piece + "{" + BoardPos(col1, row1) + " => " + BoardPos(col2, row2) + ", " + card.Name() + "}" + winner + pass
}
