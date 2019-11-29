package numerical

//import "fmt"

import "math"

func Trapezoid(start float64, end float64, steps int, fp func(float64) float64) float64 {

	h := (end - start) / float64(steps)

	total := 0.0
	fpX := fp(start)

	for i := 0; i < steps; i++ {
		x := start + h*float64(i)
		squareHeight := fpX
		fpX1 := fp(x + h)
		triangleHeight := fpX1 - squareHeight
		area := squareHeight*h + 0.5*(triangleHeight*h)
		total += area
		fpX = fpX1
	}

	return total
}

func TrapezoidFast(start float64, end float64, steps int, fp func(float64) float64) float64 {

	h := (end - start) / float64(steps)

	total := 0.5 * (fp(end) - fp(start))

	for i := 0; i < steps; i++ {
		total += fp(start + h*float64(i))
	}

	return total * h
}

func TrapezoidIterative(start float64, end float64, steps int, fp func(float64) float64) float64 {

	// adaptive method -- ignore provided 'steps'
	_ = steps
	steps = 1

	oldEstimate := 0.0
	newEstimate := trapezoidIterativeHelper(start, end, steps, fp, oldEstimate)

	for (oldEstimate-newEstimate)*(oldEstimate-newEstimate) > 1e-18 {
		steps++
		oldEstimate = newEstimate
		newEstimate = trapezoidIterativeHelper(start, end, steps, fp, oldEstimate)
	}

	return newEstimate
}

func trapezoidIterativeHelper(start float64, end float64, steps int, fp func(float64) float64, oldEstimate float64) float64 {

	if steps == 1 {
		return (fp(start) + fp(end)) * (end - start) / 2.0
	}

	n := 1 << uint(steps-2)         // number of new points
	h := (end - start) / float64(n) // spacing of new points
	x := start + h/2.0              // coord of first new point
	total := 0.0

	for i := 0; i < n; i++ {
		total += fp(x + float64(i)*h)
	}

	return (oldEstimate + h*total) / 2.0
}

func Romberg(start float64, end float64, steps int, fp func(float64) float64) (area float64) {

	MAX := 21 // iterations

	s := make([]float64, MAX)

	for k := 1; k < MAX; k++ {

		oldval := s[1]
		s[1] = trapezoidIterativeHelper(start, end, k, fp, oldval)

		for i := 2; i <= k; i++ {
			p := float64(uint(1) << (2 * uint(i-1))) // = 4 **(i-1)
			s[k] = (p*s[i-1] - oldval) / (p - 1)
			oldval = s[i]
			s[i] = s[k]
		}

		if math.Abs(area-s[k]) < 1e-9 {
			return s[k]
		}

		area = s[k]
	}

	println("max iterations reached")
	return s[MAX-1]
}

func Simpsons(start float64, end float64, steps int, fp func(float64) float64) float64 {
	total := fp(start) + 4.0*fp((start+end)/2) + fp(end)
	return (end - start) / 6.0 * total
}

func Simpsons38(start float64, end float64, steps int, fp func(float64) float64) float64 {
	total := fp(start) + 3.0*fp((2.0*start+end)/3) + 3.0*fp((2.0*end+start)/3) + fp(end)
	return (end - start) / 8.0 * total
}

func SimpsonsComposite(start float64, end float64, steps int, fp func(float64) float64) float64 {

	steps += (steps & 1)
	h := (end - start) / float64(steps)
	var totals [2]float64

	for i := 1; i < steps; i++ {
		totals[i&1] += fp(start + h*float64(i))
	}

	total := fp(start) + totals[0]*2.0 + totals[1]*4.0 + fp(end)

	return h * total / 3.0
}

func Simpsons38Composite(start float64, end float64, steps int, fp func(float64) float64) float64 {

	steps *= 3
	sections := steps / 3

	h := (end - start) / float64(steps)
	var totals [4]float64

	for i := 0; i < sections; i++ {
		x0 := start + 3.0*float64(i)*h
		totals[0] += fp(x0)
		x0 += h
		totals[1] += fp(x0)
		x0 += h
		totals[2] += fp(x0)
		x0 += h
		totals[3] += fp(x0)
	}

	total := totals[0] + 3.0*totals[1] + 3.0*totals[2] + totals[3]

	return 3.0 * h / 8.0 * total
}

// thin function -- prime the recursive pump and make the API match our other integration routines
func SimpsonsAdaptive(start float64, end float64, steps int, f func(float64) float64) float64 {
	//	epsilon := 1.0 / math.Pow(2, float64(steps))
	epsilon := 10e-8
	mid := (start + end) / 2
	h := end - start
	fStart := f(start)
	fEnd := f(end)
	fMid := f(mid)
	S := (h / 6) * (fStart + 4*fMid + fEnd)
	return simpsonsAdaptiveHelper(f, start, end, epsilon, S, fStart, fMid, fEnd, steps)
}

func simpsonsAdaptiveHelper(f func(float64) float64, start, end, epsilon, S, fStart, fMid, fEnd float64, bottom int) float64 {

	// lots of variables are declared here to prevent multiple evalutions of the function
	mid := (start + end) / 2
	h := end - start

	startmidMid := (start + mid) / 2
	fStartmid := f(startmidMid)

	endmidMid := (mid + end) / 2
	fEndmid := f(endmidMid)

	Sleft := (h / 12) * (fStart + 4*fStartmid + fMid)
	Sright := (h / 12) * (fMid + 4*fEndmid + fEnd)
	S2 := Sleft + Sright

	// we evaluate simpson's rule, then again summing the areas using the left and right midpoints.
	// If the sums are "the same", then we're done
	if bottom <= 0 || math.Abs(S2-S) <= 15*epsilon {
		return S2 + (S2-S)/15
	}

	// else recurse more
	return simpsonsAdaptiveHelper(f, start, mid, epsilon/2, Sleft, fStart, fStartmid, fMid, bottom-1) +
		simpsonsAdaptiveHelper(f, mid, end, epsilon/2, Sright, fMid, fEndmid, fEnd, bottom-1)
}

/*
double rkfp(double x, double y) {
    return (x*x) + (0.11*y*y) + 1.0;
}

double rungekutta4(double start, double end, int steps, double (*fp)(double, double), double x0, double y0)
{

   double k1, k2, k3, k4;
   double y1, x1;

   double h = (end - start) / steps;

   steps = 0;

   do {
      k1=h*fp(x0,y0);
      k2=h*fp(x0+h/2,y0+k1/2);
      k3=h*fp(x0+h/2,y0+k2/2);
      k4=h*fp(x0+h,y0+k3);
      y1=y0+0.1667*(k1+(2*k2)+(2*k3)+k4);
      x1=x0+h;
      x0=x1;
      y0=y1;
   } while (x0 <end);

   return y0;
}

*/
