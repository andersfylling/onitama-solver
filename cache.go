// +build onitama_cache

package onitamago

// The performance elbow is at Depth 7 or 8.
// However, the delta from 7 to 8 is around -560%.
// As Depth 7 is from 250-700ms, this seems like the most
// significant jump in Depth time. Therefore we do not
// cache sub-trees with a height less than 7, as this
// requires extra loop time and slows down the application
// more than I've found the cache to help.
//
//  Depth/Seconds
//
//      |     |     |     |  X  |
//  10s +-----+-----+-----+-X---+
//      |     |     |     |X    |
//      |     |     |    XX     |
//   1s +-----+-----+--XX-+-----+
//      |     |   XXXXX   |     |
//      XXXXXXXXXX  |     |     |
//      |     |     |     |     |
//   0s +-----+-----+-----+-----+
//      5     6     7     8     9
const CacheableSubTreeMinHeight = 7

func cacheableDepth(targetDepth, currentDepth uint64) bool {
	return targetDepth-currentDepth > 6
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

func (c *cacheInfo) reset() {
	c.uses = 0
	c.matches = 0
	c.Ready = false
	for i := range c.metrics {
		c.metrics[i].Reset()
	}
}

type onitamaCache struct {
	entries [10000]cacheInfo
	index   int
}

func (c *onitamaCache) match(k Key, currentDepth uint64) ([]DepthMetric, bool, bool) {
	for i := 0; i < c.index; i++ {
		if k != c.entries[i].key {
			continue
		}
		c.entries[i].matches++

		// only works when the subtree is as high or higher
		if currentDepth < c.entries[i].depth {
			continue
		}

		index := i
		if i > 0 && c.entries[i-1].matches < c.entries[i].matches {
			// move upwards to avoid being overwritten when the size limit is hit
			c.entries[i-1], c.entries[i] = c.entries[i], c.entries[i-1]
			index = i - 1
		}

		return c.entries[index].metrics[:], c.entries[index].Ready, true
	}
	return nil, false, false
}

func (c *onitamaCache) add(k Key, targetDepth, currentDepth uint64, index int) {
	if c.index == len(c.entries) {
		c.index = int(float32(len(c.entries)) * 0.8)
	}

	c.entries[c.index].reset()
	c.entries[c.index].key = k
	c.entries[c.index].depth = currentDepth
	c.entries[c.index].StopAfterIndex = index
	c.index++
}

func (c *onitamaCache) addMetrics(targetDepth, currentDepth uint64, index int, metric DepthMetric) {
	for i := c.index - 1; i >= 0; i-- {
		if c.entries[i].Ready {
			continue
		}

		if c.entries[i].StopAfterIndex > index {
			c.entries[i].Ready = true
			continue
		}

		mdepth := int(currentDepth - c.entries[i].depth)
		c.entries[i].metrics[mdepth].Increment(&metric)
	}
}
