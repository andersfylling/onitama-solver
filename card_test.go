package onitamago

import (
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
}
