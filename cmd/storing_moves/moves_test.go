package storing_moves

import (
	"github.com/andersfylling/onitamago"
	ikea "github.com/ikkerens/ikeapack"
	"math/rand"
	"testing"
)

var moves Moves

func init() {
	moves.Init()

	for i := range moves.moves {
		for j := range moves.moves[i] {
			v := rand.Uint32()
			//v := ^uint16(0)
			moves.moves[i][j] = onitamago.Move(v)
		}
	}
}

func TestPacking(t *testing.T) {
	moves.compress()

	moves2 := Moves{}
	moves2.Init()
	err := ikea.Unpack(moves.buffer, &moves2)
	if err != nil {
		panic(err)
	}

	for i := range moves.moves {
		for j := range moves.moves[i] {
			if moves.moves[i][j] != moves2.moves[i][j] {
				panic("diff values")
			}
		}
	}
}

func BenchmarkCompression(b *testing.B) {
	for i := 0; i < b.N; i++ {
		moves.compress()
	}
}
