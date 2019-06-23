package onitamago

// first node is considered root. You can think of it as it requires no moves to reach it, since it is the starting point.
type MoveNode struct {
	Instances uint64
	Depth     uint8
	Paths     map[Move]*MoveNode
	Parent    *MoveNode
}

func ConvertPathsToTree(paths [][]Move) (tree *MoveNode, nrOfNodes uint) {
	tree = &MoveNode{
		Instances: 1,
		Paths:     map[Move]*MoveNode{},
	}

	for p := range paths {
		node := tree

		var exists bool
		for i := range paths[p] {
			m := paths[p][i]
			if _, exists = node.Paths[m]; !exists {
				node.Paths[m] = &MoveNode{
					Depth:     node.Depth + 1,
					Instances: 1,
					Paths:     map[Move]*MoveNode{},
					Parent:    node,
				}
				nrOfNodes++
			} else {
				node.Paths[m].Instances++
			}
			node = node.Paths[m]
		}
	}

	return tree, nrOfNodes
}
