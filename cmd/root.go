package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yukithm/rfunc/utils"
)

var defaultConfigFiles = []string{
	"~/.config/rfunc/rfunc.toml",
	"~/.rfunc.toml",
}
var configfile string
var configOpts *GlobalOptions

var flagOpts = &FlagOptions{
	GlobalOptions: &GlobalOptions{
		Addr: "127.0.0.1:8299",
		EOL:  "NATIVE",
	},
}

var globalOpts *GlobalOptions

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

	rootCmd.Flags().SortFlags = false
	pf := rootCmd.PersistentFlags()
	pf.SortFlags = false
	pf.StringVarP(&configfile, "conf", "c", configfile, "configuration file")
	pf.StringVarP(&flagOpts.Addr, "addr", "a", flagOpts.Addr, "address and port")
	pf.StringVarP(&flagOpts.Sock, "sock", "s", flagOpts.Sock, "unix domain socket path")
	pf.StringVarP(&flagOpts.Logfile, "logfile", "l", flagOpts.Logfile, "logfile")
	pf.BoolVarP(&flagOpts.Quiet, "quiet", "q", flagOpts.Quiet, "suppress outputs (except paste content")
	pf.Var(&flagOpts.EOL, "eol", "line ending (LF|CRLF|NATIVE|PASS)")
	flagOpts.Flags = pf

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(serverCmd)
}

func Execute() (code int) {
	defer func() {
		if logdev != nil {
			logdev.Close()
		}
	}()

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
	var confErr error
	if cmd.Flag("conf").Changed {
		configOpts, confErr = loadConfig(configfile)
	} else {
		configOpts, confErr = loadDefaultConfig()
	}

	if confErr != nil {
		configOpts = &GlobalOptions{}
	}
	globalOpts = MergeFlagOptions(configOpts, flagOpts)
	globalOpts.AbsPaths()

	logdev, err = newLogDevice(globalOpts, cmd.Flags())
	logger = log.New(logdev, "", log.LstdFlags)
	cmd.Root().SetOutput(logdev)

	if confErr != nil {
		return confErr
	}

	return
}

func newLogDevice(opts *GlobalOptions, flags *pflag.FlagSet) (io.WriteCloser, error) {
	// Keep logging to the file even if specified quiet option
	if opts.Quiet && (opts.Logfile == "" || opts.Logfile == "-") {
		return utils.NopWriteCloser(ioutil.Discard), nil
	}

	switch opts.Logfile {
	case "":
		if flags.Changed("logfile") {
			return utils.NopWriteCloser(ioutil.Discard), nil
		}
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

func loadConfig(conf string) (*GlobalOptions, error) {
	if conf == "" {
		return &GlobalOptions{}, nil
	}

	path, err := utils.ExpandPath(conf)
	if err != nil {
		return nil, err
	}
	return LoadConfig(path)
}

func loadDefaultConfig() (*GlobalOptions, error) {
	for _, file := range defaultConfigFiles {
		path, err := utils.ExpandPath(file)
		if err != nil {
			return nil, err
		}

		opts, err := LoadConfig(path)
		if err == nil {
			return opts, nil
		}

		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return &GlobalOptions{}, nil
}
