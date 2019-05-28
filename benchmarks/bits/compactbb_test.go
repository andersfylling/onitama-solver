package bits

import (
	. "github.com/andersfylling/onitamago"
	"testing"
)

// BenchmarkCompactBitboard
//
//                            25bit        32bit
// name                       old time/op  new time/op  delta
// CompactBitboard/compact-8  10.5ns ± 4%  0.2ns ± 0%   -97.71%  (p=0.000 n=20+17)
func BenchmarkCompactBitboard(b *testing.B) {
	board := Bitboard(0x885300a409904245)

	b.Run("25bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MakeCompactBoard(board)
		}
	})

	b.Run("32bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MakeSemiCompactBoard(board)
		}
	})
}
