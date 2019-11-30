package main

import (
	"TFirewall"
	"crypto/tls"
	"fmt"
	"io"
	"log"
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
func Server(listen *net.TCPListener, s5listen *net.TCPListener) {
	for {
		s5conn, err := s5listen.Accept()
		if err != nil {
			fmt.Println("Error on accept socks5 connect : ", err.Error())
			continue
		}
		fmt.Println("Socks5 new socket from : ", s5conn.RemoteAddr().String())
		defer s5conn.Close()

		controlconn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error on accept control connect :", err.Error())
			continue
		}
		fmt.Println("Control new socket from : ", controlconn.RemoteAddr().String())
		defer controlconn.Close()

		go handle(controlconn, s5conn)
	}
}
func TLSServer(listen net.Listener, s5listen *net.TCPListener) {
	for {
		s5conn, err := s5listen.Accept()
		if err != nil {
			fmt.Println("Error on accept socks5 connect : ", err.Error())
			continue
		}
		fmt.Println("Socks5 new socket from : ", s5conn.RemoteAddr().String())
		defer s5conn.Close()

		controlconn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error on accept control connect :", err.Error())
			continue
		}
		fmt.Println("Control new socket from : ", controlconn.RemoteAddr().String())
		defer controlconn.Close()

		go handle(controlconn, s5conn)
	}
}

func handle(sconn net.Conn, dconn net.Conn) {
	defer sconn.Close()
	defer dconn.Close()
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	dconn.Close()
}

func ErrHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage:\n server check\n server check 20-22,80-90,22")
		fmt.Println("usage:\n server socks5 80 1080")
		os.Exit(1)
	}

	if strings.EqualFold(os.Args[1], "check") {
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
		fmt.Println("Check Server listening: ", portslice)
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
	} else if strings.EqualFold(os.Args[1], "socks5") {
		if len(os.Args) == 4 {
			var ip = "0.0.0.0"
			port, _ := strconv.Atoi(os.Args[2])
			s5port, _ := strconv.Atoi(os.Args[3])

			lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
			ErrHandler(err)
			defer lis.Close()
			fmt.Println("Control Listening: ", port)
			s5lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), s5port, ""})
			ErrHandler(err)
			defer s5lis.Close()
			fmt.Println("Socks5 Listening: ", s5port)

			Server(lis, s5lis)
		} else if len(os.Args) == 5 {
			var ip = "0.0.0.0"
			port, _ := strconv.Atoi(os.Args[2])
			s5port, _ := strconv.Atoi(os.Args[3])
			log.Printf(os.Args[0])
			cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
			if err != nil {
				log.Println(err)
				return
			}
			config := &tls.Config{Certificates: []tls.Certificate{cert}}

			lis, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", ip, port), config)
			if err != nil {
				log.Println(err)
				return
			}

			ErrHandler(err)
			defer lis.Close()
			fmt.Println("Control Listening: ", port)
			s5lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), s5port, ""})
			ErrHandler(err)
			defer s5lis.Close()
			fmt.Println("Socks5 Listening: ", s5port)

			TLSServer(lis, s5lis)

		} else {
			fmt.Println("usage:\n server socks5 80 1080 (80 is control port,1080 is socks5 port)")
		}
	}
}
