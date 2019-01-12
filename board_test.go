package onitamago

import "testing"

func TestRotateBoard(t *testing.T) {
	if BoardMask != RotateBoard(BoardMask) {
		t.Errorf("playing board was not correctly rotate. Got %d, wants %d", RotateBoard(BoardMask), BoardMask)
	}
}
