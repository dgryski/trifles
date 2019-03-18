// +build amd64

package tsip

//go:generate go run asm.go -out tsip_amd64.s
//go:noescape

func HashASM(k0, k1 uint64, p []byte) uint64
