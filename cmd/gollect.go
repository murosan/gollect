package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/murosan/gollect/gollect"
)

var (
	configFile = flag.String("config", "", "configuration filepath. if specified, all other cli option will be ignored")
	inputFile  = flag.String("f", "main.go", "filepath of main.go or glob for main package files")
	out        = flag.String("out", "stdout", "output filepath. filepath, 'stdout' and 'clipboard' are available")
)

func main() {
	flag.Parse()

	var config *gollect.Config
	if *configFile == "" {
		config = &gollect.Config{
			InputFile:   *inputFile,
			OutputPaths: []string{*out},
		}
	} else {
		config = gollect.LoadConfig(*configFile)
	}

	if err := config.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gollect.Main(config)
}
