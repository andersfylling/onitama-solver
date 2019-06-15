// +build onitama_store_wins

package buildtag

//go:inline
func Onitama_store_wins(cb func()) {
	cb()
}

var _ fOnitama_store_wins = Onitama_store_wins
