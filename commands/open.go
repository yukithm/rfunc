package commands

import (
	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/client"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open specified URLs by the server's browser",
	RunE:  runOpenCmd,
}

func runOpenCmd(cmd *cobra.Command, args []string) error {
	config := &client.Config{
		EOL: globalOpts.EOLCode(),
		TLS: &globalOpts.TLS,
	}
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), config, func(rfunc *client.RFunc) error {
		return rfunc.OpenURL(args...)
	})
}
