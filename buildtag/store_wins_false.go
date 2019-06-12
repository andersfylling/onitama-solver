// +build !onitama_store_wins

package buildtag

func Onitama_store_wins(cb func()) {}

var _ fOnitama_store_wins = Onitama_store_wins
