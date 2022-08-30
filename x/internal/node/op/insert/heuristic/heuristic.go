package heuristic

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

func Estimate(n *node.N, aabb hyperrectangle.R) float64 {
	return inherited(n, aabb) + bhr.V(aabb)
}

func Actual(n *node.N, aabb hyperrectangle.R) float64 {
	return inherited(n, aabb) + direct(n, aabb)
}

func direct(n *node.N, aabb hyperrectangle.R) float64 {
	if n == nil {
		return bhr.V(aabb)
	}

	return bhr.V(bhr.Union(n.AABB(), aabb))
}

func inherited(n *node.N, aabb hyperrectangle.R) float64 {
	if n == nil {
		return 0
	}

	return bhr.V(bhr.Union(n.AABB(), aabb)) - bhr.V(aabb) + inherited(n.Parent(), aabb)
}
