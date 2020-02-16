package gollect

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func getAstFiles(fset *token.FileSet, glob string) map[string][]*ast.File {
	g, err := filepath.Glob(glob)
	if err != nil {
		log.Fatalf("parse glob: %v", err)
	}

	res := make(map[string][]*ast.File)
	processed := make(map[string]struct{})

	type pathInfo struct{ packagePath, filepath string }
	targets := make([]*pathInfo, len(g))
	for i, p := range g {
		targets[i] = &pathInfo{
			packagePath: "main",
			filepath:    p,
		}
	}

	for ; len(targets) > 0; targets = targets[1:] {
		info := targets[0]
		processed[info.filepath] = struct{}{}

		f, err := parser.ParseFile(fset, info.filepath, nil, 0)
		if err != nil {
			log.Fatalf("parse file: %v", err)
		}

		res[info.packagePath] = append(res[info.packagePath], f)

		for _, i := range f.Imports {
			pp := trimQuotes(i.Path.Value)
			if _, ok := processed[pp]; ok || isBuiltinPackage(pp) {
				continue
			}

			for _, pfp := range findPackageFilePaths(pp) {
				targets = append(targets, &pathInfo{
					packagePath: pp,
					filepath:    pfp,
				})
			}
		}
	}

	return res
}

func findPackageFilePaths(names ...string) (paths []string) {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, names...)
	if err != nil {
		log.Fatalf("load: %v\n", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		paths = append(paths, pkg.GoFiles...)
	}
	return
}
