package main

import    "fmt"

const ConstA  = 100
const constA  = 200

const (
	constB int =           300
		constC  = 400
	_ = 900
	constD = 1000
)

 var VarA = 500
var     varA = 600
	var (
	varB = 700
	varC  int = 800
)

type A struct{}

var (
	one,
	_,
	three A
)

func Nums() []int {
	return []int{10, 20, 30}
}

func Unused() []int {
	return []int{10, 20, 30}
}

func main() {
	fmt.Println(
		ConstA,
		constB,
		constD,
		varA,
		varB,
		varC,
		Nums(),
		one,
		)
}
