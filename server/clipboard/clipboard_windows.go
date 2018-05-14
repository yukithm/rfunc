// +build windows

package clipboard

import cb "github.com/atotto/clipboard"

func GetClipboard() (Clipboard, error) {
	return WinClipboard(0), nil
}

type WinClipboard int

func (c WinClipboard) CopyText(text string) error {
	return cb.WriteAll(text)
}

func (c WinClipboard) PasteText() (string, error) {
	return cb.ReadAll()
}
