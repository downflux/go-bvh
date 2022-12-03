package metrics

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	AilaC C = C{
		Internal: 1.2,
		Leaf:     1,
		Object:   0,
	}

	// The internally cached node heuristic may differ from the heuristic as
	// strictly defined by the MacDonald and Booth paper. We are making the
	// tree heuristic here independent of the cached heuristic in order to
	// have comparable tree costs across different node heuristic
	// strategies.
	SA = hyperrectangle.SA
)

func SAH(n node.N) float64 { return AilaC.SAH(n) }

type C struct {
	Internal float64
	Leaf     float64
	Object   float64
}

// SAH returns the surface area heuristic as defined in MacDonald and Booth
// 1990.
//
// The total heuristic is comprised of three separate components -- the cost of
// the internal nodes, the cost of the leaves, and the cost of testing for
// intersections. We use track these via ci, cl, and co respectively.
//
// Per Aila et al., a "normal" SAH value is around 100.
//
// N.B.: SAH assumes the local subtree has up-to-date AABB bounding boxes and
// heuristic caches.
func (c C) SAH(n node.N) float64 {
	var ci, cl, co float64
	util.PreOrder(n, func(n node.N) {
		if !n.IsLeaf() {
			ci += SA(n.AABB().R())
		} else {
			cl += SA(n.AABB().R())
			co += SA(n.AABB().R()) * float64(len(n.Leaves()))
		}
	})
	return (c.Internal*ci + c.Leaf*cl + c.Object*co) / SA(n.AABB().R())
}
