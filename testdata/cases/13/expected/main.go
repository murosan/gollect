package main

type Vector[T any] = []T
type StringMap[V any] = map[string]V

func main() {
	var v Vector[int]
	var m StringMap[string]
	_ = v
	_ = m
}
