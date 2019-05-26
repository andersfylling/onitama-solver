package perft

import (
	"fmt"
	oni "github.com/andersfylling/onitamago"
	"time"
)

func createMetric(depth, activePlayer int, moves []oni.Move) (metric oni.DepthMetric) {
	metric = oni.DepthMetric{
		Depth:          depth,
		ActivePlayer:   activePlayer,
		GeneratedMoves: uint64(len(moves)),
	}

	// TODO: win paths
	var actions oni.Move
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

func Perft(cards []oni.Card, depth int) (metrics []oni.DepthMetric, leafs uint64, moves uint64, duration time.Duration) {
	if depth == 0 {
		return nil, 1, 0, 0
	}
	stack := oni.Stack{}

	st := oni.State{}
	st.CreateGame(cards)
	start := time.Now()

	// caching - use build tag onitama_cache
	cache := onitamaCache{}
	targetDepth := uint64(depth)

	// metrics
	buildtag_onitama_metrics(func() {
		metrics = make([]oni.DepthMetric, depth+1)
	})

	// prepare stack and move indexing
	st.GenerateMoves()
	moves = uint64(st.MovesLen())
	buildtag_onitama_metrics(func() {
		metric := createMetric(1, st.NextPlayer(), st.Moves())
		metrics[1].Increment(&metric)
	})
	if depth == 1 {
		return metrics, moves, moves, time.Now().Sub(start)
	}

	// populate stack with some work
	stack.Push(oni.MoveUndo)
	stack.PushMany(st.Moves())
	var move oni.Move
	var anyWins bool
	for {
		if move = stack.Pop(); move == oni.MoveUndo {
			// finished processing node children will yield a skip move
			// to signify we need to go one depth back up
			for ; move == oni.MoveUndo && stack.Size() > 0; move = stack.Pop() {
				st.UndoMove()
			}
			if stack.Size() == 0 {
				break
			}
		}

		st.ApplyMove(move)
		var skip bool
		buildtag_onitama_cache(targetDepth, st.Depth(), func() { // build tag "onitama_cache"
			var key oni.CacheKey
			key.Encode(&st)

			buildtag_onitama_noinfinity(func() { // build tag "onitama_noinfinity"
				if oni.InfinityBranch(&st) {
					skip = true
				}
			})
			if skip {
				return
			}

			if ms, ready, ok := cache.match(key, st.Depth()); ok {
				//fmt.Println(st.Depth(), key.String())
				if !ready {
					buildtag_onitama_metrics(func() {
						fmt.Println(st.IsParentCacheKey(key), "was not ready..")
					})
					return
				} else {
					//fmt.Println("ready")
				}
				buildtag_onitama_metrics(func() {
					for i := 0; i+int(st.Depth()) <= depth; i++ {
						m := ms[i]
						metrics[int(st.Depth())+i].Increment(&m)
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
		moves += uint64(st.MovesLen())

		buildtag_onitama_metrics(func() { // build tag "onitama_metrics"
			// populate game metrics for the cached entries
			cdepth := int(st.Depth() + 1)
			metric := createMetric(cdepth, st.NextPlayer(), st.Moves())
			metrics[cdepth].Increment(&metric)
			cache.addMetrics(targetDepth, uint64(cdepth), stack.Size(), metric)
		})

		if int(st.Depth()+1) >= depth {
			leafs += uint64(st.MovesLen())
			st.UndoMove()
		} else {
			stack.Push(oni.MoveUndo) // identify a new depth
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

	return metrics, leafs, moves, time.Now().Sub(start)
}
