package main

import "io"

type A struct{ B }

func (A) a() {}

type B struct{ C }

func (B) b() {}

type C struct {
	io.Writer
	// unused field, but this should be left.
	unused func()
}

func (C) c() {}

type T struct{ I }
type I interface {
	a()
	b()
}

func main() {
	var a A
	a.b()
	a.Write([]byte(""))

	var t T
	t.a() // this causes panic.
}
