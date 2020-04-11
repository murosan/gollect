package main

func main() {
	newT()
}

type reader interface {
	Read(p []byte) (n int, err error)
}

type T interface {
	Do(r reader)
}

func newT() T {
	return &t{}
}

// gollect: keep methods
type t struct{}

func (*t) Do(r reader) {}

func (*t) Do2() {}
