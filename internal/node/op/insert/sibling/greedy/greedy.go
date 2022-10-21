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
)

const epsilon = 1e-20

type candidate struct {
	n *node.N

	// ci is the cached inherited cost, starting from the root.
	ci float64
}

func Execute(n *node.N, aabb hyperrectangle.R) *node.N {
	q := pq.New[*node.N](0, pq.PMin)
	q.Push(n, 0)

	var opt *node.N
	h := math.Inf(0)

	buf := hyperrectangle.New(
		make([]float64, aabb.Min().Dimension()),
		make([]float64, aabb.Min().Dimension()),
	).M()

	for q.Len() > 0 {
		m, ci := q.Pop()
		if ci+heuristic.H(aabb) >= h {
			break
		}

		buf.Copy(m.AABB())
		buf.Union(aabb)
		cd := heuristic.H(buf.R())
		c := ci + cd
		if c < h {
			h = c
			opt = m
		}

		ci = c - heuristic.H(m.AABB())
		if !m.IsLeaf() && ci+heuristic.H(aabb) < h {
			q.Push(m.Left(), ci)
			q.Push(m.Right(), ci)
		}
	}
	return opt
}
