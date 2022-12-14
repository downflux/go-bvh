package candidate

import (
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/heuristic"
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
		l, r := m.Left(), m.Right()

		buf.Copy(aabb)
		buf.Union(l.AABB().R())

		lh := heuristic.H(buf.R())
		dlh := lh - l.Heuristic()

		buf.Copy(aabb)
		buf.Union(r.AABB().R())

		rh := heuristic.H(buf.R())
		drh := rh - r.Heuristic()

		// Choose an appropriate node to search. Use the node with the
		// least change in cost; in the case of a tie, use the node with
		// the smallest resultant area.
		if dlh < drh || epsilon.Within(dlh, drh) && lh < rh {
			m = l
		} else {
			m = r
		}
	}

	return m
}
