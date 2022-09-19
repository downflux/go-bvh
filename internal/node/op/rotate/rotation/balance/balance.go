// Package balance generates a rotation based on tree height. This is the
// current Box2D rotation strategy (github.com/erincatto/box2d) and notably
// varies from the Catto 2019 slides, which uses a rotation heuristic based on
// node surface area.
package balance

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/rotation"
	"github.com/downflux/go-bvh/internal/node/op/rotate/subtree"
)

// Generate creates a rotation based on tree height.
//
// N.B.: The reference Box2D implementation rotates a subtree from
//
//	  A
//	 / \
//	B   C
//	   / \
//	  F   G
//
// to
//
//	  C
//	 / \
//	F   A
//	   / \
//	  B   G
//
// Where in the unrotated tree, the following conditions are met --
//
// 1. C.Height() > B.Height() and
// 2. F.Height() > G.Height()
//
// Note that since both A and C are internal nodes, this is identical to
//
//	  A
//	 / \
//	F   C
//	   / \
//	  B   G
//
// We will instead use this latter rotation notation going forward to keep
// consistent with the swap operation as described in the Catto 2019 slides.
func Generate(n *node.N) rotation.R {
	t := subtree.New(n)
	r := rotation.R{}

	if t.A.IsLeaf() || t.A.Height() < 2 {
		return r
	}

	if t.B.Height() == t.C.Height() {
		return r
	} else if t.C.Height() > t.B.Height() {
		r.B = t.B
		r.C = t.C
		if t.F.Height() > t.G.Height() {
			r.F = t.F
			r.G = t.G
		} else {
			r.F = t.G
			r.G = t.F
		}
	} else {
		r.B = t.C
		r.C = t.B
		if t.D.Height() > t.E.Height() {
			r.F = t.D
			r.G = t.E
		} else {
			r.F = t.E
			r.G = t.D
		}
	}
	return r
}
