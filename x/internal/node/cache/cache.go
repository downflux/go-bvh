package cache

import (
	"fmt"
)

type ID int

type C[T any] struct {
	cache []T
	freed map[ID]bool
}

func New[T any]() *C[T] {
	return &C[T]{
		cache: make([]T, 0, 128),
		freed: map[ID]bool{},
	}
}

func (c *C[T]) Get(x ID) T {
	if int(x) >= len(c.cache) || c.freed[x] {
		panic(fmt.Sprintf("invalid cache ID %v", x))
	}

	return c.cache[x]
}

func (c *C[T]) Insert(t T) ID {
	if len(c.freed) > 0 {
		for x := range c.freed {
			delete(c.freed, x)
			c.cache[x] = t
			return x
		}
	}

	c.cache = append(c.cache, t)
	return ID(len(c.cache) - 1)
}

func (c *C[T]) Remove(x ID) {
	if int(x) >= len(c.cache) || c.freed[x] {
		panic(fmt.Sprintf("invalid cache ID %v", x))
	}

	var blank T
	c.cache[x] = blank
	c.freed[x] = true
}
