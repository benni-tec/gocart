package utils

import "sync"

type Set[T comparable] struct {
	sync.Mutex
	data map[T]bool
}

func NewSet[T comparable](values ...T) *Set[T] {
	set := &Set[T]{
		data: make(map[T]bool),
	}

	for _, v := range values {
		set.Add(v)
	}

	return set
}

func (set *Set[T]) Has(value T) bool {
	_, ok := set.data[value]
	return ok
}

func (set *Set[T]) Add(value T) {
	set.data[value] = true
}

func (set *Set[T]) Remove(value T) {
	delete(set.data, value)
}
