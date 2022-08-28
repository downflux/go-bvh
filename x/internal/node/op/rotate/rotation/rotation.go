package rotation

import (
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate/subtree"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

type R struct {
	B, C, F, G *node.N
}

func Generate(n *node.N) []R {
	t := subtree.New(n)

	var rs []R
	if !t.A.IsLeaf() { // t.B and t.C are non-nil
		if !t.B.IsLeaf() { // t.F and t.G are non-nil
			rs = append(rs, R{
				B: t.B, C: t.C, F: t.F, G: t.G,
			}, R{
				B: t.B, C: t.C, F: t.F, G: t.G,
			})
		}
		if !t.C.IsLeaf() { // t.E and t.D are non-nil
			rs = append(rs, R{
				B: t.C, C: t.B, F: t.E, G: t.D,
			}, R{
				B: t.C, C: t.B, F: t.D, G: t.E,
			})
		}
	}

	return rs
}

// Optimal finds the optimal rotation for a given ancester node n. The returned
// rotation object may be empty i.e. R{}, which indicates the existing rotation
// is already optimal.
func Optimal(n *node.N) R {
	var optimal R

	var h float64
	if !n.IsLeaf() {
		// The ancester node n will have the same AABB volume, so we
		// will need to check the decomposed volume of the children
		// instead.
		h = heuristic.H(n.Left().AABB()) + heuristic.H(n.Right().AABB())
	}

	for _, r := range Generate(n) {
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
			optimal = r
		}
	}

	return optimal
}
