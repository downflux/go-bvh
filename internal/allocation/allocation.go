package allocation

import (
	"fmt"
	"math/rand"

	"github.com/downflux/go-bvh/internal/allocation/id"
)

type C[T comparable] map[id.ID]T

func New[T comparable]() *C[T] {
	c := C[T](map[id.ID]T{})
	return &c
}

func (c C[T]) Allocate() id.ID {
	var i id.ID
	found := true
	for ; found; i = id.ID(rand.Int()) {
		_, found = c[i]
	}

	var blank T
	c[i] = blank

	return i
}

func (c C[T]) Insert(i id.ID, n T) error {
	m, ok := c[i]

	// id.ID must be allocated first.
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

func (c C[T]) Remove(i id.ID) error {
	n := c[i]
	var blank T
	if n == blank {
		return fmt.Errorf("cannot remove non-existent node: %v", i)
	}

	delete(c, i)
	return nil
}
