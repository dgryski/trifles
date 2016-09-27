package shellsort

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

type bigarray []int

func (b bigarray) Generate(rand *rand.Rand, size int) reflect.Value {
	nv := make([]int, size*1000)

	for i := range nv {
		nv[i] = rand.Intn(1000 * size)
	}

	return reflect.ValueOf(nv)
}

func testSort(t *testing.T, sort func(a []int), which string) {

	rand.Seed(0)

	f := func(a bigarray) bool {
		sort(a)
		for i := 1; i < len(a); i++ {
			if a[i-1] > a[i] {
				return false
			}
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("%s: %v", which, err)
	}
}

func TestShell(t *testing.T) { testSort(t, Shell, "shell") }
func TestRand(t *testing.T)  { testSort(t, RandShell, "rand") }
