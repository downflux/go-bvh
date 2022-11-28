package candidate

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// DHConnelly searches for a candidate node using the approach as outlined in
// the Rtree.chooseNode function.
func DHConnelly(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	m := dhconnellyRO(c, n, aabb)
	if !m.IsLeaf() {
		m = unsafe.Expand(c, m)
	}
	return m
}

func dhconnellyRO(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	return dhconnellyRecursive(n, aabb, buf)
}

func dhconnellyRecursive(n node.N, aabb hyperrectangle.R, buf hyperrectangle.M) node.N {
	if n.IsLeaf() {
		return n
	}

	buf.Copy(aabb)
	buf.Union(n.Left().AABB().R())

	lh := heuristic.H(buf.R())

	diff := lh - n.Left().Heuristic()
	opt := n.Left()

	buf.Copy(aabb)
	buf.Union(n.Right().AABB().R())

	rh := heuristic.H(buf.R())

	if d := rh - n.Right().Heuristic(); d < diff || (epsilon.Within(d, diff) && n.Left().Height() < opt.Height()) {
		opt = n.Right()
	}

	return dhconnellyRecursive(opt, aabb, buf)
}
