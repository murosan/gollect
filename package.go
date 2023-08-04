// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/types"
)

// Package represents analyzing information.
type Package struct {
	path    string                 // package path
	files   []*ast.File            // container of ast files
	objects map[string]*ast.Object // map of package-level objects
	info    *types.Info            // uses info
}

// NewPackage returns new Package.
func NewPackage(path string) *Package {
	return &Package{
		path:    path,
		files:   nil,
		objects: make(map[string]*ast.Object),
		info: &types.Info{
			Uses: make(map[*ast.Ident]types.Object),
			// Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
		},
	}
}

// Path returns package path.
func (pkg *Package) Path() string { return pkg.path }

// InitObjects compiles all files' objects into one map.
// This is called after parsing all ast files and before
// start analyzing dependencies.
func (pkg *Package) InitObjects() {
	for _, file := range pkg.files {
		for k, v := range file.Scope.Objects {
			k, v := k, v
			pkg.objects[k] = v
		}
	}
}

// GetObject gets and returns object which scope is package-level.
func (pkg *Package) GetObject(key string) (*ast.Object, bool) {
	o, ok := pkg.objects[key]
	return o, ok
}

// PushAstFile push ast.File to files.
func (pkg *Package) PushAstFile(f *ast.File) {
	pkg.files = append(pkg.files, f)
}

// UsesInfo gets types.Object from types.Info and returns it.
func (pkg *Package) UsesInfo(i *ast.Ident) (types.Object, bool) {
	o, ok := pkg.info.Uses[i]
	return o, ok
}

// SelInfo gets types.Selection from types.Info and returns it.
func (pkg *Package) SelInfo(expr *ast.SelectorExpr) (*types.Selection, bool) {
	s, ok := pkg.info.Selections[expr]
	return s, ok
}

// DefInfo gets types.Object from types.Info and returns it.
func (pkg *Package) DefInfo(i *ast.Ident) (types.Object, bool) {
	def, ok := pkg.info.Defs[i]
	return def, ok
}

func (pkg *Package) ObjectOf(i *ast.Ident) types.Object {
	return pkg.info.ObjectOf(i)
}

// PackageSet is a map of Package.
type PackageSet map[string]*Package

// Add sets the Package to set.
func (p PackageSet) Add(path string, pkg *Package) { p[path] = pkg }

// Get gets a Package from set.
func (p PackageSet) Get(path string) (*Package, bool) {
	pkg, ok := p[path]
	return pkg, ok
}
