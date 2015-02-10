package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type roll int

const (
	sideHome roll = iota
	sideFirst
	sideSecond
	sideEndOfLine
)

func rollDie() roll {

	switch rand.Intn(6) {
	case 0, 1:
		return sideHome
	case 2:
		return sideFirst
	case 3:
		return sideSecond
	case 4, 5:
		return sideEndOfLine
	}

	panic("unreached")
}

type queue []int

func (q *queue) moveToEnd(pos int) {
	e := (*q)[pos]
	copy((*q)[pos:], (*q)[pos+1:])
	(*q)[len(*q)-1] = e
}

func (q *queue) moveHome(player int) {
	for i, v := range *q {
		if v == player {
			q.moveToEnd(i)
			*q = (*q)[:len(*q)-1]
			return
		}
	}
}

func playGame(players int) int {
	var q queue

	for i := 0; i < 5; i++ {
		q = append(q, 0)
		for p := players - 1; p > 0; p-- {
			q = append(q, p)
		}
	}

	scores := make([]int, players)

	var player int

	for {
		switch rollDie() {

		case sideHome:
			q.moveHome(player)
			scores[player]++
			if scores[player] == 5 {
				return player
			}

		case sideFirst:
			p := q[0]
			q.moveHome(p)
			scores[p]++
			if scores[p] == 5 {
				return p
			}

		case sideSecond:
			p := q[1]
			q.moveHome(p)
			scores[p]++
			if scores[p] == 5 {
				return p
			}

		case sideEndOfLine:
			q.moveToEnd(0)
		}

		player = (player + 1) % players
	}
}

func main() {

	players := flag.Int("p", 4, "players")
	runs := flag.Int("r", 1000, "games to simulate")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	scores := make([]int, *players)

	for i := 0; i < *runs; i++ {
		scores[playGame(*players)]++
	}

	fmt.Println(scores)
}
