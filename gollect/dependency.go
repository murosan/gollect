package gollect

import (
	"fmt"
	"go/ast"
)

type Dependencies map[string]*Dependency

func (deps Dependencies) Get(name string) (*Dependency, bool) {
	d, ok := deps[name]
	return d, ok
}

func (deps Dependencies) Set(name string, d *Dependency) {
	deps[name] = d
}

func (deps Dependencies) GetOrCreate(name string) *Dependency {
	d, ok := deps.Get(name)
	if !ok {
		d = NewDependency(name)
		deps.Set(name, d)
	}
	return d
}

func (deps Dependencies) SetInternal(caller, target string) {
	c, t := deps.GetOrCreate(caller), deps.GetOrCreate(target)
	c.SetInternal(t)
}

func (deps Dependencies) SetExternal(caller *ast.Ident, path, target string) {
	deps.GetOrCreate(caller.Name).SetExternal(path, target)
}

func (deps Dependencies) SetImport(caller *ast.Ident, i *Import) {
	deps.GetOrCreate(caller.Name).SetImport(i)
}

func (deps Dependencies) Use(key string) []ExternalDependencySet {
	if d, ok := deps[key]; !ok {
		panic("no such identity. name = " + key)
	} else {
		return d.Use()
	}
}

func (deps Dependencies) IsUsed(key string) bool {
	v, ok := deps.Get(key)
	return ok && v.IsUsed()
}

func (deps Dependencies) String() (s string) {
	for _, dep := range deps {
		var internal, external, imports string
		for _, d := range dep.internal {
			internal += "\n|     " + d.name
		}
		for ed := range dep.external {
			external += "\n|     " + ed.String()
		}
		for _, i := range dep.imports {
			imports += "\n|     " + i.String()
		}

		s += fmt.Sprintf(`| Ident: %s (used = %v)
|   [internal]%s
|   [external]%s
|   [import]%s
`,
			dep.name,
			dep.IsUsed(),
			internal, external, imports,
		)
	}
	return
}

type Dependency struct {
	name     string
	imports  ImportSet
	internal map[string]*Dependency
	external ExternalDependencySet

	used      bool
	forceUsed bool
}

func NewDependency(name string) *Dependency {
	return &Dependency{
		name:     name,
		imports:  make(ImportSet),
		internal: make(map[string]*Dependency),
		external: make(ExternalDependencySet),
		used:     false,
	}
}

func (d *Dependency) SetInternal(dep *Dependency) {
	if _, ok := d.internal[dep.name]; !ok {
		d.internal[dep.name] = dep
	}
}

func (d *Dependency) SetExternal(path, target string) {
	d.external.Add(path, target)
}

func (d *Dependency) SetImport(i *Import) { d.imports.Add(i) }

func (d *Dependency) String() string { return d.name }

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
	return append(v, d.external)
}

func (d *Dependency) IsUsed() bool { return d.forceUsed || d.used }

type (
	ExternalDependency struct {
		path, name string
	}

	ExternalDependencySet map[ExternalDependency]struct{}
)

func NewExternalDependency(path, name string) ExternalDependency {
	return ExternalDependency{
		path: path,
		name: name,
	}
}

func (ed ExternalDependency) String() string {
	return fmt.Sprintf("%s -> %s", ed.path, ed.name)
}

func (eds ExternalDependencySet) Add(path, name string) {
	eds[NewExternalDependency(path, name)] = struct{}{}
}

func (eds ExternalDependencySet) Get(path, name string) (ExternalDependency, bool) {
	ed := NewExternalDependency(path, name)
	_, ok := eds[ed]
	return ed, ok
}
