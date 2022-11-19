package candidate

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// C finds an appropriate leaf node to add an AABB object. This node may need to
// be split after adding the AABB.
type C func(c *cache.C, n node.N, aabb hyperrectangle.R) node.N
