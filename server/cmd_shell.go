// +build linux darwin

package server

import (
	"errors"
	"os/exec"
	"runtime"
)

var DefaultShell = map[string][]CmdShell{
	"darwin": {
		{
			openCmd: OpenCmd{path: "open"},
		},
	},
	"linux": {
		{
			openCmd: OpenCmd{path: "xdg-open"},
		},
	},
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

func GetShell() (Shell, error) {
	if defs, ok := DefaultShell[runtime.GOOS]; ok {
		return findShell(defs)
	}

	return nil, errors.New("Unsupported OS")
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

	return nil, errors.New("no available open commands")
}
