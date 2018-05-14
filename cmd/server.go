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

func init() {
	f := serverCmd.Flags()
	f.BoolVar(&globalOpts.Server.Daemon, "daemon", globalOpts.Server.Daemon, "daemonize")
	f.StringSliceVar(&globalOpts.Server.AllowCmds, "allow-commands", nil, "allow only specified commands")
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	if globalOpts.Server.Daemon {
		godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	}

	lis, err := server.NewServerConn(globalOpts.Network(), globalOpts.Address())
	if err != nil {
		return err
	}
	defer lis.Close()

	logger.Println("Server started at", lis.Addr())
	defer logger.Println("Server stopped")
	s := &server.RFuncServer{
		Config: &server.Config{
			EOL:       globalOpts.EOLCode(),
			AllowCmds: globalOpts.Server.AllowCommands(),
			TLS:       &globalOpts.TLS,
		},
		Logger: logger,
	}
	return s.Serve(lis)
}
