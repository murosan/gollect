package gollect

/*
This test is used for only confirmation when the Go version changes.

Packages differ depending on the environment, but we ignore it
because they are not so important in competition programming.
*/

/*
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
*/
