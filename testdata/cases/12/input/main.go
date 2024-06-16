package main

func main() {
	(&S[int, string]{}).f1()
}

type S[T, U any] struct {
	t T
	u U
}

func (s S[T, U]) f1() {}

func (s *S[T, U]) f2() {}
