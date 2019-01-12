package onitamago

type Move = uint16

const MovePositionMask Move = 0x3f
const MoveMaskTo Move = MovePositionMask << 0
const MoveMaskFrom Move = MovePositionMask << 6
const MoveMaskWin Move = 0x1 << 12
const MoveMaskAttack Move = 0x2 << 13
const MoveMaskCardIndex Move = 0x1 << 15
// 0 bits left

func setMoveTo(m Move, pos BoardIndex) Move {
	return ((m | MoveMaskTo) ^ MoveMaskTo) | Move(pos)
}

func setMoveFrom(m Move, pos BoardIndex) Move {
	return ((m | MoveMaskFrom) ^ MoveMaskFrom) | Move(pos << 6)
}

func setMoveWin(m Move, hit Board) Move {
	return unsetMoveWin(m) | Move(hit << 12)
}

func unsetMoveWin(m Move) Move {
	return (m | MoveMaskWin) ^ MoveMaskWin
}

func setMoveAttack(m Move, pos BoardIndex) Move {
	return ((m | MoveMaskAttack) ^ MoveMaskAttack) | Move(pos << 13)
}

func setCardIndex(m Move, index BoardIndex) Move {
	return ((m | MoveMaskAttack) ^ MoveMaskAttack) | Move(index << 15)
}

func getMoveCardIndex(m Move) BoardIndex {
	return BoardIndex(m & MoveMaskAttack) >> 15
}

func encodeMove(fromIndex, toIndex, cardIndex BoardIndex, opponents, friends Board, temple Board) (move Move) {
	move = setMoveFrom(move, fromIndex)
	move = setMoveTo(move, toIndex)

	to := boardIndexToBoard(toIndex)
	attack := to & opponents
	attackIndex := LSB(attack)
	move = setMoveAttack(move, attackIndex)
	// TODO: check if a master is killed, cause that's a win

	templeHit := temple & to
	move = setMoveWin(move, templeHit >> 11)
	move = setMoveWin(move, templeHit >> 43)

	move = setCardIndex(move, cardIndex)

	return move
}