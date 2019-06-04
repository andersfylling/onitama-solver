package onitamago

import (
	"math/rand"
	"time"
)

const NrOfActiveCards = 5
const NrOfPlayerCards = 2

// HighestNrOfMoves holds teh highest number of moves on any card
const HighestNrOfMoves = 4

const (
	// Tiger card
	//  _  _  X  _  _
	//  _  _  _  _  _
	//  _  _  O  _  _
	//  _  _  X  _  _
	Tiger        Card = 0x2000202000000000
	TigerRotated Card = 0x20200020000000

	// Dragon card
	//  _  _  _  _  _
	//  X  _  _  _  X
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Dragon        Card = 0x88205000000000
	DragonRotated Card = 0x50208800000000

	// Frog card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  X  _  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Frog        Card = 0x40A01000000000
	FrogRotated Card = 0x40281000000000

	// Rabbit card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  O  _  X  _
	//  X  _  _  _  _
	//  _  _  _  _  _
	Rabbit        Card = 0x10284000000000
	RabbitRotated Card = 0x10A04000000000

	// Crab card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  X  _  O  _  X
	//  _  _  _  _  _
	//  _  _  _  _  _
	Crab        Card = 0x20a80000000000
	CrabRotated Card = 0xa82000000000

	// Elephant card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  X  O  X  _
	//  _  _  _  _  _
	//  _  _  _  _  _
	Elephant        Card = 0x50700000000000
	ElephantRotated Card = 0x705000000000

	// Goose card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  X  O  X  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Goose        Card = 0x40701000000000
	GooseRotated Card = Goose

	// Rooster card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  _  X  O  X  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Rooster        Card = 0x10704000000000
	RoosterRotated Card = Rooster

	// Monkey card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Monkey        Card = 0x50205000000000
	MonkeyRotated Card = Monkey

	// Mantis card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  _  _  O  _  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Mantis        Card = 0x50202000000000
	MantisRotated Card = Crane

	// Horse card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  X  O  _  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Horse        Card = 0x20602000000000
	HorseRotated Card = Ox

	// Ox card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  X  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Ox        Card = 0x20302000000000
	OxRotated Card = Horse

	// Crane card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  _  _
	//  _  X  _  X  _
	//  _  _  _  _  _
	Crane        Card = 0x20205000000000
	CraneRotated Card = Mantis

	// Boar card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  X  O  X  _
	//  _  _  _  _  _
	//  _  _  _  _  _
	Boar        Card = 0x20700000000000
	BoarRotated Card = 0x702000000000

	// Eel card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  _  O  X  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Eel        Card = 0x40304000000000
	EelRotated Card = Cobra

	// Cobra card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  _  X  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Cobra        Card = 0x10601000000000
	CobraRotated Card = Eel

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
	CardOffset BitboardPos = 45
)

// //go:generate stringer -type=Card
type Card Bitboard

func (c Card) Name() string {
	if str, ok := _Card_map[c]; ok {
		return str
	}

	panic("no such card")
}

func (c Card) Bitboard() Bitboard {
	return Bitboard(c)
}

func (c Card) Heatmap() (b BitboardHeatmap) {
	b.AddCard(c)
	return b
}

func (c *Card) Rotate() {
	switch *c {
	case Tiger:
		*c = TigerRotated
	case TigerRotated:
		*c = Tiger
	case Dragon:
		*c = DragonRotated
	case DragonRotated:
		*c = Dragon
	case Frog:
		*c = FrogRotated
	case FrogRotated:
		*c = Frog
	case Rabbit:
		*c = RabbitRotated
	case RabbitRotated:
		*c = Rabbit
	case Elephant:
		*c = ElephantRotated
	case ElephantRotated:
		*c = Elephant
	case Crab:
		*c = CrabRotated
	case CrabRotated:
		*c = Crab
	case Boar:
		*c = BoarRotated
	case BoarRotated:
		*c = Boar

	case Horse: // opposite of Ox
		*c = HorseRotated
	case Ox: // opposite of Horse
		*c = OxRotated

	case Crane: // opposite of Mantis
		*c = Mantis
	case Mantis: // opposite of Crane
		*c = Crane

	case Eel: // Opposite of Cobra
		*c = Cobra
	case Cobra: // Opposite of Eel
		*c = Eel
	}
}

// TODO: colours

func DrawCards() (selection []Card) {
	cards := []Card{
		Rooster, Rabbit, Ox, Cobra,
		Horse, Goose, Frog, Eel,
		Tiger, Dragon, Crab, Elephant, Monkey, Mantis, Crane, Boar,
	}

	selection = make([]Card, 0, NrOfActiveCards)

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

const _Card_name = "RabbitCobraRoosterCraneOxHorseBoarCrabEelGooseFrogMantisMonkeyElephantDragonTiger"

var _Card_map = map[Card]string{
	Rabbit:   _Card_name[0:6],
	Cobra:    _Card_name[6:11],
	Rooster:  _Card_name[11:18],
	Crane:    _Card_name[18:23],
	Ox:       _Card_name[23:25],
	Horse:    _Card_name[25:30],
	Boar:     _Card_name[30:34],
	Crab:     _Card_name[34:38],
	Eel:      _Card_name[38:41],
	Goose:    _Card_name[41:46],
	Frog:     _Card_name[46:50],
	Mantis:   _Card_name[50:56],
	Monkey:   _Card_name[56:62],
	Elephant: _Card_name[62:70],
	Dragon:   _Card_name[70:76],
	Tiger:    _Card_name[76:81],
}
