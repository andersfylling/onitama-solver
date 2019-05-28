package bits

import (
	. "github.com/andersfylling/onitamago"
	"testing"
)

// BenchmarkNLSB
//
//				reference	 copy (two returns)
// name         old time/op  new time/op  delta
// NLSB/NLSB-8  42.5ns ± 2%  42.6ns ± 2%   ~     (p=0.284 n=20+19)
func BenchmarkNLSB(b *testing.B) {
	board := Bitboard(0x885300a409904245)

	var nlsb2 = func(x Bitboard, i BitboardPos) (Bitboard, BitboardPos) {
		x ^= 1 << i
		return x, LSB(x)
	}

	b.Run("NLSB-refence", func(b *testing.B) {
		var boardcp Bitboard
		for i := 0; i < b.N; i++ {
			boardcp = board
			for i := LSB(boardcp); i != 64; i = NLSB(&boardcp, i) {
			}
		}
	})

	b.Run("NLSB-copy", func(b *testing.B) {
		var boardcp Bitboard
		for i := 0; i < b.N; i++ {
			boardcp = board
			for i := LSB(boardcp); i != 64; boardcp, i = nlsb2(boardcp, i) {
			}
		}
	})
}
