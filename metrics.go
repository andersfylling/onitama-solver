package onitamago

func FetchMetrics(maxDepth int) {

}

type DepthMetric struct {
	StudentsKilled  [2]int
	MastersKilled   [2]int
	TemplesTaken    [2]int
	NonViolentMoves int

	depth          int
	ActivePlayer   int
	GeneratedMoves int
}
