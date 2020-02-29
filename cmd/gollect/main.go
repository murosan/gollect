// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/murosan/gollect"
)

var (
	cnf   = flag.String("config", "", "configuration filepath. if specified, all other cli option will be ignored")
	input = flag.String("main", "main.go", "filepath of main.go or glob for main package files")
	out   = flag.String("out", "stdout", "output filepath. filepath, 'stdout' and 'clipboard' are available")

	config *gollect.Config
)

func main() {
	flag.Parse()

	if *cnf == "" {
		config = &gollect.Config{
			InputFile:   *input,
			OutputPaths: []string{*out},
		}
	} else {
		config = gollect.LoadConfig(*cnf)
	}

	if err := gollect.Main(config); err != nil {
		panic(err)
	}
}
