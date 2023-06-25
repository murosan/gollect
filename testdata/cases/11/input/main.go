package main

import (
	"fmt"
	"os"

	"github.com/murosan/gollect/testdata/cases/11/input/io"
)

func main() {
	s := []int{1, 2, 3}

	a, b := integer(100), integer(200)

	fmt.Println(
		sum1(s),
		sum1[int](s),
		sum2(s),
		sum2[int](s),
		min(a, b),
		min[integer](a, b),
		mapVector(s, func(n int) string { return fmt.Sprint(n) }),
		mapVector[int](s, func(n int) string { return fmt.Sprint(n) }),
		mapVector[int, string](s, func(n int) string { return fmt.Sprint(n) }),
	)

	vec := vector[string]{"abc", "def", "ghi"}
	vec.each(func(s string) {
		fmt.Println(s)
	})
	vec.print()
	(&vec).print()

	w := &io.Writer[int]{Writer: os.Stdout}
	w.PrintSlice([]int{1, 2, 3})
	io.Fprintln(w, "aaa")
}

func sum1[V int | float64](s []V) (n V) {
	for _, v := range s {
		n += v
	}
	return n
}

type number interface{ int | float64 }
type number2 interface{ int | float32 | float64 }
type ordering interface{ ~int | ~float64 }
type integer int

func sum2[V number](s []V) (n V) {
	for _, v := range s {
		n += v
	}
	return n
}

func min[V ordering](a, b V) V {
	if a < b {
		return a
	}
	return b
}

type vector[T any] []T

func (v vector[T]) each(f func(T)) {
	for _, t := range v {
		f(t)
	}
}

func (v *vector[T]) print() {
	for _, t := range *v {
		fmt.Print(t)
	}
}

func mapVector[T, U any](v vector[T], f func(T) U) vector[U] {
	vec := make(vector[U], len(v))
	for i, t := range v {
		vec[i] = f(t)
	}
	return vec
}
