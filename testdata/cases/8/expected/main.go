package main

import "io"

type A struct{ B }

type B struct{ C }

func (B) b() {}

type C struct {
	io.Writer
	// unused field, but this should be left.
	unused func()
}

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
