package stack

import "fmt"

type RobotStack[T any] []T

func (s *RobotStack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *RobotStack[T]) Pop() (T, error) {
	if len(*s) == 0 {
		var t T
		return t, fmt.Errorf("stack is empty")
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, nil
}
