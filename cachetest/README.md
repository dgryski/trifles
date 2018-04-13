Test different caching algorithms.

Currently implemented:

* random eviction
* lru (external, from https://github.com/golang/groupcache )
* clock
* slru
* s4lru (external, from http://github.com/dgryski/go-s4lru )
* clock-pro (external, from https://github.com/dgryski/go-clockpro )
* arc (external, from https://github.com/dgryski/go-arc )
* tinylfu (external, from https://github.com/dgryski/go-tinylfu )
* tworand -- power of two-random-choices
* randmark -- randomized marking

usage:

    $  ./cachetest -n=5000 -alg=s4lru <trace.txt

