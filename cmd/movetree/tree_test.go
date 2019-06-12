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

	_, wins, d := SearchExhaustive(cards, 8)
	fmt.Println(d)
	fmt.Println(len(wins))

	heatmap := BitboardHeatmap{}

	for i := range wins {
		board := Bitboard(0)
		bullet := BitboardPos(wins[i][len(wins[i])-1] & MoveMaskTo)
		board |= Bitboard(1) << bullet
		heatmap.AddBoard(board)
	}
}
