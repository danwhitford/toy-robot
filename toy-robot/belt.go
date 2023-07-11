package toyrobot

import "fmt"

type Belt[T any] struct {
	input []T
	ptr   int
	size  int
}

func NewBelt[T any](input []T) *Belt[T] {
	return &Belt[T]{
		input: input,
		ptr:   0,
		size:  len(input),
	}
}

func (b *Belt[T]) GetNext() (T, error) {
	var t T
	if b.ptr >= b.size {
		return t, fmt.Errorf("out of bounds")
	}
	token := b.input[b.ptr]
	b.ptr++
	return token, nil
}

func (b *Belt[T]) Peek() (T, error) {
	var t T
	if b.ptr >= b.size {
		return t, fmt.Errorf("out of bounds")
	}
	return b.input[b.ptr], nil
}

func (b *Belt[T]) HasNext() bool {
	return b.ptr < b.size
}
