// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/ast"
	"go/types"
	"sync"
)

type (
	// Package represents analyzing information.
	Package struct {
		sync.Mutex

		path    string                 // package path
		files   []*ast.File            // container of ast files
		imports *ImportSet             // shared in global
		objects map[string]*ast.Object // map of package-level objects
		info    *types.Info            // uses info
		deps    Dependencies           // pairs of ident name and Dependency
	}

	// Packages is a map of Package.
	Packages map[string]*Package
)

// NewPackage returns new Package.
func NewPackage(path string, imports *ImportSet) *Package {
	return &Package{
		path:    path,
		files:   nil,
		imports: imports,
		objects: make(map[string]*ast.Object),
		info: &types.Info{
			Uses: make(map[*ast.Ident]types.Object),
			// Types:      make(map[ast.Expr]types.TypeAndValue),
			// Defs:       make(map[*ast.Ident]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
		},
		deps: make(Dependencies),
	}
}

// InitObjects compiles all files' objects into one map.
// This is called after parsing all ast files and before
// start analyzing dependencies.
func (pkg *Package) InitObjects() {
	for _, file := range pkg.files {
		for k, v := range file.Scope.Objects {
			pkg.objects[k] = v
		}
	}
}

// Dependencies returns dependencies.
func (pkg *Package) Dependencies() Dependencies { return pkg.deps }

// PushAstFile push ast.File to files.
func (pkg *Package) PushAstFile(f *ast.File) {
	pkg.Lock()
	pkg.files = append(pkg.files, f)
	pkg.Unlock()
}

// Set sets the Package to set.
func (p Packages) Set(path string, pkg *Package) { p[path] = pkg }

// Get gets a Package from set.
func (p Packages) Get(path string) (*Package, bool) {
	pkg, ok := p[path]
	return pkg, ok
}
