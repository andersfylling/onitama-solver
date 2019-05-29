package onitamago

import "github.com/andersfylling/onitamago/oniconst"

const MaxDepth = oniconst.MaxDepth

type MoveStack interface {
	Pop() Move
	Push(Move)
}

type Stack struct {
	i     int
	stack [HighestNrOfMoves * NrOfPlayerPieces * NrOfPlayerCards * MaxDepth]Move
}

var _ MoveStack = (*Stack)(nil)

func (s *Stack) Size() int {
	return s.i
}

func (s *Stack) Pop() Move {
	s.i--
	return s.stack[s.i]
}

func (s *Stack) Push(m Move) {
	s.stack[s.i] = m
	s.i++
}

func (s *Stack) PushMany(m []Move) {
	copy(s.stack[s.i:s.i+len(m)], m)
	s.i += len(m)
}
