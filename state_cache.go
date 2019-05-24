// +build onitama_cache

package onitamago

type previousCacheKeys struct {
	entries [MaxDepth]CacheKey
}

func (st *State) setCacheKey(k CacheKey) {
	// called automatically from the CacheKey.Encode method
	st.previousCacheKeys.entries[st.currentDepth] = k
}

func (st *State) removeLastCacheKey() {}

func (st *State) IsParentCacheKey(k CacheKey) bool {
	// assumption1: the previous move can never be the same.
	//				This also allows the InfinityBranch check to work.
	for i := st.currentDepth - 1; i > 0; i-- {
		if k == st.previousCacheKeys.entries[i] {
			return true
		}
	}

	return false
}
