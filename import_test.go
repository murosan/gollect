// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/exp/maps"
)

func TestImport_ToSpec(t *testing.T) {
	cases := []struct {
		in   *Import
		want *ast.ImportSpec
	}{
		{
			in: NewImport("f", "fmt", "fmt"),
			want: &ast.ImportSpec{
				Doc:  nil,
				Name: ast.NewIdent("f"),
				Path: &ast.BasicLit{
					ValuePos: 0,
					Kind:     0,
					Value:    "\"fmt\"",
				},
				Comment: nil,
				EndPos:  0,
			},
		},
		{
			in: NewImport("", "fmt", "fmt"),
			want: &ast.ImportSpec{
				Doc:  nil,
				Name: nil,
				Path: &ast.BasicLit{
					ValuePos: 0,
					Kind:     0,
					Value:    "\"fmt\"",
				},
				Comment: nil,
				EndPos:  0,
			},
		},
	}

	for i, c := range cases {
		v := c.in.ToSpec()

		if !reflect.DeepEqual(v, c.want) {
			t.Errorf("at: %d, want: %v, actual: %v", i, c.want, v)
		}
	}
}

func TestImport_Use(t *testing.T) {
	i := NewImport("", "fmt", "fmt")
	if i.used {
		t.Errorf("wrong initial state")
	}
	i.Use()
	if !i.used {
		t.Errorf("fail")
	}
}

func TestImport_IsUsed(t *testing.T) {
	i := NewImport("", "fmt", "fmt")
	if i.used {
		t.Errorf("wrong initial state")
	}
	i.Use()
	if !i.used {
		t.Errorf("fail")
	}
}

func TestImport_IsBuiltin(t *testing.T) {
	i1 := NewImport("", "fmt", "fmt")
	if !i1.IsBuiltin() {
		t.Errorf("should be builtin")
	}

	i2 := NewImport("", "fmt", "github.com/murosan/abc")
	if i2.IsBuiltin() {
		t.Errorf("should not be builtin")
	}
}

func TestImportSet(t *testing.T) {
	set := NewImportSet()

	name := "fmt"
	i1 := NewImport("", name, "fmt")

	exists := func(i *Import) bool {
		for _, v := range set.set {
			if v == i {
				return true
			}
		}
		return false
	}

	if exists(i1) {
		t.Errorf("wrong initial state")
	}

	set.AddAndGet(i1)
	if !exists(i1) {
		t.Errorf("failing to set")
	}

	if v := set.GetOrCreate(i1.alias, i1.name, i1.path); v != i1 {
		t.Errorf("should return without create")
	}

	i2 := NewImport("f", "fmt", "fmt")
	v := set.GetOrCreate(i2.alias, i2.name, i2.path)

	if v == i1 {
		t.Errorf("should create new one because the alias name is different")
	}
	if i2.alias != v.alias || i2.name != v.name || i2.path != v.path {
		t.Errorf("want: %v, actual: %v", i2, v)
	}
}

func TestImportSet_ToDecl(t *testing.T) {
	set := NewImportSet()
	i1 := NewImport("", "fmt", "fmt")
	i2 := NewImport("f", "fmt", "fmt")
	i3 := NewImport("unused", "fmt", "fmt")
	i4 := NewImport("", "abc", "github.com/murosan/abc")
	i5 := NewImport("abcv4", "abc", "github.com/murosan/abc/v2")

	i1.Use()
	i2.Use()
	// i3 is not used
	i4.Use()
	i5.Use()

	set.AddAndGet(i1)
	set.AddAndGet(i2)
	set.AddAndGet(i3)
	set.AddAndGet(i4)
	set.AddAndGet(i5)

	want := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1,
		Specs:  []ast.Spec{i1.ToSpec(), i2.ToSpec()}, // used && builtin only
	}
	actual := set.ToDecl()

	if !eqImportGenDecl(t, want, actual) {
		t.Errorf("\nwant:   %v\nactual: %v", want, actual)
	}
}

func TestImportSet_ToDecl2(t *testing.T) {
	set := NewImportSet()
	i1 := NewImport("", "fmt", "fmt")
	i1.Use()
	set.AddAndGet(i1)

	want := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 0,
		Specs:  []ast.Spec{i1.ToSpec()},
	}
	actual := set.ToDecl()

	if !eqImportGenDecl(t, want, actual) {
		t.Errorf("\nwant:   %v\nactual: %v", want, actual)
	}
}

func eqImportGenDecl(t *testing.T, a, b *ast.GenDecl) bool {
	t.Helper()
	if !reflect.DeepEqual(a.Doc, b.Doc) ||
		a.TokPos != b.TokPos ||
		a.Tok != b.Tok ||
		a.Lparen != b.Lparen ||
		a.Rparen != b.Rparen ||
		len(a.Specs) != len(b.Specs) ||
		((a.Specs == nil) != (b.Specs == nil)) {
		return false
	}

	//lint:ignore U1000 (staticcheck) tup is used
	type tup struct{ name, path string }
	aset := make(map[tup]struct{})
	bset := make(map[tup]struct{})

	for i := range a.Specs {
		aa, ok1 := a.Specs[i].(*ast.ImportSpec)
		bb, ok2 := b.Specs[i].(*ast.ImportSpec)

		// must be *ast.ImportSpec
		if !ok1 || !ok2 {
			return false
		}

		var aname, bname string
		if aa.Name != nil {
			aname = aa.Name.Name
		}
		if bb.Name != nil {
			bname = bb.Name.Name
		}
		aset[tup{name: aname, path: aa.Path.Value}] = struct{}{}
		bset[tup{name: bname, path: bb.Path.Value}] = struct{}{}
	}

	return maps.Equal(aset, bset)
}
