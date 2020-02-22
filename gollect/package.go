package gollect

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
)

type (
	Package struct {
		path    string                 // package path
		files   []*ast.File            // container of ast files
		imports ImportSet              // shared in global
		objects map[string]*ast.Object // map of package-level objects

		info *types.Info  // uses info
		deps Dependencies // pairs of ident name and Dependency
	}

	Packages map[string]*Package
)

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

func (pkg *Package) InitObjects() {
	for _, file := range pkg.files {
		for k, v := range file.Scope.Objects {
			pkg.objects[k] = v
		}
	}
}

func (pkg *Package) Dependencies() Dependencies { return pkg.deps }

func NewPackages(
	fset *token.FileSet,
	imports ImportSet,
	packagePath,
	glob string,
) Packages {
	paths, err := filepath.Glob(glob)
	if err != nil {
		panic(fmt.Errorf("parse glob: %w", err))
	}

	packages := make(Packages)
	ParseAll(
		packages,
		fset,
		imports,
		packagePath,
		paths...,
	)

	return packages
}

func (p Packages) Set(path string, pkg *Package) {
	p[path] = pkg
}

func (p Packages) Get(path string) (*Package, bool) {
	pkg, ok := p[path]
	return pkg, ok
}
