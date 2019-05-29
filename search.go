package onitamago

import (
	"fmt"
	"github.com/andersfylling/onitamago/buildtag"
	"time"
)

func createMetric(depth, activePlayer int, moves []Move) (metric DepthMetric) {
	metric = DepthMetric{
		Depth:          depth,
		ActivePlayer:   activePlayer,
		GeneratedMoves: uint64(len(moves)),
	}

	// TODO: win paths
	var actions Move
	for _, move := range moves {
		actions = 0
		actions = (move >> 12) & 7

		metric.MastersKilled += uint64(actions & 1)
		metric.TemplesTaken += uint64((actions & (actions >> 2)) & 1)

		if actions == 0 || actions == 2 {
			metric.StudentsKilled += 1
		}

		if actions == 4 || actions == 6 {
			// the pieces attacked none nor the temple
			metric.NonViolentMoves += 1
		}
	}

	return metric
}

// SearchExhaustive uses depth first search to goes through the
// entire game tree generated from the card configuration until
// the target depth is reached.
// However, when a parent generates children moves that causes a win,
// that parent branch is no longer explored as it's assumed the player
// prioritizes a win instead of continuing the game.
func SearchExhaustive(cards []Card, targetDepth uint64) (metrics []DepthMetric, winPaths [][]Move, duration time.Duration) {
	if targetDepth == 0 {
		return nil, nil, 0
	}
	stack := Stack{}

	st := State{}
	st.CreateGame(cards)

	// caching - use build tag onitama_cache
	cache := onitamaCache{}

	// metrics
	buildtag.Onitama_metrics(func() {
		metrics = make([]DepthMetric, targetDepth+1)
	})

	// prepare stack and move indexing
	start := time.Now()
	st.GenerateMoves()
	buildtag.Onitama_metrics(func() {
		metric := createMetric(1, st.NextPlayer(), st.Moves())
		metrics[1].Increment(&metric)
	})
	if targetDepth == 1 {
		return metrics, nil, time.Now().Sub(start)
	}

	// populate stack with some work
	stack.Push(MoveUndo)
	stack.PushMany(st.Moves())
	var move Move
	var anyWins bool
	var key Key
	for {
		if move = stack.Pop(); move == MoveUndo {
			// finished processing node children will yield a skip move
			// to signify we need to go one depth back up
			for ; move == MoveUndo && stack.Size() > 0; move = stack.Pop() {
				st.UndoMove()
			}
			if stack.Size() == 0 {
				break
			}
		}

		st.ApplyMove(move)
		var skip bool
		buildtag.Onitama_cache(func() {
			// TODO: look ups can still take place
			// but new cache elements can not...
			if targetDepth-st.Depth() < CacheableSubTreeMinHeight {
				return
			}
			key.Encode(&st)
			//buildtag_onitama_noinfinity(func() {
			//	if oni.InfinityBranch(&st) {
			//		skip = true
			//	}
			//})
			//if skip {
			//	return
			//}

			if ms, ready, ok := cache.match(key, st.Depth()); ok {
				//fmt.Println(st.Depth(), key.String())
				if !ready {
					buildtag.Onitama_metrics(func() {
						fmt.Println(st.IsParentCacheKey(key), "was not ready..")
					})
					return
				} else {
					//fmt.Println("ready")
				}
				buildtag.Onitama_metrics(func() {
					for i := Number(0); i+st.Depth() <= targetDepth; i++ {
						m := ms[i]
						metrics[st.Depth()+i].Increment(&m)
					}
				})
				skip = true
			} else {
				cache.add(key, targetDepth, st.Depth(), stack.Size())
			}
		})
		if skip {
			st.UndoMove()
			continue
		}
		st.GenerateMoves()

		buildtag.Onitama_metrics(func() { // build tag "onitama_metrics"
			// populate game metrics for the cached entries
			cdepth := int(st.Depth() + 1)
			metric := createMetric(cdepth, st.NextPlayer(), st.Moves())
			metrics[cdepth].Increment(&metric)

			buildtag.Onitama_cache(func() {
				cache.addMetrics(targetDepth, uint64(cdepth), stack.Size(), metric)
			})
		})

		if st.Depth()+1 >= targetDepth {
			st.UndoMove()
		} else {
			stack.Push(MoveUndo) // identify a new depth
			anyWins = false
			for i := range st.Moves() {
				if (st.Moves()[i] & (1 << 12)) > 0 {
					anyWins = true
					break
				}
			}
			if !anyWins {
				stack.PushMany(st.Moves())
			}
		}
	}

	return metrics, winPaths, time.Now().Sub(start)
}
