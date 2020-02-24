package gollect

import (
	"fmt"
	"go/ast"
	"go/format"
	"io"
)

func Write(w io.Writer, program *Program) error {
	fset, iset, packages := program.FileSet(), program.ImportSet(), program.Packages()

	// get head of main package's ast files
	// treat this as base ast
	mainPackage := packages["main"]
	main := mainPackage.files[0]

	// delete unused codes and all imports from base ast
	main.Decls = FilterDecls(mainPackage.Dependencies(), main.Decls)
	RemoveExternalIdents(main, iset)

	// build new import decl and push it to head of decls
	main.Decls = append([]ast.Decl{iset.ToDecl()}, main.Decls...)

	if err := format.Node(w, fset, main); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	for path, pkg := range packages {
		for _, file := range pkg.files {
			if file != main {
				decls := FilterDecls(packages[path].Dependencies(), file.Decls)
				for _, d := range decls {
					RemoveExternalIdents(d, iset)
				}

				if len(decls) != 0 {
					if _, err := w.Write([]byte("\n")); err != nil {
						return err
					}
					if err := format.Node(w, fset, decls); err != nil {
						return fmt.Errorf("format: %w", err)
					}
					if _, ok := decls[len(decls)-1].(*ast.GenDecl); !ok {
						if _, err := w.Write([]byte("\n")); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
