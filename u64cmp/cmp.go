package u64cmp

import (
	"bytes"
	"reflect"
	"unsafe"
)

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

	for len(a) >= 32 && len(b) >= 32 && xor == 0 {
		xor |= a[0+0] ^ b[0+0]
		xor |= a[0+1] ^ b[0+1]
		xor |= a[0+2] ^ b[0+2]
		xor |= a[0+3] ^ b[0+3]
		xor |= a[0+4] ^ b[0+4]
		xor |= a[0+5] ^ b[0+5]
		xor |= a[0+6] ^ b[0+6]
		xor |= a[0+7] ^ b[0+7]

		xor |= a[8+0] ^ b[8+0]
		xor |= a[8+1] ^ b[8+1]
		xor |= a[8+2] ^ b[8+2]
		xor |= a[8+3] ^ b[8+3]
		xor |= a[8+4] ^ b[8+4]
		xor |= a[8+5] ^ b[8+5]
		xor |= a[8+6] ^ b[8+6]
		xor |= a[8+7] ^ b[8+7]

		xor |= a[16+0] ^ b[16+0]
		xor |= a[16+1] ^ b[16+1]
		xor |= a[16+2] ^ b[16+2]
		xor |= a[16+3] ^ b[16+3]
		xor |= a[16+4] ^ b[16+4]
		xor |= a[16+5] ^ b[16+5]
		xor |= a[16+6] ^ b[16+6]
		xor |= a[16+7] ^ b[16+7]

		xor |= a[24+0] ^ b[24+0]
		xor |= a[24+1] ^ b[24+1]
		xor |= a[24+2] ^ b[24+2]
		xor |= a[24+3] ^ b[24+3]
		xor |= a[24+4] ^ b[24+4]
		xor |= a[24+5] ^ b[24+5]
		xor |= a[24+6] ^ b[24+6]
		xor |= a[24+7] ^ b[24+7]

		a, b = a[32:], b[32:]
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

func bytesEq(a, b []uint64) bool {

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&a)).Data)

	var aBytes []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&aBytes))
	hdr.Data = uintptr(p)
	hdr.Len = len(a) * 8
	hdr.Cap = len(a) * 8

	p = unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)

	var bBytes []byte
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&bBytes))
	hdr.Data = uintptr(p)
	hdr.Len = len(b) * 8
	hdr.Cap = len(b) * 8

	return bytes.Equal(aBytes, bBytes)
}
