// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"go/ast"
	"strings"
)

// DeclType represents declaration type.
type DeclType int

func (t DeclType) String() string { return "DeclType(" + fmt.Sprint(int(t)) + ")" }

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
	fmt.Stringer
}

const sep = ";"

func makeID(p *Package, s ...string) string {
	return p.Path() + sep + strings.Join(s, sep)
}

func nameForUnderscore(id *ast.Ident) string {
	name := id.Name
	if name == "_" {
		name += fmt.Sprint(int(id.NamePos))
	}
	return name
}

// NewDecl return new Decl
func NewDecl(t DeclType, pkg *Package, ids ...string) Decl {
	switch t {
	case DecCommon:
		return NewCommonDecl(pkg, ids...)
	case DecType:
		return NewTypeDecl(pkg, ids...)
	case DecMethod:
		return NewMethodDecl(pkg, ids...)
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

// NewCommonDecl returns new CommonDecl.
// The length of ids must be one, and the value must be its name.
func NewCommonDecl(pkg *Package, ids ...string) *CommonDecl {
	return &CommonDecl{
		id:   makeID(pkg, ids...),
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
func (d *CommonDecl) UsesImport(i *Import) { d.uses.imports.AddAndGet(i) }

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

func (d *CommonDecl) String() string { return declToString(d) }

// TypeDecl represents Type declaration.
type TypeDecl struct {
	*CommonDecl
	methods struct {
		mset map[string]*MethodDecl
		keep bool
	}
}

// NewTypeDecl returns new TypeDecl.
// The length of ids must be one, and the value must be the type name.
func NewTypeDecl(pkg *Package, ids ...string) *TypeDecl {
	return &TypeDecl{
		CommonDecl: NewCommonDecl(pkg, ids...),
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

// UsesDecls returns list of decls that d uses.
func (d *TypeDecl) UsesDecls() []Decl {
	return d.uses.decls.Values()
}

// Methods returns methods as a slice.
func (d *TypeDecl) Methods() []*MethodDecl {
	v := make([]*MethodDecl, len(d.methods.mset))
	i := 0
	for _, m := range d.methods.mset {
		v[i] = m
		i++
	}
	return v
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
	name     string
	tpe      *TypeDecl
	embedded bool
}

// NewMethodDecl returns new MethodDecl. Length of ids must be two,
// the head value is its receiver's type name and second value is func name.
func NewMethodDecl(pkg *Package, ids ...string) *MethodDecl {
	return &MethodDecl{
		CommonDecl: NewCommonDecl(pkg, ids...),
		name:       ids[1],
		tpe:        nil,
		embedded:   false,
	}
}

// Name returns method name.
func (d *MethodDecl) Name() string { return d.name }

// Type returns TypeDecl this method belongs to. The field is initialized lazily.
func (d *MethodDecl) Type() *TypeDecl { return d.tpe }

// SetType sets given TypeDecl to field.
func (d *MethodDecl) SetType(t *TypeDecl) { d.tpe = t }

// IsEmbedded returns true if it is embedded method.
func (d *MethodDecl) IsEmbedded() bool { return d.embedded }

// SetEmbedded change its embedded field to true.
func (d *MethodDecl) SetEmbedded(b bool) { d.embedded = b }

// Use change this, its dependencies' and its methods' used field to true.
func (d *MethodDecl) Use() {
	if d.IsUsed() {
		return
	}
	d.CommonDecl.Use()
	d.Type().Use()
}

func declToString(decl Decl) string {
	var tpe string
	switch decl.(type) {
	case *CommonDecl:
		tpe = "CommonDecl"
	case *TypeDecl:
		tpe = "TypeDecl"
	case *MethodDecl:
		tpe = "MethodDecl"
	}

	uses := fmt.Sprintf(
		"{dset:%s,iset:%s}",
		decl.(*CommonDecl).uses.decls.String(),
		decl.(*CommonDecl).uses.imports.String(),
	)
	s := fmt.Sprintf(
		`%s{id:"%s",used:%t,uses:%s`,
		tpe, decl.ID(), decl.IsUsed(), uses,
	)

	switch decl := decl.(type) {
	case *TypeDecl:
		var ids []string
		for _, m := range decl.methods.mset {
			ids = append(ids, m.ID())
		}
		s += fmt.Sprintf(
			",methods{%s},methods.keep:%t",
			strings.Join(ids, ","),
			decl.methods.keep,
		)
	case *MethodDecl:
		s += fmt.Sprintf(",embedded:%t", decl.IsEmbedded())
	}

	s += "}"
	return s
}

// DeclSet is a set of Decl
type DeclSet struct {
	dset map[string]Decl
}

// NewDeclSet returns new DeclSet
func NewDeclSet() *DeclSet {
	return &DeclSet{dset: make(map[string]Decl)}
}

// Get gets Decl from set.
func (s *DeclSet) Get(pkg *Package, key ...string) (Decl, bool) {
	d, ok := s.dset[makeID(pkg, key...)]
	return d, ok
}

// GetOrCreate gets Decl from set if exists, otherwise create new one and add to set
// then returns it.
func (s *DeclSet) GetOrCreate(dtype DeclType, pkg *Package, key ...string) Decl {
	if d, ok := s.dset[makeID(pkg, key...)]; ok {
		return d
	}

	d := NewDecl(dtype, pkg, key...)
	s.dset[d.ID()] = d
	return d
}

// Add adds Decl to set.
func (s *DeclSet) Add(d Decl) { s.dset[d.ID()] = d }

// Values creates a slice of values of set and returns it.
func (s *DeclSet) Values() []Decl {
	i := 0
	a := make([]Decl, len(s.dset))
	for _, v := range s.dset {
		v := v
		a[i] = v
		i++
	}
	return a
}

func (s *DeclSet) String() string {
	var v []string
	for _, d := range s.dset {
		v = append(v, fmt.Sprintf(`"%s"`, d.ID()))
	}
	return "DeclSet{" + strings.Join(v, ",") + "}"
}
