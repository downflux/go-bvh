package cache

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/vector"
)

type C struct {
	epsilon float64
	k       vector.D

	data  []*N
	freed []ID
}

type O struct {
	Epsilon float64
	K       vector.D
}

func New(o O) *C {
	return &C{
		epsilon: o.Epsilon,
		k:       o.K,

		data:  make([]*N, 0, 128),
		freed: make([]ID, 0, 128),
	}
}

func (c *C) K() vector.D {
	return c.k
}

func (c *C) Insert(p, l, r ID) ID {
	var x ID
	// Reuse a node if available -- this avoids additional allocs.
	if len(c.freed) > 0 {
		x, c.freed = c.freed[len(c.freed)-1], c.freed[:len(c.freed)-1]
	} else {
		c.data = append(c.data, nil)
		x = ID(len(c.data) - 1)
	}
	c.data[x] = c.data[x].allocateOrLoad(c, x, p, l, r)
	return x
}

// Get returns a node data struct.
func (c *C) Get(x ID) (*N, bool) {
	if !x.IsValid() || int(x) >= len(c.data) {
		return nil, false
	}

	n := c.data[x]
	return n, n.IsAllocated()
}

func (c *C) GetOrDie(x ID) *N {
	n, ok := c.Get(x)
	if !ok {
		panic(fmt.Sprintf("cannot find node %v", x))
	}
	return n
}

// Delete returns a given node to the available pool.
func (c *C) Delete(x ID) bool {
	if !x.IsValid() || int(x) >= len(c.data) {
		return false
	}
	n := c.data[x]
	if !n.IsAllocated() {
		return false
	}

	n.free()
	c.freed = append(c.freed, x)

	return true
}

func (c *C) DeleteOrDie(x ID) {
	if ok := c.Delete(x); !ok {
		panic(fmt.Sprintf("cannot find node %v", x))
	}
}
