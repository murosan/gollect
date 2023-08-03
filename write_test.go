// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/murosan/gollect/testdata"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
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

		pkg := NewPackage("main")
		d, _ := program.DeclSet().Get(pkg, "main")
		d.Use()

		err := Write(&buf, program)

		if err != nil {
			t.Fatal(err)
		}

		if buf.String() != c.want {
			diff := dmp.New().DiffMain(c.want, buf.String(), true)
			t.Errorf(
				"\n[at] %d\n[diff]\n%s",
				i,
				colorDiff(diff),
			)
		}
	}
}
