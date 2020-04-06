package main

import (
	"github.com/murosan/gollect/testdata/cases/6/input/pkg"
)

func main() {
	a := pkg.NewA(&pkg.B{})
	a.Do()
}
