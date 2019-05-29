// +build !onitama_cache

package buildtag

func Onitama_cache(cb func()) {}

var _ fOnitama_cache = Onitama_cache
