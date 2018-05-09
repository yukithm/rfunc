package utils

import (
	"io"
	"os/user"
	"path/filepath"
	"strings"
)

type nopWriteCloser struct {
	io.Writer
}

func (w *nopWriteCloser) Close() error {
	return nil
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

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

func FindString(list []string, value string) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}
	return -1
}
