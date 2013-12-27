package main

import (
	"code.google.com/p/rsc/qr"
	"flag"
	"fmt"
	"log"
)

func main() {

	level := flag.String("level", "m", "QR level: L/M/Q/H")

	flag.Parse()

	var qrLevel qr.Level

	switch *level {
	case "L", "l":
		qrLevel = qr.L
	case "M", "m":
		qrLevel = qr.M
	case "H", "h":
		qrLevel = qr.H
	case "Q", "q":
		qrLevel = qr.Q
	default:
		log.Fatal("unknown QR level:", *level)
	}

	data := flag.Arg(0)

	code, _ := qr.Encode(data, qrLevel)

	black := "\033[30;40m"
	white := "\033[30;47m"
	reset := "\033[0m"

	fmt.Print(reset)

	// two pixel border
	for i := 0; i < 2; i++ {
		fmt.Print(white + "    ")
		for x := 0; x < code.Size; x++ {
			fmt.Print("  ")
		}
		fmt.Println("    " + reset)
	}

	for x := 0; x < code.Size; x++ {
		fmt.Print(white + "    ") // two pixel border
		for y := 0; y < code.Size; y++ {
			if code.Black(x, y) {
				fmt.Print(black + "  ")
			} else {
				fmt.Print(white + "  ")
			}
		}
		fmt.Println(white + "    " + reset) // two pixel border
	}

	// two pixel border
	for i := 0; i < 2; i++ {
		fmt.Print(white + "    ")
		for x := 0; x < code.Size; x++ {
			fmt.Print("  ")
		}
		fmt.Println(white + "    " + reset)
	}
}
