package main

func main() {
	t{}.Do()

	a := &A{b: &B{}}
	a.m()
}

type t struct{}

func (t) Do() {}

type (
	A struct {
		b *B
	}
	B struct{}
)

func (a *A) m() *B {
	return a.b
}
