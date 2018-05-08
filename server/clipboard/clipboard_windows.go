// +build windows

package clipboard

// TODO: implements using Clipboard API

func GetClipboard() (Clipboard, error) {
	return nil, ErrUnsupported
}
