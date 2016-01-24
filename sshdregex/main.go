package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {

	fname := flag.String("f", "/code/perl/auth.log", "input filename")
	flag.Parse()

	file, err := os.Open(*fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var i int

	// Compile regular expression first
	// r, _ := regexp.Compile(`sshd\[\d{5}\]:\s*Failed`)

	for scanner.Scan() {
		if matchSSHD(scanner.Bytes()) {
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)
}
