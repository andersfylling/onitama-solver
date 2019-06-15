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
	DeckOriginal = iota
	DeckSenseisPath

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

// DeckOriginal cards
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
)

// DeckSenseisPath cards
const (
	// Turtle card
	//  _  _  _  _  _
	//  _  _  _  _  _
	//  X  _  O  _  X
	//  _  X  _  X  _
	//  _  _  _  _  _
	Turtle        Card = 0xa85000000000
	TurtleRotated Card = Pheonix

	// Pheonix card
	//  _  _  _  _  _
	//  _  X  _  X  _
	//  X  _  O  _  X
	//  _  _  _  _  _
	//  _  _  _  _  _
	Pheonix        Card = 0x50a80000000000
	PheonixRotated Card = Turtle

	// Otter card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  _  O  _  X
	//  _  _  _  X  _
	//  _  _  _  _  _
	Otter        Card = 0x40281000000000
	OtterRotated Card = 0x40a01000000000

	// Iguana card
	//  _  _  _  _  _
	//  X  _  X  _  _
	//  _  _  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Iguana        Card = 0xa0201000000000
	IguanaRotated Card = 0x40202800000000

	// Sable card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  X  _  O  _  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Sable        Card = 0x10a04000000000
	SableRotated Card = 0x10284000000000

	// Panda card
	//  _  _  _  _  _
	//  _  _  X  X  _
	//  _  _  O  _  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Panda        Card = 0x30204000000000
	PandaRotated Card = 0x10206000000000

	// Bear card
	//  _  _  _  _  _
	//  _  X  X  _  _
	//  _  _  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Bear        Card = 0x60201000000000
	BearRotated Card = 0x40203000000000

	// Fox card
	//  _  _  _  _  _
	//  _  _  _  X  _
	//  _  _  O  X  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Fox        Card = 0x10301000000000
	FoxRotated Card = Dog

	// Giraffe card
	//  _  _  _  _  _
	//  X  _  _  _  X
	//  _  _  O  _  _
	//  _  _  X  _  _
	//  _  _  _  _  _
	Giraffe        Card = 0x88202000000000
	GiraffeRotated Card = 0x20208800000000

	// Kirin card
	//  _  X  _  X  _
	//  _  _  _  _  _
	//  _  _  O  _  _
	//  _  _  _  _  _
	//  _  _  X  _  _
	Kirin        Card = 0x5000200020000000
	KirinRotated Card = 0x2000200050000000

	// Rat card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  X  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Rat        Card = 0x20601000000000
	RatRotated Card = 0x40302000000000

	// Tanuki card
	//  _  _  _  _  _
	//  _  _  X  _  X
	//  _  _  O  _  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Tanuki        Card = 0x28204000000000
	TanukiRotated Card = 0x1020a000000000

	// Mouse card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  X  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Mouse        Card = 0x20304000000000
	MouseRotated Card = 0x10602000000000

	// Viper card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  X  _  O  _  _
	//  _  _  _  X  _
	//  _  _  _  _  _
	Viper        Card = 0x20a01000000000
	ViperRotated Card = 0x40282000000000

	// Sea Snake card
	//  _  _  _  _  _
	//  _  _  X  _  _
	//  _  _  O  _  X
	//  _  X  _  _  _
	//  _  _  _  _  _
	SeaSnake        Card = 0x20284000000000
	SeaSnakeRotated Card = 0x10a02000000000

	// Dog card
	//  _  _  _  _  _
	//  _  X  _  _  _
	//  _  X  O  _  _
	//  _  X  _  _  _
	//  _  _  _  _  _
	Dog        Card = 0x40604000000000
	DogRotated Card = Fox
)

// //go:generate stringer -type=Card
type Card Bitboard

func (c Card) Name() string {
	if str, ok := _card_name[c]; ok {
		return str
	}

	panic("no such card")
}

