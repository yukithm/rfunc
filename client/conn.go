package client

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

func NewClientConn(network, addr string) (*grpc.ClientConn, error) {
	if network == "unix" {
		return grpc.Dial(
			addr,
			grpc.WithInsecure(),
			grpc.WithDialer(unixDomainSocketDialer),
		)
	}

	return grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
}

func unixDomainSocketDialer(addr string, t time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, t)
}
