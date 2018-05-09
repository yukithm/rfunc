package cmd

import (
	"strings"

	"github.com/VividCortex/godaemon"
	"github.com/spf13/cobra"
	"github.com/yukithm/rfunc/server"
	"github.com/yukithm/rfunc/utils"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start rfunc server",
	RunE:  runServerCmd,
}

var allCommands = []string{
	"copy", "paste", "open",
}

type ServerOpts struct {
	Daemon    bool
	AllowCmds []string
}

func (o *ServerOpts) AllowCommands() []string {
	if len(o.AllowCmds) == 0 {
		return allCommands[:]
	}

	res := []string{}
	for _, name := range o.AllowCmds {
		name = strings.ToLower(name)
		if name == "all" {
			return allCommands[:]
		} else if utils.FindString(allCommands, name) != -1 && utils.FindString(res, name) == -1 {
			res = append(res, name)
		}
	}

	return res
}

var serverOpts = ServerOpts{}

func init() {
	f := serverCmd.Flags()
	f.BoolVar(&serverOpts.Daemon, "daemon", serverOpts.Daemon, "daemonize")
	f.StringSliceVar(&serverOpts.AllowCmds, "allow-commands", nil, "allow only specified commands")
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
			EOL:       globalOpts.EOL.Code(),
			AllowCmds: serverOpts.AllowCommands(),
		},
		Logger: logger,
	}
	return s.Serve(lis)
}
