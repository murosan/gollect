package gollect

import (
	"go/token"
	"os"
)

func Main(glob string) {
	fset := token.NewFileSet()
	imports := make(ImportSet)
	main := "main"

	// parse ast files and check dependencies
	packages := NewPackages(fset, imports, main, glob)
	AnalyzeForeach(fset, packages)

	// mark all used declarations
	next := []ExternalDependencySet{{}}
	next[0].Add(main, main)
	UseAll(packages, next)

	if err := Write(os.Stdout, fset, packages, imports); err != nil {
		panic(err)
	}
}

func init() {
	initBuiltinPackages()
}
