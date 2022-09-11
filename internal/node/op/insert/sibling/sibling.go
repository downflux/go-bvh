// Package sibling is a faster way to select insertion node candidates. This
// is the method used in the reference Box2D implementation
// (github.com/erincatto/box2d).  This differs from the Catto 2019 slides, which
// recommends using a priority queue -- while this may give us overall a better
// quality tree, the branch and bound method is rather slow.
package sibling

import (
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

func direct(n *node.N, aabb hyperrectangle.R) float64 {
	return heuristic.H(bhr.Union(n.AABB(), aabb))
}

// Execute finds a sibling node at which the input AABB should be inserted. In
// the case that the node is a leaf node, the caller should try to insert
// directly into the node first (i.e. if the leaf has LeafSize > 1).
func Execute(n *node.N, aabb hyperrectangle.R) *node.N {
	if n.IsLeaf() {
		return n
	}

	h := direct(n, aabb)
	c := 2 * h

	inherited := 2 * (h - heuristic.H(n.AABB()))
	l := inherited + direct(n.Left(), aabb)
	r := inherited + direct(n.Right(), aabb)

	// Calculate the lower-bound heuristics for the current node -- here, we
	// assume that the minimum heuristic is if the object is merged directly
	// into the child.
	if !n.Left().IsLeaf() {
		l -= heuristic.H(n.Left().AABB())
	}
	if !n.Right().IsLeaf() {
		r -= heuristic.H(n.Right().AABB())
	}

	if c < l && c < r {
		return n
	}

	if l < r {
		return Execute(n.Left(), aabb)
	}
	return Execute(n.Right(), aabb)
}
