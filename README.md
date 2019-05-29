# Onitamago
Go package to calculate onitama moves N steps ahead. Features:
 - Caching
 - Terminates branch exploration on wins
 - Metrics
 - Infinity branch detection
 - Heatmap generation for cards and pieces on the board
 - Iterative depth first search
 
Also holds benchmarks for several concepts. Also benchmarks of exhaustive search done for with and without caching, with and without metrics fetching. 

Missing heuristics to make educated guesses about best moves, alpha-beta pruning, min-max/nega-max.

## Build tags
 - onitama_cache
 - onitama_metrics
 - onitama_noinfinity (currently depends on onitama_cache)