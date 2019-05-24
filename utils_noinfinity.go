// +build onitama_noinfinity,onitama_cache

package onitamago

// InfinityBranch assumes that if a node in relation to the current node has the
// same cache key (unique identifier), the game is in a loop and the current node
// is considered the start of an infinity branch.
func InfinityBranch(st *State) bool {
	current := st.previousCacheKeys.entries[st.currentDepth]
	return st.IsParentCacheKey(current)
}
