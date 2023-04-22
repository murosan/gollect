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

	numD = 200
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

const (
	_ A = iota
	o
	p
)

var v2, v3 int

func main() {
	_ = two
	_ = numB
	_ = numD
	_ = d
	_ = g
	_ = l
	_ = m
	_ = v2
	_ = v3
}
