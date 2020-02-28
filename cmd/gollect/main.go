// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/murosan/gollect"
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
