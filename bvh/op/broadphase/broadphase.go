package broadphase

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

// BroadPhase checks a BVH tree for the query rectangle and returns a list of
// objects which touch the query AABB.
func BroadPhase(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, q hyperrectangle.R) []id.ID {
	n, ok := c.Get(root)
	if !ok {
		return []id.ID{}
	}

	open := make([]node.N, 0, 128)
	open = append(open, n)

	candidates := make([]id.ID, 0, 128)

	var m node.N
	for len(open) > 0 {
		m, open = open[len(open)-1], open[:len(open)-1]
		if m.IsLeaf() {
			for x := range m.Leaves() {
				candidates = append(candidates, x)
			}
		} else {
			l, r := m.Left(), m.Right()
			if !hyperrectangle.Disjoint(q, l.AABB().R()) {
				open = append(open, l)
			}
			if !hyperrectangle.Disjoint(q, r.AABB().R()) {
				open = append(open, r)
			}
		}
	}

	ids := make([]id.ID, 0, len(candidates))
	for _, x := range candidates {
		if !hyperrectangle.Disjoint(q, data[x]) {
			ids = append(ids, x)
		}
	}

	return ids
}
