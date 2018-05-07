package cmd

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
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), func(rfunc *client.RFunc) error {
		text, err := rfunc.Paste()
		if err != nil {
			return err
		}

		os.Stdout.WriteString(text)
		return nil
	})
}
