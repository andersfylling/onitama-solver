// +build onitama_metrics_infinity

package buildtag

//go:inline
func Onitama_metrics_infinity(cb func()) {
	cb()
}

var _ fOnitama_metrics_infinity = Onitama_metrics_infinity
