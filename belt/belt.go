package belt

import "fmt"

type Belt[T any] struct {
	input []T
	Ptr   int
	size  int
}

func NewBelt[T any](input []T) *Belt[T] {
	return &Belt[T]{
		input: input,
		Ptr:   0,
		size:  len(input),
	}
}

func (b *Belt[T]) GetNext() (T, error) {
	var t T
	if b.Ptr >= b.size {
		return t, fmt.Errorf("out of bounds")
	}
	token := b.input[b.Ptr]
	b.Ptr++
	return token, nil
}

func (b *Belt[T]) Peek() (T, error) {
	var t T
	if b.Ptr >= b.size {
		return t, fmt.Errorf("out of bounds")
	}
	return b.input[b.Ptr], nil
}

func (b *Belt[T]) HasNext() bool {
	return b.Ptr < b.size
}
