package main

import (
	"fmt"

	"github.com/andersfylling/onitamago"
)

func main() {
	deck := onitamago.Deck(onitamago.DeckOriginal, onitamago.DeckSenseisPath)
	configurations := onitamago.GenCardConfigs(deck)
	fmt.Println(len(configurations))
}
