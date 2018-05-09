package cmd

import (
	"github.com/VividCortex/godaemon"
	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start rfunc server",
	RunE:  runServerCmd,
}

var serverOpts = struct {
	Daemon bool
}{}

func init() {
	f := serverCmd.Flags()
	f.BoolVar(&serverOpts.Daemon, "daemon", serverOpts.Daemon, "daemonize")
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	if serverOpts.Daemon {
		godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	}

	lis, err := server.NewServerConn(globalOpts.Network(), globalOpts.Address())
	if err != nil {
		return err
	}
	defer lis.Close()

	logger.Println("Server started at", lis.Addr())
	defer logger.Println("Server stopped")
	s := &server.Server{
		Config: &server.Config{
			EOL: globalOpts.EOL.Code(),
		},
		Logger: logger,
	}
	return s.Serve(lis)
}
