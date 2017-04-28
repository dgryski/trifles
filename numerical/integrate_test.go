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
		{"simpson", Simpsons},
		{"simpson38", Simpsons38},
		{"simpsoncomposite", SimpsonsComposite},
		{"simpson38composite", Simpsons38Composite},
		{"simpson_adaptive", SimpsonsAdaptive},
		{"trapezoid", Trapezoid},
		{"trapezoidfast", TrapezoidFast},
		{"trapezoiditerative", TrapezoidIterative},
		{"romberg", Romberg},
	}

	for _, algorithm := range algorithms {

		totalTicks := uint64(0)

		var area float64

		for i := 0; i < samples; i++ {
			ticks := rdtsc()
			area = algorithm.method(0, math.SqrtPi, 100, f)
			totalTicks += rdtsc() - ticks
		}

		fmt.Printf("integrate: %-20s\t=> %.10f (avg. ticks=%d)\n", algorithm.name, area, totalTicks/samples)
	}
}
