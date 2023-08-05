// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type (
	// Import represents import.
	Import struct {
		alias, name, path string
		used              bool
	}

	// DotImport struct {
	// 	pkg  *types.Package
	// 	used bool
	// }

	// ImportSet is a set of Import.
	ImportSet struct {
		set map[isetKey]*Import
		// dots map[string]*DotImport
	}

	isetKey struct{ alias, name, path string }
)

// NewImport returns new Import.
func NewImport(alias, name, path string) *Import {
	return &Import{
		alias: alias,
		name:  name,
		path:  path,
		used:  false,
	}
}

// ToSpec creates and returns ast.ImportSpec.
func (i *Import) ToSpec() *ast.ImportSpec {
	var s ast.ImportSpec
	if i.alias != "" {
		s.Name = ast.NewIdent(i.alias)
	}
	s.Path = &ast.BasicLit{Value: strconv.Quote(i.path)}
	return &s
}

// Use changes used state to true.
func (i *Import) Use() { i.used = true }

// IsBuiltin returns if the import's path is Go language's builtin or not.
func (i *Import) IsBuiltin() bool { return isBuiltinPackage(i.path) }

func (i *Import) key() isetKey {
	return isetKey{alias: i.alias, name: i.name, path: i.path}
}

func (i *Import) String() string {
	return fmt.Sprintf("{alias: %s, name: %s, path: %s}",
		strconv.Quote(i.alias),
		strconv.Quote(i.name),
		strconv.Quote(i.path),
	)
}

// func NewDotImport(pkg *types.Package) *DotImport {
// 	return &DotImport{pkg: pkg}
// }
// func (i *DotImport) Use()            { i.used = true }
// func (i *DotImport) IsBuiltin() bool { return isBuiltinPackage(i.pkg.Path()) }
// func (i *DotImport) ToSpec() *ast.ImportSpec {
// 	return &ast.ImportSpec{
// 		Name: ast.NewIdent("."),
// 		Path: &ast.BasicLit{Value: strconv.Quote(i.pkg.Path())},
// 	}
// }
// func (i *DotImport) String() string {
// 	return fmt.Sprintf("{alias: %s, name: %s, path: %s}",
// 		strconv.Quote("."),
// 		strconv.Quote(i.pkg.Name()),
// 		strconv.Quote(i.pkg.Path()),
// 	)
// }

// NewImportSet returns new ImportSet.
func NewImportSet() *ImportSet {
	return &ImportSet{
		set: make(map[isetKey]*Import),
		// dots: make(map[string]*DotImport),
	}
}

// AddAndGet gets an Import form set if exists, otherwise
// creates new one and returns it.
func (s *ImportSet) AddAndGet(i *Import) *Import {
	key := i.key()
	v, ok := s.set[key]
	if ok {
		return v
	}

	s.set[key] = i
	return i
}

// GetOrCreate gets an Import form set if exists, otherwise
// creates new one and returns it.
func (s *ImportSet) GetOrCreate(alias, name, path string) *Import {
	return s.AddAndGet(NewImport(alias, name, path))
}

// func (s *ImportSet) AddDotImport(pkg *types.Package) {
// 	i := NewDotImport(pkg)
// 	if _, ok := s.dots[i.pkg.Path()]; !ok {
// 		s.dots[i.pkg.Path()] = i
// 	}
// }
//
// func (s *ImportSet) EachDotImports(f func(i *DotImport) bool) {
// 	for _, i := range s.dots {
// 		if !f(i) {
// 			break
// 		}
// 	}
// }

// ToDecl creates ast.GenDecl and returns it.
func (s *ImportSet) ToDecl() *ast.GenDecl {
	d := &ast.GenDecl{Tok: token.IMPORT}

	for _, i := range s.set {
		if i.used && i.IsBuiltin() {
			d.Specs = append(d.Specs, i.ToSpec())
		}
	}
	// for _, i := range s.dots {
	// 	if i.used && i.IsBuiltin() {
	// 		d.Specs = append(d.Specs, i.ToSpec())
	// 	}
	// }

	if len(d.Specs) > 1 {
		// if there is one import and Lparen value is 0,
		// generated import will be in a single line.
		// ex. import "fmt"
		//
		// if there are multiple imports, Lparen value should not be 0
		// to sort them by format.Node().
		d.Lparen = 1
	}

	return d
}

func (s *ImportSet) String() string {
	var v []string
	for _, i := range s.set {
		v = append(v, i.String())
	}
	// for _, i := range s.dots {
	// 	v = append(v, i.String())
	// }
	return fmt.Sprint(v)
}
