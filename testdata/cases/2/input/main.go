package main

import    "fmt"

const ConstA  = 100
const constA  = 200

const (
	constB int =           300
		constC  = 400
)

 var VarA = 500
var     varA = 600
	var (
	varB = 700
	varC  int = 800
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
		varA,
		varB,
		varC,
		Nums(),
		)
}
