// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func TestFilterDecls(t *testing.T) {
	deps := make(Dependencies)
	d1 := NewDependency("d1")
	d2 := NewDependency("d2")
	d3 := NewDependency("d3")
	d4 := NewDependency("d4")
	d5 := NewDependency("d5")
	f1 := NewDependency(d4.name + "." + "funcA")
	f2 := NewDependency(d4.name + "." + "funcB")
	f3 := NewDependency(d3.name + "." + "funcC")
	f4 := NewDependency(d3.name + "." + "funcD")
	deps.Set(d1)
	deps.Set(d2)
	deps.Set(d3)
	deps.Set(d4)
	deps.Set(d5)
	deps.Set(f1)
	deps.Set(f2)
	deps.Set(f3)
	deps.Set(f4)
	d4.SetMethod(f1)
	d4.SetMethod(f2)
	d1.Use()
	d2.Use()
	d4.Use()
	f1.Use()
	f2.Use()

	sut := []ast.Decl{
		&ast.GenDecl{
			// imports will be removed
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{Name: ast.NewIdent(d1.name)},
				&ast.ImportSpec{Name: ast.NewIdent("fake")},
			},
		},
		&ast.GenDecl{
			// unused codes will be removed
			// if Specs becomes empty, GenDecl will also be removed
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{ast.NewIdent(d3.name)},
					Values: []ast.Expr{ast.NewIdent(d3.name)},
				},
			},
		},
		&ast.GenDecl{
			// used codes will be left
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{ast.NewIdent(d2.name)},
					Values: []ast.Expr{ast.NewIdent(d2.name)},
				},
				&ast.ValueSpec{
					Names:  []*ast.Ident{ast.NewIdent(d3.name)},
					Values: []ast.Expr{ast.NewIdent(d3.name)},
				},
			},
		},
		&ast.FuncDecl{Name: ast.NewIdent(d4.name)}, // used
		&ast.FuncDecl{Name: ast.NewIdent(d5.name)}, // unused
		// used
		&ast.FuncDecl{
			Name: ast.NewIdent("funcA"),
			Recv: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent(d4.name)}}},
		},
		// used
		&ast.FuncDecl{
			Name: ast.NewIdent("funcB"),
			Recv: &ast.FieldList{List: []*ast.Field{
				{Type: &ast.StarExpr{X: ast.NewIdent(d4.name)}},
			}},
		},
		// unused
		&ast.FuncDecl{
			Name: ast.NewIdent("funcC"),
			Recv: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent(d5.name)}}},
		},
		// unused
		&ast.FuncDecl{
			Name: ast.NewIdent("funcD"),
			Recv: &ast.FieldList{List: []*ast.Field{
				{Type: &ast.StarExpr{X: ast.NewIdent(d5.name)}},
			}},
		},
		// used
		&ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent(d1.name)}},
		},
		// unused
		&ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent(d3.name)}},
		},
	}

	actual := FilterDecls(deps, sut)

	if len(actual) != 5 {
		t.Errorf("length is not 5. actual: %v", actual)
	}

	if actual[0] != sut[2] {
		t.Errorf("[want]\n%v\n[actual]\n%v", sut[2], actual[0])
	}

	specs := sut[2].(*ast.GenDecl).Specs
	if len(specs) != 1 {
		t.Errorf("length is not 1. actual: %v", specs)
	}

	if n := specs[0].(*ast.ValueSpec).Names[0].Name; n != d2.name {
		t.Errorf("want: %s, actual: %s", d2.name, n)
	}

	if actual[1] != sut[3] {
		t.Errorf("[want]\n%vn[actual]\n%v", sut[3], actual[1])
	}

	if n := actual[1].(*ast.FuncDecl).Name.Name; n != d4.name {
		t.Errorf("want: %s, actual: %s", d4.name, n)
	}

	if actual[2] != sut[5] {
		t.Errorf("[want]\n%vn[actual]\n%v", sut[5], actual[2])
	}

	if actual[3] != sut[6] {
		t.Errorf("[want]\n%vn[actual]\n%v", sut[6], actual[3])
	}

	if actual[4] != sut[9] {
		t.Errorf("[want]\n%vn[actual]\n%v", sut[6], actual[3])
	}
}

func TestRemoveExternalIdents(t *testing.T) {
	pkgName := func(id *ast.Ident, name, path string) *types.PkgName {
		return types.NewPkgName(0, nil, "", types.NewPackage(path, name))
	}

	pkg := NewPackage("main", nil)

	id1, name1 := ast.NewIdent("fmt"), "Println"
	id2, name2 := ast.NewIdent("ext"), "SomeFunc"

	pkg.info.Uses[id1] = pkgName(id1, "fmt", "fmt")
	pkg.info.Uses[id1] = pkgName(id2, "ext", "github.com/murosan/abc")

	funcDecl := func(ident *ast.Ident, name string) *ast.FuncDecl {
		return &ast.FuncDecl{
			Name: ident,
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ident,
								Sel: ast.NewIdent(name),
							},
						},
					},
				},
			},
		}
	}

	decls := []ast.Decl{
		funcDecl(id1, name1),
		funcDecl(id2, name2),
	}
	file := &ast.File{Decls: decls}
	RemoveExternalIdents(file, pkg)

	if len(file.Decls) != 2 || !reflect.DeepEqual(decls, file.Decls) {
		t.FailNow()
	}

	{
		// no panic
		head := decls[0].(*ast.FuncDecl).Body.List[0]
		cl := head.(*ast.ExprStmt).X.(*ast.CallExpr)
		n := cl.Fun.(*ast.Ident).Name

		if n == id1.Name {
			t.Errorf("want: %s, actual: %s", id1.Name, n)
		}
	}

	{
		// no panic
		head := decls[1].(*ast.FuncDecl).Body.List[0]
		cl := head.(*ast.ExprStmt).X.(*ast.CallExpr)
		n := cl.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name

		if n == name2 {
			t.Errorf("want: %s, actual: %s", name2, n)
		}
	}
}