func (c Card) String() string {
	var tmp string
	add := func(mark bool) {
		if mark {
			tmp += "X"
		} else {
			tmp += "_"
		}
	}

	p := 63
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			bb := Bitboard(1) << BitboardPos(p-8*y-x)
			add(bb&c.Bitboard() > 0)
		}
	}

	// add spaces and new lines
	var b []byte
	for i := range tmp {
		b = append(b, tmp[i])
		if (i+1)%5 == 0 {
			b = append(b, '\n')
		} else {
			b = append(b, []byte("  ")...)
		}
	}

	return string(b[:len(b)-1])
}

func (c Card) Bitboard() Bitboard {
	return Bitboard(c)
}

func (c Card) Heatmap() (b BitboardHeatmap) {
	b.AddCard(c)
	return b
}

var _rotated = map[Card]Card{
	Tiger:           TigerRotated,
	TigerRotated:    Tiger,
	Rooster:         RoosterRotated, // Rooster == RoosterRotated
	Goose:           GooseRotated,   // Goose == GooseRotated
	Monkey:          MonkeyRotated,  // Monkey == MonkeyRotated
	Dragon:          DragonRotated,
	DragonRotated:   Dragon,
	Frog:            Otter,
	Rabbit:          Sable,
	Elephant:        ElephantRotated,
	ElephantRotated: Elephant,
	Crab:            CrabRotated,
	CrabRotated:     Crab,
	Boar:            BoarRotated,
	BoarRotated:     Boar,
	Horse:           HorseRotated,
	Ox:              OxRotated,
	Crane:           Mantis,
	Mantis:          Crane,
	Eel:             Cobra,
	Cobra:           Eel,

	// expansion sensei's path
	Turtle:          Pheonix,
	Pheonix:         Turtle,
	Otter:           Frog,
	Iguana:          IguanaRotated,
	IguanaRotated:   Iguana,
	Sable:           Rabbit,
	Bear:            BearRotated,
	BearRotated:     Bear,
	Panda:           PandaRotated,
	PandaRotated:    Panda,
	Giraffe:         GiraffeRotated,
	GiraffeRotated:  Giraffe,
	Kirin:           KirinRotated,
	KirinRotated:    Kirin,
	Rat:             RatRotated,
	RatRotated:      Rat,
	Tanuki:          TanukiRotated,
	TanukiRotated:   Tanuki,
	Mouse:           MouseRotated,
	MouseRotated:    Mouse,
	Viper:           ViperRotated,
	ViperRotated:    Viper,
	SeaSnake:        SeaSnakeRotated,
	SeaSnakeRotated: SeaSnake,
	Dog:             Fox,
	Fox:             Dog,
}

func (c *Card) Rotate() {
	if r, ok := _rotated[*c]; ok {
		*c = r
		return
	}

	panic("unknown rotation\n" + c.String())
}

func (c Card) Rotated() Card {
	if r, ok := _rotated[c]; ok {
		return r
	}

	panic("unknown rotation")
}

func Deck(deckTypes ...uint) (cards []Card) {
	if len(deckTypes) == 0 {
		deckTypes = append(deckTypes, DeckOriginal)
	}

	decks := [][]Card{
		/* original */ {
			Rooster, Rabbit, Ox, Cobra,
			Horse, Goose, Frog, Eel,
			Tiger, Dragon, Crab, Elephant,
			Monkey, Mantis, Crane, Boar,
		},
		/* sensei's path */ {
			Turtle, Pheonix, Otter, Iguana,
			Sable, Panda, Bear, Fox,
			Giraffe, Kirin, Rat, Tanuki,
			Mouse, Viper, SeaSnake, Dog,
		},
	}

	for _, deck := range deckTypes {
		cards = append(cards, decks[deck]...)
	}
	return cards
}

// CardConfig create a card configuration with awareness of which players holds which cards
// and what the idle card is.
func CardConfig(blue [2]Card, brown [2]Card, idle Card) []Card {
	return []Card{
		brown[0], brown[1],
		blue[0], blue[1],
		idle,
	}
}

