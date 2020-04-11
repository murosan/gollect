package main

type A int

const (
	zero = iota
	one
	two
)

const (
	numA = 100
	numB
	numC
	numD = 200
	numE
)

const (
	a = 10
	b = iota
	c
	d = "d"
	e
)

const (
	f, g = iota, iota
	h, i
)

const (
	_ = iota
	j
	_
	l
)

const (
	_ A = iota
	m
	n
)

func main() {
	_ = two
	_ = numB
	_ = numD
	_ = d
	_ = g
	_ = l
	_ = m
}
