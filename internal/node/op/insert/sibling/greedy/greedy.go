// Package greedy finds an insertion sibling candidate, as per Bittner et al.
// 2013. This is the original algorithm described in the Catto 2019 slides, and
// aims to decrease the overall SAH value of the resultant tree.
package greedy

import (
	"math"

	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

const epsilon = 1e-20

func inherited(n *node.N, aabb hyperrectangle.R) float64 {
	if n.IsRoot() {
		return 0
	}
	return inherited(n.Parent(), aabb) + heuristic.H(bhr.Union(n.AABB(), aabb)) - heuristic.H(n.AABB())
}

type candidate struct {
	n *node.N
	ci float64
}

func Execute(n *node.N, aabb hyperrectangle.R) *node.N {
	q := pq.New[candidate](0, pq.PMax)
	q.Push(candidate{n: n, ci: 0}, 1 / epsilon)

	var opt *node.N
	h := math.Inf(0)

	for q.Len() > 0 {
		v, _ := q.Pop()
		if v.ci+heuristic.H(aabb) >= h {
			break
		}

		cd := heuristic.H(bhr.Union(v.n.AABB(), aabb))
		c := v.ci + cd
		if c < h {
			h = c
			opt = v.n
		}

		ci := c - heuristic.H(v.n.AABB())
		if !v.n.IsLeaf() && ci+heuristic.H(aabb) < h {
			q.Push(candidate{n: v.n.Left(), ci: ci}, 1 / (ci + epsilon))
			q.Push(candidate{n: v.n.Right(), ci: ci}, 1 / (ci + epsilon))
		}
	}
	return opt
}
