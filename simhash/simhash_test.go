package simhash

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func simhashString(s string) uint64 {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	return SimHash(scanner)
}

func TestSimHash(t *testing.T) {

	h1 := simhashString("Now is the winter of our discontent and also the time for all good people to come to the aid of the party")
	fmt.Printf("h=%016x\n", h1)

	h2 := simhashString("Now is the winter of our discontent and also the time for all good people to come and party")
	fmt.Printf("h=%016x\n", h2)

	h3 := simhashString("The more we get together together together the more we get together the happier we'll be")
	fmt.Printf("h=%016x\n", h3)

	fmt.Printf("d(h1,h2)=%d\n", hammingDistance(h1, h2))
	fmt.Printf("d(h1,h3)=%d\n", hammingDistance(h1, h3))
	fmt.Printf("d(h2,h3)=%d\n", hammingDistance(h2, h3))

}
