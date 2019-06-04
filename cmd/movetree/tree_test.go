package movetree

import (
	"fmt"
	. "github.com/andersfylling/onitamago"
	"github.com/disintegration/imaging"
	"testing"
)

func TestWinPath(t *testing.T) {

	cards := []Card{
		Rooster, Rabbit, Cobra, Mantis, Boar,
	}

	_, wins, d := SearchExhaustive(cards, 10)
	fmt.Println(d)
	fmt.Println(len(wins))

	heatmap := BitboardHeatmap{}

	for i := range wins {
		board := Bitboard(0)
		bullet := BitboardPos(wins[i][len(wins[i])-1] & MoveMaskTo)
		board |= Bitboard(1) << bullet
		heatmap.AddBoard(board)
	}

	imaging.Save(heatmap.Render(0.05), "wins.png")

	win := wins[15000]

	st := State{}
	st.CreateGame(cards)
	fmt.Println(st)

	for i := range win {
		st.ApplyMove(win[i])
		fmt.Println(st)
	}
}
