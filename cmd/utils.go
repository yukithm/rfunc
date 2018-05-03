package cmd

import (
	"os/user"
	"path/filepath"
	"strings"
)

func ExpandPath(path string) (string, error) {
	// ~user form is not supported because of user.Lookup() requires cgo.

	if strings.HasPrefix(path, "~/") {
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		path = u.HomeDir + path[1:]
	}
	epath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return epath, nil
}
