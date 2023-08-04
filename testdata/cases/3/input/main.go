package main

import (
	"sort"
)

func main() {
	newT()

	s := &S{}
	sort.Sort(s)
	println(&U{})
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

type S struct {
	sort.Interface
	data []int
}

func (s *S) Len() int           { return len(s.data) }
func (s *S) Less(i, j int) bool { return s.data[i] < s.data[j] }
func (s *S) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s *S) A()                 {}
func (s *S) Unused()            {}

type I interface {
	A()
}

type U struct {
	I
	*S
}
