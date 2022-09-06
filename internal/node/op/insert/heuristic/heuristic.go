package heuristic

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// B is the lower bound of the insertion cost for the node.
func B(n *node.N, aabb hyperrectangle.R) float64 {
	return inherited(n, aabb) + bhr.V(aabb)
}

// F is the actual cost for inserting the AABB into the node.
func F(n *node.N, aabb hyperrectangle.R) float64 {
	return inherited(n, aabb) + direct(n, aabb)
}

// direct calculates (per Catto 2019) "the surface area of the new internal node
// that will be created for the siblings."
func direct(n *node.N, aabb hyperrectangle.R) float64 {
	return bhr.V(bhr.Union(n.AABB(), aabb))
}

// inherited calculates (per Catto 2019) "the increased surface area caused by
// refitting the ancestor's boxes."
func inherited(n *node.N, aabb hyperrectangle.R) float64 {
	if n.IsRoot() {
		return 0
	}

	return bhr.V(bhr.Union(n.AABB(), aabb)) - bhr.V(aabb) + inherited(n.Parent(), aabb)
}
