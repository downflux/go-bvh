package op

import (
	"github.com/downflux/go-bvh/x/bvh/balance"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/candidate"
	"github.com/downflux/go-bvh/x/internal/cache/op/split"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func Insert(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cid.ID, x id.ID, tolerance float64) (node.N, map[id.ID]cid.ID) {
	updates := make(map[id.ID]cid.ID, c.LeafSize())

	var n node.N
	var ok bool
	if n, ok = c.Get(root); !ok {
		n = c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, false))
	} else {
		n = candidate.BrianNoyama(c, n, data[x])
	}

	if n.IsFull() {
		n.Leaves()[x] = struct{}{}
		m := unsafe.Expand(c, n)
		split.GuttmanLinear(c, data, n, m)

		for y := range m.Leaves() {
			updates[y] = m.ID()
		}
		if _, ok := updates[x]; !ok {
			updates[x] = n.ID()
		}
	} else {
		n.Leaves()[x] = struct{}{}

		updates[x] = n.ID()
	}

	var r node.N
	for m := n; m != nil; m = m.Parent() {
		if !m.IsLeaf() {
			node.SetAABB(m, data, tolerance)
			node.SetHeight(m)

			m = balance.B(m, data, tolerance)

			node.SetAABB(m, data, tolerance)
			node.SetHeight(m)
		}

		if m.Parent() == nil {
			r = m
		}
	}

	return r, updates
}
