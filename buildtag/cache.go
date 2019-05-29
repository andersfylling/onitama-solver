// +build onitama_cache

package buildtag

func Onitama_cache(cb func()) {
	cb()
}

var _ fOnitama_cache = Onitama_cache
