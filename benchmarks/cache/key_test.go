package cache

import (
	oni "github.com/andersfylling/onitamago"
	"testing"
)

// ➜ benchstat uint64-cardIDs-loop.txt uint64-MakeSemiCompactBoard.txt
//
//name               old time/op    new time/op    delta
//CacheKey_Encode-8  59.8ns ± 2%    32.3ns ± 2%    -45.92%  (p=0.000 n=20+20)
//
//name               old alloc/op   new alloc/op   delta
//CacheKey_Encode-8  0.00B          0.00B          ~     (all equal)
//
//name               old allocs/op  new allocs/op  delta
//CacheKey_Encode-8  0.00           0.00           ~     (all equal)
func BenchmarkCacheKey_Encode(b *testing.B) {
	b.ReportAllocs()
	st := oni.NewState()
	st.CreateGame(nil)

	moves := []oni.Move{
		50005,
		52003,
		50533,
		51419,
	}
	for _, move := range moves {
		st.ApplyMove(move)
	}

	var key oni.Key
	for i := 0; i < b.N; i++ {
		key.Encode(&st)
	}
}
