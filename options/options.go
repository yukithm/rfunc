package options

import (
	"net"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yukithm/rfunc/utils"
)

type Options struct {
	Addr    string        `toml:"addr"`
	Sock    string        `toml:"sock"`
	Logfile string        `toml:"logfile"`
	Quiet   bool          `toml:"quiet"`
	EOL     string        `toml:"eol"`
	Server  ServerOptions `toml:"server"`
}

func (o *Options) Clone() *Options {
	newOpts := &Options{}
	*newOpts = *o
	newOpts.Server = *o.Server.Clone()
	return newOpts
}

func (o *Options) Fill(other *Options) {
	if o.Addr == "" {
		o.Addr = other.Addr
	}

	if o.Sock == "" {
		o.Sock = other.Sock
	}

	if o.Logfile == "" {
		o.Logfile = other.Logfile
	}

	if o.Quiet == false {
		o.Quiet = other.Quiet
	}

	if o.EOL == "" {
		o.EOL = other.EOL
	}

	o.Server.Fill(&other.Server)
}

func (o *Options) AbsPaths() {
	if o.Logfile != "" && o.Logfile != "-" {
		o.Logfile = o.abs(o.Logfile)
	}
	o.Sock = o.abs(o.Sock)
}

func (o *Options) abs(path string) string {
	if path == "" {
		return ""
	}
	if ap, err := filepath.Abs(path); err == nil {
		return ap
	}
	return path
}

func (o *Options) Network() string {
	if o.Sock != "" {
		return "unix"
	}
	return "tcp"
}

func (o *Options) Address() string {
	if o.Sock != "" {
		return o.Sock
	}
	return o.Addr
}

func (o *Options) NetAddr() (net.Addr, error) {
	if o.Sock != "" {
		return net.ResolveUnixAddr("unix", o.Sock)
	}
	return net.ResolveTCPAddr("tcp", o.Addr)
}

var allCommands = []string{
	"copy", "paste", "open",
}

func (o *Options) EOLCode() string {
	switch strings.ToUpper(o.EOL) {
	case "", "PASS":
		return ""

	case "CR", "\r":
		return "\r"

	case "LF", "\n":
		return "\n"

	case "CRLF", "\r\n":
		return "\r\n"

	case "NATIVE":
		if runtime.GOOS == "windows" {
			return "\r\n"
		}
		return "\n"

	default:
		return ""
	}
}

type ServerOptions struct {
	Daemon    bool     `toml:"daemon"`
	AllowCmds []string `toml:"allow-commands"`
}

func (o *ServerOptions) Clone() *ServerOptions {
	newOpts := &ServerOptions{}
	*newOpts = *o
	if o.AllowCmds != nil {
		newOpts.AllowCmds = make([]string, len(o.AllowCmds))
		copy(newOpts.AllowCmds, o.AllowCmds)
	}
	return newOpts
}

func (o *ServerOptions) Fill(other *ServerOptions) {
	if o.Daemon == false {
		o.Daemon = other.Daemon
	}

	if o.AllowCmds == nil || len(o.AllowCmds) == 0 {
		if other.AllowCmds != nil {
			o.AllowCmds = make([]string, len(other.AllowCmds))
			copy(o.AllowCmds, other.AllowCmds)
		} else {
			o.AllowCmds = make([]string, 0)
		}
	}
}

func (o *ServerOptions) AllowCommands() []string {
	if o.AllowCmds == nil || len(o.AllowCmds) == 0 {
		return allCommands[:]
	}

	res := []string{}
	for _, name := range o.AllowCmds {
		name = strings.ToLower(name)
		if name == "all" {
			return allCommands[:]
		} else if utils.FindString(allCommands, name) != -1 && utils.FindString(res, name) == -1 {
			res = append(res, name)
		}
	}

	return res
}
