package onitamago

import (
	"fmt"
	"testing"
)

func check(t *testing.T, got, wants Card) {
	if got != wants {
		t.Errorf("incorrect result. Got %d, wants %d", got, wants)
	}
}

func TestRotateCard(t *testing.T) {
	var card Card
	var wants Card
	var got Card

	card = Elephant
	wants = 0x705000000000

	got = card
	got.Rotate()

	check(t, got, wants)

	names := []Card{
		Tiger, Goose, Rooster, Rabbit, Horse,
	}
	for _, card := range names {
		fmt.Print(uint64(card), ",")
	}

	cards := Deck(DeckOriginal, DeckSenseisPath)
	for i := range cards {
		cards[i].Rotate() // panic on unsupported cards
	}
}

func TestCard_PrettyString(t *testing.T) {
	roosterStr := `_  _  _  _  _
_  _  _  X  _
_  X  X  X  _
_  X  _  _  _
_  _  _  _  _`
	if roosterStr != Rooster.String() {
		t.Error("different strings")
	}
}

func TestCard_Name(t *testing.T) {
	cards := Deck(DeckSenseisPath, DeckOriginal)

	// make sure name can be retrieved for each type
	for _, card := range cards {
		card.Name()
	}
}
