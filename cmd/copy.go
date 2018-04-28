package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/client"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy text to server's clipboard",
	RunE:  runCopyCmd,
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func runCopyCmd(cmd *cobra.Command, args []string) error {
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), func(rfunc *client.RFunc) error {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, os.Stdin); err != nil {
			return err
		}
		return rfunc.Copy(buf.String())
	})
}
