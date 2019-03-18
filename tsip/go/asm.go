// +build ignore

package main

import (
	"strconv"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func sipround(v0, v1 Register) {
	ADDQ(v1, v0)

	ROLQ(Imm(13), v1)
	XORQ(v0, v1)

	ROLQ(Imm(35), v0)
	ADDQ(v1, v0)

	ROLQ(Imm(17), v1)
	XORQ(v0, v1)

	ROLQ(Imm(21), v0)
}

func main() {
	Package("github.com/dgryski/trifles/tsip/go/")

	TEXT("HashASM", NOSPLIT, "func(k0, k1 uint64, p []byte) uint64")

	reg_v0 := GP64()
	reg_v1 := GP64()

	Load(Param("k0"), reg_v0)
	Load(Param("k1"), reg_v1)

	reg_magic := GP64()
	MOVQ(Imm(0x736f6d6570736575), reg_magic)
	XORQ(reg_magic, reg_v0)
	MOVQ(Imm(0x646f72616e646f6d), reg_magic)
	XORQ(reg_magic, reg_v1)

	reg_p := GP64()
	reg_p_len := GP64()

	Load(Param("p").Base(), reg_p)
	Load(Param("p").Len(), reg_p_len)

	reg_b := GP64()
	MOVQ(reg_p_len, reg_b)
	SHLQ(Imm(56), reg_b)

	reg_m := GP64()

	loop_end := "loop_end"
	loop_begin := "loop_begin"
	CMPQ(reg_p_len, Imm(8))
	JL(LabelRef(loop_end))

	Label(loop_begin)
	MOVQ(Mem{Base: reg_p}, reg_m)
	XORQ(reg_m, reg_v1)
	sipround(reg_v0, reg_v1)
	XORQ(reg_m, reg_v0)

	ADDQ(Imm(8), reg_p)
	SUBQ(Imm(8), reg_p_len)
	CMPQ(reg_p_len, Imm(8))
	JGE(LabelRef(loop_begin))
	Label(loop_end)

	var labels []string
	for i := 0; i < 8; i++ {
		labels = append(labels, "sw"+strconv.Itoa(i))
	}

	for i := 0; i < 7; i++ {
		CMPQ(reg_p_len, Imm(uint64(i)))
		JE(LabelRef(labels[i]))
	}

	char := GP64()
	for i := 7; i > 0; i-- {
		Label(labels[i])
		MOVBQZX(Mem{Base: reg_p, Disp: i - 1}, char)
		SHLQ(Imm(uint64(i-1)*8), char)
		ORQ(char, reg_b)

	}

	Label(labels[0])

	XORQ(reg_b, reg_v1)
	sipround(reg_v0, reg_v1)
	XORQ(reg_b, reg_v0)

	XORQ(Imm(0xff), reg_v1)
	sipround(reg_v0, reg_v1)
	ROLQ(Imm(32), reg_v1)
	sipround(reg_v0, reg_v1)
	ROLQ(Imm(32), reg_v1)

	XORQ(reg_v1, reg_v0)

	Store(reg_v0, ReturnIndex(0))
	RET()

	Generate()
}
