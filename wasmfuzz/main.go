package main

import (
	"fmt"
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
	i32Const
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
	case i32Const:
		return "i32.const"
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

func (op i32Op) eval(args []Node) int32 {
	switch op {
	case i32Add:
		return args[0].i32 + args[1].i32
	case i32And:
		return args[0].i32 & args[1].i32
	case i32Const:
		return args[0].i32
	case i32Mul:
		return args[0].i32 * args[1].i32
	case i32Or:
		return args[0].i32 | args[1].i32
	case i32RemS:
		return args[0].i32 % args[1].i32
	case i32Sub:
		return args[0].i32 - args[1].i32
	case i32Xor:
		return args[0].i32 ^ args[1].i32
	}
	panic("unreached:" + op.String())
}

func (n *Node) Write(w io.Writer) {
	if n.op == i32Const {
		io.WriteString(w, n.op.String()+" "+strconv.FormatInt(int64(n.i32), 10)+"\n")
		return
	}

	for _, arg := range n.args {
		arg.Write(w)
	}

	io.WriteString(w, n.op.String()+"\n")
}

func (op i32Op) arity() int {
	if op == i32Const {
		return 0
	}

	return 2
}

type generator struct {
	fuel int
}

func (g *generator) i32() Node {

	if g.fuel == 0 || rand.Intn(5) == 0 {
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
			args = append(args, g.i32())
		}

		if op == i32RemS && args[1].i32 == 0 {
			// insert another node to add 1 to avoid division by zero
			constOne := Node{
				op:  i32Const,
				i32: 1, // TODO(dgryski): add random here?
			}
			addOne := Node{
				op:   i32Add,
				args: []Node{args[1], constOne},
				i32:  args[1].i32 + constOne.i32,
			}
			args[1] = addOne
		}

		return Node{
			op:   op,
			args: args,
			i32:  op.eval(args),
		}

	case i32Const:
		return Node{
			op:  i32Const,
			i32: rand.Int31(), // TODO(dgryski): select edge case sometimes
		}
	}

	panic("generate.i32: " + op.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(`(module (func (export "_start")`)
	fuel := rand.Intn(100)
	g := generator{fuel: fuel}
	n := g.i32()
	n.Write(os.Stdout)
	fmt.Println("i32.const " + strconv.FormatInt(int64(n.i32), 10))
	fmt.Println("i32.ne")
	fmt.Println("(if (then unreachable))")
	fmt.Println(") )")
}
