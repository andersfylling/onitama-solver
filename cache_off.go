// +build !onitama_cache

package onitamago

const CacheableSubTreeMinHeight = 7

func cacheableDepth(targetDepth, currentDepth uint64) bool {
	return false
}

type cacheInfo struct {
	depth   uint64
	uses    int
	matches int
	key     Key

	StopAfterIndex int
	Ready          bool

	metrics [MaxDepth]DepthMetric
}

func (c *cacheInfo) reset() {}

type onitamaCache struct {
	entries [10000]cacheInfo
	index   int
}

func (c *onitamaCache) match(k Key, currentDepth uint64) ([]DepthMetric, bool, bool) {
	return nil, false, false
}
func (c *onitamaCache) add(k Key, targetDepth, currentDepth uint64, index int)                     {}
func (c *onitamaCache) addMetrics(targetDepth, currentDepth uint64, index int, metric DepthMetric) {}
