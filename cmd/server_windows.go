// +build windows

package cmd

import "errors"

func daemonize() error {
	return errors.New("Not support daemonize on Windows")
}
