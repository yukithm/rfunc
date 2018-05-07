package cmd

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type GlobalOptions struct {
	Addr    string `json:"addr"`
	Sock    string `json:"sock"`
	Logfile string `json:"logfile"`
	Quiet   bool   `json:"quiet"`
}

func (o *GlobalOptions) Clone() *GlobalOptions {
	newOpts := &GlobalOptions{}
	*newOpts = *o
	return newOpts
}

func (o *GlobalOptions) Merge(other *GlobalOptions) {
	if other.Addr != "" {
		o.Addr = other.Addr
	}

	if other.Sock != "" {
		o.Sock = other.Sock
	}

	if other.Logfile != "" {
		o.Logfile = other.Logfile
	}

	o.Quiet = other.Quiet
}

func (o *GlobalOptions) AbsPaths() {
	if o.Logfile != "" && o.Logfile != "-" {
		o.Logfile = o.abs(o.Logfile)
	}
	o.Sock = o.abs(o.Sock)
}

func (o *GlobalOptions) abs(path string) string {
	if path == "" {
		return ""
	}
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

type FlagOptions struct {
	*GlobalOptions
	Flags *pflag.FlagSet
}

func MergeFlagOptions(opts *GlobalOptions, fopts *FlagOptions) *GlobalOptions {
	o := opts.Clone()

	if fopts.Flags.Changed("addr") || o.Addr == "" {
		o.Addr = fopts.Addr
	}
	if fopts.Flags.Changed("sock") || o.Sock == "" {
		o.Sock = fopts.Sock
	}
	if fopts.Flags.Changed("logfile") || o.Logfile == "" {
		o.Logfile = fopts.Logfile
	}
	if fopts.Flags.Changed("quiet") || o.Quiet == false {
		o.Quiet = fopts.Quiet
	}

	return o
}

func LoadConfig(file string) (*GlobalOptions, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf GlobalOptions
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
