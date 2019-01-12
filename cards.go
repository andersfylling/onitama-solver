package onitamago

import (
	"math/rand"
	"time"
)

const NrOfActiveCards = 5
const NrOfPlayerCards = 2

// HighestNrOfMoves holds teh highest number of moves on any card
const HighestNrOfMoves = 4

// //go:generate stringer -type=Card
type Card = Board

const (
	// Tiger card
	//  _  _  X  _  _
	//  _  _  _  _  _
	//  _  _  O  _  _
	//  _  _  X  _  _
	Tiger Card = 0x400040400000000 << 3 // TODO...

	// Dragon card
	//  _  _  _  _  _
	//  X  _  _  _  X
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Dragon Card = 0x11040a00000000 << 3 // TODO...

	// Frog card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  X  _  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Frog Card = 0x8140200000000 << 3 // TODO...

	// Rabbit card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  O  _  X  _
	//  X  _  _  _  _
	//  _  _  _  _  _
	Rabbit Card = 0x2050800000000 << 3 // TODO...

	// Crab card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  X  _  O  _  X
	//  _  _  _  _  _
	//  _  _  _  _  _
	Crab Card = 0x4150000000000 << 3 // TODO...

	// Elephant card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  X  O  X  _
	//  _  _  _  _  _
	//  _  _  _  _  _
	Elephant Card = 0xa0e0000000000 << 3 // TODO...

	// Goose card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  X  O  X  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Goose Card = 0x80e0200000000 << 3 // TODO...

	// Rooster card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  _  X  O  X  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Rooster Card = 0x20e0800000000 << 3 // TODO...

	// Monkey card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Monkey Card = 0xa040a00000000 << 3 // TODO...

	// Mantis card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  _  O  _  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Mantis Card = 0xa040400000000 << 3 // TODO...

	// Horse card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  X  O  _  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Horse Card = 0x40c0400000000 << 3 // TODO...

	// Ox card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  X  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Ox Card = 0x4060400000000 << 3 // TODO...

	// Crane card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Crane Card = 0x4040a00000000 << 3 // TODO...

	// Boar card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  X  O  X  _
	//  _  _  _  _  _
	//  _  _  _  _  _
	Boar Card = 0x40e0000000000 << 3 // TODO...

	// Eel card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  _  O  X  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Eel Card = 0x8060800000000 << 3 // TODO...

	// Cobra card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  _  X  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Cobra Card = 0x20c0200000000 << 3 // TODO...

	// CardOffset is how many bit position the initial card masks are shifted
	// remember that offset is number of bit positions. Note that every card
	// has their center at bit position 45.
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	45	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	//  _	_	_	_	_	_	_	_
	CardOffset BoardIndex = 0x2d
)

func rotateCard(card Card) Card {
	// Rotate the move card and shift it into the original position
	card = RotateBoard(card)
	card = card << 3     // columns to the left
	card = card << 8 * 3 // rows up

	return card
}

// TODO: colours

func DrawCards() (selection []Card) {
	cards := []Card{
		Rooster, Rabbit, Ox, Cobra,
		Horse, Goose, Frog, Eel,
		Tiger, Dragon, Crab, Elephant, Monkey, Mantis, Crane, Boar,
	}

	for {
	pick:
		rand.Seed(time.Now().UnixNano())
		card := cards[rand.Int()%len(cards)]

		// check if it already exists
		for i := range selection {
			if card == selection[i] {
				goto pick
			}
		}

		// add card
		selection = append(selection, card)
		if len(selection) == NrOfActiveCards {
			break
		}
	}

	return
}
