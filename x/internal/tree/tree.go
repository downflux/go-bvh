package tree

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type T struct {
	c *cache.C

	nodes map[id.ID]cache.ID
	data  map[id.ID]hyperrectangle.R
}

type O struct {
	K        vector.D
	LeafSize int
}

func New(o O) *T {
	return &T{
		c: cache.New(cache.O{
			K:        o.K,
			LeafSize: o.LeafSize,
		}),

		nodes: make(map[id.ID]cache.ID, 1024),
		data:  make(map[id.ID]hyperrectangle.R, 1024),
	}
}

func (t *T) K() vector.D { return t.c.K() }

// When inserting a hyperrectangle, we are traccking the leaf size in the tree
// itself, not the cache. So in order to add a new node, we need to
//
// 1. Find the best sibling from the cache with the given AABB.
// 2. Add a new cache node, or split an existing node.
// 3. Update the tree lookup table to move nodes around.
// 3a. Update the node data cache with the AABB.
// 4. Update the cache nodes with new AABBs.
// 5. Walk up the cache node, balancing along the way.
