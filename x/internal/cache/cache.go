package cache

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

type N interface {
	node.N

	Allocate(parent cid.ID, left cid.ID, right cid.ID)
	Free()
	IsAllocated() bool
}

type C struct {
	k        vector.D
	leafSize int

	data  []N
	freed []cid.ID
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
		panic(fmt.Sprintf("invalid impl leaf size %v", o.LeafSize))
	}

	return &C{
		k:        o.K,
		leafSize: o.LeafSize,

		data:  make([]N, 0, 128),
		freed: make([]cid.ID, 0, 128),
	}
}

func (c *C) K() vector.D   { return c.k }
func (c *C) LeafSize() int { return c.leafSize }

// IsAllocated checks if the given impl is tracked by the cache. This function
// returns false if the impl is in the freed pool.
func (c *C) IsAllocated(x cid.ID) bool {
	_, ok := c.Get(x)
	return ok
}

// Get returns a impl data struct.
func (c *C) Get(x cid.ID) (node.N, bool) {
	if !x.IsValid() || int(x) >= len(c.data) {
		return nil, false
	}

	n := c.data[x]
	if !n.IsAllocated() {
		return nil, false
	}
	return n, true
}

func (c *C) GetOrDie(x cid.ID) node.N {
	n, ok := c.Get(x)
	if !ok {
		panic(fmt.Sprintf("cannot find impl %v", x))
	}
	return n
}

func (c *C) Insert(p, l, r cid.ID, validate bool) cid.ID {
	if validate {
		if p.IsValid() && !c.IsAllocated(p) {
			panic(fmt.Sprintf("cannot set new impl with invalid parent %v", p))
		}
		if l.IsValid() && !c.IsAllocated(l) {
			panic(fmt.Sprintf("cannot set new impl with invalid left child %v", l))
		}
		if r.IsValid() && !c.IsAllocated(r) {
			panic(fmt.Sprintf("cannot set new impl with invalid right child %v", r))
		}
	}

	var x cid.ID
	// Reuse a impl if available -- this avoids additional allocs.
	if len(c.freed) > 0 {
		x, c.freed = c.freed[len(c.freed)-1], c.freed[:len(c.freed)-1]
	} else {
		x = cid.ID(len(c.data))
		c.data = append(c.data, impl.New(c, x))
	}
	c.data[x].Allocate(p, l, r)
	return x
}

// Delete returns a given impl to the available pool.
func (c *C) Delete(x cid.ID) bool {
	if !x.IsValid() || int(x) >= len(c.data) {
		return false
	}
	n := c.data[x]
	if !n.IsAllocated() {
		return false
	}

	n.Free()
	c.freed = append(c.freed, x)

	return true
}

func (c *C) DeleteOrDie(x cid.ID) {
	if ok := c.Delete(x); !ok {
		panic(fmt.Sprintf("cannot find impl %v", x))
	}
}
