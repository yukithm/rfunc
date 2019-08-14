package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/client"
)

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste text from server's clipboard",
	RunE:  runPasteCmd,
}

func runPasteCmd(cmd *cobra.Command, args []string) error {
	config := &client.Config{
		EOL: globalOpts.EOLCode(),
		TLS: &globalOpts.TLS,
	}
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), config, func(rfunc *client.RFunc) error {
		text, err := rfunc.Paste()
		if err != nil {
			return err
		}

		os.Stdout.WriteString(text)
		return nil
	})
}
