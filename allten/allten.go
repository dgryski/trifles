package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type op int

const (
	opAdd op = iota
	opSub
	opDiv
	opMul
	opPush
)

func randop() op {
	return op(rand.Intn(int(opPush)))
}

type operation struct {
	n int
	o op
}

var opStrs = []string{"+", "-", "/", "*"}

func (o operation) String() string {
	if o.o == opPush {
		return strconv.Itoa(o.n)
	}

	return opStrs[o.o]
}

type stack[T any] []T

func (s *stack[T]) push(t T) {
	*s = append(*s, t)
}

func (s *stack[T]) pop() T {
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t
}

type expr []operation

func (ops expr) String() string {
	var s stack[string]

	for _, op := range ops {
		if op.o == opPush {
			s.push(op.String())
			continue
		}
		a := s.pop()
		b := s.pop()
		s.push("(" + b + op.String() + a + ")")
	}

	return s.pop()
}

func (ops expr) eval() float64 {
	var s stack[float64]

	for _, op := range ops {
		if op.o == opPush {
			s.push(float64(op.n))
			continue
		}
		a := s.pop()
		b := s.pop()
		switch op.o {
		case opAdd:
			s.push(b + a)
		case opSub:
			s.push(b - a)
		case opMul:
			s.push(b * a)
		case opDiv:
			s.push(b / a)
		}
	}

	return s.pop()
}

const eps = 0.0000001

func f64eq(a, b float64) bool {
	return math.Abs(a-b) < eps
}

func randelt(in []int) (int, []int) {
	idx := rand.Intn(len(in))
	n := in[idx]
	in[idx] = in[len(in)-1]
	in = in[:len(in)-1]
	return n, in
}

func solve(numbers []int, target float64) expr {
	in := make([]int, 0, len(numbers))

	for i := 0; i < 10000; i++ {
		// copy in our numbers
		var n int
		in = append(in, numbers...)

		var ops expr
		var nnum, nops int

		for len(in) > 1 {
			// choose a number
			n, in = randelt(in)
			ops = append(ops, operation{n, opPush})
			nnum++

			// plus maybe some operations
			for n, i := rand.Intn(nnum-nops), 0; i < n; i++ {
				ops = append(ops, operation{0, randop()})
				nops++
			}
		}

		// last number
		n, in = randelt(in)
		ops = append(ops, operation{n, opPush})
		nnum++

		// remaining operations
		for nops < (nnum - 1) {
			ops = append(ops, operation{n, randop()})
			nops++
		}

		if v := ops.eval(); f64eq(v, target) {
			return ops
		}
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	target := flag.Int("t", 0, "target value")
	solns := flag.Int("c", 1, "number of solutions")
	flag.Parse()

	var numbers []int
	for _, arg := range flag.Args() {
		n, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("unable to parse number %q: %v", arg, err)
			return
		}

		numbers = append(numbers, int(n))
	}

	targets := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if *target != 0 {
		targets = []int{*target}
	}

	for _, t := range targets {
		var solnFound bool
		for i := 0; i < *solns; i++ {
			ops := solve(numbers, float64(t))
			if ops != nil {
				solnFound = true
				fmt.Println(ops, "=", t)
			}
		}
		if !solnFound {
			fmt.Println(";", t, ": no solution found")
		}
	}
}
