package numerical

import (
	"fmt"
	"math"
	"testing"
)

func f1(x float64) float64 {
	return x - math.Tan(x)
}

func TestRootSolve(t *testing.T) {

	const samples = 200

	algorithms := []struct {
		name   string
		method func(func(float64) float64, float64, float64) float64
	}{
		{"bisect", Bisect},
		{"secant", Secant},
		{"falsi", Falsi},
		{"ridder", Ridder},
	}

	for _, algorithm := range algorithms {

		totalTicks := uint64(0)

		var r float64

		for i := 0; i < samples; i++ {
			ticks := rdtsc()
			r = algorithm.method(f1, 4.4, 4.6)
			totalTicks += rdtsc() - ticks

		}

		fmt.Printf("solve: %-20s\t=> %.10f (ticks=%d)\n", algorithm.name, r, totalTicks/samples)
	}
}
