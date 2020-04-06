package main

import "fmt"

const ConstA = 100

const constB int = 300

var varA = 600
var (
	varB     = 700
	varC int = 800
)

func Nums() []int {
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
