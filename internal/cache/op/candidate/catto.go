package candidate

import (
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/op/unsafe"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// Catto finds or creates a leaf node which will result in a minimal increase in
// the heuristic.
//
// This is the method utilized in the Box2D implementatiom.
func Catto(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	m := cattoRO(c, n, aabb)

	// It is possible that the "optimal" candidate is an internal node; in
	// this case, we need to create a new leaf node to pass back to the
	// caller.
	//
	//    P
	//   /
	//  N
	//
	// to
	//      P
	//     /
	//    Q
	//   / \
	//  N   M
	//
	if !m.IsLeaf() {
		m = unsafe.Expand(c, m)
	}
	return m
}

func cattoRO(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	g := heuristic.H(aabb)

	var m node.N
	for m = n; !m.IsLeaf(); {
		buf.Copy(aabb)
		buf.Union(m.AABB().R())
		combined := heuristic.H(buf.R())

		h := 2 * combined
		inherited := 2 * (combined - g)

		var lh float64
		var rh float64

		buf.Copy(aabb)
		buf.Union(m.Left().AABB().R())
		if m.Left().IsLeaf() {
			lh = heuristic.H(buf.R()) + inherited
		} else {
			lh = heuristic.H(buf.R()) - m.Left().Heuristic() + inherited
		}

		buf.Copy(aabb)
		buf.Union(m.Right().AABB().R())
		if m.Right().IsLeaf() {
			rh = heuristic.H(buf.R()) + inherited
		} else {
			rh = heuristic.H(buf.R()) - m.Right().Heuristic() + inherited
		}

		if h < lh && h < rh {
			break
		}

		if lh < rh {
			m = m.Left()
		} else {
			m = m.Right()
		}
	}

	return m
}
