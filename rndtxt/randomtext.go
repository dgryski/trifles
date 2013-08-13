// Package rndtxt generates random text strings
package rndtxt

// Generate generates ln bytes of random text
func Generate(ln int) []byte {
	// yoinked from a go-nuts post by rsc
	var x = 0xFFFFFFFF
	text := make([]byte, ln)
	for i := range text {
		x += x
		if int(x) < 0 {
			x ^= 0x88888EEF
		}
		v := byte(uint(x) % 95)
		if v == 0 {
			text[i] = '\n'
		} else {
			text[i] = v + ' '
		}
	}
	return text
}
