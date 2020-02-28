// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"reflect"
	"testing"

	"github.com/murosan/gollect/testdata"
)

func TestResolveDependency(t *testing.T) {
	program := NewProgram(testdata.FilePaths.A)
	ParseAll(program)
	AnalyzeForeach(program)

	{
		pp := "main"
		want := make(Dependencies)
		caller := "main"
		want.SetInternal(caller, "num")
		want.SetExternal(caller, testdata.PackagePaths.Pkg1, "TypeA")
		want.SetExternal(caller, testdata.PackagePaths.Pkg1, "NumA")
		want.SetExternal(caller, testdata.PackagePaths.Pkg1, "NumC")
		want.SetExternal(caller, testdata.PackagePaths.Pkg1, "PrintMax")
		want.SetImport(caller, NewImport("", "pkg1", testdata.PackagePaths.Pkg1))
		want.SetImport(caller, NewImport("", "fmt", "fmt"))

		pkg, _ := program.Packages().Get(pp)
		actual := pkg.Dependencies()

		if !reflect.DeepEqual(want, actual) {
			t.Errorf("\npackage:%s\nwant\n%v\nactual\n%v", pp, want, actual)
		}
	}

	{
		pp := testdata.PackagePaths.Pkg1
		want := make(Dependencies)
		want.Set(NewDependency("NumA"))
		want.Set(NewDependency("NumB"))
		want.Set(NewDependency("NumC"))
		want.Set(NewDependency("TypeA"))
		want.SetExternal("PrintMax", testdata.PackagePaths.Pkg2, "Max")
		want.SetImport("PrintMax", NewImport("p", "pkg2", testdata.PackagePaths.Pkg2))
		want.SetImport("PrintMax", NewImport("", "fmt", "fmt"))
		want.SetImport("TypeA", NewImport("f", "fmt", "fmt"))

		pkg, _ := program.Packages().Get(pp)
		actual := pkg.Dependencies()

		if !reflect.DeepEqual(want, actual) {
			t.Errorf("\npackage:%s\nwant\n%v\nactual\n%v", pp, want, actual)
		}
	}

	{
		pp := testdata.PackagePaths.Pkg2
		want := make(Dependencies)
		want.Set(NewDependency("Max"))

		pkg, _ := program.Packages().Get(pp)
		actual := pkg.Dependencies()

		if !reflect.DeepEqual(want, actual) {
			t.Errorf("\npackage:%s\nwant\n%v\nactual\n%v", pp, want, actual)
		}
	}
}
