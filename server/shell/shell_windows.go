// +build windows

package shell

var DefaultShell = []CmdShell{
	{
		openCmd: OpenCmd{path: "rundll32.exe", args: []string{"url.dll,FileProtocolHandler"}},
	},
}

func GetShell() (Shell, error) {
	return findShell(DefaultShell)
}
