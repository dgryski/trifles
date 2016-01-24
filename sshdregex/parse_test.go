package main

import (
	"regexp"
	"testing"
)

var data = []byte(`Jan 18 06:41:30 corecompute sshd[42327]: Failed keyboard-interactive/pam for root from 112.100.68.182 port 48803 ssh2`)

func BenchmarkRagel(b *testing.B) {
	var hits int
	for i := 0; i < b.N; i++ {
		if matchSSHD(data) {
			hits++
		}
	}
}

func BenchmarkRegex(b *testing.B) {

	// Compile regular expression first
	r, _ := regexp.Compile(`sshd\[\d{5}\]:\s*Failed`)

	b.ResetTimer()

	var hits int

	for i := 0; i < b.N; i++ {
		if r.Match(data) {
			hits++
		}
	}
}
