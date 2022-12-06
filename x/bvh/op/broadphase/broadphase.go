package broadphase

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/stack"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// BroadPhase checks a BVH tree for the query rectangle and returns a list of
// objects which touch the query AABB.
func BroadPhase(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, q hyperrectangle.R) []id.ID {
	n, ok := c.Get(root)
	if !ok {
		return []id.ID{}
	}

	open := stack.New(make([]node.N, 0, 128))
	open.Push(n)

	ids := make([]id.ID, 0, 128)

	for m, ok := open.Pop(); ok; m, ok = open.Pop() {
		if m.IsLeaf() {
			for x := range m.Leaves() {
				if !hyperrectangle.Disjoint(q, data[x]) {
					ids = append(ids, x)
				}
			}
		} else {
			l, r := m.Left(), m.Right()
			if !hyperrectangle.Disjoint(q, l.AABB().R()) {
				open.Push(l)
			}
			if !hyperrectangle.Disjoint(q, r.AABB().R()) {
				open.Push(r)
			}
		}
	}

	return ids
}
