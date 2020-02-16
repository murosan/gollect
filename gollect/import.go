package gollect

import (
	"fmt"
	"go/ast"
	"strconv"
)

type (
	Import struct {
		alias, name, path string
		used              bool
	}

	ImportSet map[string]*Import
)

func NewImport(alias, name, path string) *Import {
	return &Import{
		alias: alias,
		name:  name,
		path:  path,
		used:  false,
	}
}

func (i *Import) AliasOrName() string {
	if i.alias != "" {
		return i.alias
	}
	return i.name
}

func (i *Import) ToSpec() *ast.ImportSpec {
	var name *ast.Ident
	if i.alias != "" {
		name = ast.NewIdent(i.alias)
	}
	return &ast.ImportSpec{
		Name: name,
		Path: &ast.BasicLit{Value: strconv.Quote(i.path)},
	}
}

func (i *Import) Use() { i.used = true }

func (i *Import) IsUsed() bool { return i.used }

func (i *Import) IsBuiltin() bool { return isBuiltinPackage(i.path) }

func (i *Import) String() string {
	return fmt.Sprintf("{alias: %s, name: %s, path: %s}",
		strconv.Quote(i.alias),
		strconv.Quote(i.name),
		strconv.Quote(i.path),
	)
}

func (s ImportSet) Add(i *Import) { s[i.AliasOrName()] = i }

func (s ImportSet) Get(key string) (*Import, bool) {
	v, ok := s[key]
	return v, ok
}

func (s ImportSet) GetOrCreate(alias, name, path string) *Import {
	i := NewImport(alias, name, path)
	if v, ok := s.Get(i.AliasOrName()); ok {
		return v
	}
	s.Add(i)
	return i
}
