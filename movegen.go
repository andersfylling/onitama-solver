package onitamago

//go:inline
func pieces(player BitboardPos, board [NrOfPlayers * NrOfPieceTypes]Bitboard) Bitboard {
	// bind the function to our compile time assumption
	if NrOfPieceTypes != 3 {
		panic("NrOfPieceTypes is no longer 2")
	}
	return board[player*NrOfPieceTypes] | board[player*NrOfPieceTypes+1]
}

//go:inline
func addMovePositionsToPieceBB(card Card, pos BitboardPos, friendlyPieces Bitboard) (moves Bitboard) {
	moves = card.Bitboard() >> (CardOffset - pos)
	moves ^= moves & friendlyPieces // remove moves that hits a friendly warrior
	moves &= BoardMask              // ignore positions outside the board

	return
}

func generateMoves(st *State) (moveIndex Number) {
	friends := pieces(st.activePlayer, st.board)
	var moves Bitboard
	var pieces Bitboard
	for c := st.activePlayer * NrOfPlayerCards; c < (st.activePlayer*NrOfPlayerCards + NrOfPlayerCards); c++ {
		card := st.playerCards[c]
		// TODO: remove if sentence
		// add some virtual layer
		if st.currentDepth%2 == 1 {
			card.Rotate()
		}
		pieces = friends
		for i := LSB(pieces); i != 64; i = NLSB(&pieces, i) {
			moves = addMovePositionsToPieceBB(card, i, friends)

			for j := LSB(moves); j != 64; j = NLSB(&moves, j) {
				st.generatedMoves[moveIndex] = encodeMove(st, i, j, c)
				moveIndex++
			}
		}
	}

	// No moves could be generated. But the player must swap card
	if moveIndex == 0 {
		bb := CurrentMasterBitboard(st)
		from := LSB(bb)

		m := MovePassBase
		m.addFrom(from)
		m.addTo(from)

		st.generatedMoves[0] = m

		m.addCardIndex(1)
		st.generatedMoves[1] = m

		moveIndex = 2
	}

	st.generatedMovesLen = int(moveIndex)
	return moveIndex
}
