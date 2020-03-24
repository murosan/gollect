// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"sync"
)

type (
	// Import represents import.
	Import struct {
		alias, name, path string
		used              bool
	}

	// ImportSet is a set of Import.
	ImportSet struct {
		sync.RWMutex
		m map[string]*Import
	}
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
	return &ImportSet{m: make(map[string]*Import)}
}

// Len returns length of map.
func (s *ImportSet) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

// Foreach apply block to each element.
func (s *ImportSet) Foreach(block func(*Import)) {
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.m {
		block(v)
	}
}

// Add adds the Import to set.
func (s *ImportSet) Add(i *Import) {
	s.Lock()
	s.m[i.AliasOrName()] = i
	s.Unlock()
}

// Get gets an Import from set.
func (s *ImportSet) Get(key string) (*Import, bool) {
	s.RLock()
	v, ok := s.m[key]
	s.RUnlock()
	return v, ok
}

// GetOrCreate gets an Import form set, if the set has no matched value
// creates new one.
func (s *ImportSet) GetOrCreate(alias, name, path string) *Import {
	i := NewImport(alias, name, path)
	if v, ok := s.Get(i.AliasOrName()); ok {
		return v
	}
	s.Add(i)
	return i
}

// ToDecl creates ast.GenDecl and returns it.
func (s *ImportSet) ToDecl() *ast.GenDecl {
	d := &ast.GenDecl{Tok: token.IMPORT}

	s.RLock()
	defer s.RUnlock()

	for _, i := range s.m {
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
