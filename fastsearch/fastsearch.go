package fastsearch

// https://github.com/python/cpython/blob/master/Objects/stringlib/fastsearch.h

type bloom uint64

func (b *bloom) add(ch byte) {
	*b |= 1 << (ch & 63)
}

func (b bloom) get(ch byte) bool {
	return (b & (1 << (ch & 63))) != 0
}

func fastsearch(s, p []byte) int {

	n := len(s)
	m := len(p)

	w := n - m

	if w < 0 {
		return -1
	}

	var mask bloom

	mlast := m - 1
	skip := mlast - 1

	ss := s[m-1:]
	pp := p[m-1:]

	/* create compressed boyer-moore delta 1 table */

	/* process pattern[:-1] */
	for i, v := range p[:mlast] {
		mask.add(v)
		if v == p[mlast] {
			skip = mlast - i - 1
		}
	}
	/* process pattern[-1] outside the loop */
	mask.add(p[mlast])

	for i := 0; i <= w; i++ {
		/* note: using mlast in the skip path slows things down on x86 */
		if ss[i] == pp[0] {
			/* candidate match */
			var j int
			for j = 0; j < mlast; j++ {
				if s[i+j] != p[j] {
					break
				}
			}
			if j == mlast {
				return i
			}
			/* miss: check if next character is part of pattern */
			if !mask.get(ss[i+1]) {
				i = i + m
			} else {
				i = i + skip
			}
		} else {
			/* skip: check if next character is part of pattern */
			if i+1 >= len(ss) {
				continue
			}
			if !mask.get(ss[i+1]) {
				i = i + m
			}
		}
	}

	return -1
}
