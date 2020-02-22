package gollect

import (
	"golang.org/x/tools/go/packages"
)

var builtinPackages = make(map[string]interface{})

func initBuiltinPackages() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		builtinPackages[p.PkgPath] = struct{}{}
	}
}

func isBuiltinPackage(path string) bool {
	_, ok := builtinPackages[path]
	return ok
}
