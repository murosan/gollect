package main

import "fmt"

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

func main() {
	fmt.Println(two, numB, numD, d, g, l, m)
}
