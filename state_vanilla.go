// +build !onitama_cache

package onitamago

type previousCacheKeys struct{}

func (st *State) setCacheKey(k CacheKey) {}

func (st *State) removeLastCacheKey() {}

func (st *State) IsParentCacheKey(k CacheKey) bool { return false }
