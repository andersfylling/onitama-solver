package onitamago

import (
	"fmt"
	"testing"
)

func TestOnitama(t *testing.T) {
	st := NewGame()
	fmt.Println(st)

	st.GenerateMoves()
	moves := st.generatedMoves
	cards := st.playerCards[:]
	player := st.activePlayer
	for i := range moves {
		move := moves[i]
		if move == 0 {
			continue
		}

		fmt.Println("moved:", explainMove(move, player, cards))
		st.ApplyMove(move)
		fmt.Println(st)
		st.UndoMove()
	}
}
