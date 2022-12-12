package candidate

import (
	"math"

	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/op/unsafe"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-pq/pq"
)

// Bittner traverses down a tree and gets the insertion sibling candidate, as per
// Bittner et al.  2013.  This is the original algorithm described in the Catto
// 2019 slides, and aims to decrease the overall SAH value of the resultant
// tree.
func Bittner(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	m := bittnerRO(c, n, aabb)
	if !m.IsLeaf() {
		m = unsafe.Expand(c, m)
	}
	return m
}

func bittnerRO(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	l := heuristic.H(aabb)

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

	for q.Len() > 0 {
		m, induced := q.Pop()

		// Prune the current node and any other nodes in the queue if
		// the mininum incurred penalty is greater than the current
		// optimal lower bound, as l is the highest allowable additional
		// penalty on top of the current induced cost.
		lower := induced + l
		if lower >= h {
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
			induced = cost - m.Heuristic()
			if bound := induced + l; bound < h {
				q.Push(m.Left(), induced)
				q.Push(m.Right(), induced)
			}
		}
	}

	return opt
}
