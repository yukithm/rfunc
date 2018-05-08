// +build linux

package clipboard

var DefaultClipboard = []CmdClipboard{
	{
		copyCmd:  Cmd{path: "xclip", args: []string{"-in", "-selection", "clipboard"}},
		pasteCmd: Cmd{path: "xclip", args: []string{"-out", "-selection", "clipboard"}},
	},
	{
		copyCmd:  Cmd{path: "xsel", args: []string{"--input", "--clipboard"}},
		pasteCmd: Cmd{path: "xsel", args: []string{"--output", "--clipboard"}},
	},
}

func GetClipboard() (Clipboard, error) {
	return findCmdClipboard(DefaultClipboard)
}
