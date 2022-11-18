package sibling

import (
	"math"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-pq/pq"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// FindBox2D implements a find sibling algorithm as per the Box2D
// implementation. See
// https://github.com/erincatto/box2d/blob/9dc24a6fd4f32442c4bcf80791de47a0a7d25afb/src/collision/b2_dynamic_tree.cpp#L197
// for more information.
func FindBox2D(c *cache.C, x cid.ID, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(vector.V(make([]float64, c.K())), vector.V(make([]float64, c.K()))).M()

	var n node.N
	for n = c.GetOrDie(x); !n.IsLeaf(); {
		h := heuristic.H(n.AABB().R())
		nlh := heuristic.H(n.Left().AABB().R())
		nrh := heuristic.H(n.Right().AABB().R())

		buf.Copy(n.AABB().R())
		buf.Union(aabb)

		ch := heuristic.H(buf.R())

		// Find the cost of creating a new parent for this node and the new leaf.
		direct := 2 * ch

		// Find the minimum cost of pushing the AABB further down the
		// tree.
		inherited := 2 * (ch - h)

		var lh, rh float64

		buf.Copy(n.Left().AABB().R())
		buf.Union(aabb)
		if n.Left().IsLeaf() {
			lh = inherited + heuristic.H(buf.R())
		} else {
			lh = inherited + (nlh - heuristic.H(buf.R()))
		}

		buf.Copy(n.Right().AABB().R())
		buf.Union(aabb)
		if n.Left().IsLeaf() {
			rh = inherited + heuristic.H(buf.R())
		} else {
			rh = inherited + (nrh - heuristic.H(buf.R()))
		}

		if direct < lh && direct < rh {
			break
		}

		if lh < rh {
			n = n.Left()
		} else {
			n = n.Right()
		}
	}

	return n
}

// Find traverses down a tree and gets the insertion sibling candidate, as per
// Bittner et al.  2013.  This is the original algorithm described in the Catto
// 2019 slides, and aims to decrease the overall SAH value of the resultant
// tree.
func Find(c *cache.C, x cid.ID, aabb hyperrectangle.R) node.N {
	l := heuristic.H(aabb)

	n := c.GetOrDie(x)

	// The priority queue q weight is the induced cost, i.e. the cost of
	// traveling to the associated node. We pop the lowest-cost nodes first,
	// as those are probably the nodes with the highest chance of being
	// "optimal", which may allow us to skip a vast majority of later node
	// expansions.
	q := pq.New[node.N](0, pq.PMin)
	q.Push(n, 0)

	var opt node.N
	// h tracks the minimum heuristic penalty which would be incurred if the
	// input AABB is inserted here.
	h := math.Inf(1)

	buf := hyperrectangle.New(
		vector.V(make([]float64, aabb.Min().Dimension())),
		vector.V(make([]float64, aabb.Min().Dimension())),
	).M()
	for q.Len() > 0 {
		m, induced := q.Pop()

		// Prune the current node and any other nodes in the queue if
		// the mininum incurred penalty is greater than the current
		// optimal lower bound, as l is the highest allowable additional
		// penalty on top of the current induced cost.
		if induced+l >= h {
			break
		}

		// We define the direct cost to be the cost of merging the input
		// AABB with the current node -- that is, the cost of creating a
		// new node which contains both the AABB and the current node.
		// Note that this heuristic does not care about how much we are
		// expanding the node itself, but rather just the cost of the
		// end result. The penalty due to how much we expand the node
		// is included in the induced cost when adding the node to the
		// queue.
		buf.Copy(m.AABB().R())
		buf.Union(aabb)

		direct := heuristic.H(buf.R())
		cost := induced + direct

		if cost < h {
			h = cost
			opt = m
		}

		// Calculate the induced cost of the children. Here, we take
		// into account how much we will have to expand the current node
		// m to fully accomodate the input AABB. That is, this allows us
		// to prefer nodes which contains more (or are closer to) the
		// input AABB.
		if !m.IsLeaf() {
			induced = cost - heuristic.H(m.AABB().R())
			if bound := induced + l; bound < h {
				q.Push(m.Left(), induced)
				q.Push(m.Right(), induced)
			}
		}
	}
	return opt
}
