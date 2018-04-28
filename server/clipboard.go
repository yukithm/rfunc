package server

type Clipboard interface {
	CopyText(text string) error
	PasteText() (string, error)
}
