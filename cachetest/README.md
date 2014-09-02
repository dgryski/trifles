Test different caching algorithms.

Currently implemented:

* random eviction
* lru (external, from https://github.com/golang/groupcache )
* lfu (external, from https://github.com/calmh/lfucache )
* clock
* slru
* s4lru (external, from http://github.com/dgryski/go-s4lru )
* clock-pro (external, from https://github.com/dgryski/go-clockpro )
* arc (external, from https://github.com/dgryski/go-arc )

usage:

    $  ./cachetest -n=5000 -alg=s4lru <trace.txt

