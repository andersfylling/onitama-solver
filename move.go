package onitamago

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

// Move holds from, to, actions, the player, and which card index was used.
type Move uint16

func (m *Move) Reset() {
	*m = 0
}

func (m Move) To() BitboardPos {
	return BitboardPos(m & MoveMaskTo)
}

func (m Move) From() BitboardPos {
	return BitboardPos(m&MoveMaskFrom) >> 6
}

func (m Move) Action() Number {
	return BitboardPos(m&MoveMaskAction) >> 12
}

func (m Move) Win() bool {
	return (m.Action() & 1) == 1
}

func (m Move) WinByTemple() bool {
	return m.Action() == 7
}

func (m Move) Pass() bool {
	return m.Action() == 2 && m.From() == m.To()
}

func (m Move) PieceType() Piece {
	return Piece((m.Action() & 2) >> 1) // Master == 1, Student == 0
}

func (m Move) BoardIndex() Number {
	return Number(m.PieceType())
}

func (m Move) HostileBoardIndex() Number {
	action := m.Action()
	master := action & 1
	temple := action & 4
	return (temple >> 1) | (((temple >> 2) | master) ^ (temple >> 2)) // TODO: simplify
}

func (m Move) CardIndex() Number {
	return Number(m&MoveMaskCardIndex) >> 15
}

func (m *Move) Encode(st *State, fromIndex, toIndex, cardIndex BitboardPos) {
	m.Reset()

	///////////////////
	// From
	///////////////////
	m.addFrom(fromIndex)

	///////////////////
	// To
	///////////////////
	m.addTo(toIndex)

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
	m.addAction(action)

	///////////////////
	// Card Selection
	// -> relative to current state. Card is either 0 or 1
	///////////////////
	m.addCardIndex(cardIndex)
}

func (m *Move) addTo(p BitboardPos) {
	*m |= Move(p)
}

func (m *Move) addFrom(p BitboardPos) {
	*m |= Move(p << 6)
}

func (m *Move) addAction(a Number) {
	a = 7 & a // 0b111 & 0bxxx
	a = a << 12
	*m |= Move(a)
}

func (m *Move) addCardIndex(i BitboardPos) {
	*m |= Move(i) << 15
}

func (m Move) String() string {
	from := bitboardIndexToOnitamaIndex(m.From())
	to := bitboardIndexToOnitamaIndex(m.To())

	row1 := from / 5
	col1 := 4 - (from % 5) // reversed, due to string
	row2 := to / 5
	col2 := 4 - (to % 5) // reversed, due to string

	var piece string
	if m.PieceType() == Master {
		piece = "master"
	} else {
		piece = "student"
	}

	var winner string
	if m.Win() {
		winner = ", WINNER"
	}

	var pass string
	if m.Pass() {
		pass = " (pass)"
	}

	return piece + "{" + BoardPos(col1, row1) + " => " + BoardPos(col2, row2) + "}" + winner + pass
}

func (m Move) Card(playerIndex Number, cardsBeforeMove []Card) Card {
	i := playerIndex*NrOfPlayerCards + m.CardIndex()
	return cardsBeforeMove[i]
}

// getMoveActionAttackedPieceBoard if a piece was attacked, this returns the
// board where the related bit is set.
// This should be used together with getMoveHostileBoardIndex to handle
// temple attacks.
// func getMoveActionAttackedPieceBoard(m Move) Bitboard {
// 	action := getMoveAction(m)
// 	if (action & 4) > 0 {
// 		return 0
// 	}
//
// 	return boardIndexToBoard(getMoveTo(m))
// }

// deprecated
func encodeMove(st *State, fromIndex, toIndex, cardIndex BitboardPos) (m Move) {
	m.Encode(st, fromIndex, toIndex, cardIndex)
	return
}
