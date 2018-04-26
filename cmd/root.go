package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	Addr    string
	Sock    string
	Logfile string
}

var options = Options{
	Addr: "127.0.0.1:3299",
}

var rootCmd = &cobra.Command{
	Use:   "rfunc",
	Short: "rfunc is a utility functions over the network",
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&options.Addr, "addr", options.Addr, "address and port")
	pf.StringVar(&options.Sock, "sock", options.Sock, "unix domain socket path")
	pf.StringVar(&options.Logfile, "logfile", options.Logfile, "logfile")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
}
