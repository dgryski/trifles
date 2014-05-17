package main

import (
	"log"
	"math"

	"github.com/dgryski/trifles/inthash"
	"github.com/willf/bitset"
)

func main() {

	b := bitset.New(math.MaxUint32)

	start := uint32(0)
	i32 := start
	steps := 0

	for {
		if b.Test(uint(i32)) && i32 != start {
			log.Println("mismatch for start=", start, "at i32=", i32)
			break
		}
		b.Set(uint(i32))
		i32 = inthash.Jenkins32(i32)

		if i32 == start {
			log.Println("cycle for", start, "len=", steps)
			steps = 0
			for start < math.MaxUint32 && b.Test(uint(start)) {
				start++
			}
			if start == math.MaxUint32 {
				log.Println("all done")
				break
			}
			log.Println("start=", start)
			i32 = start
		}

		steps++
		if steps%(32*1024*1024) == 0 {
			log.Println("start=", start, "steps=", steps)
		}
	}

}
