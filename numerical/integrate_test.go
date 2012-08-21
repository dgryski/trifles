package numerical

import (
	"fmt"
	"math"
	"testing"
)

func f(x float64) float64 {
	return 6 + 2*x*x*math.Cos(x*x)
}

func TestIntegrators(t *testing.T) {

	const samples = 200

	algorithms := []struct {
		name   string
		method func(float64, float64, int, func(float64) float64) float64
	}{
		{"simpson", simpsons},
		{"simpson38", simpsons38},
		{"simpsoncomposite", simpsons_composite},
		{"simpson38composite", simpsons38_composite},
		{"simpson_adaptive", simpsons_adaptive},
		{"trapezoid", trapezoid},
		{"trapezoidfast", trapezoid_fast},
		{"trapezoiditerative", trapezoid_iterative},
		{"romberg", romberg},
	}

	for _, algorithm := range algorithms {

		total_ticks := uint64(0)

		var area float64

		for i := 0; i < samples; i++ {
			ticks := rdtsc()
			area = algorithm.method(0, math.SqrtPi, 100, f)
			total_ticks += rdtsc() - ticks
		}

		fmt.Printf("integrate: %-20s\t=> %.10f (avg. ticks=%d)\n", algorithm.name, area, total_ticks/samples)
	}
}
