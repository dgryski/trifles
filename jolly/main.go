package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	n := flag.Int("n", 1, "games to simulate")

	flag.Parse()

	for i := 0; i < *n; i++ {
		t := oneGame()
		fmt.Println("turns", t)
	}
}

func oneGame() int {
	var turns int
	var square = 1

	var missturns int

	for square < 46 {
		turns++

		var turnsmissed bool
		if missturns > 0 {
			turns += missturns
			missturns = 0
			turnsmissed = true
		}

		switch square {
		case 1, 2, 4, 5, 8, 9, 11, 12, 14, 16, 17, 18, 20, 22, 24, 26, 28, 29, 30, 32, 33, 36, 38, 39, 41, 42, 44, 45:
			square += roll()

		case 3:
			square = 1

		case 6:
			square = 18

		case 7, 19, 27, 35, 43:
			if turnsmissed {
				square += roll()
			} else {
				missturns = 1
			}

		case 10, 23:
			if d := roll(); d == 6 {
				square += d
			}

		case 13:
			if turnsmissed {
				square += roll()
			} else {
				missturns = 2
			}

		case 15:
			square = 18
		case 21:
			square = 24
		case 25:
			square += 4
		case 31:
			square -= 3
		case 34:
			square = 36
		case 37:
			square = 41
		case 40:
			square = 17
		}
	}

	return turns
}

func roll() int {
	return rand.Intn(6) + 1
}
