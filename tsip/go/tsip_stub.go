// +build amd64

package tsip

//go:generate python -m peachpy.x86_64 tsip.py -S -o tsip_amd64.s -mabi=goasm
//go:noescape

func HashASM(k0, k1 uint64, p []byte) uint64
