// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"reflect"
	"testing"
)

func TestDependencies_Get(t *testing.T) {
	deps := make(Dependencies)
	name := "name"

	// should be empty
	if d, ok := deps.Get(name); ok || d != nil {
		t.Errorf("expected empty, but got something")
	}

	// also test set method
	deps.Set(NewDependency(name))
	if d, ok := deps.Get(name); !ok || d == nil {
		t.Errorf("expected non empty, but got nothing")
	}
}

func TestDependencies_GetOrCreate(t *testing.T) {
	deps := make(Dependencies)
	name := "name"

	// should be empty
	if d, ok := deps.Get(name); ok || d != nil {
		t.Errorf("expected empty, but got something")
	}

	if d := deps.GetOrCreate(name); d == nil || d.name != name {
		t.Errorf("expected non empty, but got nothing")
	}

	// ensure created dependency is set
	if d, ok := deps.Get(name); !ok || d == nil {
		t.Errorf("expected empty, but got something")
	}
}

func TestDependencies_SetInternal(t *testing.T) {
	deps := make(Dependencies)
	caller, target := "caller", "target"

	deps.SetInternal(caller, target)

	clr, ok1 := deps.Get(caller)
	tgt, ok2 := deps.Get(target)

	if !ok1 || !ok2 {
		t.Errorf("could not find object. GetOrCreate might not be called correctly")
	}

	tgt2, ok3 := clr.internal[target]
	if !ok3 || tgt != tgt2 {
		t.Errorf("caller is failing to set dependency. want:%v, actual:%v", tgt, tgt2)
	}
}

func TestDependencies_SetExternal(t *testing.T) {
	deps := make(Dependencies)

	caller := "name"
	path := "github.com/murosan/abc"
	target := "target"

	deps.SetExternal(caller, path, target)

	d, ok := deps.Get(caller)
	if !ok || d == nil {
		t.Errorf("could not find object. GetOrCreate might not be called correctly")
	}

	if _, ok := d.external.Get(path, target); !ok {
		t.Errorf("external dependency is not set correctly")
	}
}

func TestDependencies_SetImport(t *testing.T) {
	deps := make(Dependencies)

	caller := "callerName"

	alias, name, path := "alias", "name", "path"
	i := NewImport(alias, name, path)

	deps.SetImport(caller, i)

	if d, ok := deps.Get(caller); !ok || d == nil {
		t.Errorf("GetOrCreate might not be called correctly")
	} else if im, ok := d.imports.Get(i.AliasOrName()); !ok || im != i {
		t.Errorf("Import is not set correctly. found:%t, isSame:%t", ok, im == i)
	}
}

func TestDependencies_Use(t *testing.T) {
	deps := make(Dependencies)
	key := "key"
	d := NewDependency(key)

	shouldNotPanic(
		t,
		func() { deps.Use(key) },
		"expected no panic, but recovered",
	)

	deps.Set(d)

	shouldNotPanic(
		t,
		func() { deps.Use(key) },
		"expected no panic, but recovered",
	)
}

func TestDependencies_IsUsed(t *testing.T) {
	deps := make(Dependencies)
	key := "key"
	d := NewDependency(key)

	if deps.IsUsed(key) {
		t.Errorf("should be false when unknown key was passed")
	}

	deps.Set(d)

	if deps.IsUsed(key) {
		t.Errorf("expected is not be used but used")
	}

	deps.Use(key)

	if !deps.IsUsed(key) {
		t.Errorf("expected is used but not used")
	}
}

func TestDependency_SetInternal(t *testing.T) {
	dep := NewDependency("name")
	tgt := NewDependency("target")

	if _, ok := dep.internal[tgt.name]; ok {
		t.Errorf("should be empty")
	}

	dep.SetInternal(tgt)
	if tgt2, ok := dep.internal[tgt.name]; !ok || tgt != tgt2 {
		t.Errorf("target was not set correctly")
	}
}

func TestDependency_SetExternal(t *testing.T) {
	dep := NewDependency("name")
	path, name := "github.com/murosan/abc", "def"

	if _, ok := dep.external.Get(path, name); ok {
		t.Errorf("should be empty")
	}

	dep.SetExternal(path, name)
	if ed, ok := dep.external.Get(path, name); !ok || ed.name != name || ed.path != path {
		t.Errorf(
			"want:(%t, %s, %s), actual: (%t, %s, %s)",
			true, name, path, ok, ed.name, ed.path,
		)
	}
}

func TestDependency_SetImport(t *testing.T) {
	dep := NewDependency("name")
	im := NewImport("alias", "name", "path")

	if _, ok := dep.imports.Get(im.AliasOrName()); ok {
		t.Errorf("should be empty")
	}

	dep.SetImport(im)
	if im2, ok := dep.imports.Get(im.AliasOrName()); !ok || im != im2 {
		t.Errorf("want: (%t, %v), actual: (%t, %v)", true, im, ok, im2)
	}
}

func TestDependency_Use(t *testing.T) {
	dep := NewDependency("name")
	dep.used = true

	if res := dep.Use(); len(res) != 0 {
		t.Errorf("want empty, actual: %v", res)
	}
}

