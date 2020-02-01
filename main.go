package main

import "github.com/shihanng/tfvar/cmd"

func main() {
	c, sync := cmd.New()
	_ = c.Execute()

	sync()
}
