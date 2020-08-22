package gollect

import (
	"reflect"
	"strings"
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
		p := pkg.PkgPath
		if !strings.Contains(p, "internal/") {
			m[pkg.PkgPath] = struct{}{}
		}
	}

	eq := reflect.DeepEqual(m, builtinPackages)
	if !eq {
		t.Error("please update builtinPackages")
	}
}
