// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

type D[T comparable] struct {
	ID   T
	AABB hyperrectangle.R
}

type N[T comparable] struct {
	Parent *N[T]
	Left   *N[T]
	Right  *N[T]

	Data             []D[T]
	AABBCacheIsValid bool
	AABBCache        hyperrectangle.R
}

func (n *N[T]) Leaf() bool { return len(n.Data) > 0 }
func (n *N[T]) AABB() hyperrectangle.R {
	if n.AABBCacheIsValid {
		return n.AABBCache
	}

	n.AABBCacheIsValid = true
	if n.Leaf() {
		rs := make([]hyperrectangle.R, len(n.Data))
		for i := 0; i < len(n.Data); i++ {
			rs[i] = n.Data[i].AABB
		}
		n.AABBCache = bhr.AABB(rs)
	} else {
		n.AABBCache = bhr.Union(n.Left.AABB(), n.Right.AABB())
	}

	return n.AABBCache
}
