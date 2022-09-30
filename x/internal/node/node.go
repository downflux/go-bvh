package node

import (
	"github.com/downflux/go-bvh/x/internal/cache"
)

type N struct {
	id     cache.ID
	parent cache.ID
	left   cache.ID
	right  cache.ID
}

type O struct {
	Parent cache.ID
	Left   cache.ID
	Right  cache.ID
}

func New(c *cache.C[*N], o O) *N {
	n := &N{
		parent: o.Parent,
		left:   o.Left,
		right:  o.Right,
	}
	x := c.Insert(n)

	n.id = x
	return n
}

func (n *N) ID() cache.ID { return n.id }

func (n *N) IsRoot(c *cache.C[*N]) bool { return n.Parent(c) == nil }
func (n *N) IsLeaf(c *cache.C[*N]) bool { return n.Left() == nil || n.Right() == nil }

func (n *N) Parent(c *cache.C[*N]) *N {
	m, ok := c.Get(n.parent)
	if !ok {
		return nil
	}

	return m
}

func (n *N) Left(c *cache.C[*N]) *N {
	m, ok := c.Get(n.left)
	if !ok {
		return nil
	}

	return m
}

func (n *N) Right(c *cache.C[*N]) *N {
	m, ok := c.Get(n.right)
	if !ok {
		return nil
	}

	return m
}
