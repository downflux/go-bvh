package node

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/cache"
)

type Branch int

func (b Branch) Sibling() Branch { return b ^ 1 }
func (b Branch) IsValid() bool {
	return b == BranchLeft || b == BranchRight
}

const (
	BranchLeft Branch = iota
	BranchRight

	BranchInvalid
)

type N struct {
	id       cache.ID
	parent   cache.ID
	children [2]cache.ID
}

type O struct {
	Parent cache.ID
	Left   cache.ID
	Right  cache.ID
}

func New(c *cache.C[*N], o O) *N {
	n := &N{
		parent:   o.Parent,
		children: [2]cache.ID{o.Left, o.Right},
	}
	x := c.Insert(n)

	n.id = x
	return n
}

func (n *N) ID() cache.ID { return n.id }

func (n *N) IsRoot(c *cache.C[*N]) bool { return n.Parent(c) == nil }
func (n *N) IsLeaf(c *cache.C[*N]) bool { return n.Left(c) == nil || n.Right(c) == nil }

func (n *N) Parent(c *cache.C[*N]) *N {
	m, ok := c.Get(n.parent)
	if !ok {
		return nil
	}

	return m
}

func (n *N) Branch(child cache.ID) Branch {
	if n.children[BranchLeft] == child {
		return BranchLeft
	}
	if n.children[BranchRight] == child {
		return BranchRight
	}
	return BranchInvalid
}

// Child is a convenience function for programatic tree explorations -- instead
// of calling
//
//	n.Left(c)
//
// we can instead call
//
//	n.Child[c, BranchLeft]
func (n *N) Child(c *cache.C[*N], b Branch) *N {
	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch option %v", b))
	}

	m, ok := c.Get(n.children[b])
	if !ok {
		return nil
	}

	return m
}

func (n *N) Left(c *cache.C[*N]) *N  { return n.Child(c, BranchLeft) }
func (n *N) Right(c *cache.C[*N]) *N { return n.Child(c, BranchRight) }
