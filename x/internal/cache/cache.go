package cache

import (
	"fmt"
)

type C struct {
	data  []*N
	freed map[ID]bool
}

func New() *C {
	return &C{
		data:  make([]*N, 0, 128),
		freed: map[ID]bool{},
	}
}

func (c *C) Insert(p, l, r ID) ID {
	var x ID
	// Reuse a node if available -- this avoids additional allocs.
	if len(c.freed) > 0 {
		for i := range c.freed {
			delete(c.freed, i)
			x = ID(i)
			break
		}
	} else {
		c.data = append(c.data, nil)
		x = ID(len(c.data) - 1)
	}
	c.data[x] = c.data[x].allocateOrLoad(x, p, l, r)
	return x
}

// Get returns a node data struct.
func (c *C) Get(x ID) (*N, bool) {
	if !x.IsValid() || int(x) >= len(c.data) || c.freed[x] {
		return nil, false
	}

	return c.data[x], true
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
	if !x.IsValid() || int(x) >= len(c.data) || c.freed[x] {
		return false
	}

	c.data[x].free()
	c.freed[x] = true
	return true
}

func (c *C) DeleteOrDie(x ID) {
	if ok := c.Delete(x); !ok {
		panic(fmt.Sprintf("cannot find node %v", x))
	}
}
