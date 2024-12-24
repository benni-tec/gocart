package utils

import "github.com/google/uuid"

type Cache[T any] struct {
	data map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		data: make(map[string]T),
	}
}

func (c *Cache[T]) NewKey() string {
	id := uuid.New().String()
	if c.IsKey(id) {
		return c.NewKey()
	}

	return id
}

func (c *Cache[T]) IsKey(id string) bool {
	if id == "" {
		return false
	}

	_, ok := c.data[id]
	return ok
}

func (c *Cache[T]) Set(key string, value T) {
	c.data[key] = value
}

func (c *Cache[T]) Get(key string) T {
	return c.data[key]
}
