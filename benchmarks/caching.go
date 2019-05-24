// +build !onitama_cache

package main

import oni "github.com/andersfylling/onitamago"

type cacheInfo struct {
	depth          uint64
	uses           int
	matches        int
	key            oni.CacheKey
	StopAfterIndex int
}

func (c *cacheInfo) reset() {}

type onitamaCache struct {
	entries [100000]cacheInfo
	index   int
}

func (c *onitamaCache) match(k oni.CacheKey, currentDepth uint64) ([]oni.DepthMetric, bool) {
	return nil, false
}

func (c *onitamaCache) add(k oni.CacheKey, targetDepth, depth uint64, index int) {}
func (c *onitamaCache) addMetrics(depth int, index int, metric oni.DepthMetric)  {}

func doCache(targetDepth, currentDepth uint64, cb func()) {}
