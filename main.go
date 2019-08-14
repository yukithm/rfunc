package main

import (
	"os"

	"github.com/yukithm/rfunc/commands"
)

func main() {
	code := commands.Execute()
	if code != 0 {
		os.Exit(code)
	}
}
