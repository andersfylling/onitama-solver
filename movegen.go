package onitamago

func pieces(player BoardIndex, board [NrOfPlayers * NrOfPieceTypes]Board) Board {
	// bind the function to our compile time assumption
	if NrOfPieceTypes != 3 {
		panic("NrOfPieceTypes is no longer 2")
	}
	return board[player*NrOfPieceTypes] | board[player*NrOfPieceTypes+1]
}

func generateMoves(st *State) (moveIndex Index) {
	friends := pieces(st.activePlayer, st.board)
	for c := st.activePlayer * NrOfPlayerCards; c < (st.activePlayer*NrOfPlayerCards + NrOfPlayerCards); c++ {
		moves := st.playerCards[c]
		// TODO: remove if sentence
		// add some virtual layer
		if st.activePlayer == OppositePlayer {
			moves = rotateCard(moves)
		}
		pieces := friends
		for i := LSB(pieces); i != 64; i = NLSB(&pieces, i) {
			move := moves >> (CardOffset - i)
			move ^= move & friends // remove moves that hits a friendly warrior
			move &= BoardMask       // ignore positions outside the board

			for j := LSB(move); j != 64; j = NLSB(&move, j) {
				st.generatedMoves[moveIndex] = encodeMove(st, i, j, c)
				moveIndex++
			}
		}
	}

	st.generatedMovesLen = int(moveIndex)
	return moveIndex
}
