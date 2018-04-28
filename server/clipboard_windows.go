// +build windows

package server

import "errors"

// TODO: implements using Clipboard API

func GetClipboard() (Clipboard, error) {
	return nil, errors.New("Unsupported OS")
}
