package cmd

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
)

type LineEnding string

const (
	EOL_PASS   = LineEnding("")
	EOL_LF     = LineEnding("\n")
	EOL_CRLF   = LineEnding("\r\n")
	EOL_NATIVE = LineEnding("NATIVE")
)

func (e *LineEnding) String() string {
	switch *e {
	case EOL_PASS:
		return "PASS"
	case EOL_LF:
		return "LF"
	case EOL_CRLF:
		return "CRLF"
	case EOL_NATIVE:
		return "NATIVE"
	}
	return "Invalid LineEnding value!"
}

func (e *LineEnding) Type() string {
	return "string"
}

func (e *LineEnding) Set(value string) error {
	switch strings.ToUpper(value) {
	case "PASS":
		*e = EOL_PASS
	case "LF":
		*e = EOL_LF
	case "CRLF":
		*e = EOL_CRLF
	case "NATIVE":
		*e = EOL_NATIVE
	default:
		return errors.New("Unsupported")
	}

	return nil
}

func (e *LineEnding) Code() string {
	if *e != EOL_NATIVE {
		return string(*e)
	}

	switch runtime.GOOS {
	case "windows":
		return string(EOL_CRLF)
	default:
		return string(EOL_LF)
	}
}

type GlobalOptions struct {
	Addr    string     `json:"addr"`
	Sock    string     `json:"sock"`
	Logfile string     `json:"logfile"`
	Quiet   bool       `json:"quiet"`
	EOL     LineEnding `json:"eol"`
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
	o.EOL = other.EOL
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
	if fopts.Flags.Changed("eol") || o.EOL == "" {
		o.EOL = fopts.EOL
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
