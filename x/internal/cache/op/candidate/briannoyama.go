package candidate

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// BrianNoyama returns a leaf node to be used for containing the AABB. See the
// briannoyama implementation for more information.
func BrianNoyama(c *cache.C, n node.N, aabb hyperrectangle.R) node.N {
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	var m node.N
	for m = n; !m.IsLeaf(); {
		buf.Copy(aabb)
		buf.Union(m.Left().AABB().R())

		lh := heuristic.H(buf.R())

		buf.Copy(aabb)
		buf.Union(m.Right().AABB().R())

		rh := heuristic.H(buf.R())

		if lh < rh {
			m = m.Left()
		} else {
			m = m.Right()
		}
	}

	return m
}
