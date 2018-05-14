// +build linux freebsd darwin

package cmd

import "github.com/VividCortex/godaemon"

func daemonize() error {
	_, _, err := godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	return err
}
