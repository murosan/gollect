package gollect

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"

	"github.com/atotto/clipboard"
)

func Write(w io.Writer, program *Program) error {
	var buf bytes.Buffer

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

	if err := format.Node(&buf, fset, main); err != nil {
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
					buf.Write([]byte("\n"))
					if err := format.Node(&buf, fset, decls); err != nil {
						return fmt.Errorf("format: %w", err)
					}
					if _, ok := decls[len(decls)-1].(*ast.GenDecl); !ok {
						buf.Write([]byte("\n"))
					}
				}
			}
		}
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	// TODO: refactor
	if !clipboard.Unsupported {
		if err := clipboard.WriteAll(buf.String()); err != nil {
			return err
		}
	}
	return nil
}
