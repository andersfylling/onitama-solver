package main

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

	for depth := 1; depth <= 10; depth++ {
		// sequential only!
		title := fmt.Sprintf("depth(%d)", depth)
		b.Run(title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _, d := perft(cards, depth)
				durations = append(durations, d)
			}
		})
		fmt.Println("duration", average(durations))
		durations = durations[:0] // reset length
	}
}

func TestFetchMetrics(t *testing.T) {
	//return
	cards := []oni.Card{
		oni.Frog, oni.Eel,
		oni.Dragon, oni.Crab,
		oni.Tiger,
	}

	var metrics [][]oni.DepthMetric
	base := 1
	for depth := base; depth <= 9; depth++ {
		// sequential only!
		m, _, _, _ := perft(cards, depth)
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
		b.Write([]byte("\t{ // perft(" + strconv.Itoa(i+base) + ") \n"))
		for j := range metrics[i] {
			b.WriteString("\t\t" + strconv.FormatUint(metrics[i][j].GeneratedMoves, 10) + ",\n")
		}
		b.Write([]byte("\t},\n"))
	}
	b.Write([]byte("}\n"))

	fmt.Print(b.String())
}

var vanilla_metrics = [][]int{
	{ // perft(1)
		0,
		11,
	},
	{ // perft(2)
		0,
		11,
		88,
	},
	{ // perft(3)
		0,
		11,
		88,
		992,
	},
	{ // perft(4)
		0,
		11,
		88,
		992,
		11159,
	},
	{ // perft(5)
		0,
		11,
		88,
		992,
		11159,
		126293,
	},
	{ // perft(6)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
	},
	{ // perft(7)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
	},
	{ // perft(8)
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
	{ // perft(9)
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
	{ // perft(1)
		0,
		11,
	},
	{ // perft(2)
		0,
		11,
		88,
	},
	{ // perft(3)
		0,
		11,
		88,
		992,
	},
	{ // perft(4)
		0,
		11,
		88,
		992,
		11159,
	},
	{ // perft(5)
		0,
		11,
		88,
		992,
		11159,
		126293,
	},
	{ // perft(6)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
	},
	{ // perft(7)
		0,
		11,
		88,
		992,
		11159,
		126293,
		1207255,
		13172583,
	},
	{ // perft(8)
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
	{ // perft(9)
		0,
		11,
		88,
		992,
		11154,
		127878,
		1192006,
		13007030,
		130994968,
		1403905796,
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
				t.Errorf("cache and vanilla for perft(%d) at depth %d was different. Cache got %d, expected %d. diff %d", i+1, j, b, a, a-b)
			}
		}
	}
}
