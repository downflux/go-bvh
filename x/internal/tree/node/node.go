package node

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/cache"
)

type N struct {
	cache  *cache.C
	id     cache.ID
	parent cache.ID
	branch cache.B
}

func New(c *cache.C, x cache.ID) *N {
	n := &N{}
	n.load(c, x)
	return n
}

func (n *N) load(c *cache.C, x cache.ID) {
	n.cache = c
	n.id = x
	n.parent = cache.IDInvalid
	n.branch = cache.BInvalid

	cn := n.cache.GetOrDie(x)
	if p, ok := n.cache.Get(cn.Parent()); ok {
		n.parent = p.ID()
		if x == p.Left() {
			n.branch = cache.BLeft
		} else {
			n.branch = cache.BRight
		}
	}
}

func (n *N) Branch() cache.B { return n.branch }
func (n *N) ID() cache.ID    { return n.id }

func (n *N) IsRoot() bool {
	_, ok := n.cache.Get(n.parent)
	return !ok
}

func (n *N) IsLeaf() bool {
	l := n.Left()
	if (l == nil) != (n.Right() == nil) {
		panic(fmt.Sprintf("invalid node %v: dangling child", n.id))
	}
	return l == nil
}

func (n *N) IterParent() *N {
	n.load(n.cache, n.cache.GetOrDie(n.id).Parent())
	return n
}

func (n *N) IterChild(b cache.B) *N {
	if !b.IsValid() {
		panic(fmt.Sprintf("cannot iterate on invalid branch %v", b))
	}

	x := n.cache.GetOrDie(n.id).Child(b)
	n.load(n.cache, x)

	return n
}

func (n *N) IterLeft() *N  { return n.IterChild(cache.BLeft) }
func (n *N) IterRight() *N { return n.IterChild(cache.BRight) }

func (n *N) Parent() *N {
	x := n.cache.GetOrDie(n.id).Parent()
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}

	return New(n.cache, x)
}

func (n *N) Child(b cache.B) *N {
	if !b.IsValid() {
		return nil
	}

	x := n.cache.GetOrDie(n.id).Child(b)
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}
	return New(n.cache, x)
}

func (n *N) Right() *N { return n.Child(cache.BRight) }
func (n *N) Left() *N  { return n.Child(cache.BLeft) }
