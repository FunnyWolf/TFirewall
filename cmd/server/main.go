package main

import (
	"TFirewall"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var c chan os.Signal

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func recvUDP(conn *net.UDPConn) {
	var buf [20]byte
	for {
		n, radder, err := conn.ReadFromUDP(buf[0:])
		if err == nil {
			if string(buf[0:n]) == TFirewall.KeySend {
				fmt.Printf("RecvUDP On %s From %s:%d\n", conn.LocalAddr().String(), radder.IP, radder.Port)
			}
		}
		_, err = conn.WriteToUDP([]byte(TFirewall.KeySend), radder)
		//checkError(err)
	}
}

func listenTCP(listener *net.TCPListener) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: %s", err.Error())
		}
		defer conn.Close()
		var buf [20]byte
		n, err := conn.Read(buf[0:])
		if err == nil {
			if string(buf[0:n]) == TFirewall.KeySend {
				fmt.Printf("RecvTCP On %s From %s \n", conn.LocalAddr().String(), conn.RemoteAddr().String())
			}
		}
		_, err = conn.Write([]byte(TFirewall.KeySend))
	}

}

func main() {
	portslice := []int{}
	if len(os.Args) > 1 {
		portListStrs := strings.Split(os.Args[1], ",")
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
	fmt.Println("Server listening: ", portslice)
	for _, port := range portslice {
		tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
		checkError(err)
		tcpListener, err := net.ListenTCP("tcp", tcpAddr)
		defer tcpListener.Close()
		go listenTCP(tcpListener)
		checkError(err)
	}

	for _, port := range portslice {
		udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
		checkError(err)
		udpconn, err := net.ListenUDP("udp", udpAddr)
		defer udpconn.Close()
		checkError(err)
		go recvUDP(udpconn)
	}

	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)
	//阻塞直到有信号传入
	s := <-c
	fmt.Println("exit : ", s)
}
