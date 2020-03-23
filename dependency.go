// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
)

// Dependencies is a map of Dependency.
// The key is ident name of the Dependency.
type Dependencies map[string]*Dependency

// Get returns Dependency.
func (deps Dependencies) Get(name string) (*Dependency, bool) {
	d, ok := deps[name]
	return d, ok
}

// Set sets Dependency.
func (deps Dependencies) Set(d *Dependency) { deps[d.name] = d }

// GetOrCreate returns Dependency.
// If no Dependency found, creates new one.
func (deps Dependencies) GetOrCreate(name string) *Dependency {
	d, ok := deps.Get(name)
	if !ok {
		d = NewDependency(name)
		deps.Set(d)
	}
	return d
}

// SetInternal lets caller set internal dependency.
func (deps Dependencies) SetInternal(caller, target string) {
	c, t := deps.GetOrCreate(caller), deps.GetOrCreate(target)
	c.SetInternal(t)
}

// SetExternal lets caller set external dependency.
func (deps Dependencies) SetExternal(caller, path, target string) {
	deps.GetOrCreate(caller).SetExternal(path, target)
}

// SetImport lets caller set import.
func (deps Dependencies) SetImport(caller string, i *Import) {
	deps.GetOrCreate(caller).SetImport(i)
}

// SetMethod lets caller set method.
func (deps Dependencies) SetMethod(caller string, target string) {
	c, t := deps.GetOrCreate(caller), deps.GetOrCreate(target)
	c.SetMethod(t)
}

// TurnOnKeepMethodOption lets caller set keepMethod option to true
func (deps Dependencies) TurnOnKeepMethodOption(caller string) {
	deps.GetOrCreate(caller).TurnOnKeepMethodOption()
}

// Use set dependency's use state to true.
// Panics if deps has no key. Therefore, make sure
// to set all dependencies before.
func (deps Dependencies) Use(key string) []ExternalDependencySet {
	if d, ok := deps[key]; !ok {
		panic("no such identity. name = " + key)
	} else {
		return d.Use()
	}
}

// IsUsed gets Dependency from map and returns whether it is used or not.
// If deps has no key, returns false.
func (deps Dependencies) IsUsed(key string) bool {
	v, ok := deps.Get(key)
	return ok && v.IsUsed()
}

func (deps Dependencies) String() (s string) {
	for _, dep := range deps {
		var internal, external, imports, methods string
		for _, d := range dep.internal {
			internal += "\n|     " + d.name
		}
		for ed := range dep.external {
			external += "\n|     " + ed.String()
		}
		for _, i := range dep.imports {
			imports += "\n|     " + i.String()
		}
		for _, m := range dep.methods {
			methods += "\n|     " + m.String()
		}

		s += fmt.Sprintf(`| Ident: %s (used = %v)
|   [internal]%s
|   [external]%s
|   [import]%s
|   [methods]%s
`,
			dep.name,
			dep.IsUsed(),
			internal, external, imports, methods,
		)
	}
	return
}

// Dependency represents what the identifier is depending on.
type Dependency struct {
	name     string
	imports  ImportSet
	internal map[string]*Dependency // Dependencies inside same package
	external ExternalDependencySet  // Dependencies of external package

	// TODO: refactor
	methods map[string]*Dependency

	used bool // lazily changed to true by Use method

	// an option for type declaration. if true, all methods will be kept,
	// otherwise only used methods directly from code will be kept.
	keepMethods bool
}

// NewDependency returns new Dependency.
func NewDependency(name string) *Dependency {
	return &Dependency{
		name:        name,
		imports:     make(ImportSet),
		internal:    make(map[string]*Dependency),
		external:    make(ExternalDependencySet),
		methods:     make(map[string]*Dependency),
		used:        false,
		keepMethods: false,
	}
}

// SetInternal sets Dependency in same package.
func (d *Dependency) SetInternal(dep *Dependency) {
	if _, ok := d.internal[dep.name]; !ok {
		d.internal[dep.name] = dep
	}
}

// SetExternal sets Dependency of external package.
func (d *Dependency) SetExternal(path, target string) {
	d.external.Add(path, target)
}

// SetImport sets import dependency.
func (d *Dependency) SetImport(i *Import) { d.imports.Add(i) }

// SetMethod sets the Dependency to methods map.
func (d *Dependency) SetMethod(dep *Dependency) {
	if _, ok := d.methods[dep.name]; !ok {
		d.methods[dep.name] = dep
	}
}

// TurnOnKeepMethodOption sets keepMethods option to true
func (d *Dependency) TurnOnKeepMethodOption() { d.keepMethods = true }

func (d *Dependency) String() string { return d.name }

// Use changes it's use state to true and
// calls Use method of internal Dependencies.
// This returns all external Dependencies they are depending on.
func (d *Dependency) Use() (v []ExternalDependencySet) {
	if d.used {
		return
	}
	d.used = true

	for _, i := range d.imports {
		i.Use()
	}
	for _, dep := range d.internal {
		v = append(v, dep.Use()...)
	}
	if d.keepMethods {
		for _, method := range d.methods {
			v = append(v, method.Use()...)
		}
	}
	return append(v, d.external)
}

// IsUsed returns used state.
func (d *Dependency) IsUsed() bool { return d.used }

type (
	// ExternalDependency represents external package's dependency information.
	ExternalDependency struct {
		path, name string
	}

	// ExternalDependencySet is a set of ExternalDependency.
	ExternalDependencySet map[ExternalDependency]struct{}
)

// NewExternalDependency returns new ExternalDependency.
func NewExternalDependency(path, name string) ExternalDependency {
	return ExternalDependency{
		path: path,
		name: name,
	}
}

func (ed ExternalDependency) String() string {
	return fmt.Sprintf("%s -> %s", ed.path, ed.name)
}

// Add adds new ExternalDependency to set.
func (eds ExternalDependencySet) Add(path, name string) {
	eds[NewExternalDependency(path, name)] = struct{}{}
}

// Get gets ExternalDependency from set.
func (eds ExternalDependencySet) Get(path, name string) (ExternalDependency, bool) {
	ed := NewExternalDependency(path, name)
	_, ok := eds[ed]
	return ed, ok
}

// UseAll calls Use to each dependency across all package.
func UseAll(packages Packages, next []ExternalDependencySet) {
	for ; len(next) > 0; next = next[1:] {
		for ed := range next[0] {
			if pkg, ok := packages.Get(ed.path); ok {
				deps := pkg.Dependencies()
				next = append(next, deps.Use(ed.name)...)
			} else {
				panic("unknown package. path=" + ed.path)
			}
		}
	}
}
