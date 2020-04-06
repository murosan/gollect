package pkg

type (
	// A is a type. this comment will be left.
	A struct {
		b *B
	}

	// B decl.
	// This line will also be left.
	B struct{ n int }

	// This comment will be removed.
)

func NewA(b *B) *A {
	return &A{b: b}
}

func (a *A) Do() *B     { return a.getb() }
func (a *A) Unused() *B { return a.getb() }

func (a *A) getb() *B { return a.b }
