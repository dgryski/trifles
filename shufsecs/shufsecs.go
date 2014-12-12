package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func outputShuf(currentSecond string, requests []string) {
	perm := rand.Perm(len(requests))
	for _, v := range perm {
		fmt.Print(currentSecond, "\t", requests[v], "\n")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var requests []string
	var currentSecond string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.SplitN(line, "\t", 3)
		if currentSecond != fields[0] {
			outputShuf(currentSecond, requests)
			currentSecond = fields[0]
			requests = requests[:0]
		}

		count := 1
		if len(fields) == 3 {
			count, _ = strconv.Atoi(fields[2])
		}
		for i := 0; i < count; i++ {
			requests = append(requests, fields[1])
		}
	}

	outputShuf(currentSecond, requests)
}
