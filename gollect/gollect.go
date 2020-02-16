package gollect

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
)

func Main(glob string) {
	fset := token.NewFileSet()
	packages := make(map[string]*Package)
	imports := make(ImportSet)

	astFiles := getAstFiles(fset, glob)

	for path, files := range astFiles {
		pkg := NewPackage(path, files, imports)
		packages[path] = pkg

		if err := ExecCheck(fset, pkg); err != nil {
			log.Fatalln(err)
		}
		pkg.InitObjects()
		ResolveDependency(pkg)
	}

	main := "main"
	next := []ExternalDependencySet{make(ExternalDependencySet)}
	next[0].Add(main, main)
	for ; len(next) > 0; next = next[1:] {
		for ed := range next[0] {
			deps := packages[ed.path].Dependencies()
			next = append(next, deps.Use(ed.name)...)
		}
	}

	astMain := astFiles[main][0]
	astMain.Decls = FilterDecls(packages[main].Dependencies(), astMain.Decls)

	RemoveExternalIdents(astMain, imports)

	importDecl := &ast.GenDecl{Tok: token.IMPORT}
	for _, i := range imports {
		if i.IsUsed() && i.IsBuiltin() {
			importDecl.Specs = append(importDecl.Specs, i.ToSpec())
		}
	}
	astMain.Decls = append([]ast.Decl{importDecl}, astMain.Decls...)
	ast.SortImports(fset, astMain)

	var buf bytes.Buffer
	format.Node(&buf, fset, astMain)

	for path, files := range astFiles {
		if path != main {
			for _, file := range files {
				decls := FilterDecls(packages[path].Dependencies(), file.Decls)
				for _, d := range decls {
					RemoveExternalIdents(d, imports)
				}

				buf.Write([]byte("\n"))
				format.Node(&buf, fset, decls)
				buf.Write([]byte("\n"))
			}
		}
	}

	buf.WriteTo(os.Stdout)
}

func init() {
	initBuiltinPackages()
}
