package perft

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)
import oni "github.com/andersfylling/onitamago"

func BenchmarkPERFT(b *testing.B) {
	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	average := func(ds []time.Duration) (avg time.Duration) {
		for _, d := range ds {
			avg += d
		}
		avg = avg / time.Duration(len(ds))

		return
	}
	durations := []time.Duration{}

	for depth := 10; depth <= 10; depth++ {
		// sequential only!
		title := fmt.Sprintf("depth(%d)", depth)
		b.Run(title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, d := Perft(cards, depth)
				durations = append(durations, d)
			}
		})
		fmt.Println("duration", average(durations))
		durations = durations[:0] // reset length
	}
}

func TestFetchMetrics(t *testing.T) {
	skip := true
	buildtag_onitama_metrics(func() {
		skip = false
	})
	if skip {
		return
	}

	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	var metrics [][]oni.DepthMetric
	base := 1
	for depth := base; depth <= 9; depth++ {
		// sequential only!
		m, _, _, _ := Perft(cards, depth)
		metrics = append(metrics, m)
		var moves uint64
		if len(m) > 0 {
			moves = m[len(m)-1].GeneratedMoves
		}
		fmt.Println("depth", depth, "=", moves)
	}

	b := strings.Builder{}
	b.Write([]byte(" = [][]int{\n"))
	for i := range metrics {
		b.Write([]byte("\t{ // Perft(" + strconv.Itoa(i+base) + ") \n"))
		for j := range metrics[i] {
			b.WriteString("\t\t" + strconv.FormatUint(metrics[i][j].GeneratedMoves, 10) + ",\n")
		}
		b.Write([]byte("\t},\n"))
	}
	b.Write([]byte("}\n"))

	fmt.Print(b.String())
}

func TestPERFTCacheAcc(t *testing.T) {
	skip := true
	var showCachePrune bool
	buildtag_onitama_metrics(func() {
		skip = false
		buildtag_onitama_cache(100, 0, func() {
			showCachePrune = true
		})
	})
	if skip {
		return
	}

	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	for depth := 9; depth <= 9; depth++ {
		// sequential only!
		m, leafs, _, _ := Perft(cards, depth)
		mleafs := m[len(m)-1].GeneratedMoves
		fmt.Println(depth, leafs)

		if showCachePrune {
			fmt.Println("cache pruned", mleafs-leafs)
		}
	}
}

var vanilla_metrics = [][]int{
	{ // Perft(1)
		0,
		11,
	},
	{ // Perft(2)
		0,
		11,
		88,
	},
	{ // Perft(3)
		0,
		11,
		88,
		992,
	},
	{ // Perft(4)
		0,
		11,
		88,
		992,
		11159,
	},
	{ // Perft(5)
		0,
		11,
		88,
		992,
		11159,
		126293,
	},
	{ // Perft(6)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
	},
	{ // Perft(7)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
	},
	{ // Perft(8)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
		133192301,
	},
	{ // Perft(9)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
		133192301,
		1438696527,
	},
}
var cache_metrics = [][]int{
	{ // Perft(1)
		0,
		11,
	},
	{ // Perft(2)
		0,
		11,
		88,
	},
	{ // Perft(3)
		0,
		11,
		88,
		992,
	},
	{ // Perft(4)
		0,
		11,
		88,
		992,
		11159,
	},
	{ // Perft(5)
		0,
		11,
		88,
		992,
		11159,
		126293,
	},
	{ // Perft(6)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
	},
	{ // Perft(7)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
	},
	{ // Perft(8)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
		133192301,
	},
	{ // Perft(9)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
		133192301,
		1438696527,
	},
}

func TestCompareMetricsToCache(t *testing.T) {
	if len(vanilla_metrics) != len(cache_metrics) {
		t.Fatal("cache and vanilla metrics should be the same length")
	}

	for i := range vanilla_metrics {
		if len(vanilla_metrics[i]) != len(cache_metrics[i]) {
			t.Fatal("cache and vanilla depth metrics should be the same length")
		}

		for j := range vanilla_metrics[i] {
			a, b := vanilla_metrics[i][j], cache_metrics[i][j]
			if a != b {
				t.Errorf("cache and vanilla for Perft(%d) at depth %d was different. Cache got %d, expected %d. diff %d", i+1, j, b, a, a-b)
			}
		}
	}
}
