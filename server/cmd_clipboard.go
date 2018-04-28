// +build linux darwin

package server

import (
	"errors"
	"os/exec"
	"runtime"
)

var DefaultClipboard = map[string][]CmdClipboard{
	"darwin": {
		{
			copyCmd:  ClipCmd{path: "pbcopy"},
			pasteCmd: ClipCmd{path: "pbpaste"},
		},
	},
	"linux": {
		{
			copyCmd:  ClipCmd{path: "xclip", args: []string{"-in", "-selection", "clipboard"}},
			pasteCmd: ClipCmd{path: "xclip", args: []string{"-out", "-selection", "clipboard"}},
		},
		{
			copyCmd:  ClipCmd{path: "xsel", args: []string{"--input", "--clipboard"}},
			pasteCmd: ClipCmd{path: "xsel", args: []string{"--output", "--clipboard"}},
		},
	},
}

type ClipCmd struct {
	path string
	args []string
}

type CmdClipboard struct {
	copyCmd  ClipCmd
	pasteCmd ClipCmd
}

func (c *CmdClipboard) CopyText(text string) error {
	return nil
}

func (c *CmdClipboard) PasteText() (string, error) {
	return "", nil
}

func GetClipboard() (Clipboard, error) {
	if defs, ok := DefaultClipboard[runtime.GOOS]; ok {
		return findClipboard(defs)
	}

	return nil, errors.New("Unsupported OS")
}

func findClipboard(candidates []CmdClipboard) (*CmdClipboard, error) {
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
			copyCmd:  ClipCmd{cpath, candidate.copyCmd.args},
			pasteCmd: ClipCmd{ppath, candidate.pasteCmd.args},
		}, nil
	}

	return nil, errors.New("no available clipboard commands")
}
