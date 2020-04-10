// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"strings"
	"sync"
)

// DeclType represents declaration type.
type DeclType int

func (t DeclType) String() string { return fmt.Sprint(int(t)) }

const (
	_ DeclType = iota
	// DecCommon represents common declaration type.
	// e.g, Var, Const, Func
	DecCommon
	// DecType represents type declaration type.
	// e.g, Struct, Interface
	DecType
	// DecMethod represents method declaration.
	DecMethod
)

// Decl represents a declaration.
type Decl interface {
	ID() string
	Node() ast.Node
	SetNode(ast.Node)
	Pkg() *Package
	IsUsed() bool
	Use()
	Uses(Decl)
	UsesImport(*Import)
}

const sep = ";"

func makeID(p *Package, s ...string) string {
	return p.Path() + sep + strings.Join(s, sep)
}

// NewDecl return new Decl
func NewDecl(t DeclType, pkg *Package, id ...string) Decl {
	switch t {
	case DecCommon:
		return NewCommonDecl(pkg, id...)
	case DecType:
		return NewTypeDecl(pkg, id...)
	case DecMethod:
		return NewMethodDecl(pkg, id...)
	default:
		panic(fmt.Sprintf("unknown DeclType %s", t))
	}
}

// CommonDecl represents one of Var, Const, Func declaration.
type CommonDecl struct {
	id   string
	node ast.Node
	pkg  *Package
	used bool
	uses struct {
		decls   *DeclSet
		imports *ImportSet
	}
}

// NewCommonDecl returns new CommonDecl
func NewCommonDecl(pkg *Package, id ...string) *CommonDecl {
	return &CommonDecl{
		id:   makeID(pkg, id...),
		node: nil,
		pkg:  pkg,
		used: false,
		uses: struct {
			decls   *DeclSet
			imports *ImportSet
		}{
			decls:   NewDeclSet(),
			imports: NewImportSet(),
		},
	}
}

// ID returns id made by makeID(package-path, declName).
func (d *CommonDecl) ID() string { return d.id }

// Node returns ast.Node. Its field is initialized lazily.
func (d *CommonDecl) Node() ast.Node { return d.node }

// SetNode sets node.
func (d *CommonDecl) SetNode(n ast.Node) { d.node = n }

// Pkg returns Package.
func (d *CommonDecl) Pkg() *Package { return d.pkg }

// IsUsed returns true if it is used from main package.
func (d *CommonDecl) IsUsed() bool { return d.used }

// Uses sets given decl to dependency map.
func (d *CommonDecl) Uses(decl Decl) { d.uses.decls.Add(decl) }

// UsesImport sets given import to dependency map.
func (d *CommonDecl) UsesImport(i *Import) { d.uses.imports.Add(i) }

// Use change this and its dependencies' used field to true.
func (d *CommonDecl) Use() {
	if d.IsUsed() {
		return
	}
	d.used = true
	for _, i := range d.uses.imports.Values() {
		i.Use()
	}
	for _, d := range d.uses.decls.Values() {
		d.Use()
	}
}

// TypeDecl represents Type declaration.
type TypeDecl struct {
	*CommonDecl
	methods struct {
		mset map[string]*MethodDecl
		keep bool
	}
}

// NewTypeDecl returns new TypeDecl.
func NewTypeDecl(pkg *Package, id ...string) *TypeDecl {
	return &TypeDecl{
		CommonDecl: NewCommonDecl(pkg, id...),
		methods: struct {
			mset map[string]*MethodDecl
			keep bool
		}{
			mset: make(map[string]*MethodDecl),
			keep: false,
		},
	}
}

// Use change this, its dependencies' and its methods' used field to true.
func (d *TypeDecl) Use() {
	if d.IsUsed() {
		return
	}
	d.CommonDecl.Use()
	if d.methods.keep {
		for _, m := range d.methods.mset {
			m.Use()
		}
	}
}

// SetMethod sets given method to methods set.
func (d *TypeDecl) SetMethod(m *MethodDecl) { d.methods.mset[m.ID()] = m }

// KeepMethod set true its keep method option.
// When the field is true, all methods will not removed even the method
// is not used from main.
func (d *TypeDecl) KeepMethod() { d.methods.keep = true }

// MethodDecl represents method declaration.
type MethodDecl struct {
	*CommonDecl
	tpe         *TypeDecl
	inheritFrom *MethodDecl
	embedded    bool
}

// NewMethodDecl returns new MethodDecl.
func NewMethodDecl(pkg *Package, id ...string) *MethodDecl {
	return &MethodDecl{
		CommonDecl:  NewCommonDecl(pkg, id...),
		tpe:         nil,
		embedded:    false,
		inheritFrom: nil,
	}
}

// Type returns TypeDecl this method belongs to. The field is initialized lazily.
func (d *MethodDecl) Type() *TypeDecl { return d.tpe }

// SetType sets given TypeDecl to field.
func (d *MethodDecl) SetType(t *TypeDecl) { d.tpe = t }

// IsEmbedded returns true if it is embedded method.
// TODO: this is not used for now.
func (d *MethodDecl) IsEmbedded() bool { return d.embedded }

// SetEmbedded change its embedded field to true.
// TODO: this is not used for now.
func (d *MethodDecl) SetEmbedded(b bool) { d.embedded = b }

// SetInheritFrom sets given method to inheritFrom field.
// TODO: this is not used for now.
func (d *MethodDecl) SetInheritFrom(m *MethodDecl) { d.inheritFrom = m }

// Use change this, its dependencies' and its methods' used field to true.
func (d *MethodDecl) Use() {
	if d.IsUsed() {
		return
	}
	d.CommonDecl.Use()
	d.Type().Use()
	if d.IsEmbedded() {
		d.inheritFrom.Use()
	}
}

// DeclSet is a set of Decl
type DeclSet struct {
	mux  sync.RWMutex
	dset map[string]Decl
}

// NewDeclSet returns new DeclSet
func NewDeclSet() *DeclSet {
	return &DeclSet{dset: make(map[string]Decl)}
}

// Get gets Decl from set.
func (s *DeclSet) Get(pkg *Package, key ...string) (Decl, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	d, ok := s.dset[makeID(pkg, key...)]
	return d, ok
}

// GetOrCreate gets Decl from set if exists, otherwise create new one and add to set
// then returns it.
func (s *DeclSet) GetOrCreate(dtype DeclType, pkg *Package, key ...string) Decl {
	s.mux.Lock()
	defer s.mux.Unlock()

	if d, ok := s.dset[makeID(pkg, key...)]; ok {
		return d
	}

	d := NewDecl(dtype, pkg, key...)
	s.dset[d.ID()] = d
	return d
}

// Add adds Decl to set.
func (s *DeclSet) Add(d Decl) {
	s.mux.Lock()
	s.dset[d.ID()] = d
	s.mux.Unlock()
}

// Values creates a slice of values of set and returns it.
func (s *DeclSet) Values() []Decl {
	s.mux.RLock()
	defer s.mux.RUnlock()
	i := 0
	a := make([]Decl, len(s.dset))
	for _, v := range s.dset {
		v := v
		a[i] = v
		i++
	}
	return a
}
