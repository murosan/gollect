package gollect

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

func FilterDecls(deps Dependencies, decls []ast.Decl) (res []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			FilterGenDecl(deps, decl)
			if len(decl.Specs) != 0 {
				res = append(res, decl)
			}

		case *ast.FuncDecl:
			if IsUsedFuncDecl(deps, decl) {
				res = append(res, decl)
			}
		}
	}
	return
}

func FilterGenDecl(deps Dependencies, node *ast.GenDecl) {
	switch node.Tok {
	case token.VAR, token.CONST, token.TYPE:
		node.Specs = FilterSpecs(deps, node.Specs)

	case token.IMPORT:
		// remove all imports to add unique ones later
		node.Specs = nil
	}
}

func FilterSpecs(deps Dependencies, specs []ast.Spec) (res []ast.Spec) {
	for _, spec := range specs {
		switch spec := spec.(type) {
		case *ast.ValueSpec:
			FilterValueSpec(deps, spec)
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

func FilterValueSpec(deps Dependencies, spec *ast.ValueSpec) {
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

func IsUsedFuncDecl(deps Dependencies, decl *ast.FuncDecl) bool {
	id := decl.Name

	if decl.Recv != nil {
		switch expr := decl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			id = expr
		case *ast.StarExpr:
			id = expr.X.(*ast.Ident)
		}
	}

	if id == nil {
		return false
	}

	return deps.IsUsed(id.Name)
}

func RemoveExternalIdents(node ast.Node, pkg *Package) {
	iset, uses := pkg.imports, pkg.info.Uses

	astutil.Apply(node, func(cr *astutil.Cursor) bool {
		switch n := cr.Node().(type) {
		case nil:
			return false

		case *ast.SelectorExpr:
			if i, ok := n.X.(*ast.Ident); ok && i != nil {
				if _, ok := uses[i].(*types.PkgName); ok {
					ip, ok := iset.Get(i.Name)
					if !ok || !isBuiltinPackage(ip.path) {
						cr.Replace(n.Sel)
					}
				}
			}
		}

		return true
	}, nil)
}
