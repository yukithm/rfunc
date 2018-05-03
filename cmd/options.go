package cmd

import (
	"net"
	"path/filepath"
)

type GlobalOptions struct {
	Addr    string
	Sock    string
	Logfile string
}

func (o *GlobalOptions) AbsPaths() {
	if o.Logfile != "" && o.Logfile != "-" {
		o.Logfile = o.abs(o.Logfile)
	}
	o.Sock = o.abs(o.Sock)
}

func (o *GlobalOptions) abs(path string) string {
	if ap, err := filepath.Abs(path); err == nil {
		return ap
	}
	return path
}

func (o *GlobalOptions) Network() string {
	if o.Sock != "" {
		return "unix"
	}
	return "tcp"
}

func (o *GlobalOptions) Address() string {
	if o.Sock != "" {
		return o.Sock
	}
	return o.Addr
}

func (o *GlobalOptions) NetAddr() (net.Addr, error) {
	if o.Sock != "" {
		return net.ResolveUnixAddr("unix", o.Sock)
	}
	return net.ResolveTCPAddr("tcp", o.Addr)
}
