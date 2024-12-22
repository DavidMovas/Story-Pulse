package net

import (
	"fmt"
	"net"
)

func CheckPortAvailable(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}

	_ = listener.Close()
	return true
}

func FindFreePort(host string, startPort int) (int, error) {
	port := startPort
	for {
		if CheckPortAvailable(host, port) {
			return port, nil
		}

		port++
		if port > 65535 {
			return 0, fmt.Errorf("port %d is out of range", port)
		}
	}
}