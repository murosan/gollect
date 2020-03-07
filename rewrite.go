// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

// FilterDecls returns new slice that consists of used declarations.
// All unused declaration will be removed.
// Be careful this method manipulates decls directly.
func FilterDecls(deps Dependencies, decls []ast.Decl) (res []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			filterGenDecl(deps, decl)
			if l := len(decl.Specs); l != 0 {
				if l == 1 {
					decl.Lparen, decl.Rparen = 0, 0 // delete '(' and ')'
				}
				res = append(res, decl)
			}

		case *ast.FuncDecl:
			if isUsedFuncDecl(deps, decl) {
				res = append(res, decl)
			}
		}
	}
	return
}

func filterGenDecl(deps Dependencies, node *ast.GenDecl) {
	switch node.Tok {
	case token.VAR, token.CONST, token.TYPE:
		node.Specs = filterSpecs(deps, node.Specs)

	case token.IMPORT:
		// remove all imports to add unique ones later
		node.Specs = nil
	}
}

func filterSpecs(deps Dependencies, specs []ast.Spec) (res []ast.Spec) {
	for _, spec := range specs {
		switch spec := spec.(type) {
		case *ast.ValueSpec:
			filterValueSpec(deps, spec)
			if len(spec.Names) != 0 {
				res = append(res, spec)
			}

		case *ast.TypeSpec:
			if deps.IsUsed(spec.Name.Name) {
				res = append(res, spec)
			}
		}
	}
	return
}

func filterValueSpec(deps Dependencies, spec *ast.ValueSpec) {
	var names []*ast.Ident
	var values []ast.Expr

	for i, id := range spec.Names {
		if deps.IsUsed(id.Name) {
			names = append(names, id)
			values = append(values, spec.Values[i])
		}
	}

	spec.Names = names
	spec.Values = values
}

func isUsedFuncDecl(deps Dependencies, decl *ast.FuncDecl) bool {
	id := decl.Name

	if decl.Recv != nil {
		switch expr := decl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			id = expr
		case *ast.StarExpr:
			id = expr.X.(*ast.Ident)
		}
	}

	return id != nil && deps.IsUsed(id.Name)
}

// RemoveExternalIdents removes external package's selectors.
//
//   fmt.Println() → fmt.Println() // keep builtin packages
//   mypkg.SomeFunc() → SomeFunc() // remove package selector
//
func RemoveExternalIdents(node ast.Node, pkg *Package) {
	astutil.Apply(node, func(cr *astutil.Cursor) bool {
		switch n := cr.Node().(type) {
		case nil:
			return false

		case *ast.SelectorExpr:
			if i, ok := n.X.(*ast.Ident); ok && i != nil {
				if pn, ok := pkg.info.Uses[i].(*types.PkgName); ok {
					if !isBuiltinPackage(pn.Imported().Path()) {
						cr.Replace(n.Sel)
					}
				}
			}
		}

		return true
	}, nil)
}
