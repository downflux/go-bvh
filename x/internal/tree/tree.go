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

func (t *T) Insert(x id.ID, aabb hyperrectangle.R) error { return nil }
func (t *T) Update(x id.ID, aabb hyperrectangle.R) error { return nil }
func (t *T) Remove(x id.ID) error                        { return nil }
func (t *T) BroadPhase(q hyperrectangle.R) []id.ID       { return nil }

// When inserting an AABB object, remember we also need to update the leaf node
// n.Data() as well as the caches tracked by the tree struct.
//
// 0. update T.data
// 1. find the best sibling node with the given input AABB;
// 2. add a new cache node, or
//    a. split an existing node;
// 3. update cache node n.Data() with AABB
//    a. expand AABB by some factor for both new and existing nodes
// 4. update T.nodes with new AABB, and
//    a. update other nodes which may have moved during split;
// 5. walk up the cache node (i.e. n.Parent()) and balance the tree
