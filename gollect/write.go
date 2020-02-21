package gollect

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

func Write(
	w io.Writer,
	fset *token.FileSet,
	packages map[string]*Package,
	imports ImportSet,
) error {
	var buf bytes.Buffer

	// get head of main package's ast files
	// treat this as base ast
	mainPackage := packages["main"]
	main := mainPackage.files[0]

	// delete unused codes and all imports from base ast
	main.Decls = FilterDecls(mainPackage.Dependencies(), main.Decls)
	RemoveExternalIdents(main, imports)

	// build new import decl and push it to head of decls
	i := importDecl(imports)
	main.Decls = append([]ast.Decl{i}, main.Decls...)

	if err := format.Node(&buf, fset, main); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	for path, pkg := range packages {
		for _, file := range pkg.files {
			if file != main {
				decls := FilterDecls(packages[path].Dependencies(), file.Decls)
				for _, d := range decls {
					RemoveExternalIdents(d, imports)
				}

				buf.Write([]byte("\n"))
				if err := format.Node(&buf, fset, decls); err != nil {
					return fmt.Errorf("format: %w", err)
				}
				if len(decls) != 0 {
					if _, ok := decls[len(decls)-1].(*ast.GenDecl); !ok {
						buf.Write([]byte("\n"))
					}
				}
			}
		}
	}

	_, err := buf.WriteTo(w)
	return err
}

func importDecl(imports ImportSet) *ast.GenDecl {
	d := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1, // if zero, imports will not be sorted
	}
	for _, i := range imports {
		if i.IsUsed() && i.IsBuiltin() {
			d.Specs = append(d.Specs, i.ToSpec())
		}
	}
	return d
}
