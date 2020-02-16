package main

import (
	"flag"
	"os"

	"github.com/murosan/gollect/gollect"
)

var (
	glob = flag.String("f", "main.go", "input path of main package file")
)

func main() {
	flag.Parse()
	if *glob == "" {
		flag.Usage()
		os.Exit(1)
	}

	gollect.Main(*glob)
}
