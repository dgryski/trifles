// Package wtflog provides more intuitive logging levels
/*
Original proposal https://twitter.com/jbarnette/status/144791640589078530
*/
package wtflog

import (
	"log"
	"os"
)

var (
	Lwtf = log.New(os.Stderr, "WTF ", log.LstdFlags)
	Lomg = log.New(os.Stderr, "OMG ", log.LstdFlags)
	Lfyi = log.New(os.Stderr, "FYI ", log.LstdFlags)
	Lbtw = log.New(os.Stderr, "BTW ", log.LstdFlags)
)

func WTF(args ...interface{}) { Lwtf.Println(args...) }
func OMG(args ...interface{}) { Lomg.Println(args...) }
func FYI(args ...interface{}) { Lfyi.Println(args...) }
func BTW(args ...interface{}) { Lbtw.Println(args...) }