// DrawCards draws five random cards from the original 16 card deck
func DrawCards() (selection []Card) {
	cards := Deck()

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

// GenCardConfigs generates all possible unique card configurations that respects
// the ordered set ({a1, a2}, {a3, a4}, a5)
func GenCardConfigs(selection []Card) (configs [][]Card) {
	rmDuplicates := func(s []Card) (uniques []Card) {
		for i := range s {
			var exists bool
			for j := range uniques {
				if s[i] == uniques[j] {
					exists = true
					break
				}
			}

			if exists {
				continue
			}

			uniques = append(uniques, s[i])
		}

		return uniques
	}

	genTuples := func(s []Card) (tuples [][2]Card) {
		for i := range s {
			for j := i + 1; j < len(s); j++ {
				tuples = append(tuples, [2]Card{s[i], s[j]})
			}
		}

		return tuples
	}

	missingCards := func(s, taken []Card) (missing []Card) {
		for i := range s {
			var match bool
			for j := range taken {
				if s[i] == taken[j] {
					match = true
					break
				}
			}

			if !match {
				missing = append(missing, s[i])
			}
		}

		return missing
	}

	selection = rmDuplicates(selection)
	tuples := genTuples(selection)

	// generate the first part of the ordered set ({a1, a2}, {a3, a4}, a5);
	// two first sub-set a1-a4.
	var bases [][]Card
	for i := range tuples {
		for j := i + 1; j < len(tuples); j++ {
			p1 := append(tuples[i][:], tuples[j][:]...)
			bases = append(bases, p1)

			p2 := append(tuples[j][:], tuples[i][:]...)
			bases = append(bases, p2)
		}
	}

	// combine a5 and the first sub sets
	for i := range bases {
		b := bases[i]
		options := missingCards(selection, b)

		for j := range options {
			config := append(b, options[j])
			configs = append(configs, config)
		}
	}

	// remove duplicates in each card config
	prev := configs
	configs = configs[:0]
	for i := range prev {
		config := rmDuplicates(prev[i])
		if len(config) == len(prev[i]) {
			configs = append(configs, config)
		}
	}

	return configs
}

var _card_name = map[Card]string{
	Tiger:           "Tiger",
	TigerRotated:    "TigerRotated",
	Rooster:         "Rooster", // Rooster == RoosterRotated
	Goose:           "Goose",   // Goose == GooseRotated
	Monkey:          "Monkey",  // Monkey == MonkeyRotated
	Dragon:          "Dragon",
	DragonRotated:   "DragonRotated",
	Frog:            "Frog",
	Rabbit:          "Rabbit",
	Elephant:        "Elephant",
	ElephantRotated: "ElephantRotated",
	Crab:            "Crab",
	CrabRotated:     "CrabRotated",
	Boar:            "Boar",
	BoarRotated:     "BoarRotated",
	Horse:           "Horse",
	Ox:              "Ox",
	Crane:           "Crane",
	Mantis:          "Mantis",
	Eel:             "Eel",
	Cobra:           "Cobra",

	// expansion sensei's path
	Turtle:          "Turtle",
	Pheonix:         "Pheonix",
	Otter:           "Otter",
	Iguana:          "Iguana",
	IguanaRotated:   "IguanaRotated",
	Sable:           "Sable",
	Bear:            "Bear",
	BearRotated:     "BearRotated",
	Panda:           "Panda",
	PandaRotated:    "PandaRotated",
	Giraffe:         "Giraffe",
	GiraffeRotated:  "GiraffeRotated",
	Kirin:           "Kirin",
	KirinRotated:    "KirinRotated",
	Rat:             "Rat",
	RatRotated:      "RatRotated",
	Tanuki:          "Tanuki",
	TanukiRotated:   "TanukiRotated",
	Mouse:           "Mouse",
	MouseRotated:    "MouseRotated",
	Viper:           "Viper",
	ViperRotated:    "ViperRotated",
	SeaSnake:        "SeaSnake",
	SeaSnakeRotated: "SeaSnakeRotated",
	Dog:             "Dog",
	Fox:             "Fox",
}
