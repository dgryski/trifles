package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	var z big.Int
	var a, b, c big.Int
	var r big.Int

	two := big.NewInt(2)

	var n_add, n_sub, n_mul2, n_div2, n_gcd, n_sqr, n_invmod, n_exptmod int

	for scanner.Scan() {
		cmd := scanner.Text()

		fmt.Println(n_add, n_sub, n_mul2, n_div2, n_gcd, n_sqr, n_invmod)

		switch cmd {

		case "add_d", "add":
			n_add++
			next(scanner, &a)
			next(scanner, &b)
			next(scanner, &r)

			z.Add(&a, &b)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		case "sub_d", "sub":
			n_sub++
			next(scanner, &a)
			next(scanner, &b)
			next(scanner, &r)

			z.Sub(&a, &b)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		case "mul2":
			n_mul2++
			next(scanner, &a)
			next(scanner, &r)

			z.Mul(&a, two)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		case "div2":
			n_div2++
			next(scanner, &a)
			next(scanner, &r)

			z.Div(&a, two)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}
		case "sqr":
			n_sqr++
			next(scanner, &a)
			next(scanner, &r)

			z.Mul(&a, &a)
			if false && z.Cmp(&r) != 0 {
				fmt.Printf("sqr: %s %s %s\n", a.String(), r.String(), z.String())
			}

		case "gcd":
			n_gcd++
			next(scanner, &a)
			next(scanner, &b)
			next(scanner, &r)

			z.GCD(nil, nil, &a, &b)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		case "invmod":
			n_invmod++
			next(scanner, &a)
			next(scanner, &b)
			next(scanner, &r)

			z.ModInverse(&a, &b)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		case "exptmod":
			n_exptmod++
			next(scanner, &a)
			next(scanner, &b)
			next(scanner, &c)
			next(scanner, &r)

			z.Exp(&a, &b, &c)
			if z.Cmp(&r) != 0 {
				fmt.Println(a, b, r, z)
			}

		default:
			// nothing
		}
	}
}

func next(scanner *bufio.Scanner, n *big.Int) {

	ok := scanner.Scan()
	if !ok {
		panic("end of input")
	}

	txt := scanner.Text()

	n.SetString(txt, 10)
}
