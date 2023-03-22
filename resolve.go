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
	fset, dset := program.FileSet(), program.DeclSet()
	iset, pset := program.ImportSet(), program.PackageSet()

	for _, pkg := range pset {
		ExecCheck(fset, pkg)
		pkg.InitObjects()
		NewDeclFinder(dset, iset, pkg).Files(pkg.files)
	}

	dsetValues := dset.Values()
	for _, d := range dsetValues {
		wg.Add(1)
		go func(d Decl) {
			NewDependencyResolver(dset, iset, pset).Check(d)
			wg.Done()
		}(d)
	}
	wg.Wait()

	for _, d := range dsetValues {
		wg.Add(1)
		go func(d Decl) {
			NewDependencyResolver(dset, iset, pset).CheckEmbedded(d)
			wg.Done()
		}(d)
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

// DeclFinder find package-level declarations and set it to DeclSet.
type DeclFinder struct {
	dset *DeclSet
	iset *ImportSet
	pkg  *Package
}

// NewDeclFinder returns new DeclFinder
func NewDeclFinder(dset *DeclSet, iset *ImportSet, pkg *Package) *DeclFinder {
	return &DeclFinder{
		dset: dset,
		iset: iset,
		pkg:  pkg,
	}
}

// Files finds package-level declarations foreach file.decls concurrently.
func (f *DeclFinder) Files(files []*ast.File) {
	var wg sync.WaitGroup

	for _, file := range files {
		for _, decl := range file.Decls {
			wg.Add(1)

			go func(decl ast.Decl) {
				f.Decl(decl)
				wg.Done()
			}(decl)
		}
	}

	wg.Wait()
}

// Decl finds package-level declarations from ast.Decl.
func (f *DeclFinder) Decl(decl ast.Decl) {
	switch decl := decl.(type) {
	case *ast.GenDecl:
		f.GenDecl(decl)

	case *ast.FuncDecl:
		f.FuncDecl(decl)
	}
}

// GenDecl finds package-level declarations from ast.GenDecl.
func (f *DeclFinder) GenDecl(decl *ast.GenDecl) {
	switch decl.Tok {
	case token.CONST:
		f.varSpecs(decl, true)
	case token.VAR:
		f.varSpecs(decl, false)
	case token.TYPE:
		f.typeSpecs(decl)
	}
}

func (f *DeclFinder) varSpecs(decl *ast.GenDecl, isConst bool) {
	var prev Decl
	iota := isConst && f.hasIota(decl)

	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, id := range spec.Names {
			name := nameForUnderscore(id)
			d := f.dset.GetOrCreate(DecCommon, f.pkg, name)
			d.SetNode(spec)
			if prev != nil {
				if iota {
					d.Uses(prev)
					prev.Uses(d)
				} else if isConst && len(spec.Values) == 0 {
					d.Uses(prev)
				}
			}
			prev = d
		}
	}
}

func (f *DeclFinder) hasIota(decl *ast.GenDecl) (has bool) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, v := range spec.Values {
			ast.Inspect(v, func(node ast.Node) bool {
				id, ok := node.(*ast.Ident)
				if !ok || id.Name != "iota" {
					return true
				}

				uses, ok := f.pkg.UsesInfo(id)
				if !ok || uses.Parent() != types.Universe {
					return true
				}

				has = true
				return false
			})

			if has {
				return
			}
		}
	}
	return
}

func (f *DeclFinder) typeSpecs(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		id := spec.Name
		tdecl := f.dset.GetOrCreate(DecType, f.pkg, id.Name).(*TypeDecl)
		tdecl.SetNode(spec)

		if decl.Doc != nil {
			for _, doc := range decl.Doc.List {
				if strings.HasPrefix(doc.Text, keepMethods.String()) {
					tdecl.KeepMethod()
				}
			}
		}

		def, ok := f.pkg.DefInfo(id)
		if !ok {
			continue
		}

		// fill methods
		// https://pkg.go.dev/go/types?tab=doc#example-MethodSet
		for _, t := range []types.Type{def.Type(), types.NewPointer(def.Type())} {
			mset := types.NewMethodSet(t)
			for i := 0; i < mset.Len(); i++ {
				m := mset.At(i)
				mdecl := f.dset.GetOrCreate(
					DecMethod,
					f.pkg,
					id.Name,
					m.Obj().Name(),
				).(*MethodDecl)
				if len(m.Index()) >= 2 {
					mdecl.SetEmbedded(true)
				}
				mdecl.SetType(tdecl)
				tdecl.SetMethod(mdecl)
			}
		}
	}
}

