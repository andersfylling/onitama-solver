package onitamago

import (
	"math/rand"
	"strconv"
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

const _Card_name = "RabbitCobraRoosterCraneOxHorseBoarCrabEelGooseFrogMantisMonkeyElephantDragonTiger"

var _Card_map = map[Card]string{
	4547854970388480:    _Card_name[0:6],
	4609221463113728:    _Card_name[6:11],
	4627019807588352:    _Card_name[11:18],
	9042727224213504:    _Card_name[18:23],
	9060113251827712:    _Card_name[23:25],
	9112889809960960:    _Card_name[25:30],
	9130344557051904:    _Card_name[30:34],
	9191917208207360:    _Card_name[34:38],
	18067449945522176:   _Card_name[38:41],
	18137612531269632:   _Card_name[41:46],
	18190389089402880:   _Card_name[46:50],
	22553319947894784:   _Card_name[50:56],
	22553526106324992:   _Card_name[56:62],
	22641143439163392:   _Card_name[62:70],
	38316124802121728:   _Card_name[70:76],
	2305878331024736256: _Card_name[76:81],
}

func CardName(card Card) string {
	if str, ok := _Card_map[card]; ok {
		return str
	}
	return "Card(" + strconv.FormatInt(int64(card), 10) + ")"
}
