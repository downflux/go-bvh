// Package singlepass is a faster way to select insertion node candidates. This
// is the method used in the dyntree implementation (github.com/imVexed/dyntree)
// as well as the reference Box2D implementation (github.com/erincatto/box2d).
// This differs from the Catto 2019 slides, which recommends using a priority
// queue -- while this may give us overall a better quality tree, the branch and
// bound method is rather slow.
package singlepass

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// Execute finds a sibling node at which the input AABB should be inserted. In
// the case that the ndoe is a leaf node, the caller should try to insert
// directly into the node first (i.e. if the leaf has LeafSize > 1).
func Execute(n *node.N, aabb hyperrectangle.R) *node.N {
	if n.IsLeaf() {
		return n
	}

	// Calculate the lower-bound heuristics for the current node -- here, we
	// assume that the minimum heuristic is if the object is merged directly
	// into the child.
	l := hyperrectangle.SA(bhr.Union(n.Left().AABB(), aabb)) + hyperrectangle.SA(n.Right().AABB())
	r := hyperrectangle.SA(bhr.Union(n.Right().AABB(), aabb)) + hyperrectangle.SA(n.Left().AABB())

	// Calculate the cost of adding a new leaf node at the current node.
	c := hyperrectangle.SA(n.AABB()) + hyperrectangle.SA(aabb)

	if c < l && c < r {
		return n
	}

	if l < r {
		return Execute(n.Left(), aabb)
	}
	return Execute(n.Right(), aabb)
}
