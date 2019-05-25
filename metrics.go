package onitamago

func FetchMetrics(maxDepth int) {}

type DepthMetric struct {
	StudentsKilled  uint64
	MastersKilled   uint64
	TemplesTaken    uint64
	NonViolentMoves uint64

	Depth          int
	ActivePlayer   int
	GeneratedMoves uint64
}

func (m *DepthMetric) Reset() {
	m.StudentsKilled = 0
	m.MastersKilled = 0
	m.TemplesTaken = 0
	m.NonViolentMoves = 0
	m.Depth = 0
	m.ActivePlayer = 0
	m.GeneratedMoves = 0
}
func (m *DepthMetric) Increment(metric *DepthMetric) {
	m.StudentsKilled += metric.StudentsKilled
	m.MastersKilled += metric.MastersKilled
	m.TemplesTaken += metric.TemplesTaken
	m.GeneratedMoves += metric.GeneratedMoves

	// TODO: do not set these for every update..
	m.Depth = metric.Depth
	m.ActivePlayer = metric.ActivePlayer
}
