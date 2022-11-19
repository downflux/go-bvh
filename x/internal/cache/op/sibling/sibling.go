package sibling

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// S finds an appropriate leaf node to add an AABB object. This node may need to
// be split after adding the AABB.
type S func(n node.N, aabb hyperrectangle.R) node.N
