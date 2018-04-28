package server

type Shell interface {
	OpenURL(url ...string) error
}
