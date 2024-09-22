#!/bin/bash

n=20

for t in new; do
    echo "build $t compiler, then press enter"
    read
    tinygo build -o tinyjsonbench.$t
    echo -n "warmup $t: "
    for i in $(seq 3); do ./tinyjsonbench.$t >/dev/null; echo -n "."; done; echo
    echo -n "benchmark $t: "
    for i in $(seq $n); do time ./tinyjsonbench.$t >/dev/null ; echo -n "."; done 2> $t.out; echo
    awk '/real/{sub(/0m/, "", $2); sub(/s/, "", $2); print $2}' $t.out > $t.times
    dist < $t.times
done

tinystat old.times new.times
