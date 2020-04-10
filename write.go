// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"go/format"
	"io"
)

// Write writes filtered and formatted code to io.Writer.
func Write(w io.Writer, program *Program) error {
	fset, dset := program.FileSet(), program.DeclSet()
	iset, pset := program.ImportSet(), program.PackageSet()

	// get head of main package's ast files
	// treat this as base ast
	mainPackage := pset["main"]
	main := mainPackage.files[0]

	filter := NewFilter(dset, mainPackage)

	// delete unused codes and all imports from base ast
	main.Decls = filter.Decls(main.Decls)
	filter.PackageSelectorExpr(main)

	// build new import decl and push it to head of decls
	if ispec := iset.ToDecl(); len(ispec.Specs) != 0 {
		main.Decls = append([]ast.Decl{ispec}, main.Decls...)
	}

	if err := format.Node(w, fset, main); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	for path, pkg := range pset {
		for _, file := range pkg.files {
			if file == main {
				continue
			}

			filter := NewFilter(dset, pset[path])
			decls := filter.Decls(file.Decls)
			for _, d := range decls {
				filter.PackageSelectorExpr(d)
			}

			if len(decls) != 0 {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
				if err := format.Node(w, fset, decls); err != nil {
					return fmt.Errorf("format: %w", err)
				}
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
