// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
)

// ParseAll parses all ast files and sets to Program's map.
// This also parses external imported package's ast.
func ParseAll(program *Program) {
	paths := []string{"main"}
	for ; len(paths) > 0; paths = paths[1:] {
		pp := paths[0]
		if _, ok := program.Packages().Get(pp); ok {
			continue
		}

		var fp []string
		if pp == "main" {
			fp = program.FilePaths()
		} else {
			fp = FindFilePaths(pp)
		}

		pkg := NewPackage(pp, program.ImportSet())
		program.Packages().Set(pp, pkg)

		ParseAst(program.FileSet(), pkg, fp...)
		if len(pkg.files) == 0 {
			panic(fmt.Sprintf("there are no files. paths=%v", fp))
		}
		paths = append(paths, NextPackagePaths(pkg)...)
	}
}

// ParseAst parses ast and pushes to files slice.
func ParseAst(fset *token.FileSet, p *Package, paths ...string) {
	for _, path := range paths {
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(fmt.Errorf("parse file (path = %s): %w", path, err))
		}

		p.files = append(p.files, f)
	}
}

// FindFilePaths finds filepaths from package path.
// https://pkg.go.dev/golang.org/x/tools/go/packages?tab=doc#example-package
func FindFilePaths(path string) (paths []string) {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(fmt.Errorf("load: %w", err))
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		paths = append(paths, pkg.GoFiles...)
	}
	return
}

// NextPackagePaths returns list of imported package paths.
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
