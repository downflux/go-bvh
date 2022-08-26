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
		n.AABBCache = n.Data[0].AABB
		for _, p := range n.Data[1:] {
			n.AABBCache = bhr.Union(n.AABBCache, p.AABB)
		}
	} else {
		n.AABBCache = bhr.Union(n.Left.AABB(), n.Right.AABB())
	}

	return n.AABBCache
}
