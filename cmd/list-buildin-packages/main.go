package main

import (
	"os"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

/*
https://github.com/murosan/list-builtin-packages
```
var builtinPackages = map[string]interface{}{
	"archive/tar": struct{}{},
	"archive/zip": struct{}{},
	"bufio": struct{}{},
...
```
*/

var builtinPackages []string

func main() {
	tmpl, err := template.New("a").Parse(`var builtinPackages = map[string]interface{}{{"{"}}{{ range . }}
	"{{ . }}": struct{}{},{{ end }}
}
`)

	if err != nil {
		panic(err)
	}

	sort.Strings(builtinPackages)

	err = tmpl.Execute(os.Stdout, builtinPackages)
	if err != nil {
		panic(err)
	}
}

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		p := pkg.PkgPath
		if !strings.Contains(p, "internal/") {
			builtinPackages = append(builtinPackages, p)
		}
	}
}
