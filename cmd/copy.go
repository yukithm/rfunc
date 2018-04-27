package cmd

import "github.com/spf13/cobra"

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy text to server's clipboard",
	RunE:  runCopyCmd,
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func runCopyCmd(cmd *cobra.Command, args []string) error {
	return nil
}
