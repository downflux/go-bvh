package op

import (
	"math"

	"github.com/downflux/go-bvh/x/bvh/balance"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/candidate"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// Split partitions a leaf node which may be full and adds some non-zero amount
// of objects into a separate leaf node.
//
// This algorithm is described in Guttman 1984, section 3.5.3.
func Split(c *cache.C, data map[id.ID]hyperrectangle.R, from node.N, to node.N) {
	if c.LeafSize() == 1 {
		for x := range from.Leaves() {
			to.Leaves()[x] = struct{}{}
			delete(from.Leaves(), x)
			return
		}
	}

	node.SetAABB(from, data, 1)
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()
	buf.Copy(from.AABB().R())

	// Reset the leaves within the source node, as data will be copied into
	// here.
	nodes := make([]id.ID, 0, len(from.Leaves()))
	for x := range from.Leaves() {
		nodes = append(nodes, x)
		delete(from.Leaves(), x)
	}

	separation := math.Inf(-1)
	var r id.ID
	var l id.ID

	// Pick node seeds -- one AABB will go into the source node, and one the
	// destination.
	for i := vector.D(0); i < c.K(); i++ {
		var left id.ID
		var right id.ID

		min := math.Inf(1)
		max := math.Inf(-1)

		for _, x := range nodes {
			aabb := data[x]

			if aabb.Min().X(i) > max {
				right = x
				max = aabb.Min().X(i)
			}
			if aabb.Max().X(i) < min && right != x {
				left = x
				min = aabb.Max().X(i)
			}
		}

		if s := (min - max) / (buf.Max().X(i) - buf.Min().X(i)); s > separation {
			separation = s
			r = right
			l = left
		}
	}

	from.Leaves()[l] = struct{}{}
	to.Leaves()[r] = struct{}{}

	// Set AABBs based on the smallest net increase in node size.
	for _, x := range nodes {
		if x == l || x == r {
			continue
		}

		node.SetAABB(from, data, 1)
		buf.Copy(from.AABB().R())
		buf.Union(data[x])
		dlh := heuristic.H(buf.R()) - heuristic.H(from.AABB().R())

		node.SetAABB(to, data, 1)
		buf.Copy(to.AABB().R())
		buf.Union(data[x])
		drh := heuristic.H(buf.R()) - heuristic.H(to.AABB().R())

		if dlh < drh {
			from.Leaves()[x] = struct{}{}
		} else {
			to.Leaves()[x] = struct{}{}
		}
	}
}

func Insert(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cid.ID, x id.ID, tolerance float64) (cid.ID, map[id.ID]cid.ID) {
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
		Split(c, data, n, m)

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
		}

		if m.Parent() == nil {
			r = m
		}
	}

	return r.ID(), updates
}
