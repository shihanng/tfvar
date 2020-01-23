package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shihanng/tfvar/pkg/tfvar"
)

func main() {
	names, err := tfvar.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(names)
}
