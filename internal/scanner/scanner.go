package scanner

import (
	"fmt"
	"net"
	"time"
)

func GeneratePorts(initialPort int, finalPort int) []int {
	var ports []int
	for i := initialPort; i <= finalPort; i++ {
		ports = append(ports, i)
	}
	return ports
}

func Scan(address string, portsToScan []int, timeout int, channel chan int) {
	for i := 0; i < len(portsToScan); i++ {
		port := portsToScan[i]
		target := fmt.Sprintf("%s:%d", address, port)
		conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
		if err != nil {
			channel <- 0
			continue
		}
		conn.Close()
		channel <- port
	}
}
