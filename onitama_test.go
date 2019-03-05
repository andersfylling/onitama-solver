package onitamago

import (
	"fmt"
	"math/rand"
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

func Perft(cards []Card, depth int) {
	//game := Game{}

	st := &State{
		otherPlayer:  1,
		activePlayer: 0,
	}
	st.CreateGame(cards)

	var nodes uint64
	var totalNodes uint64 // inc. root
	var wins uint64

	//var current int
	//var next int

	var perft func(int) uint64
	perft = func(depth int) uint64 {
		st.GenerateMoves()
		if st.generatedMovesLen == 0 {
			return 0
		}
		if depth == 1 {
			for i := range st.generatedMoves[:st.generatedMovesLen] {
				if getMoveWin(st.generatedMoves[i]) == 1 {
					wins++
				}
			}
			totalNodes += uint64(st.generatedMovesLen)
			return uint64(st.generatedMovesLen)
		}
		var nodes uint64

		moves := make([]Move, st.generatedMovesLen)
		copy(moves, st.generatedMoves[:st.generatedMovesLen])
		for i := range moves {
			st.ApplyMove(moves[i])
			totalNodes++
			if st.hasWon {
				wins++
				st.UndoMove() //moves[i])
				break
			} else {
				nodes += perft(depth - 1)
			}
			st.UndoMove() //moves[i])
		}
		return nodes
	}

	start := time.Now()
	nodes++ // root node
	if depth > 0 {
		nodes = perft(depth)
	}
	fmt.Println("depth", depth, ", time:", time.Now().Sub(start), ", moves", nodes, ", wins", wins, ", total", totalNodes)

	//
	//for {
	//	st.GenerateMoves()
	//	for i := 0; i < st.generatedMovesLen; i++ {
	//		totalNodes++
	//		if getMoveWin(st.generatedMoves[i]) == 1 {
	//			// skip winner nodes
	//			wins++
	//			//st.ApplyMove(st.generatedMoves[i])
	//			//PrintWinnerPath(st)
	//			//st.UndoMove()
	//		} else {
	//			next++
	//			game.Tree[next] = st.generatedMoves[i]
	//		}
	//	}
	//
	//	genDepth := int(st.currentDepth) + 1
	//	if genDepth == depth {
	//		nodes += uint64(st.generatedMovesLen) //includes winner nodes
	//		current--
	//		next = current
	//		st.UndoMove()
	//	} else if genDepth < depth {
	//		if current > next {
	//			// if none was generated, then it's a win and should never been seen as a
	//			// node of interest
	//			panic("node generated 0 children")
	//		}
	//		current = next
	//		next-- // to overwrite the move in the game tree, such that it is not re-applied
	//	} else {
	//		panic("went below given depth")
	//	}
	//
	//	if current <= 0 {
	//		// done iterating the game tree
	//		break
	//	}
	//
	//	st.ApplyMove(game.Tree[current])
	//}

}

func perftIterative(cards []Card, depth int) (leafs uint64, moves uint64, duration time.Duration) {
	stack := Stack{}

	st := State{}
	st.CreateGame(cards)
	skipMove := ^Move(0)

	start := time.Now()

	// prepare stack and move indexing
	st.GenerateMoves()
	moves = uint64(st.generatedMovesLen)
	if int(st.currentDepth+1) >= depth {
		return moves, moves, time.Now().Sub(start)
	}

	// populate stack with some work
	stack.Push(skipMove)
	stack.PushMany(st.generatedMoves[:st.generatedMovesLen])
	var move Move
	var anyWins Move
	for {
		if move = stack.Pop(); move == skipMove {
			// finished processing node children
			for ; move == skipMove && stack.Size() > 0; move = stack.Pop() {
				st.UndoMove()
			}
			if stack.Size() == 0 {
				break
			}
		}

		st.ApplyMove(move)
		st.GenerateMoves()
		moves += uint64(st.generatedMovesLen)

		if int(st.currentDepth+1) >= depth {
			leafs += uint64(st.generatedMovesLen)
			st.UndoMove()
		} else {
			stack.Push(skipMove) // identify a new depth
			anyWins = 0
			for i := range st.generatedMoves[:st.generatedMovesLen] {
				anyWins |= st.generatedMoves[i] & (1 << 12)
			}
			if anyWins == 0 {
				stack.PushMany(st.generatedMoves[:st.generatedMovesLen])
			}
		}
	}

	return leafs, moves, time.Now().Sub(start)
}

func TestPerft(t *testing.T) {
	cards := []Card{
		Frog, Eel,
		Dragon, Crab,
		Tiger,
	}

	const depth = 4
	for i := 1; i <= depth; i++ {
		leafs, moves, duration := perftIterative(cards, i)
		perf := float64(moves) / float64(duration.Seconds())
		fmt.Println("depth", i, ",leafs", leafs, fmt.Sprintf(",moves/sec %0.2f", perf), ",duration", duration)
	}
}

func BenchmarkState_GenerateMoves(b *testing.B) {
	st := State{}
	st.CreateGame(nil)
	for i := 0; i < b.N; i++ {
		st.GenerateMoves()
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
		//fmt.Println(st)
	}

	st.GenerateMoves()
	if st.generatedMovesLen > 0 {
		t.Errorf("expected 0 children to be generated. Got %d", st.generatedMovesLen)
	}
	//for i := range st.generatedMoves[:st.generatedMovesLen] {
	//	fmt.Println(i, st.generatedMoves[i], explainMove(st.generatedMoves[i], st.activePlayer, st.playerCards[:]))
	//	st.ApplyMove(st.generatedMoves[i])
	//	fmt.Println(st)
	//	st.UndoMove()
	//}
	//
	//fmt.Println("REVERSING")
	//for {
	//	move := st.previousMoves[st.currentDepth]
	//	st.UndoMove()
	//	fmt.Println(explainMove(move, st.activePlayer, st.playerCards[:]))
	//	fmt.Println(st)
	//	if st.currentDepth == 0 {
	//		break
	//	}
	//}

}

func TestRandomSampling(t *testing.T) {
	st := State{}
	st.CreateGame(nil)

	stCopy := st

	var moves []Move
	//fmt.Println(st)
	for i := 0; i < 20; i++ {
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
			cards += CardName(stCopy.playerCards[i]) + ", "
		}
		cards += CardName(stCopy.suspendedCard)
		t.Errorf("apply and undo move create different roots. Card config: %s\n", cards)
		fmt.Println(st)
		fmt.Println(stCopy)
	}
}
