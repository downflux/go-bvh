package node

import (
	"math"

	"github.com/downflux/go-bvh/y/id"
	"github.com/downflux/go-bvh/y/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// Update makes the node up-to-date. This function assumes the child nodes are
// also up-to-date.
func Update(n *N, data map[id.ID]hyperrectangle.R) {
	if n.IsLeaf() {
		n.heightCache = 0
	} else {
		h := 0
		for _, n := range n.children {
			if g := n.Height(); g > h {
				h = g
			}
		}

		n.heightCache = h + 1
	}

	aabb := n.AABB().M()
	var init bool
	if !n.IsLeaf() {
		for _, m := range n.children {
			if !init {
				init = true
				aabb.Copy(m.AABB())
			} else {
				aabb.Union(m.AABB())
			}
		}
	} else {
		aabb.Copy(data[n.leaf])

		f := math.Pow(n.tolerance, 1/float64(n.k))

		tmin, tmax := aabb.Min(), aabb.Max()
		for i := vector.D(0); i < n.k; i++ {
			delta := tmax[i] - tmin[i]
			offset := delta * (f - 1) / 2
			tmin[i] = tmin[i] - offset
			tmax[i] = tmax[i] + offset
		}
	}
	heuristic.N(aabb.R())
}
