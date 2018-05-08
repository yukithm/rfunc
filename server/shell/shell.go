package shell

import (
	"errors"
	"os/exec"
)

var (
	ErrUnsupported = errors.New("shell functions are not supported")
	ErrCmdNotFound = errors.New("No available commands")
)

type Shell interface {
	OpenURL(url ...string) error
}

type OpenCmd struct {
	path string
	args []string
}

func (c *OpenCmd) Run(args []string) error {
	cargs := append([]string{}, c.args...)
	cargs = append(cargs, args...)
	cmd := exec.Command(c.path, cargs...)
	return cmd.Run()
}

type CmdShell struct {
	openCmd OpenCmd
}

func (c *CmdShell) OpenURL(url ...string) error {
	return c.openCmd.Run(url)
}

func findShell(candidates []CmdShell) (*CmdShell, error) {
	for _, candidate := range candidates {
		var path string
		var err error

		if path, err = exec.LookPath(candidate.openCmd.path); err != nil {
			continue
		}

		return &CmdShell{
			openCmd: OpenCmd{path, candidate.openCmd.args},
		}, nil
	}

	return nil, ErrCmdNotFound
}
