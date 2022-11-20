package split

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// S will move objects from an overly-full source node n into an empty node m.
// This function will not validate the leaf nodes (i.e. set AABB).
//
// The source node tracks exactly c.LeafSize() + 1 objects.
type S func(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N, m node.N)
