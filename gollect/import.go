package gollect

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type (
	// Import represents import.
	Import struct {
		alias, name, path string
		used              bool
	}

	// ImportSet is a set of Import.
	ImportSet map[string]*Import
)

// NewImport returns new Import.
func NewImport(alias, name, path string) *Import {
	return &Import{
		alias: alias,
		name:  name,
		path:  path,
		used:  false,
	}
}

// AliasOrName returns alias if the alias is non-empty otherwise returns name.
// The return value is expected to be used as a key of ImportSet.
func (i *Import) AliasOrName() string {
	if i.alias != "" {
		return i.alias
	}
	return i.name
}

// ToSpec creates and returns ast.ImportSpec.
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

// Use changes used state to true.
func (i *Import) Use() { i.used = true }

// IsUsed returns used state.
func (i *Import) IsUsed() bool { return i.used }

// IsBuiltin returns if the import's path is Go language's builtin or not.
func (i *Import) IsBuiltin() bool { return isBuiltinPackage(i.path) }

func (i *Import) String() string {
	return fmt.Sprintf("{alias: %s, name: %s, path: %s}",
		strconv.Quote(i.alias),
		strconv.Quote(i.name),
		strconv.Quote(i.path),
	)
}

// Add adds the Import to set.
func (s ImportSet) Add(i *Import) { s[i.AliasOrName()] = i }

// Get gets an Import from set.
func (s ImportSet) Get(key string) (*Import, bool) {
	v, ok := s[key]
	return v, ok
}

// GetOrCreate gets an Import form set, if the set has no matched value
// creates new one.
func (s ImportSet) GetOrCreate(alias, name, path string) *Import {
	i := NewImport(alias, name, path)
	if v, ok := s.Get(i.AliasOrName()); ok {
		return v
	}
	s.Add(i)
	return i
}

// ToDecl creates ast.GenDecl and returns it.
func (s ImportSet) ToDecl() *ast.GenDecl {
	d := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1, // if zero, imports will not be sorted
	}

	for _, i := range s {
		if i.IsUsed() && i.IsBuiltin() {
			d.Specs = append(d.Specs, i.ToSpec())
		}
	}

	return d
}
