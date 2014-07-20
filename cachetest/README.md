Test different caching algorithms.

Currently implemented:
    - random eviction
    - lru (external, from github.com/golang/groupcache/lru)
    - lfu (external, from https://github.com/calmh/lfucache )
    - clock
    - slru
    - s4lru (described in http://www.cs.cornell.edu/~qhuang/papers/sosp_fbanalysis.pdf )


usage:

    $  ./cachetest -n=5000 -alg=s4lru <trace.txt


