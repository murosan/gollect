// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/murosan/gollect/testdata"
)

func TestWrite(t *testing.T) {
	cases := []struct {
		path,
		want string
	}{
		{
			path: testdata.FilePaths.Write1,
			want: `package main

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
`,
		},
		{
			path: testdata.FilePaths.Write2,
			want: `package main

func main() {
	FuncA()
	FuncB()
}

func FuncA() {}

func FuncB() {

}
`,
		},
	}

	for i, c := range cases {
		var buf bytes.Buffer

		program := NewProgram()
		paths, _ := filepath.Glob(c.path)
		ParseAll(program, "main", paths)
		AnalyzeForeach(program)

		next := []ExternalDependencySet{{}}
		next[0].Add("main", "main")
		UseAll(program.Packages(), next)

		err := Write(&buf, program)

		if err != nil {
			t.Fatal(err)
		}

		if buf.String() != c.want {
			t.Errorf(
				"\n[at] %d\n[want]\n%s\n[actual]\n%s",
				i,
				c.want,
				buf.String(),
			)
		}
	}
}
