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

// TODO:
// - add float ops
// - add i64 ops
// - add opcodes for tables/...
// - add simd opcodes

type opcode int

const (
	opI32Add = iota
	opI32And
	opI32Clz
	opI32Const
	opI32Ctz
	opI32Mul
	opI32Popcnt
	opI32Or
	opI32RemS
	opI32Rotl
	opI32Rotr
	opI32Shl
	opI32ShrS
	opI32ShrU
	opI32Sub
	opI32Xor

	opLocalGet
	opLocalSet
	opCall
	opSelect

	opMAXOPCODE
)

func (op opcode) String() string {
	switch op {
	case opI32Add:
		return "i32.add"
	case opI32And:
		return "i32.and"
	case opI32Clz:
		return "i32.clz"
	case opI32Const:
		return "i32.const"
	case opI32Ctz:
		return "i32.ctz"
	case opI32Mul:
		return "i32.mul"
	case opI32Popcnt:
		return "i32.popcnt"
	case opI32Or:
		return "i32.or"
	case opI32RemS:
		return "i32.rem_s"
	case opI32Rotl:
		return "i32.rotl"
	case opI32Rotr:
		return "i32.rotr"
	case opI32Shl:
		return "i32.shl"
	case opI32ShrS:
		return "i32.shr_s"
	case opI32ShrU:
		return "i32.shr_u"
	case opI32Sub:
		return "i32.sub"
	case opI32Xor:
		return "i32.xor"
	case opSelect:
		return "select"
	case opCall:
		return "call"
	case opLocalGet:
		return "local.get"
	case opLocalSet:
		return "local.set"

	case opMAXOPCODE:
		return "i32MAXOPCODE"
	}

	panic("unreached: " + strconv.Itoa(int(op)))
}

type Node struct {
	op   opcode
	args []Node
	i32  int32
}

func (n *Node) Write(w io.Writer) {

	for _, arg := range n.args {
		arg.Write(w)
	}

	if n.op == opI32Const || n.op == opLocalGet || n.op == opLocalSet || n.op == opCall {
		io.WriteString(w, n.op.String()+" "+strconv.FormatInt(int64(n.i32), 10)+"\n")
		return
	}

	io.WriteString(w, n.op.String()+"\n")
}

func (op opcode) arity() int {
	if op == opI32Const {
		return 0
	}

	if op == opI32Clz || op == opI32Ctz || op == opI32Popcnt {
		return 1
	}

	if op == opSelect {
		return 3
	}

	return 2
}

type Func struct {
	name    string
	nargs   int
	nlocals int
	body    Node
	idx     int32
}

type generator struct {
	funcs []Func
	fuel  int
}

func (g *generator) i32(f *Func) Node {

	if g.fuel == 0 || rand.Intn(5) == 0 {

		// sometimes generate a local.get
		if rand.Intn(3) == 0 && (f.nargs+f.nlocals) > 0 {
			return Node{
				op:  opLocalGet,
				i32: int32(rand.Intn(f.nargs + f.nlocals)),
			}
		}

		// generate a regular i32.const
		n := Node{
			op:  opI32Const,
			i32: rand.Int31(), // TODO(dgryski): generate edge values some of the time
		}

		if rand.Intn(2) == 0 {
			n.i32 = -n.i32
		}
		return n
	}

	// TODO:dgryski add semantic preserving operations here:
	// call an identity function, + 0, * 1,  xor 0, xor N / xor N, and -1,
	// if 0 unreachable else X, if 1 X else unreachable ; while(1) X break;

	g.fuel--

	op := opcode(rand.Intn(opMAXOPCODE))
	switch op {

	case opCall:
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
					op:   opCall,
					args: args,
					i32:  idx,
				}

				return n
			}
		}

		// generate a new function instead
		nargs := rand.Intn(16)
		nlocals := rand.Intn(16)
		newf := Func{
			nargs:   nargs,
			nlocals: nlocals,
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

	case opLocalGet:
		if f.nargs == 0 {
			// no arguments to current function; generate a random const instead
			return Node{
				op:  opI32Const,
				i32: rand.Int31(), // TODO(dgryski): generate edge values some of the time
			}
		}
		return Node{
			op:  op,
			i32: int32(rand.Intn(f.nargs + f.nlocals)),
		}

	case opLocalSet:
		if (f.nargs + f.nlocals) == 0 {
			// no arguments to current function; generate a random const instead
			return Node{
				op:  opI32Const,
				i32: rand.Int31(), // TODO(dgryski): generate edge values some of the time
			}
		}

		// generate two arguments; one for the set, one to leave on the stack as the generated i32
		return Node{
			op:   op,
			args: []Node{g.i32(f), g.i32(f)},
			i32:  int32(rand.Intn(f.nargs + f.nlocals)),
		}

	case opI32Const:
		return Node{
			op:  opI32Const,
			i32: rand.Int31(), // TODO(dgryski): select edge case sometimes
		}
	default:
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
		for i := 0; i < f.nlocals; i++ {
			io.WriteString(w, "(local i32) ")
		}
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
