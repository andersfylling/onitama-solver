package main

import (
	"fmt"
	oni "github.com/andersfylling/onitamago"
	perft "github.com/andersfylling/onitamago/benchmarks/perft"
)

func main() {
	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	// perft(11):
	//  - vanilla 32m41.239508707s
	//  - cached 34m7.398049079s
	_, _, _, d := perft.Perft(cards, 11)
	fmt.Println("duration", d)
}
