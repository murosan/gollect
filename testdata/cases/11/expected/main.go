package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"os"
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
		max(a, b),
		max[integer](a, b),
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

	w := &Writer[int]{Writer: os.Stdout}
	w.PrintSlice([]int{1, 2, 3})
	Fprintln(w, "aaa")
}

func sum1[V int | float64](s []V) (n V) {
	for _, v := range s {
		n += v
	}
	return n
}

type number interface{ int | float64 }

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

func max[V constraints.Ordered](a, b V) V {
	if a > b {
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

type Printable interface {
	~int | ~string
}

type Writer[T Printable] struct {
	io.Writer
}

func (w *Writer[T]) PrintSlice(s []T) { PrintSlice(s) }

func PrintSlice[T any](s []T) {
	if len(s) != 0 {
		fmt.Print(s[0])
		for _, v := range s[1:] {
			fmt.Print(" ", v)
		}
	}
	fmt.Println()
}

func Fprintln[T any](w io.Writer, v T) { fmt.Fprintln(w, v) }
