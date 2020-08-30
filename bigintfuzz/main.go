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

	var a, b, c big.Int

	ca := bigintCreate();
	cb := bigintCreate();
	cc := bigintCreate();


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

		ca = bigintAssignString(ca, a.String())
		cb = bigintAssignString(cb, b.String())

		switch rnd.Intn(3) {

		case 0:
			bigintAdd(ca, cb, cc)
			s := bigintFormat(cc)

			c.Add(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("+", &a, &b, &c, s)
			}

		case 1:
			r := bigintSub(ca, cb, cc)
			s := bigintFormat(r)

			c.Sub(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("-", &a, &b, &c, s)
			}

		case 2:
			r := bigintMul(ca, cb, cc)
			s := bigintFormat(r)

			c.Mul(&a, &b)
			if cstr := (&c).String(); s != cstr {
				log("*", &a, &b, &c, s)
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

func bigintCreate() *C.bigint {
	return C.bigint_create()
}

func bigintAssignString(a *C.bigint, s string) *C.bigint {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return C.bigint_assign_string(a, cs)
}

func bigintAdd(a, b, c *C.bigint) *C.bigint {
	return C.bigint_add(a, b, c)
}

func bigintSub(a, b, c *C.bigint) *C.bigint {
	return C.bigint_sub(a, b, c)
}

func bigintMul(a, b, c *C.bigint) *C.bigint {
	return C.bigint_mul(a, b, c)
}

func bigintFormat(a *C.bigint) string {
	var s [1000]byte
	C.bigint_format(a, (*C.char)(unsafe.Pointer(&s[0])))
	i := bytes.IndexByte(s[:], 0)
	return string(s[:i])
}
