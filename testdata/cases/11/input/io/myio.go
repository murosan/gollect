package io

import (
	"fmt"
	"io"
)

type Printable interface {
	~int | ~string
}

type Writer[T Printable] struct {
	io.Writer
}

func (w *Writer[T]) Println(v T)      { fmt.Println(v) }
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
