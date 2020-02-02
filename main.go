package main

import (
	"os"

	"github.com/shihanng/tfvar/cmd"
)

var version = "dev"

func main() {
	c, sync := cmd.New(os.Stdout, version)
	_ = c.Execute()

	sync()
}
