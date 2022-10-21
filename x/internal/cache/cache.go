package cache

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/vector"
)

type C struct {
	k        vector.D
	leafSize int

	data  []*N
	freed []ID
}

type O struct {
	K        vector.D
	LeafSize int
}

func New(o O) *C {
	if o.K <= 0 {
		panic(fmt.Sprintf("invalid AABB dimension %v", o.K))
	}
	if o.LeafSize <= 0 {
		panic(fmt.Sprintf("invalid node leaf size %v", o.LeafSize))
	}

	return &C{
		k:        o.K,
		leafSize: o.LeafSize,

		data:  make([]*N, 0, 128),
		freed: make([]ID, 0, 128),
	}
}

func (c *C) K() vector.D   { return c.k }
func (c *C) LeafSize() int { return c.leafSize }

func (c *C) IsAllocated(x ID) bool {
	_, ok := c.Get(x)
	return ok
}

func (c *C) Insert(p, l, r ID, validate bool) ID {
	if validate {
		if p.IsValid() && !c.IsAllocated(p) {
			panic(fmt.Sprintf("cannot set new node with invalid parent %v", p))
		}
		if l.IsValid() && !c.IsAllocated(l) {
			panic(fmt.Sprintf("cannot set new node with invalid left child %v", l))
		}
		if r.IsValid() && !c.IsAllocated(r) {
			panic(fmt.Sprintf("cannot set new node with invalid right child %v", r))
		}
	}

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
