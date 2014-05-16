package inthash

import (
	"testing"
	"testing/quick"
)

func TestHash64Inv(t *testing.T) {

	f := func(x uint64) bool {
		y := Hash64(uint64(x))
		return Hash64Inv(y) == x
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
