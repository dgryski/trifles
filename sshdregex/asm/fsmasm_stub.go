//go:build !tinygo

package asm

//go:noescape
func Match(data []byte) int
