package main

import (
	"os"

	"github.com/yukithm/rfunc/cmd"
)

func main() {
	code := cmd.Execute()
	if code != 0 {
		os.Exit(code)
	}
}
