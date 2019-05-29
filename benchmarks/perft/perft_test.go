package perft

import (
	"fmt"
	"github.com/andersfylling/onitamago/buildtag"
	"strconv"
	"strings"
	"testing"
	"time"
)
import oni "github.com/andersfylling/onitamago"

func BenchmarkPerft(b *testing.B) {
	cardsSlice := [...][]oni.Card{
		{4627019807588352, 9191917208207360, 18067449945522176, 4609221463113728, 18137612531269632},
		{18190389089402880, 9191917208207360, 4609221463113728, 9112889809960960, 22553526106324992},
		{9060113251827712, 4609221463113728, 4547854970388480, 38316124802121728, 18067449945522176},
		{22641143439163392, 4547854970388480, 9112889809960960, 2305878331024736256, 18137612531269632},
		{38316124802121728, 9060113251827712, 9130344557051904, 18067449945522176, 4627019807588352},
		/*5*/ {4627019807588352, 18137612531269632, 22553319947894784, 9130344557051904, 4547854970388480},
		{38316124802121728, 9060113251827712, 9042727224213504, 22553526106324992, 9130344557051904},
		{9060113251827712, 18190389089402880, 22553319947894784, 22641143439163392, 4609221463113728},
		{22553526106324992, 22641143439163392, 9112889809960960, 4627019807588352, 18137612531269632},
		{9042727224213504, 9191917208207360, 4547854970388480, 18067449945522176, 22553319947894784},
		{38316124802121728, 9060113251827712, 4547854970388480, 22641143439163392, 22553319947894784},
		{22553526106324992, 18067449945522176, 38316124802121728, 9191917208207360, 4609221463113728},
		{22553526106324992, 22641143439163392, 18190389089402880, 9112889809960960, 18137612531269632},
		{4609221463113728, 2305878331024736256, 4627019807588352, 9191917208207360, 18137612531269632},
		{9112889809960960, 9130344557051904, 22553319947894784, 4547854970388480, 18067449945522176},
		{4627019807588352, 9060113251827712, 4609221463113728, 18067449945522176, 22641143439163392},
		{9130344557051904, 18137612531269632, 9060113251827712, 22553319947894784, 18067449945522176},
		{9191917208207360, 4547854970388480, 38316124802121728, 9060113251827712, 2305878331024736256},
		{4627019807588352, 9191917208207360, 9130344557051904, 18190389089402880, 22553526106324992},
		{9060113251827712, 9191917208207360, 18137612531269632, 4609221463113728, 2305878331024736256},
	}
	cards := cardsSlice[2]
	depth := 9

	b.ReportAllocs()

	title := fmt.Sprintf("depth(%d)", depth)
	var d time.Duration
	b.Run(title, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, d = oni.SearchExhaustive(cards, uint64(depth))
		}
	})
	fmt.Println("duration", d)
}

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

	for depth := 8; depth <= 8; depth++ {
		// sequential only!
		title := fmt.Sprintf("depth(%d)", depth)
		b.Run(title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, d := oni.SearchExhaustive(cards, uint64(depth))
				durations = append(durations, d)
			}
		})
		fmt.Println("duration", average(durations))
		durations = durations[:0] // reset length
	}
}

func TestFetchMetrics(t *testing.T) {
	skip := true
	buildtag.Onitama_metrics(func() {
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
		m, _, _ := oni.SearchExhaustive(cards, uint64(depth))
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
