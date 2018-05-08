// +build !linux,!darwin,!windows

package clipboard

func GetClipboard() (Clipboard, error) {
	return nil, ErrUnsupported
}