// FuncDecl finds package-level declarations from ast.FuncDecl.
func (f *DeclFinder) FuncDecl(decl *ast.FuncDecl) {
	name := decl.Name.Name

	if decl.Recv == nil {
		d := f.dset.GetOrCreate(DecCommon, f.pkg, name)
		d.SetNode(decl)
		return
	}

	recvID := receiverID(decl.Recv.List[0].Type)
	if recvID != nil {
		md := f.dset.GetOrCreate(DecMethod, f.pkg, recvID.Name, name)
		md.SetNode(decl)
	}
}

func receiverID(expr ast.Expr) *ast.Ident {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr
	case *ast.StarExpr:
		return expr.X.(*ast.Ident)
	default:
		return nil
	}
}

// DependencyResolver provides a method for checking dependency.
type DependencyResolver struct {
	dset *DeclSet
	iset *ImportSet
	pset PackageSet
}

// NewDependencyResolver returns new DependencyResolver
func NewDependencyResolver(dset *DeclSet, iset *ImportSet, pset PackageSet) *DependencyResolver {
	return &DependencyResolver{
		dset: dset,
		iset: iset,
		pset: pset,
	}
}

// Check checks on which given decl depending on.
func (r *DependencyResolver) Check(decl Decl) {
	if decl.Node() == nil {
		return
	}

	ast.Inspect(decl.Node(), func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.SelectorExpr:
			if sel, ok := decl.Pkg().SelInfo(node); ok {
				n := named(sel.Recv())
				if n == nil {
					return true
				}

				path := n.Obj().Pkg().Path()
				if isBuiltinPackage(path) {
					return true
				}

				pkg, ok := r.pset.Get(path)
				if !ok {
					return true
				}

				d, ok := r.dset.Get(pkg, n.Obj().Name(), node.Sel.Name)
				if ok {
					decl.Uses(d)
				}

				return true
			}

			if id, ok := node.X.(*ast.Ident); ok && id != nil {
				uses, _ := decl.Pkg().UsesInfo(id)
				switch uses := uses.(type) {
				case *types.PkgName:
					imported := uses.Imported()
					alias, name, path := id.Name, imported.Name(), imported.Path()
					if name == alias {
						alias = ""
					}

					i := r.iset.GetOrCreate(alias, name, path)
					decl.UsesImport(i)
					pkg, ok := r.pset.Get(path)

					if !ok || isBuiltinPackage(path) {
						return true
					}

					d, ok := r.dset.Get(pkg, node.Sel.Name)
					if ok {
						decl.Uses(d)
					}
				}
			}

		case *ast.Ident:
			if _, ok := decl.Pkg().GetObject(node.Name); !ok {
				// break when the object is
				//   - not a package-level declaration
				//   - an external package object
				break
			}

			uses, ok := decl.Pkg().UsesInfo(node)
			if !ok || uses.Pkg() == nil || isBuiltinPackage(uses.Pkg().Path()) {
				break
			}

			switch obj := uses.(type) {
			case *types.Const, *types.Var, *types.Func, *types.TypeName:
				d, _ := r.dset.Get(decl.Pkg(), obj.Name())
				decl.Uses(d)
			}
		}

		return true
	})
}

// CheckEmbedded checks set a method inherit from to dependency set
// when the decl is embedded method.
func (r *DependencyResolver) CheckEmbedded(decl Decl) {
	mdecl, ok := decl.(*MethodDecl)
	if !ok || !mdecl.IsEmbedded() {
		return
	}

	for _, d := range mdecl.Type().UsesDecls() {
		tdecl, ok := d.(*TypeDecl)
		if !ok {
			continue
		}

		for _, d := range tdecl.Methods() {
			if mdecl.Name() == d.Name() {
				mdecl.Uses(d)
			}
		}
	}
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