func TestDependency_Use2(t *testing.T) {
	dep1 := NewDependency("dep1")
	dep2 := NewDependency("dep2")
	dep3 := NewDependency("dep3")
	dep4 := NewDependency("dep4")

	path1 := "github.com/murosan/abc"

	dep1.SetInternal(dep2)
	dep1.SetExternal(path1, dep3.name)
	dep1.SetExternal(path1, dep4.name)
	dep2.SetExternal(path1, dep4.name)

	im1 := NewImport("alias", "abc", path1)
	im2 := NewImport("", "fmt", "fmt")

	dep1.SetImport(im1)
	dep2.SetImport(im2)

	if dep1.IsUsed() || dep2.IsUsed() || dep3.IsUsed() || dep4.IsUsed() {
		t.Errorf("failing")
	}
	if im1.IsUsed() || im2.IsUsed() {
		t.Errorf("failing")
	}

	actual := dep1.Use()

	if !dep1.IsUsed() || !dep2.IsUsed() || dep3.IsUsed() || dep4.IsUsed() {
		t.Errorf("failing")
	}
	if !im1.IsUsed() || !im2.IsUsed() {
		t.Errorf("failing")
	}

	expected := []ExternalDependencySet{
		{
			NewExternalDependency(path1, dep3.name): struct{}{},
			NewExternalDependency(path1, dep4.name): struct{}{},
		},
		{
			NewExternalDependency(path1, dep4.name): struct{}{},
		},
	}

	if reflect.DeepEqual(actual, expected) {
		t.Errorf("want: %v, actual: %v", expected, actual)
	}
}

func TestDependency_IsUsed(t *testing.T) {
	dep := NewDependency("name")
	if dep.IsUsed() {
		t.Errorf("fail")
	}

	dep.used = true
	if !dep.IsUsed() {
		t.Errorf("fail")
	}
}

func TestExternalDependencySet(t *testing.T) {
	path1, name1 := "github.com/murosan/abc", "name1"
	path2, name2 := "github.com/murosan/def", "name2"
	ed1 := NewExternalDependency(path1, name1)
	ed2 := NewExternalDependency(path2, name2)

	set := make(ExternalDependencySet)

	set.Add(path1, name1)
	set.Add(path1, name1)
	set.Add(path2, name2)

	if len(set) != 2 {
		t.Errorf("want: 2, actual: %d", len(set))
	}

	if v, ok := set.Get(path1, name1); !ok || v != ed1 {
		t.Errorf("want: (%t, %v), actual: (%t, %v)", true, ed1, ok, v)
	}

	if v, ok := set.Get(path2, name2); !ok || v != ed2 {
		t.Errorf("want: (%t, %v), actual: (%t, %v)", true, ed1, ok, v)
	}

	if _, ok := set.Get("a", name1); ok {
		t.Errorf("want: %t, actual: %t", false, ok)
	}
}

func TestUseAll(t *testing.T) {
	packages := make(Packages)
	is := NewImportSet()

	p1 := NewPackage("main", is)
	p2 := NewPackage("pkg2", is)
	p3 := NewPackage("pkg3", is)
	packages.Set(p1.path, p1)
	packages.Set(p2.path, p2)
	packages.Set(p3.path, p3)

	main := NewDependency("main")
	varA := NewDependency("varA")
	varB := NewDependency("varB")
	varC := NewDependency("varC")
	numA := NewDependency("NumA")
	typeA := NewDependency("TypeA")

	p1.Dependencies().Set(main)
	p1.Dependencies().Set(varA)
	p1.Dependencies().Set(varB)
	p2.Dependencies().Set(varC)
	p2.Dependencies().Set(numA)
	p3.Dependencies().Set(typeA)

	im1 := NewImport("alias1", "name1", "path1")
	im2 := NewImport("alias2", "name2", "path2")
	im3 := NewImport("alias3", "name3", "path3")

	p1.Dependencies().SetInternal(main.name, varA.name)
	p1.Dependencies().SetExternal(main.name, p2.path, numA.name)
	p1.Dependencies().SetExternal(main.name, p3.path, typeA.name)
	p2.Dependencies().SetExternal(numA.name, p3.path, typeA.name)

	p1.Dependencies().SetImport(main.name, im1)
	p2.Dependencies().SetImport(numA.name, im1)
	p2.Dependencies().SetImport(numA.name, im2)
	p1.Dependencies().SetImport(varB.name, im3) // set to unused

	next := []ExternalDependencySet{{}}
	next[0].Add(main.name, main.name)
	UseAll(packages, next)

	for i, d := range []*Dependency{main, varA, numA, typeA} {
		if !d.IsUsed() {
			t.Errorf("%s is unused. index=%d", d.name, i)
		}
	}

	for i, d := range []*Dependency{varB, varC} {
		if d.IsUsed() {
			t.Errorf("%s is used. index=%d", d.name, i)
		}
	}

	for n, i := range []*Import{im1, im2} {
		if !i.IsUsed() {
			t.Errorf("%v is unused. index=%d", i, n)
		}
	}

	for n, i := range []*Import{im3} {
		if i.IsUsed() {
			t.Errorf("%v is unused. index=%d", i, n)
		}
	}
}

func TestUseAll2(t *testing.T) {
	packages := make(Packages)
	next := []ExternalDependencySet{{}}
	next[0].Add("main", "main")

	shouldPanic(
		t,
		func() {
			UseAll(packages, next)
		},
		"should panic if next contains unknown identity",
	)
}
