package tree

import (
	"fmt"
	"math"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-pq/pq"
)

// sibling finds an insertion sibling candidate, as per Bittner et al.  2013.
// This is the original algorithm described in the Catto 2019 slides, and aims
// to decrease the overall SAH value of the resultant tree.
func sibling(c *cache.C, x cache.ID, aabb hyperrectangle.R) cache.ID {
	l := heuristic.H(aabb)

	n := c.GetOrDie(x)

	// The priority queue q weight is the induced cost, i.e. the cost of
	// traveling to the associated node. We pop the lowest-cost nodes first,
	// as those are probably the nodes with the highest chance of being
	// "optimal", which may allow us to skip a vast majority of later node
	// expansions.
	q := pq.New[*cache.N](0, pq.PMin)
	q.Push(n, 0)

	var opt *cache.N
	// h tracks the minimum heuristic penalty which would be incurred if the
	// input AABB is inserted here.
	h := math.Inf(1)

	buf := hyperrectangle.New(
		vector.V(make([]float64, aabb.Min().Dimension())),
		vector.V(make([]float64, aabb.Min().Dimension())),
	).M()

	for q.Len() > 0 {
		m, g := q.Pop()
		fmt.Printf("DEBUG: m.ID() == %v, g == %v, h == %v\n", m.ID(), g, h)

		// Prune the current node and any other nodes in the queue if
		// the mininum incurred penalty is greater than the current
		// optimal lower bound, as l is the highest allowable additional
		// penalty on top of the current induced cost.
		if g+l >= h {
			break
		}

		// We define the direct cost to be the cost of merging the input
		// AABB with the current node -- that is, the cost of creating a
		// new node which contains both the AABB and the current node.
		buf.Copy(m.AABB().R())
		buf.Union(aabb)
		d := heuristic.H(buf.R())

		if f := g + d; f < h {
			h = f
			opt = m
		}

		f := g - heuristic.H(m.AABB().R())
		if !m.IsLeaf() && f+l < h {
			q.Push(m.Left(), f)
			q.Push(m.Right(), f)
		}
	}
	if !opt.IsAllocated() {
		panic("cannot find valid insertion sibling candidate")
	}
	return opt.ID()
}
