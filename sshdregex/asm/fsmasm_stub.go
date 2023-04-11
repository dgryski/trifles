//go:build !tinygo

package asm

//go:noescape
func MatchASM(data []byte) int
