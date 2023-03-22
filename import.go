// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	ImportSet struct{ iset []*Import }
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

// ToSpec creates and returns ast.ImportSpec.
func (i *Import) ToSpec() *ast.ImportSpec {
	var s ast.ImportSpec
	if i.alias != "" {
		s.Name = ast.NewIdent(i.alias)
	}
	s.Path = &ast.BasicLit{Value: strconv.Quote(i.path)}
	return &s
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

// NewImportSet returns new ImportSet.
func NewImportSet() *ImportSet {
	return &ImportSet{iset: make([]*Import, 0)}
}

// Len returns length of map.
func (s *ImportSet) Len() int { return len(s.iset) }

// Values apply block to each element.
func (s *ImportSet) Values() []*Import {
	a := make([]*Import, len(s.iset))
	copy(a, s.iset)
	return a
}

// AddAndGet gets an Import form set if exists, otherwise
// creates new one and returns it.
func (s *ImportSet) AddAndGet(i *Import) *Import {
	for _, v := range s.iset {
		if v.alias == i.alias &&
			v.name == i.name &&
			v.path == i.path {
			return v
		}
	}

	s.iset = append(s.iset, i)
	return i
}

// GetOrCreate gets an Import form set if exists, otherwise
// creates new one and returns it.
func (s *ImportSet) GetOrCreate(alias, name, path string) *Import {
	return s.AddAndGet(NewImport(alias, name, path))
}

// ToDecl creates ast.GenDecl and returns it.
func (s *ImportSet) ToDecl() *ast.GenDecl {
	d := &ast.GenDecl{Tok: token.IMPORT}

	for _, i := range s.iset {
		if i.IsUsed() && i.IsBuiltin() {
			d.Specs = append(d.Specs, i.ToSpec())
		}
	}

	if len(d.Specs) > 1 {
		// if there is one import and Lparen value is 0,
		// generated import will be in a single line.
		// ex. import "fmt"
		//
		// if there are multiple imports, Lparen value should not be 0
		// to sort them by format.Node().
		d.Lparen = 1
	}

	return d
}

func (s *ImportSet) String() string {
	return fmt.Sprint(s.iset)
}
