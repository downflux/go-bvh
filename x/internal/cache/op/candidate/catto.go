package candidate

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// Catto finds or creates a leaf node which will result in a minimal increase in
// the heuristic.
//
// This is the method utilized in the Box2D implementation.
func Catto(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	aabbh := heuristic.H(aabb)

	for n := n; !n.IsLeaf(); {
		buf.Copy(aabb)
		buf.Union(n.AABB().R())
		combined := heuristic.H(buf.R())

		h := 2 * combined
		inherited := 2 * (combined - aabbh)

		var lh, rh float64

		buf.Copy(aabb)
		buf.Union(n.Left().AABB().R())
		if n.Left().IsLeaf() {
			lh = heuristic.H(buf.R()) + inherited
		} else {
			lh = heuristic.H(buf.R()) - heuristic.H(n.Left().AABB().R()) + inherited
		}

		buf.Copy(aabb)
		buf.Union(n.Right().AABB().R())
		if n.Right().IsLeaf() {
			rh = heuristic.H(buf.R()) + inherited
		} else {
			rh = heuristic.H(buf.R()) - heuristic.H(n.Right().AABB().R()) + inherited
		}

		if h < lh && h < rh {
			break
		}

		if lh < rh {
			n = n.Left()
		} else {
			n = n.Right()
		}
	}

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
	if !n.IsLeaf() {
		q := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, false))
		m := c.GetOrDie(c.Insert(q.ID(), cid.IDInvalid, cid.IDInvalid, false))

		q.SetLeft(n.ID())
		q.SetRight(m.ID())

		if !n.IsRoot() {
			p := n.Parent()

			q.SetParent(p.ID())
			p.SetChild(p.Branch(n.ID()), q.ID())
			n.SetParent(q.ID())
		}

		n = m
	}

	return n
}
