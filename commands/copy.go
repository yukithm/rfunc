package commands

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

func runCopyCmd(cmd *cobra.Command, args []string) error {
	config := &client.Config{
		EOL: globalOpts.EOLCode(),
		TLS: &globalOpts.TLS,
	}
	return client.RunRFunc(globalOpts.Network(), globalOpts.Address(), config, func(rfunc *client.RFunc) error {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, os.Stdin); err != nil {
			return err
		}
		return rfunc.Copy(buf.String())
	})
}
