package main

import (
	"os"

	"github.com/terem/bom/cmd/bom/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
