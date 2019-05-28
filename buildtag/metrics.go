// +build onitama_metrics

package buildtag

//go:inline
func Onitama_metrics(cb func()) {
	cb()
}

var _ fOnitama_metrics = Onitama_metrics
