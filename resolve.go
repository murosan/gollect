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
)

// AnalyzeForeach executes analyzing dependency for each packages.
func AnalyzeForeach(program *Program, initialPkg, initialObj string) {
	fset, dset := program.FileSet(), program.DeclSet()
	iset, pset := program.ImportSet(), program.PackageSet()

	for _, pkg := range pset {
		ExecCheck(fset, pkg)
		pkg.InitObjects()
		NewDeclFinder(dset, iset, pkg).Files()
	}

	pkg, ok := pset.Get(initialPkg)
	if !ok {
		panic("No such package: " + initialPkg)
	}

	initial, ok := dset.Get(pkg, initialObj)
	if !ok {
		panic("No such decl:" + initialObj)
	}
	resolver := NewDependencyResolver(dset, iset, pset)
	resolver.CheckEach(initial)

	for _, d := range dset.ListInitOrUnderscore() {
		resolver.CheckEach(d)
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

// DeclFinder find package-level declarations and set it to DeclSet.
type DeclFinder struct {
	dset DeclSet
	iset *ImportSet
	pkg  *Package
}

// NewDeclFinder returns new DeclFinder
func NewDeclFinder(dset DeclSet, iset *ImportSet, pkg *Package) *DeclFinder {
	return &DeclFinder{
		dset: dset,
		iset: iset,
		pkg:  pkg,
	}
}

// Files finds package-level declarations foreach file.decls concurrently.
func (f *DeclFinder) Files() {
	for _, file := range f.pkg.files {
		for _, decl := range file.Decls {
			f.Decl(decl)
		}
	}
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
	case token.IMPORT:
		// f.importSpecs(decl)
	case token.CONST:
		f.varSpecs(decl, true)
	case token.VAR:
		f.varSpecs(decl, false)
	case token.TYPE:
		f.typeSpecs(decl)
	}
}

// func (f *DeclFinder) importSpecs(decl *ast.GenDecl) {
// 	for _, spec := range decl.Specs {
// 		spec, ok := spec.(*ast.ImportSpec)
// 		if ok && spec.Name != nil && spec.Name.Name == "." {
// 			pkg := f.pkg.Info().Defs[spec.Name].(*types.PkgName)
// 			f.iset.AddDotImport(pkg.Imported())
// 		}
// 	}
// }

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

				uses, ok := f.pkg.Info().Uses[id]
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

		def, ok := f.pkg.Info().Defs[id]
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

		switch tpe := spec.Type.(type) {
		case *ast.StructType:
			for _, field := range tpe.Fields.List {
				if id := fieldTypeID(field.Type); id != nil {
					if it, ok := f.pkg.Info().ObjectOf(id).Type().Underlying().(*types.Interface); ok {
						// use each interface methods.
						// when struct type has an interface type in its field list,
						// the interface methods should be left.
						// example:
						//   type S struct { sort.Interface }
						// when the case of above example,
						// method `Len`, `Less` ans `Swap` will be left.
						for i := 0; i < it.NumMethods(); i++ {
							m, ok := tdecl.GetMethodByName(it.Method(i).Name())
							if ok {
								tdecl.Uses(m)
							}
						}
					}
				}
			}
		}
	}
}

func fieldTypeID(expr ast.Expr) *ast.Ident {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr
	case *ast.SelectorExpr:
		return expr.Sel
	case *ast.IndexExpr:
		return expr.X.(*ast.Ident)
	default:
		return nil
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
	case *ast.IndexExpr:
		return expr.X.(*ast.Ident)
	case *ast.IndexListExpr:
		return expr.X.(*ast.Ident)
	case *ast.StarExpr:
		return receiverID(expr.X)
	default:
		return nil
	}
}

// DependencyResolver provides a method for checking dependency.
type DependencyResolver struct {
	dset DeclSet
	iset *ImportSet
	pset PackageSet

	queue []Decl
}

// NewDependencyResolver returns new DependencyResolver
func NewDependencyResolver(dset DeclSet, iset *ImportSet, pset PackageSet) *DependencyResolver {
	return &DependencyResolver{
		dset: dset,
		iset: iset,
		pset: pset,
	}
}

func (r *DependencyResolver) use(decl, usedBy Decl) {
	if decl.IsUsed() {
		return
	}
	decl.Use()

	switch decl := decl.(type) {
	case *MethodDecl:
		r.use(decl.Type(), nil) // should check earlier to resolve embedded methods.
		decl.GetUses().Each(func(d Decl) { r.use(d, decl) })
		r.push(decl)

	case *TypeDecl:
		r.push(decl) // should check type earlier to resolve embedded methods.
		decl.GetUses().Each(func(d Decl) { r.use(d, decl) })
		if decl.ShouldKeepMethods() {
			decl.EachMethod(func(m *MethodDecl) { r.use(m, nil) })
		}

		// use lazily for checking embedded methods
		tpe, ok := usedBy.(*TypeDecl)
		if ok && tpe != nil {
			tpe.Uses(decl)
		}

	default:
		decl.GetUses().Each(func(d Decl) { r.use(d, decl) })
		r.push(decl)
	}
}

func (r *DependencyResolver) useImport(i *Import) { r.iset.AddAndGet(i).Use() }

func (r *DependencyResolver) push(decl Decl) { r.queue = append(r.queue, decl) }

func (r *DependencyResolver) pop() (decl Decl) {
	decl = r.queue[0]
	r.queue = r.queue[1:]
	return
}

func (r *DependencyResolver) CheckEach(initial Decl) {
	r.use(initial, nil)

	for len(r.queue) > 0 {
		decl := r.pop()
		r.Check(decl)
		r.CheckEmbedded(decl)
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
			if sel, ok := decl.Pkg().Info().Selections[node]; ok {
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
					r.use(d, decl)
				}

				return true
			}

			if id, ok := node.X.(*ast.Ident); ok && id != nil {
				uses := decl.Pkg().Info().Uses[id]
				switch uses := uses.(type) {
				case *types.PkgName:
					imported := uses.Imported()
					alias, name, path := id.Name, imported.Name(), imported.Path()
					if name == alias {
						alias = ""
					}

					i := r.iset.GetOrCreate(alias, name, path)
					r.useImport(i)
					pkg, ok := r.pset.Get(path)

					if !ok || isBuiltinPackage(path) {
						return true
					}

					d, ok := r.dset.Get(pkg, node.Sel.Name)
					if ok {
						r.use(d, decl)
					}
				}
			}

		case *ast.Ident:
			if _, ok := decl.Pkg().GetObject(node.Name); !ok {
				// break when the object is
				//   - not a package-level declaration
				//   - an external package object
				return true
			}

			uses, ok := decl.Pkg().Info().Uses[node]
			if !ok || uses.Pkg() == nil || isBuiltinPackage(uses.Pkg().Path()) {
				return true
			}

			switch obj := uses.(type) {
			case *types.Const, *types.Var, *types.Func, *types.TypeName:
				d, _ := r.dset.Get(decl.Pkg(), obj.Name())
				r.use(d, decl)
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

	// dependency checking of mdecl.Type() should be done before
	mdecl.Type().GetUses().Each(func(d Decl) {
		tdecl, ok := d.(*TypeDecl)
		if !ok {
			return
		}

		m, ok := tdecl.GetMethod(mdecl)
		if ok {
			r.use(m, nil)
		}
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
