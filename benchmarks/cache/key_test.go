package cache

import (
	oni "github.com/andersfylling/onitamago"
	"testing"
)

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

	var key oni.CacheKey
	for i := 0; i < b.N; i++ {
		key.Encode(&st)
	}
}
