package aabb

import (
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/rotation"
	"github.com/downflux/go-bvh/internal/node/op/rotate/subtree"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

func list(n *node.N) []rotation.R {
	t := subtree.New(n)

	var rs []rotation.R
	if !t.A.IsLeaf() { // t.B and t.C are non-nil
		if !t.B.IsLeaf() { // t.E and t.D are non-nil
			rs = append(rs, rotation.R{
				B: t.C, C: t.B, F: t.D, G: t.E,
			}, rotation.R{
				B: t.C, C: t.B, F: t.E, G: t.D,
			})
		}
		if !t.C.IsLeaf() { // t.F and t.G are non-nil
			rs = append(rs, rotation.R{
				B: t.B, C: t.C, F: t.F, G: t.G,
			}, rotation.R{
				B: t.B, C: t.C, F: t.G, G: t.F,
			})
		}
	}

	return rs
}

// Generate finds the optimal rotation for a given ancester node n. The returned
// rotation object may be empty i.e. R{}, which indicates the existing rotation
// is already optimal.
func Generate(n *node.N) rotation.R {
	if n.IsLeaf() {
		return rotation.R{}
	}

	// The ancester node n will have the same AABB volume, so we
	// will need to check the decomposed volume of the children
	// instead.
	h := heuristic.H(n.Left().AABB()) + heuristic.H(n.Right().AABB())
	var optimal rotation.R

	for _, r := range list(n) {
		// Calculate the decomposed volume of the simulated rotation F
		// and C'.
		if g := heuristic.H(r.F.AABB()) + heuristic.H(
			// Compute the AABB volume for a simulated merge of the
			// B and G nodes into C'.
			bhr.Union(
				r.B.AABB(),
				r.G.AABB(),
			),
		); g < h {
			h = g
			optimal = r
		}
	}

	return optimal
}
