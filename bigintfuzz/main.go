package main

/*

#include <stdio.h>
#include <stdlib.h>
#include "bigint.h"

*/
import "C"

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"os"
	"time"
	"unsafe"
)

func main() {

	var ca, cb C.bigint
	var a, b, c big.Int

	bigintInit(&ca)
	bigintInit(&cb)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 1000000; i++ {
		if i % 10000 == 0 {
			fmt.Println("iterations", i)
		}
		nsize := rnd.Intn(200) + 1
		n := new(big.Int).Lsh(big.NewInt(1), uint(nsize))
		a.Rand(rnd, n)


		nsize = rnd.Intn(200)  + 1
		n = new(big.Int).Lsh(big.NewInt(1), uint(nsize))
		b.Rand(rnd, n)

		long := int32(rnd.Uint32())

		bigintAssignString(&ca, a.String())
		bigintAssignString(&cb, b.String())

		switch rnd.Intn(6) {

		case 0:
			bigintAdd(&ca, &cb)
			s := bigintFormat(&ca)

			c.Add(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("+", &a, &b, &c, s)
			}

		case 1:
			r := bigintSub(&ca, &cb)
			s := bigintFormat(r)

			c.Sub(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("-", &a, &b, &c, s)
			}

		case 2:
			r := bigintMul(&ca, &cb)
			s := bigintFormat(r)

			c.Mul(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("*", &a, &b, &c, s)
			}

		case 3:
			r := bigintAddInt(&ca, long)
			s := bigintFormat(r)

			b.SetInt64(int64(long))
			c.Add(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("l+", &a, &b, &c, s)
			}

		case 4:
			r := bigintSubInt(&ca, long)
			s := bigintFormat(r)

			b.SetInt64(int64(long))
			c.Sub(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("l-", &a, &b, &c, s)
			}

		case 5:
			r := bigintMulInt(&ca, long)
			s := bigintFormat(r)

			b.SetInt64(int64(long))
			c.Mul(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("l*", &a, &b, &c, s)
			}
		}
	}
}

func compare(a, b string) string {

	var sb strings.Builder

	if len(a) != len(b) {
		return	fmt.Sprintf("lengths differ: %v <=> %v", a, b)
	}

	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			sb.WriteByte(a[i])
		} else {
			sb.WriteString("**>")
			sb.WriteByte(a[i])
			sb.WriteByte(':')
			sb.WriteByte(b[i])
			sb.WriteString("<**")
		}
	}

	return sb.String()
}

func log(op string, a, b, c *big.Int, s string) {
	fmt.Printf("%v %s %v = %v\n", a, op, b, c)
	fmt.Printf("fail: got %q want %q\n", s, c.String())
	fmt.Printf("diff: %s\n", compare(s, c.String()))
	os.Exit(1)
}

func bigintInit(a *C.bigint) {
	C.bigint_init(a)
}

func bigintAssignString(a *C.bigint, s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.bigint_assign_string(a, cs, 10)
}

func bigintAdd(a, b *C.bigint) *C.bigint {
	return C.bigint_add_bigint(a, b)
}

func bigintAddInt(a *C.bigint, b int32) *C.bigint {
	return C.bigint_add_integer(a, C.long(b))
}

func bigintSub(a, b *C.bigint) *C.bigint {
	return C.bigint_sub_bigint(a, b)
}

func bigintSubInt(a *C.bigint, b int32) *C.bigint {
	return C.bigint_sub_integer(a, C.long(b))
}

func bigintMul(a, b *C.bigint) *C.bigint {
	return C.bigint_mul_bigint(a, b)
}

func bigintMulInt(a *C.bigint, b int32) *C.bigint {
	return C.bigint_mul_integer(a, C.long(b))
}

func bigintFormat(a *C.bigint) string {
	var s [1000]byte
	C.bigint_format(a, (*C.char)(unsafe.Pointer(&s[0])))
	i := bytes.IndexByte(s[:], 0)
	return string(s[:i])
}
