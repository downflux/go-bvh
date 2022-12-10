package unsafe

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// Remove will delete a subtree. This will re-link the sibling node and ensure
// the input cache is still valid. This method does not wory about balancing or
// updating the node-specific AABB cache.
//
//	  Q
//	 / \
//	T   P
//	   / \
//	  M   N
//
// to
//
//	  Q
//	 / \
//	T   M
//
// Return the sibling node M.
func Remove(c *cache.C, n node.N) node.N {
	var m node.N

	if !n.IsRoot() {
		p := n.Parent()
		m = p.Child(p.Branch(n.ID()).Sibling())

		m.SetParent(cid.IDInvalid)

		if !p.IsRoot() {
			q := p.Parent()

			q.SetChild(q.Branch(p.ID()), m.ID())
			m.SetParent(q.ID())

			c.Delete(p.ID())
		}
	}

	open := make([]node.N, 0, 128)
	open = append(open, n)

	var s node.N
	// Remove the entire subtree (i.e. propagate downwards).
	for len(open) > 0 {
		s, open = open[len(open)-1], open[:len(open)-1]
		if !s.IsLeaf() {
			open = append(open, s.Left(), s.Right())
		}

		c.Delete(s.ID())
	}

	return m
}
