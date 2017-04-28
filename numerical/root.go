package numerical

import (
	"math"
)

func Bisect(f func(x float64) float64, x0, x1 float64) (root float64) {

	eps := 10e-8

	fx0 := f(x0)

	for x1-x0 > eps {
		root = (x0 + x1) / 2
		froot := f(root)
		if froot*fx0 < 0 {
			x1 = root
		} else {
			x0 = root
			fx0 = froot
		}
	}

	return root
}

func Secant(f func(x float64) float64, x0, x1 float64) (root float64) {

	eps := 10e-8

	fx0 := f(x0)
	for math.Abs(x1-x0) > eps {
		fx1 := f(x1)
		x2 := x1 - fx1*((x1-x0)/(fx1-fx0))
		x1, x0 = x2, x1
		fx0 = fx1
	}

	return x1
}

// wikipedia

func Falsi(f func(x float64) float64, s, t float64) (r float64) {

	side := 0

	fs := f(s)
	ft := f(t)

	m := 100 // max iterations

	e := 5e-15 // eps

	for n := 1; n <= m; n++ {
		r = (fs*t - ft*s) / (fs - ft)
		if math.Abs(t-s) < e*math.Abs(t+s) {
			break
		}
		fr := f(r)

		if fr*ft > 0 {
			t = r
			ft = fr
			if side == -1 {
				fs /= 2
			}
			side = -1
		} else if fs*fr > 0 {
			s = r
			fs = fr
			if side == +1 {
				ft /= 2
			}
			side = +1
		} else {
			break
		}
	}
	return r
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Numerical Methods in Engineering with Python -> Go

func Ridder(f func(x float64) float64, a, b float64) (r float64) {

	tol := 1.0e-9

	fa := f(a)
	if fa == 0.0 {
		return a
	}
	fb := f(b)
	if fb == 0.0 {
		return b
	}

	// if fa*fb > 0.0: errror.err('Root is not bracketed')

	var xOld float64

	for i := 0; i < 30; i++ {
		// Compute the improved root x from Ridder's formula
		c := 0.5 * (a + b)
		fc := f(c)
		s := math.Sqrt(fc*fc - fa*fb)
		if s == 0.0 {
			return 0.0
		}
		dx := (c - a) * fc / s
		if (fa - fb) < 0.0 {
			dx = -dx
		}
		x := c + dx
		fx := f(x)
		// Test for convergence
		if i > 0 && math.Abs(x-xOld) < tol*max(math.Abs(x), 1.0) {
			return x
		}
		xOld = x
		// Re-bracket the root as tightly as possible
		if fc*fx > 0.0 {
			if fa*fx < 0.0 {
				b = x
				fb = fx
			} else {
				a = x
				fa = fx
			}
		} else {
			a = c
			b = x
			fa = fc
			fb = fx
		}

	}

	// too many iterations
	return 0

}
