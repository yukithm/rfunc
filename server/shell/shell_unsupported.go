// +build !linux,!darwin,!windows

package shell

func GetShell() (Shell, error) {
	return nil, ErrUnsupported
}
