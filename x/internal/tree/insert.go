package tree

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// insert adds a new AABB into a tree, and returns the root along with the ID of
// the newly inserted node.
//
// The input data cache is a read-only map within the insert function.
func insert(c *cache.C, root cache.ID, data map[id.ID]hyperrectangle.R, x id.ID, aabb hyperrectangle.R) (cache.ID, cache.ID) {
	// t is the new node into which we insert the AABB.
	var t *cache.N

	s := c.GetOrDie(sibling(c, root, aabb))
	if s.IsLeaf() && !s.IsFull() {
		t = s
	} else {
		q := s.Parent()

		p := c.GetOrDie(c.Insert(q.ID(), cache.IDInvalid, cache.IDInvalid, false))
		t = c.GetOrDie(c.Insert(p.ID(), cache.IDInvalid, cache.IDInvalid, false))

		q.SetChild(q.Branch(s.ID()), p.ID())
		s.SetParent(p.ID())

		p.SetLeft(s.ID())
		p.SetRight(t.ID())
	}
	t.Data()[x] = true

	if s.IsLeaf() && s.IsFull() {
		// Move nodes.
	}

	var n *cache.N
	for n = t; n != nil; n = n.Parent() {
		// Refit AABBs.
		if n.IsLeaf() {
			// Merge data nodes.
		} else {
			n.AABB().Copy(n.Left().AABB().R())
			n.AABB().Union(n.Right().AABB().R())
		}
		// Rebalance and refit.
	}

	return n.ID(), t.ID()

	// The newly inserted node t will change the bounding box values of its
	// parents up to the root. We will need to calculate that here, along
	// with rebalancing each node as we traverse up the tree.

	// rebalance ...
}
