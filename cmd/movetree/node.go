package movetree

import "github.com/andersfylling/onitamago"

// Root is depth 0, such that children moves are at depth 1
type MoveRoot struct {
	cards [5]onitamago.Card
	moves []*MoveNode
}

type MoveNode struct {
	move  onitamago.Move
	moves []*MoveNode
}
