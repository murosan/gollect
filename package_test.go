// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"testing"
)

func TestPackage_InitObjects(t *testing.T) {
	p := NewPackage("", nil)
	if len(p.objects) != 0 {
		t.FailNow()
	}

	p.files = []*ast.File{
		{Scope: &ast.Scope{Objects: map[string]*ast.Object{
			"main": {},
			"varA": {},
		}}},
		{Scope: &ast.Scope{Objects: map[string]*ast.Object{
			"varB": {},
		}}},
	}

	p.InitObjects()

	if len(p.objects) != 3 {
		t.Errorf("length must be %d. actual: %d", 3, len(p.objects))
	}

	for i, k := range []string{"main", "varA", "varB"} {
		if _, ok := p.objects[k]; !ok {
			t.Errorf("not found at: %d, key: %s", i, k)
		}
	}
}

func TestPackages(t *testing.T) {
	packages := make(Packages)
	imports := NewImportSet()
	p := NewPackage("github.com/murosan/abc", imports)

	packages.Set(p.path, p)

	if pkg, ok := packages.Get(p.path); !ok || pkg != p {
		t.Errorf("want: (%t, %v), actual: (%t, %v)", true, p, ok, pkg)
	}
}
