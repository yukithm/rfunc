// +build windows

package commands

import "errors"

func daemonize() error {
	return errors.New("Not support daemonize on Windows")
}
