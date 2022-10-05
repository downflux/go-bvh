package node

import (
	"fmt"

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
	cn := c.GetOrDie(x)
	if p, ok := c.Get(cn.Parent()); ok {
		n.parent = p.ID()
		if x == p.Left() {
			n.branch = BLeft
		} else {
			n.branch = BRight
		}
	}

	// Ensure either the node is a leaf or both children are valid.
	_, cl := c.Get(cn.Left())
	_, cr := c.Get(cn.Right())
	if cl != cr {
		panic(fmt.Sprintf("invalid node %v: dangling child node", x))
	}

	return n
}

func (n *N) Branch() B    { return n.branch }
func (n *N) ID() cache.ID { return n.id }

func (n *N) IsRoot() bool {
	_, ok := n.cache.Get(n.parent)
	return !ok
}

func (n *N) IsLeaf() bool {
	l := n.Left()
	if (l == nil) != (n.Right() == nil) {
		panic(fmt.Sprintf("invalid node %v: dangling child", n.ID()))
	}
	return l == nil
}

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
