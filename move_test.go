package onitamago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMove_MoveFriendlyBoardIndex(t *testing.T) {
	var f = func(action MoveAction) BoardIndex {
		move := setMoveAction(0, action)
		return getMoveFriendlyBoardIndex(move)
	}

	// find the piece that is being moved.
	// 1 = master, 0 = student
	assert.EqualValues(t, f(0), 0)
	assert.EqualValues(t, f(1), 0)
	assert.EqualValues(t, f(2), 1)
	assert.EqualValues(t, f(3), 1)
	assert.EqualValues(t, f(4), 0)
	assert.EqualValues(t, f(5), 0)
	assert.EqualValues(t, f(6), 1)
	assert.EqualValues(t, f(7), 1)
}

func TestMove_MoveHostileBoardIndex(t *testing.T) {
	var f = func(action MoveAction) BoardIndex {
		move := setMoveAction(0, action)
		return getMoveHostileBoardIndex(move)
	}

	// find hostile piece being attacked/destroyed.
	// 0 = student, 1 = master, 2 = temple or no attack
	assert.EqualValues(t, f(0), 0)
	assert.EqualValues(t, f(1), 1)
	assert.EqualValues(t, f(2), 0)
	assert.EqualValues(t, f(3), 1)
	assert.EqualValues(t, f(4), 2)
	assert.EqualValues(t, f(5), 2)
	assert.EqualValues(t, f(6), 2)
	assert.EqualValues(t, f(7), 2)
}

func TestMove_MoveWin(t *testing.T) {
	var f = func(action MoveAction) BoardIndex {
		move := setMoveAction(0, action)
		return getMoveWin(move)
	}

	// if the first action bit is set, it means that a template or master was taken.
	// regardless of the other bits. 0 = have not won, 1 = win
	assert.EqualValues(t, f(0), 0)
	assert.EqualValues(t, f(1), 1)
}
