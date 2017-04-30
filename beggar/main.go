package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
)

func shuffle(b []byte) []byte {
	d := make([]byte, len(b))
	p := rand.Perm(len(b))
	for i, v := range p {
		d[i] = b[v]
	}
	return d
}

type deck []byte

type table struct {
	turn    int
	players []deck
	stack   deck

	turns  int
	tricks int
	swap   int
	tlens  []int
	cards  int
}

func (t *table) deal() {
	t.cards++
	t.stack = append(t.stack, t.players[t.turn][0])
	t.players[t.turn] = t.players[t.turn][1:]
}

func (t *table) done() bool {
	return len(t.players[0]) == 0 || len(t.players[1]) == 0
}

func (t *table) swapPlayer() {
	t.turn = 1 - t.turn
	t.turns++
	fmt.Printf("t.turns = %+v\n", t.turns)
}

var needToPlay = [256]int{
	'A': 4,
	'J': 1,
	'K': 3,
	'Q': 2,
}

func (t *table) trick(remain int) {

	// player `turn` just dealt a card
	// so player `1 - turn` has `remain` chances to cover it

	t.tricks++
	t.swapPlayer()

	for ; remain > 0 && !t.done(); remain-- {
		t.deal()

		need := needToPlay[t.stack[len(t.stack)-1]]
		if need != 0 {
			// remain count is +1 because of -- at the bottom
			remain = need + 1
			t.swapPlayer()
		}
	}

	if remain == 0 {
		t.swapPlayer()

		// player `turn` didn't cover in time
		// so player 1-turn gets all the cards and gets to deal next
		t.players[t.turn] = append(t.players[t.turn], t.stack...)
		//	t.tlens = append(t.tlens, len(t.stack))
		fmt.Printf("t.players = %+v\n", t.players)
		t.stack = nil
	}
}

func (t *table) play() {

	for len(t.players[0]) > 0 && len(t.players[1]) > 0 {

		t.deal()

		need := needToPlay[t.stack[len(t.stack)-1]]
		if need != 0 {
			t.trick(need)
		} else {
			t.swapPlayer()
		}
	}
}

func oneGame() (d deck, turns int, tricks int) {

	for i := 0; i < 4; i++ {
		for j := 0; j < 9; j++ {
			d = append(d, '.')
		}

		d = append(d, 'K', 'Q', 'J', 'A')
	}

	d = shuffle(d)

	orig := make(deck, len(d))
	copy(orig, d)

	t := table{
		players: []deck{d[:26], d[26:]},
	}

	t.play()

	return orig, t.turns, t.tricks
}

type game struct {
	d      deck
	turns  int
	tricks int
}

func playGames(games int, quitch chan struct{}, gamech chan game) {

	var played int

	var maxTurns int
	var maxGame deck
	var maxTricks int

run:
	for {
		d, turns, tricks := oneGame()

		if turns > maxTurns {
			log.Printf("turns=%d deck=%s\n", turns, d)
			maxTurns = turns
			maxGame = d
			maxTricks = tricks
		}

		select {
		case <-quitch:
			break run
		default:

		}

		played++

		if games > 0 && played > games {
			break run
		}
	}

	gamech <- game{maxGame, maxTurns, maxTricks}

}

func main() {

	games := flag.Int("g", 10, "number of games")
	p := flag.Int("p", runtime.NumCPU(), "number of parallel games")

	flag.Parse()

	quitch := make(chan struct{})

	gamech := make(chan game)

	for i := 0; i < *p; i++ {
		go playGames(*games, quitch, gamech)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("quitting...")
		close(quitch)
	}()

	for i := 0; i < *p; i++ {
		g := <-gamech
		fmt.Printf("turns=%d tricks=%d deck=%s\n", g.turns, g.tricks, g.d)
	}
}
