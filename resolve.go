// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"strings"
	"sync"
)

// AnalyzeForeach executes analyzing dependency for each packages.
func AnalyzeForeach(program *Program) {
	var wg sync.WaitGroup
	for _, pkg := range program.Packages() {
		wg.Add(1)
		go func(pkg *Package) {
			ExecCheck(program.FileSet(), pkg)
			pkg.InitObjects()
			ResolveDependency(pkg)
			wg.Done()
		}(pkg)
	}
	wg.Wait()
}

// ExecCheck executes types.Config.Check
func ExecCheck(fset *token.FileSet, pkg *Package) {
	conf := &types.Config{
		Importer: importer.ForCompiler(fset, "source", nil),
	}

	if _, err := conf.Check(pkg.path, fset, pkg.files, pkg.info); err != nil {
		panic(fmt.Errorf("types.Conf check: %w", err))
	}
}

// ResolveDependency analyzes dependency for each decls.
func ResolveDependency(pkg *Package) {
	for _, file := range pkg.files {
		for _, decl := range file.Decls {
			resolve(pkg, decl)
		}
	}
}

func resolve(pkg *Package, decl ast.Decl) {
	switch node := decl.(type) {
	case *ast.GenDecl:
		switch node.Tok {
		case token.VAR, token.CONST, token.TYPE:
			for _, spec := range node.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for i, id := range spec.Names {
						name := id.Name
						if pkg.objects[name].Decl == spec {
							pkg.Dependencies().GetOrCreate(name)
							setDependency(pkg, name, spec.Type)
							if len(spec.Values) > i {
								setDependency(pkg, name, spec.Values[i])
							}
						}
					}
				case *ast.TypeSpec:
					id := spec.Name.Name
					pkg.Dependencies().GetOrCreate(id)
					setDependency(pkg, id, spec.Type)
					if node.Doc != nil {
						for _, doc := range node.Doc.List {
							if strings.HasPrefix(doc.Text, keepMethods.String()) {
								pkg.Dependencies().TurnOnKeepMethodOption(id)
							}
						}
					}
				}
			}
		}

	case *ast.FuncDecl:
		id := node.Name.Name

		if node.Recv != nil {
			var typeID *ast.Ident
			switch expr := node.Recv.List[0].Type.(type) {
			case *ast.Ident:
				typeID = expr
			case *ast.StarExpr:
				typeID = expr.X.(*ast.Ident)
			}

			if typeID != nil {
				id = typeID.Name + "." + id
				pkg.Dependencies().SetMethod(typeID.Name, id)
				pkg.Dependencies().SetInternal(id, typeID.Name)
			}
		}

		pkg.Dependencies().GetOrCreate(id)
		setDependency(pkg, id, node)
	}
}

func setDependency(pkg *Package, id string, node ast.Node) {
	if node == nil {
		return
	}

	ast.Inspect(node, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.SelectorExpr:
			if sel, ok := pkg.info.Selections[node]; ok {
				if n := named(sel.Recv()); n != nil {
					path := n.Obj().Pkg().Path()
					if !isBuiltinPackage(path) {
						pkg.Dependencies().SetExternal(id, path, n.Obj().Name()+"."+node.Sel.Name)
					}
				}
				return true
			}

			if i, ok := node.X.(*ast.Ident); ok && i != nil {
				switch uses := pkg.info.Uses[i].(type) {
				case *types.PkgName:
					p := uses.Imported()
					alias, name, path := i.Name, p.Name(), p.Path()
					if name == alias {
						alias = ""
					}

					i := pkg.imports.GetOrCreate(alias, name, path)
					pkg.Dependencies().SetImport(id, i)

					if !isBuiltinPackage(path) {
						pkg.Dependencies().SetExternal(id, path, node.Sel.Name)
					}
				}
			}

		case *ast.Ident:
			if _, ok := pkg.objects[node.Name]; !ok {
				// break when the object is
				//   - not a package-level declaration
				//   - an external package object
				break
			}

			uses, ok := pkg.info.Uses[node]
			if !ok || uses.Pkg() == nil || isBuiltinPackage(uses.Pkg().Path()) {
				break
			}

			switch obj := uses.(type) {
			case *types.Const, *types.Var, *types.Func, *types.TypeName:
				pkg.Dependencies().SetInternal(id, obj.Name())
			}
		}

		return true
	})
}

func named(expr types.Type) *types.Named {
	switch expr := expr.(type) {
	case *types.Named:
		return expr
	case *types.Pointer:
		return named(expr.Elem())
	default:
		return nil
	}
}
