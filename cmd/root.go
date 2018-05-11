package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/options"
	"github.com/yukithm/rfunc/utils"
)

var defaultOpts = &options.Options{
	Addr: "127.0.0.1:8299",
	EOL:  "NATIVE",
}
var globalOpts = &options.Options{}

var logger *log.Logger
var logdev io.WriteCloser

var rootCmd = &cobra.Command{
	Use:               "rfunc",
	Short:             "rfunc is utility functions over the network",
	Version:           Version,
	SilenceErrors:     false,
	SilenceUsage:      true,
	PersistentPreRunE: initApp,
}

func init() {
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(serverCmd)
}

func initFlags() {
	rootCmd.Flags().SortFlags = false
	pf := rootCmd.PersistentFlags()
	pf.SortFlags = false

	pf.StringP("conf", "c", options.ConfigFile, "configuration file")
	pf.StringVarP(&globalOpts.Addr, "addr", "a", globalOpts.Addr, "address and port")
	pf.StringVarP(&globalOpts.Sock, "sock", "s", globalOpts.Sock, "unix domain socket path")
	pf.StringVarP(&globalOpts.Logfile, "logfile", "l", globalOpts.Logfile, "logfile")
	pf.BoolVarP(&globalOpts.Quiet, "quiet", "q", globalOpts.Quiet, "suppress outputs (except paste content)")
	pf.StringVar(&globalOpts.EOL, "eol", globalOpts.EOL, "line ending (LF|CRLF|NATIVE|PASS)")

	// TLS options
	pf.StringVar(&globalOpts.TLS.CertFile, "tls-cert", globalOpts.TLS.CertFile, "Certificate file")
	pf.StringVar(&globalOpts.TLS.KeyFile, "tls-key", globalOpts.TLS.KeyFile, "Private key file")
	pf.StringVar(&globalOpts.TLS.CAFile, "tls-ca", globalOpts.TLS.CAFile, "CA Certificate file")
}

func Execute() (code int) {
	defer func() {
		if logdev != nil {
			logdev.Close()
		}
	}()

	configOpts, err := options.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	globalOpts.Fill(configOpts)
	globalOpts.Fill(defaultOpts)

	initFlags()

	prognameSwitch()
	if err := rootCmd.Execute(); err != nil {
		// error message is printed by command library
		return 1
	}

	return 0
}

func prognameSwitch() {
	var cmd string
	progname := strings.ToLower(filepath.Base(os.Args[0]))
	if strings.Index(progname, "open") != -1 {
		cmd = "open"
	} else if strings.Index(progname, "copy") != -1 {
		cmd = "copy"
	} else if strings.Index(progname, "paste") != -1 {
		cmd = "paste"
	}

	var args []string
	if cmd == "" {
		args = os.Args[1:]
	} else {
		args = []string{cmd}
		if len(os.Args) > 1 {
			args = append(args, os.Args[1:]...)
		}
	}

	rootCmd.SetArgs(args)
}

func initApp(cmd *cobra.Command, args []string) (err error) {
	globalOpts.AbsPaths()

	logdev, err = newLogDevice(globalOpts)
	logger = log.New(logdev, "", log.LstdFlags)
	cmd.Root().SetOutput(logdev)

	return
}

func newLogDevice(opts *options.Options) (io.WriteCloser, error) {
	// Keep logging to the file even if specified quiet option
	if opts.Quiet && (opts.Logfile == "" || opts.Logfile == "-") {
		return utils.NopWriteCloser(ioutil.Discard), nil
	}

	switch opts.Logfile {
	case "":
		return utils.NopWriteCloser(os.Stderr), nil

	case "-":
		return utils.NopWriteCloser(os.Stdout), nil

	default:
		file, err := os.OpenFile(opts.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
