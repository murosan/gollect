package gollect

import (
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestIsBuiltin(t *testing.T) {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	for _, pkg := range pkgs {
		m[pkg.PkgPath] = struct{}{}
	}

	eq := reflect.DeepEqual(m, builtinPackages)
	if !eq {
		t.Error("please update builtinPackages")
	}
}
