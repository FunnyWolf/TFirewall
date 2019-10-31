package main

import (
	"TFirewall"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func TcpCheck(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	} else {
		_, _ = conn.Write([]byte(TFirewall.KeySend))
		return true
	}
}
func UdpCheck(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return false
	} else {
		_, _ = conn.Write([]byte(TFirewall.KeySend))
		return true
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:\n client 192.168.1.100\n client 192.168.1.100 20-22,80-90,22")
		os.Exit(1)
	}
	portslice := []int{}
	if len(os.Args) > 2 {
		portListStrs := strings.Split(os.Args[2], ",")
		for i := 0; i < len(portListStrs); i++ {
			if strings.Contains(portListStrs[i], "-") {
				parts := strings.Split(portListStrs[i], "-")
				start, err := strconv.Atoi(parts[0])
				if err == nil {
					end, err := strconv.Atoi(parts[1])
					if err == nil {
						for port := start; port < end; port++ {
							if !TFirewall.Contain(portslice, port) {
								portslice = append(portslice, port)
							}
						}
					}
				}
			} else {
				intPort, err := strconv.Atoi(portListStrs[i])
				if err == nil {
					if !TFirewall.Contain(portslice, intPort) {
						portslice = append(portslice, intPort)
					}
				}
			}
		}
	} else {
		portslice = TFirewall.TcpPorts()
	}

	ip := os.Args[1]
	for _, port := range portslice {
		addr := ip + ":" + strconv.Itoa(port)
		timeout := 1 * time.Second
		TcpCheck(addr, timeout)
	}

	for _, port := range portslice {
		addr := ip + ":" + strconv.Itoa(port)
		timeout := 1 * time.Second
		UdpCheck(addr, timeout)
	}
	fmt.Println("finish")
}
