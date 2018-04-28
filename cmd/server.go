package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start rfunc server",
	RunE:  runServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	lis, err := server.NewServerConn(globalOpts.Network(), globalOpts.Address())
	if err != nil {
		return err
	}

	logger.Println("Server started at", lis.Addr())
	s := &server.Server{}
	return s.Serve(lis)
}
