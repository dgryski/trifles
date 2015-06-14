package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type network struct {
	ip   uint32
	bits uint32
}

func (n network) String() string {
	return fmt.Sprintf("%d.%d.%d.%d/%d", byte(n.ip>>24), byte(n.ip>>16), byte(n.ip>>8), byte(n.ip), n.bits)
}

func range2cidr(ipStart uint32, ipEnd uint32) []network {
	var nets []network
	newip := ipStart
	for newip <= ipEnd {
		bits, mask := uint32(1), uint32(1)
		for bits < 32 {
			newip = ipStart | mask
			if (newip > ipEnd) || ((ipStart>>bits)<<bits) != ipStart {
				bits--
				mask >>= 1
				break
			}
			bits++
			mask = (mask << 1) + 1
		}
		nets = append(nets, network{ipStart, 32 - bits})
		ipStart = (ipStart | mask) + 1
		newip = ipStart
	}
	return nets
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			fmt.Println("bad line:", scanner.Text())
			continue
		}
		ipStart, ipEnd := str2ip32(fields[0]), str2ip32(fields[1])
		if ipStart == 0 || ipEnd == 0 {
			fmt.Println("error parsing:", scanner.Text())
			continue
		}
		fmt.Printf("%s - %s: %v", fields[0], fields[1], range2cidr(ipStart, ipEnd))
	}
}

func str2ip32(s string) uint32 {
	netip := net.ParseIP(s)
	if netip == nil {
		return 0
	}
	ip4 := netip.To4()
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}
