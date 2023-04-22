package main

import "fmt"

type A struct{}
type B struct{}

var _ = A{}

func main() {}

func init() { fmt.Println(&B{}) }
