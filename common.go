package TFirewall

const KeySend string = "HelloFarmer"
const KeyRecv string = "HelloCoder"

func TcpPorts() []int {
	var TCPPORTS = []int{21, 25, 53, 80, 443, 110, 465, 3389, 8080, 8443}
	return TCPPORTS
}
func UdpPorts() []int {
	var UDPPORTS = []int{21, 25, 53, 80, 443, 110, 465, 3389, 8080, 8443}
	return UDPPORTS
}
func Contain(obj []int, target int) bool {
	for _, port := range obj {
		if port == target {
			return true
		}
	}
	return false
}
