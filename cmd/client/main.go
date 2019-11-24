package main

import (
	"TFirewall"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}

func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		return
	}

	if b[0] == 0x05 { //只处理Socks5协议
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var addr string
		switch b[3] {
		case 0x01:
			sip := sockIP{}
			if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
				log.Println("请求解析错误")
				return
			}
			addr = sip.toAddr()
		case 0x03:
			host := string(b[5 : n-2])
			var port uint16
			err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
			if err != nil {
				log.Println(err)
				return
			}
			addr = fmt.Sprintf("%s:%d", host, port)
		}

		server, err := net.Dial("tcp", addr)
		if err != nil {
			return
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}

}

func connect_socks5_server(serverip string, serverport string) {
	var RemoteConn net.Conn
	var err error
	for {
		for {
			RemoteConn, err = net.Dial("tcp", serverip+":"+serverport)
			if err == nil {
				break
			}
		}
		go handleClientRequest(RemoteConn)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage:\n client check 192.168.1.100\n client check 192.168.1.100 20-22,80-90,22")
		fmt.Println("usage:\n client socks5 192.168.1.100 80")
		os.Exit(1)
	}

	if strings.EqualFold(os.Args[1], "check") {
		portslice := []int{}
		if len(os.Args) > 3 {
			portListStrs := strings.Split(os.Args[3], ",")
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
		ip := os.Args[2]
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
		fmt.Println("check finish")
	} else if strings.EqualFold(os.Args[1], "socks5") {
		if len(os.Args) > 3 {
			connect_socks5_server(os.Args[2], os.Args[3])
		} else {
			fmt.Println("usage:\n client socks5 192.168.1.100 80")
		}
	}
}
