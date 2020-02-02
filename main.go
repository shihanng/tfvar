package main

import (
	"os"

	"github.com/shihanng/tfvar/cmd"
)

func main() {
	c, sync := cmd.New(os.Stdout)
	_ = c.Execute()

	sync()
}
