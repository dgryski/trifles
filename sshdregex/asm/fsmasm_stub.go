//go:build amd64 && !tinygo

package asm

//go:noescape
func Match(data []byte) int
