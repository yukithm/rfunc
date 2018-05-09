package cmd

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
		EOL: globalOpts.EOL.Code(),
	}
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), config, func(rfunc *client.RFunc) error {
		return rfunc.OpenURL(args...)
	})
}
