package node

import (
	"github.com/downflux/go-bvh/x/internal/cache"
)

type N struct {
	cache  *cache.C
	id     cache.ID
	parent cache.ID
	branch B
}

func New(c *cache.C, x cache.ID) *N {
	n := &N{
		cache:  c,
		id:     x,
		parent: cache.IDInvalid,
		branch: BInvalid,
	}
	if p, ok := c.Get(
		c.GetOrDie(x).Parent()); ok {
		n.parent = p.ID()
		if x == p.Left() {
			n.branch = BLeft
		} else {
			n.branch = BRight
		}
	}
	return n
}

func (n *N) Branch() B { return n.branch }

func (n *N) IsRoot() bool { return n.parent.IsValid() }
func (n *N) IsLeaf() bool { return n.Left() == nil }

func (n *N) Parent() *N {
	x := n.cache.GetOrDie(n.id).Parent()
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}

	return New(n.cache, x)
}

func (n *N) Child(b B) *N {
	if !b.IsValid() {
		return nil
	}
	var x cache.ID
	if b == BLeft {
		x = n.cache.GetOrDie(n.id).Left()
	} else {
		x = n.cache.GetOrDie(n.id).Right()
	}

	if _, ok := n.cache.Get(x); !ok {
		return nil
	}
	return New(n.cache, x)
}

func (n *N) Right() *N { return n.Child(BRight) }
func (n *N) Left() *N  { return n.Child(BLeft) }
