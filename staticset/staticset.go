package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

func search(s []uint32, nbits uint) uint32 {

	slots := 1 << nbits

	log.Printf("searching for multiplier to fit %d keys in %d slots\n", len(s), slots)

	if len(s) > slots {
		return 0
	}

	dst := make([]uint32, slots)

	for i := 1; i < math.MaxUint32; i++ {

		if i&((1<<20)-1) == 0 {
			log.Println(i)
		}

		// try with multiplier `i`
		for d := range dst {
			dst[d] = 0
		}

		var placed int
		for _, v := range s {
			h := (v * uint32(i) >> (32 - nbits))
			if dst[h] != 0 {
				break
			}
			dst[h] = v
			placed++
		}
		if placed == len(s) {
			return uint32(i)
		}
	}
	return 0
}

func u32(s string) uint32 {
	return binary.BigEndian.Uint32([]byte(s))
}

func main() {

	keys := []string{
		"SIP/",
		"INVI",
		"ACK ",
		"CANC",
		"BYE ",
		"PRAC",
		"REGI",
		"OPTI",
		"INFO",
		"UPDA",
		"SUBS",
		"NOTI",
		"MESS",
		"REFE",
		"PUBL",
	}

	var k32s []uint32

	for _, s := range keys {
		k32s = append(k32s, u32(s))
	}

	m := search(k32s, 4)

	if m == 0 {
		log.Fatalf("unable to find multiplier")
	}

	fmt.Println("m=", m)
	for _, v := range k32s {
		fmt.Println((v * m) >> 28)
	}
}
