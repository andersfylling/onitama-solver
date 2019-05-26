package main

import (
	perft "github.com/andersfylling/onitamago/benchmarks/perft"
	oni "github.com/andersfylling/onitamago"
	"fmt"
)

func main() {
	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	_, _, _, d := perft.Perft(cards, 10)
	fmt.Println("duration", d)
}
