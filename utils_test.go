package onitamago

import "testing"

func TestMakeCompactBoard(t *testing.T) {
	normal := Board(0x2a1408140200)
	compact := Board(0x1551141)

	got := MakeCompactBoard(normal)
	if got != compact {
		t.Errorf("got different compacted board. Got %d, wants %d", got, compact)
	}
}
