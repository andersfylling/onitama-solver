package movetree

import (
	"fmt"
	"testing"

	. "github.com/andersfylling/onitamago"
)

func TestWinPath(t *testing.T) {

	cards := []Card{
		Rooster, Rabbit, Cobra, Mantis, Boar,
	}

	_, wins, d := SearchExhaustive(cards, 6)
	fmt.Println(d)
	fmt.Println(wins.Instances)

	st := State{}
	st.CreateGame(cards)
	fmt.Println(st)
	node := wins
	for {
		if len(node.Paths) == 0 {
			break
		}

		var move Move
		for k, v := range node.Paths {
			move = k
			node = v
			break
		}

		st.ApplyMove(move)
		fmt.Println(st)
	}
	//
	// heatmap := BitboardHeatmap{}
	//
	// for i := range wins {
	// 	board := Bitboard(0)
	// 	bullet := BitboardPos(wins[i][len(wins[i])-1] & MoveMaskTo)
	// 	board |= Bitboard(1) << bullet
	// 	heatmap.AddBoard(board)
	// }
}

func TestWinPath2(t *testing.T) {

	cards := []Card{
		Rooster, Rabbit, Cobra, Mantis, Boar,
	}

	_, wins, d := SearchForTempleWins(cards, 6)
	fmt.Println(d)
	fmt.Println(len(wins))
}
