package unsafe

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// Expand creates a new node with s as its sibling. This will re-link any
// existing parents or siblings of s and ensure that the generated cache is
// still valid.
//
// The input node s must not be nil.
//
// The input node s is a node within the cache.
//
//	  Q
//	 / \
//	N   T
//
// to
//
//	    Q
//	   / \
//	  P   T
//	 / \
//	N   M
func Expand(c *cache.C, n node.N) node.N {
	p := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, false))
	m := c.GetOrDie(c.Insert(p.ID(), cid.IDInvalid, cid.IDInvalid, false))

	p.SetLeft(n.ID())
	p.SetRight(m.ID())

	if !n.IsRoot() {
		q := n.Parent()

		p.SetParent(q.ID())
		q.SetChild(q.Branch(n.ID()), p.ID())
		n.SetParent(p.ID())
	}

	return m
}
