package gollect

import (
	"os"
)

func Main(glob string) {
	p := NewProgram(glob)

	// parse ast files and check dependencies
	ParseAll(p)
	AnalyzeForeach(p)

	// mark all used declarations
	next := []ExternalDependencySet{{}}
	next[0].Add("main", "main")
	UseAll(p.Packages(), next)

	if err := Write(os.Stdout, p); err != nil {
		panic(err)
	}
}

func init() {
	initBuiltinPackages()
}
