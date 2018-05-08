// +build darwin

package clipboard

var DefaultClipboard = []CmdClipboard{
	{
		copyCmd:  Cmd{path: "pbcopy"},
		pasteCmd: Cmd{path: "pbpaste"},
	},
}

func GetClipboard() (Clipboard, error) {
	return findCmdClipboard(DefaultClipboard)
}
