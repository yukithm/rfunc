// +build windows

package server

import "errors"

// TODO: implements using ShellExecute API
//       or call "rundll32.exe url.dll,FileProtocolHandler URL"

func GetShell() (Shell, error) {
	return nil, errors.New("Unsupported OS")
}
