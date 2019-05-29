// +build !onitama_cache

package onitamago

type previousCacheKeys struct{}

func (st *State) setCacheKey(k Key) {}

func (st *State) removeLastCacheKey() {}

func (st *State) IsParentCacheKey(k Key) bool { return false }
