package op

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/tree/node"
)

func IsAncestor(c *cache.C, n cache.ID, m cache.ID) bool {
	for p, ok := c.Get(m); !ok; p, ok = c.Get(p.ID()) {
		if p.ID() == n {
			return true
		}
	}
	return false
}

// Swap moves two nodes in the same tree. The two nodes must not be direct
// ancestors of one another.
//
//	  A
//	 / \
//	n   B
//	   / \
//	  m   C
//
// to
//
//	  A
//	 / \
//	m   B
//	   / \
//	  n   C
func Swap(c *cache.C, from cache.ID, to cache.ID) {
	if IsAncestor(c, from, to) || IsAncestor(c, to, from) {
		panic(fmt.Sprintf("cannot swap direct ancestor nodes %v and %v", from, to))
	}

	n, m := node.New(c, from), node.New(c, to)
	// p and q will always be valid nodes, since n and m can never be the
	// root node.
	p, q := n.Parent(), m.Parent()

	c.GetOrDie(n.ID()).SetParent(q.ID())
	c.GetOrDie(m.ID()).SetParent(p.ID())
	c.GetOrDie(p.ID()).SetChild(n.Branch(), m.ID())
	c.GetOrDie(q.ID()).SetChild(m.Branch(), n.ID())
}
