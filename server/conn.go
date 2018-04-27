package server

import (
	"net"
	"os"
)

func NewServerConn(network, addr string) (net.Listener, error) {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	if network == "unix" {
		if err := os.Chmod(addr, 0700); err != nil {
			lis.Close()
			return nil, err
		}
	}

	return lis, nil
}
