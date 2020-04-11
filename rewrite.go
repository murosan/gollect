// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// Filter provides a method for filtering slice of ast.Decl.
type Filter struct {
	dset *DeclSet
	pkg  *Package
}

// NewFilter returns new Filter.
func NewFilter(dset *DeclSet, pkg *Package) *Filter {
	return &Filter{
		dset: dset,
		pkg:  pkg,
	}
}

// Decls returns new slice that consists of used declarations.
// All unused declaration will be removed.
// Be careful this method manipulates decls directly.
func (f *Filter) Decls(decls []ast.Decl) (res []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			f.genDecl(decl)
			if l := len(decl.Specs); l != 0 {
				if l == 1 {
					decl.Lparen, decl.Rparen = 0, 0 // delete '(' and ')'
				}
				res = append(res, decl)
			}

		case *ast.FuncDecl:
			if f.isUsedFuncDecl(decl) {
				res = append(res, decl)
			}
		}
	}
	return
}

func (f *Filter) genDecl(node *ast.GenDecl) {
	switch node.Tok {
	case token.VAR, token.CONST, token.TYPE:
		node.Specs = f.specs(node.Specs)

	case token.IMPORT:
		// remove all imports to add unique ones later
		node.Specs = nil
	}

	// remove gollect annotation comments
	f.annotation(node)
}

func (f *Filter) specs(specs []ast.Spec) (res []ast.Spec) {
	for _, spec := range specs {
		switch spec := spec.(type) {
		case *ast.ValueSpec:
			f.valueSpec(spec)
			if len(spec.Names) != 0 {
				res = append(res, spec)
			}

		case *ast.TypeSpec:
			if f.isUsed(spec.Name.Name) {
				res = append(res, spec)
			}
		}
	}
	return
}

func (f *Filter) valueSpec(spec *ast.ValueSpec) {
	var names []*ast.Ident
	var values []ast.Expr

	for i, id := range spec.Names {
		name := nameForUnderscore(id)
		if f.isUsed(name) {
			names = append(names, id)
			if len(spec.Values) > i {
				values = append(values, spec.Values[i])
			}
		}
	}

	spec.Names = names
	spec.Values = values
}

func (f *Filter) isUsedFuncDecl(decl *ast.FuncDecl) bool {
	if decl.Recv == nil && decl.Name.Name == "init" {
		return true
	}

	var keys []string

	if decl.Recv != nil {
		switch expr := decl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			keys = append(keys, expr.Name)
		case *ast.StarExpr:
			keys = append(keys, expr.X.(*ast.Ident).Name)
		}
	}

	keys = append(keys, decl.Name.Name)
	return f.isUsed(keys...)
}

func (f *Filter) annotation(node *ast.GenDecl) {
	if node.Doc == nil {
		return
	}
	docs := make([]*ast.Comment, len(node.Doc.List))
	i := 0
	for _, doc := range node.Doc.List {
		doc := doc
		if !strings.Contains(doc.Text, annotationPrefix) {
			docs[i] = doc
			i++
		}
	}
	node.Doc.List = docs[:i]
}

func (f *Filter) isUsed(id ...string) bool {
	b, ok := f.dset.Get(f.pkg, id...)
	return ok && b.IsUsed()
}

// PackageSelectorExpr removes external package's selectors.
//
//   fmt.Println() → fmt.Println() // keep builtin packages
//   mypkg.SomeFunc() → SomeFunc() // remove package selector
//
func (f *Filter) PackageSelectorExpr(node ast.Node) {
	astutil.Apply(node, func(cr *astutil.Cursor) bool {
		switch n := cr.Node().(type) {
		case nil:
			return false

		case *ast.SelectorExpr:
			i, ok := n.X.(*ast.Ident)
			if !ok || i == nil {
				break
			}

			uses, _ := f.pkg.UsesInfo(i)
			pn, ok := uses.(*types.PkgName)
			if !ok || isBuiltinPackage(pn.Imported().Path()) {
				break
			}

			cr.Replace(n.Sel)
		}

		return true
	}, nil)
}
