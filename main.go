package main

import (
	"os"

	"github.com/matheuzgomes/Snip/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
