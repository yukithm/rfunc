// +build linux

package shell

var DefaultShell = []CmdShell{
	{
		openCmd: OpenCmd{path: "xdg-open"},
	},
}

func GetShell() (Shell, error) {
	return findShell(DefaultShell)
}
