package onitamago

import (
	"testing"
)

func TestFilterForcedMoves(t *testing.T) {
	cards := []Card{
		Goose, Tiger, Rooster, Rabbit, Dragon,
	}

	SearchExhaustiveForForcedWins(cards, 8)
}
