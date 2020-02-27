package gollect

import (
	"go/ast"
	"go/types"
)

type (
	// Package represets analyzing information.
	Package struct {
		path    string                 // package path
		files   []*ast.File            // container of ast files
		imports ImportSet              // shared in global
		objects map[string]*ast.Object // map of package-level objects
		info    *types.Info            // uses info
		deps    Dependencies           // pairs of ident name and Dependency
	}

	// Packages is a map of Package.
	Packages map[string]*Package
)

// NewPackage returns new Package.
func NewPackage(path string, imports ImportSet) *Package {
	return &Package{
		path:    path,
		files:   nil,
		imports: imports,
		objects: make(map[string]*ast.Object),
		info: &types.Info{
			Uses: make(map[*ast.Ident]types.Object),
		},
		deps: make(Dependencies),
	}
}

// InitObjects compiles all files' objests into one map.
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

// Set sets the Package to set.
func (p Packages) Set(path string, pkg *Package) { p[path] = pkg }

// Get gets a Package from set.
func (p Packages) Get(path string) (*Package, bool) {
	pkg, ok := p[path]
	return pkg, ok
}
