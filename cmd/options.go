package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/pflag"
	"github.com/yukithm/rfunc/utils"
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
	case "PASS", EOL_PASS:
		return "PASS"
	case "LF", EOL_LF:
		return "LF"
	case "CRLF", EOL_CRLF:
		return "CRLF"
	case EOL_NATIVE:
		return "NATIVE"
	}
	return string(*e)
}

func (e *LineEnding) Type() string {
	return "string"
}

func (e *LineEnding) Set(value string) error {
	switch strings.ToUpper(value) {
	case "PASS", "":
		*e = EOL_PASS
	case "LF", "\n":
		*e = EOL_LF
	case "CRLF", "\r\n":
		*e = EOL_CRLF
	case "NATIVE":
		*e = EOL_NATIVE
	default:
		return fmt.Errorf("Unsupported value '%s'", *e)
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

type Options struct {
	Addr    string        `toml:"addr"`
	Sock    string        `toml:"sock"`
	Logfile string        `toml:"logfile"`
	Quiet   bool          `toml:"quiet"`
	EOL     LineEnding    `toml:"eol"`
	Server  ServerOptions `toml:"server"`
}

func (o *Options) Clone() *Options {
	newOpts := &Options{}
	*newOpts = *o
	newOpts.Server = *o.Server.Clone()
	return newOpts
}

func (o *Options) Merge(other *Options) {
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

type FlagOptions struct {
	*Options
	GlobalFlags *pflag.FlagSet
	ServerFlags *pflag.FlagSet
}

func MergeFlagOptions(opts *Options, fopts *FlagOptions) *Options {
	o := opts.Clone()

	if fopts.GlobalFlags.Changed("addr") || o.Addr == "" {
		o.Addr = fopts.Addr
	}
	if fopts.GlobalFlags.Changed("sock") || o.Sock == "" {
		o.Sock = fopts.Sock
	}
	if fopts.GlobalFlags.Changed("logfile") || o.Logfile == "" {
		o.Logfile = fopts.Logfile
	}
	if fopts.GlobalFlags.Changed("quiet") || o.Quiet == false {
		o.Quiet = fopts.Quiet
	}
	if fopts.GlobalFlags.Changed("eol") || o.EOL == "" {
		o.EOL = fopts.EOL
	}

	if fopts.ServerFlags.Changed("daemon") || o.Server.Daemon == false {
		o.Server.Daemon = fopts.Server.Daemon
	}
	if fopts.ServerFlags.Changed("allow-commands") || o.Server.AllowCmds == nil || len(o.Server.AllowCmds) == 0 {
		o.Server.AllowCmds = make([]string, len(fopts.Server.AllowCmds))
		copy(o.Server.AllowCmds, fopts.Server.AllowCmds)
	}

	return o
}

var allCommands = []string{
	"copy", "paste", "open",
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

func (o *ServerOptions) AllowCommands() []string {
	if len(o.AllowCmds) == 0 {
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

func LoadConfig(file string) (*Options, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Options
	decoder := toml.NewDecoder(f)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	if err := conf.EOL.Set(string(conf.EOL)); err != nil {
		return nil, fmt.Errorf("eol option error: %s in %s", err.Error(), file)
	}
	return &conf, nil
}
