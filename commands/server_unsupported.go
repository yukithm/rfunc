// +build netbsd openbsd dragonfly solaris

package commands

import "errors"

func daemonize() error {
	// See: https://github.com/VividCortex/godaemon
	return errors.New("Not support daemonize")
}
