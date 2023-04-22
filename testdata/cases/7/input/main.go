package main

import _ "github.com/murosan/gollect/testdata/cases/6/input/pkg"
import "fmt"

type A struct{}
type B struct{}

var _ = A{}

func main() {}

func init() { fmt.Println(&B{}) }
