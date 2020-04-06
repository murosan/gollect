package main

import (
	"io"
)

func main() {
	newT()
}

type T interface {
	Do(r io.Reader)
}

func newT() T {
	return &t{}
}

// gollect: keep methods
type t struct{}

func (*t) Do(r io.Reader) {}

func (*t) Do2() {}
