package main

import (
	"regexp"
	"testing"

	"github.com/dgryski/trifles/sshdregex/asm"
	"github.com/dgryski/trifles/sshdregex/c"
)

var data = []byte(`Jan 18 06:41:30 corecompute sshd[42327]: Failed keyboard-interactive/pam for root from 112.100.68.182 port 48803 ssh2`)

var sink int

func BenchmarkRagel(b *testing.B) {
	var hits int
	for i := 0; i < b.N; i++ {
		if MatchRagel(data) {
			hits++
		}
	}

	sink += hits
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

	sink += hits
}

func BenchmarkFSM(b *testing.B) {
	var hits int

	for i := 0; i < b.N; i++ {
		if Match(data) != -1 {
			hits++
		}
	}

	sink += hits
}

func BenchmarkFSMC(b *testing.B) {
	var hits int

	for i := 0; i < b.N; i++ {
		if c.Match(data) != -1 {
			hits++
		}
	}

	sink += hits
}

func BenchmarkMatchASM(b *testing.B) {
	var hits int

	for i := 0; i < b.N; i++ {
		if asm.MatchASM(data) != -1 {
			hits++
		}
	}

	sink += hits
}
