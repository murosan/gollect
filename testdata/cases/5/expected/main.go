package main

type A struct{}

func (A) m1()  {}
func (A) m2()  {}
func (A) m3()  {}
func (A) m4()  {}
func (*A) m5() {}
func (A) m6()  {}
func (A) m7()  {}
func (A) m8()  {}
func (A) m9()  {}
func (A) m10() {}
func (A) m11() {}

var a A
var s = []A{{}, {}, {}}

func main() {
	a.m1()
	A.m2(a)
	s[2].m3()
	A{}.m4()
	(&A{}).m5()
	(A{}).m6()
	(A{}).m7()
	func() A { return A{} }().m8()
	(func() A { return A{} }()).m9()
	(func() A { return A{} })().m10()
	func() (a A) { return }().m11()
}
