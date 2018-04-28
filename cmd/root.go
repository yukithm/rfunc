package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	Addr    string
	Sock    string
	Logfile string
	NoLog   bool
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

var globalOpts = GlobalOptions{
	Addr: "127.0.0.1:8299",
}

var logger *log.Logger
var logfile io.WriteCloser

var rootCmd = &cobra.Command{
	Use:           "rfunc",
	Short:         "rfunc is a utility functions over the network",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		lf := cmd.Flag("logfile")
		logger, err = newLogger(globalOpts.Logfile, lf.Changed)
		return
	},
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&globalOpts.Addr, "addr", globalOpts.Addr, "address and port")
	pf.StringVar(&globalOpts.Sock, "sock", globalOpts.Sock, "unix domain socket path")
	pf.StringVar(&globalOpts.Logfile, "logfile", globalOpts.Logfile, "logfile")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if logger != nil {
			logger.Print(err)
		}
		os.Exit(1)
	}
}

func newLogger(logfile string, explicit bool) (*log.Logger, error) {
	switch logfile {
	case "":
		if explicit {
			return log.New(ioutil.Discard, "", log.LstdFlags), nil
		}
		return log.New(os.Stderr, "", log.LstdFlags), nil

	case "-":
		return log.New(os.Stdout, "", log.LstdFlags), nil

	default:
		file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		return log.New(file, "", log.LstdFlags), nil
	}
}
