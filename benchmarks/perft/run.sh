#!/bin/sh

# without metrics
go test -bench=BenchmarkPERFT -count 10 >02-results-vanilla-perft.txt
go test -tags=onitama_cache -bench=BenchmarkPERFT -count 10 >02-results-cached-perft.txt

# with metrics
go test -tags=metrics -bench=BenchmarkPERFT -count 10 >02-results-vanilla-metrics-perft.txt
go test -tags=onitama_cache,metrics -bench=BenchmarkPERFT -count 10 >02-results-cached-metrics-perft.txt

# TODO:

benchstat vanilla-perft.txt cached-perft.txt >01-results-perft-diff.txt