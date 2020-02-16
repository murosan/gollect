package gollect

import (
	"go/ast"
	"go/types"
)

type Package struct {
	path    string
	files   []*ast.File
	imports ImportSet
	objects map[string]*ast.Object // map of package level objects

	info *types.Info
	deps Dependencies // pairs of ident name and Dependency
}

func NewPackage(path string, files []*ast.File, imports ImportSet) *Package {
	return &Package{
		path:    path,
		files:   files,
		imports: imports,
		objects: nil,
		info: &types.Info{
			Uses: map[*ast.Ident]types.Object{},
		},
		deps: make(Dependencies),
	}
}

func (pkg *Package) InitObjects() {
	objects := make(map[string]*ast.Object)
	pkg.objects = objects
	for _, file := range pkg.files {
		for k, v := range file.Scope.Objects {
			objects[k] = v
		}
	}
}

func (pkg *Package) Dependencies() Dependencies { return pkg.deps }
