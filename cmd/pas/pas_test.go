package pas

import (
	"fmt"
	"testing"

	"github.com/andersfylling/onitamago"
)

// Does there exist a situation where a user _have_ to say pass? Where the pieces can not move?

func TestMustSayPass(t *testing.T) {
	// simplify the problem to N friendly pieces, as only these can be blocking, and
	// merge the two player cards into one. Then for any piece configuration does there
	// exist a situation where the two merged cards does not allow the piece(s) to move?

	generateMergeCards := func(cards []onitamago.Card) (merged []onitamago.Card) {
		for i := range cards {
			for j := i + 1; j < len(cards); j++ {
				merged = append(merged, cards[i] | cards[j])
			}
		}

		return merged
	}

	mergedCards := generateMergeCards(onitamago.Deck())
	fmt.Println(len(mergedCards))

	board := onitamago.Bitboard(0x200)
	for i := onitamago.Bitboard(0); i < 5; i++ {
		for j := onitamago.Bitboard(0); j < 5; j++ {
			bb := (board << i) << (8*j)

			for _, card := range mergedCards {
				if ((bb | card.Bitboard()) & onitamago.BoardMask) ^ bb == 0 {
					panic("crap")
				}
			}
		}
	}
}