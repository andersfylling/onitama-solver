package onitamago

import (
	"fmt"
	"github.com/andersfylling/onitamago/buildtag"
	"github.com/andersfylling/onitamago/oniconst"

	"time"
)

func nextMove(stack *Stack, st *State) (move Move, ok bool) {
	if move = stack.Pop(); move == MoveUndo {
		// finished processing node children will yield a skip move
		// to signify we need to go one Depth back up
		for ; move == MoveUndo && stack.Size() > 0; move = stack.Pop() {
			st.UndoMove()
		}

		// the last card is a MoveUndo, so if the stack size is 0
		// we can actually assume no moves are left
		if stack.Size() == 0 {
			return 0, false
		}
	}

	return move, true
}

func createMetric(depth, activePlayer uint8, moves []Move) (metric DepthMetric) {
	metric = DepthMetric{
		Depth:          depth,
		ActivePlayer:   activePlayer,
		GeneratedMoves: uint64(len(moves)),
	}

	// TODO: win Paths
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

type moveNode struct {
	Instances uint
	Depth     uint
	Paths     map[Move]*moveNode
}

// SearchExhaustive uses Depth first search to goes through the
// entire game tree generated from the card configuration until
// the target Depth is reached.
// However, when a parent generates children moves that causes a win,
// that parent branch is no longer explored as it's assumed the player
// prioritizes a win instead of continuing the game.
//
// Build tags
// - onitama_store_wins will populate winPaths
func SearchExhaustive(cards []Card, targetDepth uint64) (metrics []DepthMetric, winPaths *moveNode /*infinityPaths [][]Move,*/, duration time.Duration) {
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

	winPaths = &moveNode{
		Instances: 1,
		Paths:     map[Move]*moveNode{},
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
			// to signify we need to go one Depth back up
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
		currentDepth := st.Depth() + 1

		buildtag.Onitama_metrics(func() { // build tag "onitama_metrics"
			// populate game metrics for the cached entries
			cdepth := uint8(currentDepth)
			metric := createMetric(cdepth, st.NextPlayer(), st.Moves())
			metrics[cdepth].Increment(&metric)

			buildtag.Onitama_cache(func() {
				cache.addMetrics(targetDepth, uint64(cdepth), stack.Size(), metric)
			})
		})

		// win Paths
		for i := range st.Moves() {
			if st.generatedMoves[i].Win() {
				anyWins = true
				if !oniconst.StoreWins {
					break
				}

				node := winPaths
				node.Instances++

				var exists bool
				for m := uint64(1); m <= st.Depth(); m++ {
					if st.previousMoves[m] == 0 {
						break
					}
					if _, exists = node.Paths[st.previousMoves[m]]; !exists {
						node.Paths[st.previousMoves[m]] = &moveNode{
							Depth:     node.Depth + 1,
							Instances: 1,
							Paths:     map[Move]*moveNode{},
						}
					} else {
						node.Paths[st.previousMoves[m]].Instances++
					}
					node = node.Paths[st.previousMoves[m]]
				}
			}
		}

		if currentDepth >= targetDepth {
			st.UndoMove()
		} else {
			stack.Push(MoveUndo) // identify a new Depth
			if !anyWins {
				stack.PushMany(st.Moves())
			}
		}
		anyWins = false
	}

	return metrics, winPaths, time.Now().Sub(start)
}

// SearchForTempleWins Stores paths of moves whenever a win by temple is achieved
func SearchForTempleWins(cards []Card, targetDepth uint64) (metrics []DepthMetric, paths [][]Move, duration time.Duration) {
	if targetDepth == 0 {
		return nil, nil, 0
	}
	stack := Stack{}
	st := State{}
	st.CreateGame(cards)

	// metrics
	buildtag.Onitama_metrics(func() {
		metrics = make([]DepthMetric, targetDepth+1)
	})

	paths = make([][]Move, 0, 1000*(targetDepth*targetDepth))

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
	for {
		if move = stack.Pop(); move == MoveUndo {
			// finished processing node children will yield a skip move
			// to signify we need to go one Depth back up
			for ; move == MoveUndo && stack.Size() > 0; move = stack.Pop() {
				st.UndoMove()
			}
			if stack.Size() == 0 {
				break
			}
		}

		st.ApplyMove(move)
		st.GenerateMoves()
		currentDepth := st.Depth() + 1

		buildtag.Onitama_metrics(func() { // build tag "onitama_metrics"
			// populate game metrics for the cached entries
			cdepth := uint8(currentDepth)
			metric := createMetric(cdepth, st.NextPlayer(), st.Moves())
			metrics[cdepth].Increment(&metric)
		})

		// win Paths
		for i := range st.Moves() {
			if st.generatedMoves[i].Win() {
				anyWins = true

				if !st.generatedMoves[i].WinByTemple() {
					continue
				}

				path := make([]Move, currentDepth)
				for m := uint64(1); m <= st.Depth(); m++ {
					if st.previousMoves[m] == 0 {
						break
					}
					path[m-1] = st.previousMoves[m]
				}
				path[currentDepth-1] = st.Moves()[i]
				paths = append(paths, path)
			}
		}

		if currentDepth >= targetDepth {
			st.UndoMove()
		} else {
			stack.Push(MoveUndo) // identify a new Depth
			if !anyWins {
				stack.PushMany(st.Moves())
			}
		}
		anyWins = false
	}

	return metrics, paths, time.Now().Sub(start)
}

func SearchExhaustiveForForcedWins(cards []Card, targetDepth uint64) (metrics []DepthMetric, paths [][]Move, duration time.Duration) {
	metrics, paths, duration = SearchForTempleWins(cards, targetDepth)
	paths = FilterForcedMoves(paths)
	return
}

// SearchExhaustiveInfinityPaths looks through a card configuration to detect if a infinity branch exists
// within the given search space. Note that pruning for wins are turned off so this has a larger search space
// then the SearchExhaustive function.
// Required build tags:
// - onitama_cache: to store the parent states as unique keys for look ups
func SearchExhaustiveInfinityPaths(cards []Card, targetDepth uint64, limitHits int) (paths [][]Move, duration time.Duration) {
	// var hasBuildTag bool
	// buildtag.Onitama_metrics_infinity(func() {
	// 	hasBuildTag = true
	// })
	// if !hasBuildTag {
	// 	panic("missing build tag onitama_metrics_infinity")
	// }
	if targetDepth < 4 {
		return nil, 0
	}
	var start time.Time
	stack := Stack{}
	st := State{}
	st.CreateGame(cards)

	endTimer := func() time.Duration {
		return time.Now().Sub(start)
	}

	// prepare stack and move indexing
	start = time.Now()
	st.GenerateMoves()

	// populate stack with initial moves
	stack.Push(MoveUndo)
	stack.PushMany(st.Moves())
	var move Move
	var key Key
	for {
		if move = stack.Pop(); move == MoveUndo {
			// finished processing node children will yield a skip move
			// to signify we need to go one Depth back up
			for ; move == MoveUndo && stack.Size() > 0; move = stack.Pop() {
				st.UndoMove()
			}
			if stack.Size() == 0 {
				break
			}
		}
		st.ApplyMove(move)

		// look for infinity games
		key.Encode(&st)
		if InfinityBranch(&st) {
			paths = append(paths, st.MoveHistory())
			if len(paths) == limitHits {
				return paths, endTimer()
			}
		}

		// generate the children. But don't prune on wins
		st.GenerateMoves()
		if st.Depth() >= targetDepth {
			st.UndoMove()
		} else {
			stack.Push(MoveUndo) // identify a new Depth
			stack.PushMany(st.Moves())
		}
	}

	return paths, endTimer()
}

func FilterForcedMoves(paths [][]Move) (forced [][]Move) {
	// TODO
	return paths
}