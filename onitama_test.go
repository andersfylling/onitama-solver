package onitamago

import (
	"fmt"
	"testing"
	"time"
)

//
//func TestOnitama(t *testing.T) {
//	st := NewGame()
//	fmt.Println(st)
//
//	st.GenerateMoves()
//
//	st.ApplyMove(st.generatedMoves[0])
//	fmt.Println(st)
//	st.GenerateMoves()
//	moves := st.generatedMoves[:st.generatedMovesLen]
//	cards := st.playerCards[:]
//	player := st.activePlayer
//	for i := range moves {
//		move := moves[i]
//		if move == 0 {
//			panic("empty")
//		}
//
//		fmt.Println("moved:", explainMove(move, player, cards))
//		st.ApplyMove(move)
//		fmt.Println(st)
//		st.UndoMove()
//	}
//}

func PrintWinnerPath(st State) {
	fmt.Println("WINNER")
	for i := len(st.previousMoves) - 1; i >= 0; i-- {
		if st.previousMoves[i] == 0 {
			continue
		}

		fmt.Println(st)
		st.UndoMove()
		fmt.Println(explainMove(st.previousMoves[i], st.activePlayer, st.playerCards[:]))
	}

	fmt.Println(st)
}

func Perft(st State, depth int) {
	game := Game{}

	start := time.Now()

	var nodes uint64
	var wins uint64

	var current int
	var next int
	for {
		st.GenerateMoves()
		for i := 0; i < st.generatedMovesLen; i++ {
			if getMoveWin(st.generatedMoves[i]) == 1 {
				// skip winner nodes
				wins++
				//st.ApplyMove(st.generatedMoves[i])
				//PrintWinnerPath(st)
				//st.UndoMove()
			} else {
				next++
				game.Tree[next] = st.generatedMoves[i]
			}
		}

		genDepth := int(st.currentDepth) + 1
		if genDepth < depth {
			current = next
			next-- // to overwrite the move in the game tree, such that it is not re-applied
		} else {
			nodes += uint64(st.generatedMovesLen) //includes winner nodes

			if current == 0 {
				// done iterating the game tree
				break
			}

			current--
			next = current
			st.UndoMove()
		}
		st.ApplyMove(game.Tree[current])
	}

	fmt.Println("depth", depth, ", time:", time.Now().Sub(start), ", moves", nodes, ", wins", wins)
}

func TestPerft(t *testing.T) {
	st := NewState()
	st.CreateGame([]Card{
		Frog, Eel,
		Dragon, Crab,
		Tiger,
	})

	const depth = 3
	for i := 1; i <= depth; i++ {
		Perft(st, i)
	}
}

func BenchmarkState_GenerateMoves(b *testing.B) {
	st := State{}
	st.CreateGame(nil)
	for i := 0; i < b.N; i++ {
		st.GenerateMoves()
	}
}
