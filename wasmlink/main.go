package main

import "fmt"

//export add
func add(a, b int32) int32

//export matchermain
func matchermain(a, b *byte) int32

func main() {
	fmt.Println(add(2, 3))

	runMatcher("meow ")
	runMatcher("woof ")
}

func runMatcher(s string) {
	buf := []byte(s)

	b := &buf[0]
	e := &buf[len(buf)-1]

	if matchermain(b, e) != -1 {
		fmt.Println("matcher hit for string", s)
	} else {
		fmt.Println("matcher miss for string", s)
	}
}
