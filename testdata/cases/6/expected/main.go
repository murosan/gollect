package main

func main() {
	a := NewA(&B{})
	a.Do()
}

type (
	// A is a type. this comment will be left.
	A struct {
		b *B
	}

	// B decl.
	// This line will also be left.
	B struct{ n int }
)

func NewA(b *B) *A {
	return &A{b: b}
}

func (a *A) Do() *B { return a.getb() }

func (a *A) getb() *B { return a.b }
