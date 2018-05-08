package clipboard

import (
	"bytes"
	"errors"
	"os/exec"
)

var (
	ErrUnsupported = errors.New("Clipboard is not supported")
	ErrCmdNotFound = errors.New("No available clipboard commands")
)

type Clipboard interface {
	CopyText(text string) error
	PasteText() (string, error)
}

type Cmd struct {
	path string
	args []string
}

func (c *Cmd) Run(input []byte) ([]byte, error) {
	cmd := exec.Command(c.path, c.args...)
	if input != nil && len(input) > 0 {
		cmd.Stdin = bytes.NewReader(input)
	}
	return cmd.Output()
}

type CmdClipboard struct {
	copyCmd  Cmd
	pasteCmd Cmd
}

func (c *CmdClipboard) CopyText(text string) error {
	_, err := c.copyCmd.Run([]byte(text))
	return err
}

func (c *CmdClipboard) PasteText() (string, error) {
	out, err := c.pasteCmd.Run(nil)
	return string(out), err
}

func findCmdClipboard(candidates []CmdClipboard) (*CmdClipboard, error) {
	for _, candidate := range candidates {
		var cpath, ppath string
		var err error

		if cpath, err = exec.LookPath(candidate.copyCmd.path); err != nil {
			continue
		}
		if ppath, err = exec.LookPath(candidate.pasteCmd.path); err != nil {
			continue
		}

		return &CmdClipboard{
			copyCmd:  Cmd{cpath, candidate.copyCmd.args},
			pasteCmd: Cmd{ppath, candidate.pasteCmd.args},
		}, nil
	}

	return nil, ErrCmdNotFound
}
