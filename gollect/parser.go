package gollect

import (
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func ParseAll(
	packages Packages,
	fset *token.FileSet,
	imports ImportSet,
	packagePath string,
	filePaths ...string,
) {
	paths := []string{packagePath}
	for ; len(paths) > 0; paths = paths[1:] {
		pp := paths[0]
		if _, ok := packages[pp]; ok {
			continue
		}

		var fp []string
		if pp == packagePath {
			fp = filePaths
		} else {
			fp = FindFilePaths(pp)
		}

		pkg := NewPackage(pp, imports)
		packages[pp] = pkg

		ParseAst(fset, pkg, fp...)
		paths = append(paths, NextPackagePaths(pkg)...)
	}
}

func ParseAst(fset *token.FileSet, p *Package, paths ...string) {
	for _, path := range paths {
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("parse file (path = %s): %v", path, err)
		}

		p.files = append(p.files, f)
	}
}

func FindFilePaths(path string) (paths []string) {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		log.Fatalf("load: %v\n", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		paths = append(paths, pkg.GoFiles...)
	}
	return
}

func NextPackagePaths(p *Package) (paths []string) {
	m := make(map[string]interface{})
	for _, f := range p.files {
		for _, i := range f.Imports {
			p := trimQuotes(i.Path.Value)
			if _, ok := m[p]; !ok && !isBuiltinPackage(p) {
				m[p] = struct{}{}
				paths = append(paths, p)
			}
		}
	}
	return
}
