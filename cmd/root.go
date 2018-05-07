package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var defaultConfigFiles = []string{
	"~/.config/rfunc/rfunc.json",
	"~/.rfunc.json",
}
var configfile string
var configOpts *GlobalOptions

var flagOpts = &FlagOptions{
	GlobalOptions: &GlobalOptions{
		Addr: "127.0.0.1:8299",
	},
}

var globalOpts *GlobalOptions

var logger *log.Logger
var logfile io.WriteCloser

var rootCmd = &cobra.Command{
	Use:               "rfunc",
	Short:             "rfunc is a utility functions over the network",
	SilenceErrors:     true,
	SilenceUsage:      true,
	PersistentPreRunE: initApp,
}

func init() {
	cobra.EnableCommandSorting = false

	rootCmd.Flags().SortFlags = false
	pf := rootCmd.PersistentFlags()
	pf.SortFlags = false
	pf.StringVar(&configfile, "conf", configfile, "configuration file")
	pf.StringVar(&flagOpts.Addr, "addr", flagOpts.Addr, "address and port")
	pf.StringVar(&flagOpts.Sock, "sock", flagOpts.Sock, "unix domain socket path")
	pf.StringVar(&flagOpts.Logfile, "logfile", flagOpts.Logfile, "logfile")
	flagOpts.Flags = pf

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(serverCmd)
}

func Execute() {
	prognameSwitch()
	if err := rootCmd.Execute(); err != nil {
		if logger != nil {
			logger.Print(err)
		}
		os.Exit(1)
	}
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

	logger, err = newLogger(globalOpts.Logfile, cmd.Flag("logfile").Changed)
	if confErr != nil {
		return confErr
	}

	return
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

func loadConfig(conf string) (*GlobalOptions, error) {
	if conf == "" {
		return &GlobalOptions{}, nil
	}

	path, err := ExpandPath(conf)
	if err != nil {
		return nil, err
	}
	return LoadConfig(path)
}

func loadDefaultConfig() (*GlobalOptions, error) {
	for _, file := range defaultConfigFiles {
		path, err := ExpandPath(file)
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
