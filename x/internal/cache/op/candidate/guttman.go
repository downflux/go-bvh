package candidate

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// Guttman finds an existing candidate leaf node to insert the AABB. This is
// based on Guttman 1984. This is the method referenced in Catto 2019 for
// finding a candidate sibling node; as the Box2D BVH implementation does not
// support multi-AABB leaves, the candidate node will always be split (and thus,
// is instead called a "sibling" instead).
func Guttman(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, aabb.Min().Dimension())),
		vector.V(make([]float64, aabb.Min().Dimension())),
	).M()

	var m node.N
	for m = n; !m.IsLeaf(); {
		buf.Copy(aabb)
		buf.Union(m.Left().AABB().R())
		lh := heuristic.H(buf.R())
		dlh := lh - heuristic.H(m.Left().AABB().R())

		buf.Copy(aabb)
		buf.Union(m.Right().AABB().R())
		rh := heuristic.H(buf.R())
		drh := rh - heuristic.H(m.Right().AABB().R())

		// Choose an appropriate node to search. Use the node with the
		// least change in cost; in the case of a tie, use the node with
		// the smallest resultant area.
		if dlh < drh || epsilon.Within(dlh, drh) && lh < rh {
			m = m.Left()
		} else {
			m = m.Right()
		}
	}

	return m
}
