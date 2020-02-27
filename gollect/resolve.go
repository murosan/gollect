package gollect

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
)

// AnalyzeForeach executes analyzing dependency for each packages.
func AnalyzeForeach(program *Program) {
	for _, pkg := range program.Packages() {
		ExecCheck(program.FileSet(), pkg)
		pkg.InitObjects()
		ResolveDependency(pkg)
	}
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
						if pkg.objects[id.Name].Decl == spec {
							pkg.Dependencies().GetOrCreate(id.Name)
							setDependency(pkg, id, spec.Values[i])
						}
					}
				case *ast.TypeSpec:
					pkg.Dependencies().GetOrCreate(spec.Name.Name)
					setDependency(pkg, spec.Name, spec.Type)
				}
			}
		}

	case *ast.FuncDecl:
		id := node.Name

		if node.Recv != nil {
			switch expr := node.Recv.List[0].Type.(type) {
			case *ast.Ident:
				id = expr
			case *ast.StarExpr:
				id = expr.X.(*ast.Ident)
			}
		}

		if id == nil {
			break
		}

		if _, ok := pkg.objects[id.Name]; ok {
			pkg.Dependencies().GetOrCreate(id.Name)
			setDependency(pkg, id, node.Body)
		}
	}
}

func setDependency(pkg *Package, id *ast.Ident, node ast.Node) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.SelectorExpr:
			if i, ok := node.X.(*ast.Ident); ok && i != nil {
				if pkgName, ok := pkg.info.Uses[i].(*types.PkgName); ok {
					p := pkgName.Imported()
					alias, name, path := i.Name, p.Name(), p.Path()
					if name == alias {
						alias = ""
					}

					i := pkg.imports.GetOrCreate(alias, name, path)
					pkg.Dependencies().SetImport(id.Name, i)

					if !isBuiltinPackage(path) {
						pkg.Dependencies().SetExternal(id.Name, path, node.Sel.Name)
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
			if !ok || uses.Pkg() == nil {
				break
			}

			switch obj := uses.(type) {
			case *types.Const, *types.Var, *types.Func, *types.TypeName:
				pkg.Dependencies().SetInternal(id.Name, obj.Name())
			}
		}

		return true
	})
}
