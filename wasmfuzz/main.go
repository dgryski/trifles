package main

import (
	"bufio"
	"flag"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type i32Op int

const (
	i32Add = iota
	i32And
	i32Call
	i32Const
	i32LocalGet
	i32Mul
	i32Or
	i32RemS
	i32Sub
	i32Xor

	i32MAXOPCODE
)

func (op i32Op) String() string {
	switch op {
	case i32Add:
		return "i32.add"
	case i32And:
		return "i32.and"
	case i32Call:
		return "call"
	case i32Const:
		return "i32.const"
	case i32LocalGet:
		return "local.get"
	case i32Mul:
		return "i32.mul"
	case i32Or:
		return "i32.or"
	case i32RemS:
		return "i32.rem_s"
	case i32Sub:
		return "i32.sub"
	case i32Xor:
		return "i32.xor"
	case i32MAXOPCODE:
		return "i32MAXOPCODE"
	}

	panic("unreached: " + strconv.Itoa(int(op)))
}

type Node struct {
	op   i32Op
	args []Node
	i32  int32
}

func (n *Node) Write(w io.Writer) {
	for _, arg := range n.args {
		arg.Write(w)
	}

	if n.op == i32Const || n.op == i32LocalGet || n.op == i32Call {
		io.WriteString(w, n.op.String()+" "+strconv.FormatInt(int64(n.i32), 10)+"\n")
		return
	}

	io.WriteString(w, n.op.String()+"\n")
}

func (op i32Op) arity() int {
	if op == i32Const {
		return 0
	}

	return 2
}

type Func struct {
	name  string
	nargs int
	body  Node
	idx   int32
}

type generator struct {
	funcs []Func
	fuel  int
}

func (g *generator) i32(f *Func) Node {

	if g.fuel == 0 || rand.Intn(5) == 0 {

		// sometimes generate a local.get
		if rand.Intn(3) == 0 && f.nargs > 0 {
			return Node{
				op:  i32LocalGet,
				i32: int32(rand.Intn(f.nargs)),
			}
		}

		// generate a regular i32.const
		n := Node{
			op:  i32Const,
			i32: rand.Int31(), // TODO(dgryski): generate edge values some of the time
		}

		if rand.Intn(2) == 0 {
			n.i32 = -n.i32
		}
		return n
	}

	g.fuel--

	op := i32Op(rand.Intn(i32MAXOPCODE))
	switch op {
	case i32Add, i32And, i32Mul, i32Or, i32RemS, i32Sub, i32Xor:
		var args []Node
		nargs := op.arity()
		for i := 0; i < nargs; i++ {
			args = append(args, g.i32(f))
		}

		// ideally we'd be able to avoid divide-by-zero here, but ditching it for now;
		// too hard to statically determine now that we've added function calls

		return Node{
			op:   op,
			args: args,
		}

	case i32Call:
		if rand.Intn(3) == 0 {
			idx := int32(rand.Intn(len(g.funcs)))
			if idx != f.idx {
				// don't allow recursive calls
				call := g.funcs[idx]
				var args []Node
				for i := 0; i < call.nargs; i++ {
					args = append(args, g.i32(f))
				}
				n := Node{
					op:   i32Call,
					args: args,
					i32:  idx,
				}

				return n
			}
		}

		// generate a new function instead
		nargs := rand.Intn(16)
		newf := Func{
			nargs: nargs,
		}
		newf.body = g.i32(&newf)
		g.funcs = append(g.funcs, newf)
		newf.idx = int32(len(g.funcs)) - 1

		// now generate a call to f
		var args []Node
		for i := 0; i < newf.nargs; i++ {
			args = append(args, g.i32(f))
		}
		return Node{
			op:   op,
			args: args,
			i32:  newf.idx,
		}

	case i32LocalGet:
		if f.nargs == 0 {
			// no arguments to current function; generate a random const instead
			return Node{
				op:  i32Const,
				i32: rand.Int31(), // TODO(dgryski): generate edge values some of the time
			}
		}
		return Node{
			op:  op,
			i32: int32(rand.Intn(f.nargs)),
		}

	case i32Const:
		return Node{
			op:  i32Const,
			i32: rand.Int31(), // TODO(dgryski): select edge case sometimes
		}
	}

	panic("generate.i32: " + op.String())
}

func (g *generator) Write(w io.Writer) {
	io.WriteString(w, "(module ")

	for _, f := range g.funcs {
		io.WriteString(w, "(func ")
		if f.name != "" {
			io.WriteString(w, `(export "`+f.name+`") `)
		}
		for i := 0; i < f.nargs; i++ {
			io.WriteString(w, "(param i32) ")
		}
		io.WriteString(w, "(result i32) ")
		io.WriteString(w, "\n")
		f.body.Write(w)
		io.WriteString(w, ")\n")
	}

	io.WriteString(w, ")\n")
}

func main() {
	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	flag.Parse()

	rand.Seed(*seed)
	fuel := rand.Intn(1000)
	start := Func{
		name:  "_start",
		nargs: 0,
		idx:   0,
	}
	g := generator{
		fuel:  fuel,
		funcs: []Func{start},
	}

	g.funcs[0].body = g.i32(&start)

	b := bufio.NewWriter(os.Stdout)
	g.Write(b)
	b.Flush()
}
