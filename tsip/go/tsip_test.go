package tsip

import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

func TestTSIP(t *testing.T) {

	var k0 uint64 = 0x0706050403020100
	var k1 uint64 = 0x0f0e0d0c0b0a0908

	f, err := os.Open("testdata/tsip.txt")
	if err != nil {
		t.Fatalf("unable to open expected results: %v", err)
	}

	scanner := bufio.NewScanner(f)

	var buf []byte

	var i int

	for scanner.Scan() {

		wantStr := scanner.Text()

		want, err := strconv.ParseUint(wantStr, 16, 64)
		if err != nil {
			t.Errorf("error parsing want: %v: %v", wantStr, err)
		}

		got := Hash(k0, k1, buf)

		if got != want {
			t.Errorf("failed: %d: got=%v, want=%v\n", i, got, want)
		}

		got = HashASM(k0, k1, buf)

		if got != want {
			t.Errorf("failed asm: %d: got=%v, want=%v\n", i, got, want)
		}

		buf = append(buf, byte(i))
		i++
	}
}
