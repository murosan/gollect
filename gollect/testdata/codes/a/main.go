package main

import (
	"fmt"

	"github.com/murosan/gollect/gollect/testdata/codes/pkg1"
)

// comment inside main should be left

var num = 1000

func main() {
	(&pkg1.TypeA{}).Do3()

	// this comment should be left
	fmt.Println(pkg1.NumA)
	fmt.Println(pkg1.NumC)
	pkg1.PrintMax(pkg1.NumA, num)
}
