package u64cmp

func naive(a, b []uint64) bool {

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func unrollxor(a, b []uint64) bool {

	if len(a) != len(b) {
		return false
	}

	var xor uint64

	for len(a) >= 8 && len(b) >= 8 && xor == 0 {
		xor |= a[0] ^ b[0]
		xor |= a[1] ^ b[1]
		xor |= a[2] ^ b[2]
		xor |= a[3] ^ b[3]
		xor |= a[4] ^ b[4]
		xor |= a[5] ^ b[5]
		xor |= a[6] ^ b[6]
		xor |= a[7] ^ b[7]
		a, b = a[8:], b[8:]
	}

	if xor != 0 {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
