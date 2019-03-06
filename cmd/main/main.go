package main

import (
	"fmt"
	"time"
)
import oni "github.com/andersfylling/onitamago"

func perftIterative(cards []oni.Card, depth int) (leafs uint64, moves uint64, duration time.Duration) {
	stack := oni.Stack{}

	st := oni.State{}
	st.CreateGame(cards)
	skipMove := ^oni.Move(0)

	start := time.Now()

	// prepare stack and move indexing
	st.GenerateMoves()
	moves = uint64(st.MovesLen())
	if int(st.Depth()+1) >= depth {
		return moves, moves, time.Now().Sub(start)
	}

	// populate stack with some work
	stack.Push(skipMove)
	stack.PushMany(st.Moves())
	var move oni.Move
	var anyWins bool
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
		moves += uint64(st.MovesLen())

		if int(st.Depth()+1) >= depth {
			leafs += uint64(st.MovesLen())
			st.UndoMove()
		} else {
			stack.Push(skipMove) // identify a new depth
			anyWins = false
			for i := range st.Moves() {
				if (st.Moves()[i] & (1 << 12)) > 0 {
					anyWins = true
					break
				}
			}
			if !anyWins {
				stack.PushMany(st.Moves())
			}
		}
	}

	return leafs, moves, time.Now().Sub(start)
}

func main() {
	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	const depth = 9
	for i := 1; i <= depth; i++ {
		leafs, moves, duration := perftIterative(cards, i)
		perf := float64(moves) / float64(duration.Seconds())
		fmt.Println("depth", i, ",leafs", leafs, fmt.Sprintf(",moves/sec %0.2f", perf), ",duration", duration)
	}
}
