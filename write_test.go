// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"bytes"
	"testing"

	"github.com/murosan/gollect/testdata"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer

	program := NewProgram(testdata.FilePaths.Write)
	ParseAll(program)
	AnalyzeForeach(program)

	next := []ExternalDependencySet{{}}
	next[0].Add("main", "main")
	UseAll(program.Packages(), next)

	err := Write(&buf, program)

	if err != nil {
		t.Fatal(err)
	}

	want := `package main

import (
	"fmt"
	f "fmt"
)

func main() {
	fmt.Println(NonFormatted1)
	Abc()
	Def()
}

var NonFormatted1 = func() int {
	return 100
}()

//  comment
func Abc() {
	fmt.Println("abc")
}

func Def() { f.Println("def") }
`

	if buf.String() != want {
		t.Errorf("\n[want]\n%s\n[actual]\n%s", want, buf.String())
	}
}
