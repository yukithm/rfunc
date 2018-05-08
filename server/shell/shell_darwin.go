// +build darwin

package shell

var DefaultShell = []CmdShell{
	{
		openCmd: OpenCmd{path: "open"},
	},
}

func GetShell() (Shell, error) {
	return findShell(DefaultShell)
}
