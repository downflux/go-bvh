package allocation

import (
	"fmt"
	"math/rand"
)

type ID int

type C[T comparable] map[ID]T

func New[T comparable](data []T) C[T] {
	c := C[T](map[ID]T{})
	for _, p := range data {
		c.Insert(c.Allocate(), p)
	}
	return c
}

func (c C[T]) Allocate() ID {
	var i ID
	found := true
	for ; found; i = ID(rand.Int()) {
		_, found = c[i]
	}

	var blank T
	c[i] = blank

	return i
}

func (c C[T]) Insert(i ID, n T) error {
	m, ok := c[i]

	// ID must be allocated first.
	if !ok {
		return fmt.Errorf("inserting an unallocated node %v", i)
	}

	var blank T
	if m != blank {
		return fmt.Errorf("duplicate node found with same index %v", i)
	}

	c[i] = n

	return nil
}
