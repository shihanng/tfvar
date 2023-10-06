package main

import (
	"os"

	"github.com/shihanng/tfvar/cmd"
)

var version = "dev"

func main() {
	c, sync := cmd.New(os.Stdout, version)
	defer sync()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
