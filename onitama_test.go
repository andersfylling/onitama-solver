package onitamago

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestState_GenerateMoves_Pass(t *testing.T) {
	cards := CardConfig([2]Card{Tiger, Ox}, [2]Card{Horse, Crane}, Dragon)
	st := State{}
	st.CreateGame(cards)

	// move all blue pieces to the right lane
	st.board[NrOfPieceTypes*BluePlayer+StudentsIndex] = 0x20200020200
	st.board[NrOfPieceTypes*BluePlayer+MasterIndex] = 0x2000000

	// remove brown piece that is in the way
	st.board[NrOfPieceTypes*BrownPlayer+StudentsIndex] ^= 0x20000000000

	// This creates a situation where the blue player can not move any pieces and must call pass
	st.GenerateMoves()
	if st.generatedMovesLen != 2 {
		t.Error("expected there to be exactly 2 moves")
	}

	for i := 0; i < st.generatedMovesLen; i++ {
		move := st.generatedMoves[i]
		if !IsPassMove(move) {
			desc := explainMove(st.generatedMoves[i], BluePlayer, cards)
			t.Errorf("expected move to be a pass move, got %d, description %s", move, desc)
		}
	}
}

func TestState_GenerateMoves_Winner(t *testing.T) {
	st := NewState()
	st.CreateGame([]Card{
		Frog, Eel,
		Dragon, Crab,
		Tiger,
	})

	moves := []Move{
		50005, // student{A1 => A2, Crab}
		52003, // student{B5 => C4, Eel}
		50533, // student{A2 => A4, Tiger}
		51419, // student{C4 => C3, Crab}
		6507,  // student{A4 => C5, Dragon}, WINNER
	}

	for i := range moves {
		st.ApplyMove(moves[i])
	}

	st.GenerateMoves()
	if st.generatedMovesLen > 0 {
		t.Errorf("expected 0 children to be generated. Got %d", st.generatedMovesLen)
	}

}

func TestRandomSampling(t *testing.T) {
	st := State{}
	st.CreateGame(nil)

	stCopy := st

	var moves []Move
	//fmt.Println(st)
	for i := 0; i < MaxDepth-1; i++ {
		st.GenerateMoves()
		if st.generatedMovesLen == 0 {
			break
		}

		//fmt.Println(explainMove(st.generatedMoves[0], st.activePlayer, st.playerCards[:]))
		r := rand.Intn(st.generatedMovesLen)
		move := st.generatedMoves[r]
		moves = append(moves, move)
		st.ApplyMove(move)
		//fmt.Println(st)
	}

	//fmt.Println("REVERSING")
	for i := len(moves) - 1; i >= 0; i-- {
		//move := st.previousMoves[st.currentDepth]
		st.UndoMove() //moves[i])
		//fmt.Println(explainMove(move, st.activePlayer, st.playerCards[:]))
		//fmt.Println(st)
		if st.currentDepth == 0 {
			break
		}
	}

	if st.String() != stCopy.String() {
		var cards string
		for i := range stCopy.playerCards {
			cards += stCopy.playerCards[i].Name() + ", "
		}
		cards += stCopy.suspendedCard.Name()
		t.Errorf("apply and undo move create different roots. Card config: %s\n", cards)
		fmt.Println(st)
		fmt.Println(stCopy)
	}
}

var sink bool

func BenchmarkSkippingWinningDepths(b *testing.B) {
	data := []int{1, 4, 1, 1, 1, 2, 1, 2, 3, 1, 2, 1, 2}

	//inline
	extract := func(i int) int {
		return (i & 4) >> 2 // (x & 0b100) >> 2 => 0 or 1
	}

	b.Run("merge", func(b *testing.B) {
		var m int
		for i := 0; i < b.N; i++ {
			m = 0
			for i := range data {
				m |= extract(data[i])
			}
			sink = m == 1
		}
	})

	b.Run("if break", func(b *testing.B) {
		var anyWins bool
		for i := 0; i < b.N; i++ {
			anyWins = false
			for i := range data {
				if extract(data[i]) > 0 {
					anyWins = true
					break
				}
			}
			if !anyWins {
				sink = true
			}
		}
	})
}

func BenchmarkState_GenerateMoves(b *testing.B) {
	st := State{}
	st.CreateGame(nil)
	for i := 0; i < b.N; i++ {
		st.GenerateMoves()
	}
}
