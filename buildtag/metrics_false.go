// +build !onitama_metrics

package buildtag

func Onitama_metrics(cb func()) {}

var _ fOnitama_metrics = Onitama_metrics
